// Package handlers содержит HTTP-обработчики для API приложения.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	apperrors "github.com/NailUsmanov/api_organization/internal/errors"
	"github.com/NailUsmanov/api_organization/internal/service"
)

// EmployeeHandler обрабатывает HTTP-запросы, связанные с сотрудниками.
type EmployeeHandler struct {
	empService service.EmployeeService
}

// NewEmployeeHandler создаёт новый экземпляр обработчика сотрудников.
func NewEmployeeHandler(empService service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{empService: empService}
}

// CreateEmployee обрабатывает POST /departments/{id}/employees - создание нового сотрудника в отделе.
func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	deptIDStr := r.PathValue("id")
	deptID, err := strconv.ParseUint(deptIDStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	var req createEmployeeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	emp, err := h.empService.Create(r.Context(), uint(deptID), req.FullName, req.Position, req.HiredAt)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrEmployeeDepartmentNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, apperrors.ErrInvalidFullName),
			errors.Is(err, apperrors.ErrInvalidPosition):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

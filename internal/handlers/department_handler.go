// Package handlers содержит HTTP-обработчики (handlers) для API приложения.
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	apperrors "github.com/NailUsmanov/api_organization/internal/errors"
	"github.com/NailUsmanov/api_organization/internal/service"
)

// DepartmentHandler обрабатывает HTTP-запросы, связанные с подразделениями.
type DepartmentHandler struct {
	depService service.DepartmentService
}

// NewDepartmentHandler создаёт новый экземпляр обработчика отделов.
func NewDepartmentHandler(depService service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{depService: depService}
}

// CreateDepartment обрабатывает POST /departments - создание нового подразделения.
func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var req createDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dept, err := h.depService.Create(r.Context(), req.Name, req.ParentID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrDepartmentNameConflict):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, apperrors.ErrParentNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dept)
}

// GetDepartment обрабатывает GET /departments/{id} - получение информации об отделе.
func (h *DepartmentHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	depth := 1
	if d := r.URL.Query().Get("depth"); d != "" {
		if val, err := strconv.Atoi(d); err == nil && val >= 1 && val <= 5 {
			depth = val
		} else {
			http.Error(w, "depth must be integer between 1 and 5", http.StatusBadRequest)
			return
		}
	}

	includeEmployees := true
	if ie := r.URL.Query().Get("include_employees"); ie != "" {
		if val, err := strconv.ParseBool(ie); err == nil {
			includeEmployees = val
		} else {
			http.Error(w, "invalid include_employees value", http.StatusBadRequest)
			return
		}
	}

	dept, err := h.depService.GetByID(r.Context(), uint(id), depth, includeEmployees)
	if err != nil {
		if errors.Is(err, apperrors.ErrDepartmentNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dept)
}

// UpdateDepartment обрабатывает PATCH /departments/{id} - обновление отдела.
func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	var req updateDepartmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	dept, err := h.depService.Update(r.Context(), uint(id), req.Name, req.ParentID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrDepartmentNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, apperrors.ErrDepartmentNameConflict),
			errors.Is(err, apperrors.ErrSelfParent),
			errors.Is(err, apperrors.ErrCycleDetected):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dept)
}

// DeleteDepartment обрабатывает DELETE /departments/{id} - удаление отдела.
func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode == "" {
		http.Error(w, "mode query parameter is required (cascade or reassign)", http.StatusBadRequest)
		return
	}

	var reassignTo *uint
	if mode == "reassign" {
		reassignStr := r.URL.Query().Get("reassign_to_department_id")
		if reassignStr == "" {
			http.Error(w, "reassign_to_department_id is required for reassign mode", http.StatusBadRequest)
			return
		}
		val, err := strconv.ParseUint(reassignStr, 10, 32)
		if err != nil {
			http.Error(w, "invalid reassign_to_department_id", http.StatusBadRequest)
			return
		}
		reassignTo = new(uint)
		*reassignTo = uint(val)
	}

	err = h.depService.Delete(r.Context(), uint(id), mode, reassignTo)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrDepartmentNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, apperrors.ErrReassignWithChildren):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, apperrors.ErrInvalidMode),
			errors.Is(err, apperrors.ErrReassignTargetRequired),
			errors.Is(err, apperrors.ErrTargetDepartmentNotFound):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// internal/handlers/employee_handler_test.go
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apperrors "github.com/NailUsmanov/api_organization/internal/errors"
	"github.com/NailUsmanov/api_organization/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmployeeService — мок для EmployeeService
type MockEmployeeService struct {
	mock.Mock
}

func (m *MockEmployeeService) Create(ctx context.Context, departmentID uint, fullName, position string, hiredAt *time.Time) (*models.Employee, error) {
	args := m.Called(ctx, departmentID, fullName, position, hiredAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func setupEmployeeTest(t *testing.T) (*MockEmployeeService, *EmployeeHandler, *http.ServeMux) {
	mockSvc := new(MockEmployeeService)
	handler := NewEmployeeHandler(mockSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /departments/{id}/employees", handler.CreateEmployee)

	return mockSvc, handler, mux
}

func TestCreateEmployee_Success(t *testing.T) {
	mockSvc, _, mux := setupEmployeeTest(t)

	hiredAt := time.Now()
	reqBody := createEmployeeRequest{
		FullName: "John Doe",
		Position: "Developer",
		HiredAt:  &hiredAt,
	}
	body, _ := json.Marshal(reqBody)

	expectedEmp := &models.Employee{
		ID:           1,
		DepartmentID: 1,
		FullName:     "John Doe",
		Position:     "Developer",
		HiredAt:      &hiredAt,
		CreatedAt:    time.Now(),
	}

	// Используем mock.MatchedBy для сравнения времени
	mockSvc.On("Create", mock.Anything, uint(1), "John Doe", "Developer", mock.MatchedBy(func(t *time.Time) bool {
		return t != nil && t.Equal(hiredAt)
	})).Return(expectedEmp, nil)

	req := httptest.NewRequest(http.MethodPost, "/departments/1/employees", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp models.Employee
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedEmp.FullName, resp.FullName)
	assert.Equal(t, expectedEmp.Position, resp.Position)

	mockSvc.AssertExpectations(t)
}

func TestCreateEmployee_DepartmentNotFound(t *testing.T) {
	mockSvc, _, mux := setupEmployeeTest(t)

	reqBody := createEmployeeRequest{
		FullName: "John Doe",
		Position: "Developer",
	}
	body, _ := json.Marshal(reqBody)

	mockSvc.On("Create", mock.Anything, uint(999), "John Doe", "Developer", (*time.Time)(nil)).Return(nil, apperrors.ErrEmployeeDepartmentNotFound)

	req := httptest.NewRequest(http.MethodPost, "/departments/999/employees", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), apperrors.ErrEmployeeDepartmentNotFound.Error())

	mockSvc.AssertExpectations(t)
}

func TestCreateEmployee_InvalidFullName(t *testing.T) {
	mockSvc, _, mux := setupEmployeeTest(t)

	reqBody := createEmployeeRequest{
		FullName: "",
		Position: "Developer",
	}
	body, _ := json.Marshal(reqBody)

	mockSvc.On("Create", mock.Anything, uint(1), "", "Developer", (*time.Time)(nil)).Return(nil, apperrors.ErrInvalidFullName)

	req := httptest.NewRequest(http.MethodPost, "/departments/1/employees", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), apperrors.ErrInvalidFullName.Error())

	mockSvc.AssertExpectations(t)
}

func TestCreateEmployee_InvalidPosition(t *testing.T) {
	mockSvc, _, mux := setupEmployeeTest(t)

	reqBody := createEmployeeRequest{
		FullName: "John Doe",
		Position: "",
	}
	body, _ := json.Marshal(reqBody)

	mockSvc.On("Create", mock.Anything, uint(1), "John Doe", "", (*time.Time)(nil)).Return(nil, apperrors.ErrInvalidPosition)

	req := httptest.NewRequest(http.MethodPost, "/departments/1/employees", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), apperrors.ErrInvalidPosition.Error())

	mockSvc.AssertExpectations(t)
}

func TestCreateEmployee_InvalidJSON(t *testing.T) {
	_, _, mux := setupEmployeeTest(t)

	req := httptest.NewRequest(http.MethodPost, "/departments/1/employees", bytes.NewReader([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid request body")
}

func TestCreateEmployee_InvalidDepartmentID(t *testing.T) {
	_, _, mux := setupEmployeeTest(t)

	reqBody := createEmployeeRequest{FullName: "John", Position: "Dev"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/departments/abc/employees", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid department id")
}

// internal/handlers/department_handler_test.go
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

type MockDepartmentService struct {
	mock.Mock
}

func (m *MockDepartmentService) Create(ctx context.Context, name string, parentID *uint) (*models.Department, error) {
	args := m.Called(ctx, name, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentService) GetByID(ctx context.Context, id uint, depth int, includeEmployees bool) (*models.Department, error) {
	args := m.Called(ctx, id, depth, includeEmployees)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentService) Update(ctx context.Context, id uint, name *string, parentID *uint) (*models.Department, error) {
	args := m.Called(ctx, id, name, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentService) Delete(ctx context.Context, id uint, mode string, reassignTo *uint) error {
	args := m.Called(ctx, id, mode, reassignTo)
	return args.Error(0)
}

func setupDepartmentTest(t *testing.T) (*MockDepartmentService, *DepartmentHandler, *http.ServeMux) {
	mockSvc := new(MockDepartmentService)
	handler := NewDepartmentHandler(mockSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /departments", handler.CreateDepartment)
	mux.HandleFunc("GET /departments/{id}", handler.GetDepartment)
	mux.HandleFunc("PATCH /departments/{id}", handler.UpdateDepartment)
	mux.HandleFunc("DELETE /departments/{id}", handler.DeleteDepartment)

	return mockSvc, handler, mux
}

// --- CREATE ---
func TestCreateDepartment_Success(t *testing.T) {
	mockSvc, _, mux := setupDepartmentTest(t)

	reqBody := createDepartmentRequest{Name: "IT", ParentID: nil}
	body, _ := json.Marshal(reqBody)

	expectedDept := &models.Department{
		ID:        1,
		Name:      "IT",
		ParentID:  nil,
		CreatedAt: time.Now(),
	}

	mockSvc.On("Create", mock.Anything, "IT", (*uint)(nil)).Return(expectedDept, nil)

	req := httptest.NewRequest(http.MethodPost, "/departments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp models.Department
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedDept.Name, resp.Name)

	mockSvc.AssertExpectations(t)
}

func TestCreateDepartment_Conflict(t *testing.T) {
	mockSvc, _, mux := setupDepartmentTest(t)

	reqBody := createDepartmentRequest{Name: "IT", ParentID: nil}
	body, _ := json.Marshal(reqBody)

	mockSvc.On("Create", mock.Anything, "IT", (*uint)(nil)).Return(nil, apperrors.ErrDepartmentNameConflict)

	req := httptest.NewRequest(http.MethodPost, "/departments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), apperrors.ErrDepartmentNameConflict.Error())

	mockSvc.AssertExpectations(t)
}

func TestCreateDepartment_ParentNotFound(t *testing.T) {
	mockSvc, _, mux := setupDepartmentTest(t)

	parentID := uint(999)
	reqBody := createDepartmentRequest{Name: "Backend", ParentID: &parentID}
	body, _ := json.Marshal(reqBody)

	mockSvc.On("Create", mock.Anything, "Backend", &parentID).Return(nil, apperrors.ErrParentNotFound)

	req := httptest.NewRequest(http.MethodPost, "/departments", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), apperrors.ErrParentNotFound.Error())

	mockSvc.AssertExpectations(t)
}

// --- GET ---
func TestGetDepartment_Success(t *testing.T) {
	mockSvc, _, mux := setupDepartmentTest(t)

	expectedDept := &models.Department{
		ID:        1,
		Name:      "IT",
		ParentID:  nil,
		CreatedAt: time.Now(),
		Employees: []models.Employee{},
		Children:  []models.Department{},
	}

	mockSvc.On("GetByID", mock.Anything, uint(1), 1, true).Return(expectedDept, nil)

	req := httptest.NewRequest(http.MethodGet, "/departments/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.Department
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, expectedDept.Name, resp.Name)

	mockSvc.AssertExpectations(t)
}

func TestGetDepartment_NotFound(t *testing.T) {
	mockSvc, _, mux := setupDepartmentTest(t)

	mockSvc.On("GetByID", mock.Anything, uint(999), 1, true).Return(nil, apperrors.ErrDepartmentNotFound)

	req := httptest.NewRequest(http.MethodGet, "/departments/999", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), apperrors.ErrDepartmentNotFound.Error())

	mockSvc.AssertExpectations(t)
}

// --- UPDATE ---
func TestUpdateDepartment_Success(t *testing.T) {
	mockSvc, _, mux := setupDepartmentTest(t)

	newName := "NewIT"
	reqBody := updateDepartmentRequest{Name: &newName}
	body, _ := json.Marshal(reqBody)

	updatedDept := &models.Department{
		ID:        1,
		Name:      newName,
		ParentID:  nil,
		CreatedAt: time.Now(),
	}

	mockSvc.On("Update", mock.Anything, uint(1), &newName, (*uint)(nil)).Return(updatedDept, nil)

	req := httptest.NewRequest(http.MethodPatch, "/departments/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.Department
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, newName, resp.Name)

	mockSvc.AssertExpectations(t)
}

func TestUpdateDepartment_NotFound(t *testing.T) {
	mockSvc, _, mux := setupDepartmentTest(t)

	newName := "NewIT"
	reqBody := updateDepartmentRequest{Name: &newName}
	body, _ := json.Marshal(reqBody)

	mockSvc.On("Update", mock.Anything, uint(999), &newName, (*uint)(nil)).Return(nil, apperrors.ErrDepartmentNotFound)

	req := httptest.NewRequest(http.MethodPatch, "/departments/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), apperrors.ErrDepartmentNotFound.Error())

	mockSvc.AssertExpectations(t)
}

// --- DELETE ---
func TestDeleteDepartment_Cascade(t *testing.T) {
	mockSvc, _, mux := setupDepartmentTest(t)

	mockSvc.On("Delete", mock.Anything, uint(1), "cascade", (*uint)(nil)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/departments/1?mode=cascade", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeleteDepartment_Reassign(t *testing.T) {
	mockSvc, _, mux := setupDepartmentTest(t)

	target := uint(2)
	mockSvc.On("Delete", mock.Anything, uint(1), "reassign", &target).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/departments/1?mode=reassign&reassign_to_department_id=2", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeleteDepartment_MissingMode(t *testing.T) {
	_, _, mux := setupDepartmentTest(t)

	req := httptest.NewRequest(http.MethodDelete, "/departments/1", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "mode query parameter is required")
}

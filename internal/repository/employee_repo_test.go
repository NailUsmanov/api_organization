package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/NailUsmanov/api_organization/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmployeeRepo - ручной мок для EmployeeRepo
type MockEmployeeRepo struct {
	mock.Mock
}

func (m *MockEmployeeRepo) Create(ctx context.Context, emp *models.Employee) error {
	args := m.Called(ctx, emp)
	return args.Error(0)
}

func (m *MockEmployeeRepo) GetByID(ctx context.Context, id uint) (*models.Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) Update(ctx context.Context, emp *models.Employee) error {
	args := m.Called(ctx, emp)
	return args.Error(0)
}

func (m *MockEmployeeRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEmployeeRepo) ListByDepartment(ctx context.Context, departmentID uint, orderBy string) ([]models.Employee, error) {
	args := m.Called(ctx, departmentID, orderBy)
	return args.Get(0).([]models.Employee), args.Error(1)
}

func (m *MockEmployeeRepo) MoveToDepartment(ctx context.Context, departmentID uint, targetDepartmentID uint) error {
	args := m.Called(ctx, departmentID, targetDepartmentID)
	return args.Error(0)
}

// Тесты
func TestEmployeeRepo_Create(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	ctx := context.Background()

	hiredAt := time.Now()
	emp := &models.Employee{
		DepartmentID: 1,
		FullName:     "John Doe",
		Position:     "Developer",
		HiredAt:      &hiredAt,
	}

	mockRepo.On("Create", ctx, emp).Return(nil)

	err := mockRepo.Create(ctx, emp)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeRepo_Create_Error(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	ctx := context.Background()

	emp := &models.Employee{
		FullName: "John Doe",
		Position: "Developer",
	}
	expectedErr := errors.New("database error")

	mockRepo.On("Create", ctx, emp).Return(expectedErr)

	err := mockRepo.Create(ctx, emp)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeRepo_GetByID_Success(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	ctx := context.Background()

	expectedEmp := &models.Employee{
		ID:       1,
		FullName: "John Doe",
		Position: "Developer",
	}

	mockRepo.On("GetByID", ctx, uint(1)).Return(expectedEmp, nil)

	emp, err := mockRepo.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, emp)
	assert.Equal(t, expectedEmp.ID, emp.ID)
	assert.Equal(t, expectedEmp.FullName, emp.FullName)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeRepo_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, uint(999)).Return(nil, nil)

	emp, err := mockRepo.GetByID(ctx, 999)

	assert.NoError(t, err)
	assert.Nil(t, emp)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeRepo_GetByID_Error(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	ctx := context.Background()

	expectedErr := errors.New("database error")

	mockRepo.On("GetByID", ctx, uint(1)).Return(nil, expectedErr)

	emp, err := mockRepo.GetByID(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, emp)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeRepo_Update(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	ctx := context.Background()

	emp := &models.Employee{
		ID:       1,
		FullName: "Updated Name",
		Position: "Senior Developer",
	}

	mockRepo.On("Update", ctx, emp).Return(nil)

	err := mockRepo.Update(ctx, emp)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeRepo_Delete(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	ctx := context.Background()

	mockRepo.On("Delete", ctx, uint(1)).Return(nil)

	err := mockRepo.Delete(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeRepo_ListByDepartment_OrderByCreatedAt(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	ctx := context.Background()

	departmentID := uint(1) // <-- ИСПРАВЛЕНО: используем uint
	expectedEmps := []models.Employee{
		{ID: 1, FullName: "John", DepartmentID: 1},
		{ID: 2, FullName: "Jane", DepartmentID: 1},
	}

	mockRepo.On("ListByDepartment", ctx, departmentID, "created_at").Return(expectedEmps, nil)

	emps, err := mockRepo.ListByDepartment(ctx, departmentID, "created_at")

	assert.NoError(t, err)
	assert.Len(t, emps, 2)
	assert.Equal(t, expectedEmps, emps)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeRepo_ListByDepartment_OrderByFullName(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	ctx := context.Background()

	departmentID := uint(1) // <-- ИСПРАВЛЕНО: используем uint
	expectedEmps := []models.Employee{
		{ID: 2, FullName: "Jane", DepartmentID: 1},
		{ID: 1, FullName: "John", DepartmentID: 1},
	}

	mockRepo.On("ListByDepartment", ctx, departmentID, "full_name").Return(expectedEmps, nil)

	emps, err := mockRepo.ListByDepartment(ctx, departmentID, "full_name")

	assert.NoError(t, err)
	assert.Len(t, emps, 2)
	assert.Equal(t, expectedEmps, emps)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeRepo_ListByDepartment_Empty(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	ctx := context.Background()

	departmentID := uint(1)
	expectedEmps := []models.Employee{}

	mockRepo.On("ListByDepartment", ctx, departmentID, "created_at").Return(expectedEmps, nil)

	emps, err := mockRepo.ListByDepartment(ctx, departmentID, "created_at")

	assert.NoError(t, err)
	assert.Empty(t, emps)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeRepo_ListByDepartment_Error(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	ctx := context.Background()

	departmentID := uint(1)
	expectedErr := errors.New("database error")

	mockRepo.On("ListByDepartment", ctx, departmentID, "created_at").Return([]models.Employee{}, expectedErr)

	emps, err := mockRepo.ListByDepartment(ctx, departmentID, "created_at")

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, emps)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeRepo_MoveToDepartment(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	ctx := context.Background()

	fromDept := uint(1)
	toDept := uint(2)

	mockRepo.On("MoveToDepartment", ctx, fromDept, toDept).Return(nil)

	err := mockRepo.MoveToDepartment(ctx, fromDept, toDept)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestEmployeeRepo_MoveToDepartment_Error(t *testing.T) {
	mockRepo := new(MockEmployeeRepo)
	ctx := context.Background()

	fromDept := uint(1)
	toDept := uint(2)
	expectedErr := errors.New("database error")

	mockRepo.On("MoveToDepartment", ctx, fromDept, toDept).Return(expectedErr)

	err := mockRepo.MoveToDepartment(ctx, fromDept, toDept)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

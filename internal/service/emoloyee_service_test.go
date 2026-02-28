// internal/service/employee_service_test.go
package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/NailUsmanov/api_organization/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEmployeeRepoForService - мок для repository.EmployeeRepo
type MockEmployeeRepoForService struct {
	mock.Mock
}

func (m *MockEmployeeRepoForService) Create(ctx context.Context, emp *models.Employee) error {
	args := m.Called(ctx, emp)
	return args.Error(0)
}

func (m *MockEmployeeRepoForService) GetByID(ctx context.Context, id uint) (*models.Employee, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Employee), args.Error(1)
}

func (m *MockEmployeeRepoForService) Update(ctx context.Context, emp *models.Employee) error {
	args := m.Called(ctx, emp)
	return args.Error(0)
}

func (m *MockEmployeeRepoForService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEmployeeRepoForService) ListByDepartment(ctx context.Context, departmentID uint, orderBy string) ([]models.Employee, error) {
	args := m.Called(ctx, departmentID, orderBy)
	return args.Get(0).([]models.Employee), args.Error(1)
}

func (m *MockEmployeeRepoForService) MoveToDepartment(ctx context.Context, departmentID uint, targetDepartmentID uint) error {
	args := m.Called(ctx, departmentID, targetDepartmentID)
	return args.Error(0)
}

// MockDepartmentRepoForEmployee - мок для repository.DepartmentRepo (нужен только GetByID)
type MockDepartmentRepoForEmployee struct {
	mock.Mock
}

func (m *MockDepartmentRepoForEmployee) Create(ctx context.Context, dept *models.Department) error {
	args := m.Called(ctx, dept)
	return args.Error(0)
}

func (m *MockDepartmentRepoForEmployee) GetByID(ctx context.Context, id uint) (*models.Department, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentRepoForEmployee) Update(ctx context.Context, dept *models.Department) error {
	args := m.Called(ctx, dept)
	return args.Error(0)
}

func (m *MockDepartmentRepoForEmployee) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDepartmentRepoForEmployee) GetChildren(ctx context.Context, parentID *uint) ([]models.Department, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]models.Department), args.Error(1)
}

func (m *MockDepartmentRepoForEmployee) GetSubTree(ctx context.Context, rootID uint, depth int) ([]models.Department, error) {
	args := m.Called(ctx, rootID, depth)
	return args.Get(0).([]models.Department), args.Error(1)
}

func (m *MockDepartmentRepoForEmployee) GetByNameAndParent(ctx context.Context, name string, parentID *uint) (*models.Department, error) {
	args := m.Called(ctx, name, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

// Вспомогательная функция для создания сервиса с моками
func setupEmployeeService(t *testing.T) (*EmpService, *MockEmployeeRepoForService, *MockDepartmentRepoForEmployee) {
	mockEmpRepo := new(MockEmployeeRepoForService)
	mockDeptRepo := new(MockDepartmentRepoForEmployee)
	service := NewEmpService(mockEmpRepo, mockDeptRepo)
	return service, mockEmpRepo, mockDeptRepo
}

// --- Тесты для Create ---

func TestEmployeeCreate_Success(t *testing.T) {
	service, mockEmpRepo, mockDeptRepo := setupEmployeeService(t)
	ctx := context.Background()

	departmentID := uint(1)
	hiredAt := time.Now()

	// Ожидаем проверку существования отдела
	mockDeptRepo.On("GetByID", ctx, departmentID).Return(&models.Department{ID: 1, Name: "IT"}, nil)

	// Ожидаем создание сотрудника
	mockEmpRepo.On("Create", ctx, mock.MatchedBy(func(emp *models.Employee) bool {
		return emp.DepartmentID == int(departmentID) &&
			emp.FullName == "John Doe" &&
			emp.Position == "Developer" &&
			emp.HiredAt == &hiredAt
	})).Return(nil)

	emp, err := service.Create(ctx, departmentID, "John Doe", "Developer", &hiredAt)

	assert.NoError(t, err)
	assert.NotNil(t, emp)
	assert.Equal(t, "John Doe", emp.FullName)
	assert.Equal(t, "Developer", emp.Position)
	assert.Equal(t, int(departmentID), emp.DepartmentID)
	assert.Equal(t, &hiredAt, emp.HiredAt)

	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertExpectations(t)
}

func TestEmployeeCreate_Success_WithoutHiredAt(t *testing.T) {
	service, mockEmpRepo, mockDeptRepo := setupEmployeeService(t)
	ctx := context.Background()

	departmentID := uint(1)

	// Ожидаем проверку существования отдела
	mockDeptRepo.On("GetByID", ctx, departmentID).Return(&models.Department{ID: 1, Name: "IT"}, nil)

	// Ожидаем создание сотрудника
	mockEmpRepo.On("Create", ctx, mock.MatchedBy(func(emp *models.Employee) bool {
		return emp.DepartmentID == int(departmentID) &&
			emp.FullName == "John Doe" &&
			emp.Position == "Developer" &&
			emp.HiredAt == nil
	})).Return(nil)

	emp, err := service.Create(ctx, departmentID, "John Doe", "Developer", nil)

	assert.NoError(t, err)
	assert.NotNil(t, emp)
	assert.Equal(t, "John Doe", emp.FullName)
	assert.Equal(t, "Developer", emp.Position)
	assert.Equal(t, int(departmentID), emp.DepartmentID)
	assert.Nil(t, emp.HiredAt)

	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertExpectations(t)
}

func TestEmployeeCreate_DepartmentNotFound(t *testing.T) {
	service, mockEmpRepo, mockDeptRepo := setupEmployeeService(t)
	ctx := context.Background()

	departmentID := uint(999)

	// Отдел не найден
	mockDeptRepo.On("GetByID", ctx, departmentID).Return(nil, nil)

	emp, err := service.Create(ctx, departmentID, "John Doe", "Developer", nil)

	assert.Error(t, err)
	assert.Equal(t, "department not found", err.Error())
	assert.Nil(t, emp)

	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertNotCalled(t, "Create")
}

func TestEmployeeCreate_DepartmentError(t *testing.T) {
	service, mockEmpRepo, mockDeptRepo := setupEmployeeService(t)
	ctx := context.Background()

	departmentID := uint(1)
	dbErr := errors.New("database error")

	// Ошибка при проверке отдела
	mockDeptRepo.On("GetByID", ctx, departmentID).Return(nil, dbErr)

	emp, err := service.Create(ctx, departmentID, "John Doe", "Developer", nil)

	assert.Error(t, err)
	assert.Equal(t, dbErr, err)
	assert.Nil(t, emp)

	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertNotCalled(t, "Create")
}

func TestEmployeeCreate_EmptyFullName(t *testing.T) {
	service, mockEmpRepo, mockDeptRepo := setupEmployeeService(t)
	ctx := context.Background()

	departmentID := uint(1)

	// Отдел существует
	mockDeptRepo.On("GetByID", ctx, departmentID).Return(&models.Department{ID: 1, Name: "IT"}, nil)

	emp, err := service.Create(ctx, departmentID, "", "Developer", nil)

	assert.Error(t, err)
	assert.Equal(t, "full_name must be not empty", err.Error())
	assert.Nil(t, emp)

	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertNotCalled(t, "Create")
}

func TestEmployeeCreate_FullNameTooLong(t *testing.T) {
	service, mockEmpRepo, mockDeptRepo := setupEmployeeService(t)
	ctx := context.Background()

	departmentID := uint(1)
	longName := string(make([]byte, 201))

	// Отдел существует
	mockDeptRepo.On("GetByID", ctx, departmentID).Return(&models.Department{ID: 1, Name: "IT"}, nil)

	emp, err := service.Create(ctx, departmentID, longName, "Developer", nil)

	assert.Error(t, err)
	assert.Equal(t, "full_name can't be longer than 200", err.Error())
	assert.Nil(t, emp)

	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertNotCalled(t, "Create")
}

func TestEmployeeCreate_EmptyPosition(t *testing.T) {
	service, mockEmpRepo, mockDeptRepo := setupEmployeeService(t)
	ctx := context.Background()

	departmentID := uint(1)

	// Отдел существует
	mockDeptRepo.On("GetByID", ctx, departmentID).Return(&models.Department{ID: 1, Name: "IT"}, nil)

	emp, err := service.Create(ctx, departmentID, "John Doe", "", nil)

	assert.Error(t, err)
	assert.Equal(t, "position must be non-empty and max 200 characters", err.Error())
	assert.Nil(t, emp)

	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertNotCalled(t, "Create")
}

func TestEmployeeCreate_PositionTooLong(t *testing.T) {
	service, mockEmpRepo, mockDeptRepo := setupEmployeeService(t)
	ctx := context.Background()

	departmentID := uint(1)
	longPosition := string(make([]byte, 201))

	// Отдел существует
	mockDeptRepo.On("GetByID", ctx, departmentID).Return(&models.Department{ID: 1, Name: "IT"}, nil)

	emp, err := service.Create(ctx, departmentID, "John Doe", longPosition, nil)

	assert.Error(t, err)
	assert.Equal(t, "position must be non-empty and max 200 characters", err.Error())
	assert.Nil(t, emp)

	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertNotCalled(t, "Create")
}

func TestEmployeeCreate_TrimSpaces(t *testing.T) {
	service, mockEmpRepo, mockDeptRepo := setupEmployeeService(t)
	ctx := context.Background()

	departmentID := uint(1)

	// Отдел существует
	mockDeptRepo.On("GetByID", ctx, departmentID).Return(&models.Department{ID: 1, Name: "IT"}, nil)

	// Ожидаем создание сотрудника с обрезанными пробелами
	mockEmpRepo.On("Create", ctx, mock.MatchedBy(func(emp *models.Employee) bool {
		return emp.FullName == "John Doe" && // пробелы обрезаны
			emp.Position == "Developer" // пробелы обрезаны
	})).Return(nil)

	emp, err := service.Create(ctx, departmentID, "  John Doe  ", "  Developer  ", nil)

	assert.NoError(t, err)
	assert.NotNil(t, emp)
	assert.Equal(t, "John Doe", emp.FullName)
	assert.Equal(t, "Developer", emp.Position)

	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertExpectations(t)
}

func TestEmployeeCreate_WithHiredAt(t *testing.T) {
	service, mockEmpRepo, mockDeptRepo := setupEmployeeService(t)
	ctx := context.Background()

	departmentID := uint(1)
	hiredAt := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)

	// Отдел существует
	mockDeptRepo.On("GetByID", ctx, departmentID).Return(&models.Department{ID: 1, Name: "IT"}, nil)

	// Ожидаем создание сотрудника с hiredAt
	mockEmpRepo.On("Create", ctx, mock.MatchedBy(func(emp *models.Employee) bool {
		return emp.FullName == "John Doe" &&
			emp.Position == "Developer" &&
			emp.HiredAt != nil &&
			emp.HiredAt.Equal(hiredAt)
	})).Return(nil)

	emp, err := service.Create(ctx, departmentID, "John Doe", "Developer", &hiredAt)

	assert.NoError(t, err)
	assert.NotNil(t, emp)
	assert.Equal(t, "John Doe", emp.FullName)
	assert.Equal(t, "Developer", emp.Position)
	assert.Equal(t, hiredAt, *emp.HiredAt)

	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertExpectations(t)
}

func TestEmployeeCreate_CreateError(t *testing.T) {
	service, mockEmpRepo, mockDeptRepo := setupEmployeeService(t)
	ctx := context.Background()

	departmentID := uint(1)
	dbErr := errors.New("database error")

	// Отдел существует
	mockDeptRepo.On("GetByID", ctx, departmentID).Return(&models.Department{ID: 1, Name: "IT"}, nil)

	// Ошибка при создании
	mockEmpRepo.On("Create", ctx, mock.Anything).Return(dbErr)

	emp, err := service.Create(ctx, departmentID, "John Doe", "Developer", nil)

	// Проверяем, что вернулась ошибка
	assert.Error(t, err)
	assert.Equal(t, dbErr, err)
	assert.Nil(t, emp)

	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertExpectations(t)
}

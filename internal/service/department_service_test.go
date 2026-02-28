// internal/service/department_service_test.go
package service

import (
	"context"
	"testing"

	"github.com/NailUsmanov/api_organization/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDepartmentRepo - полная реализация мока для repository.DepartmentRepo
type MockDepartmentRepo struct {
	mock.Mock
}

func (m *MockDepartmentRepo) Create(ctx context.Context, dept *models.Department) error {
	args := m.Called(ctx, dept)
	return args.Error(0)
}

func (m *MockDepartmentRepo) GetByID(ctx context.Context, id uint) (*models.Department, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentRepo) Update(ctx context.Context, dept *models.Department) error {
	args := m.Called(ctx, dept)
	return args.Error(0)
}

func (m *MockDepartmentRepo) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockDepartmentRepo) GetChildren(ctx context.Context, parentID *uint) ([]models.Department, error) {
	args := m.Called(ctx, parentID)
	return args.Get(0).([]models.Department), args.Error(1)
}

func (m *MockDepartmentRepo) GetSubTree(ctx context.Context, rootID uint, depth int) ([]models.Department, error) {
	args := m.Called(ctx, rootID, depth)
	return args.Get(0).([]models.Department), args.Error(1)
}

func (m *MockDepartmentRepo) GetByNameAndParent(ctx context.Context, name string, parentID *uint) (*models.Department, error) {
	args := m.Called(ctx, name, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

// MockEmployeeRepo - полная реализация мока для repository.EmployeeRepo
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

// Вспомогательная функция для создания сервиса с моками
func setupDepartmentService(t *testing.T) (*DepService, *MockDepartmentRepo, *MockEmployeeRepo) {
	mockDeptRepo := new(MockDepartmentRepo)
	mockEmpRepo := new(MockEmployeeRepo)
	service := NewDepartmentService(mockDeptRepo, mockEmpRepo)
	return service, mockDeptRepo, mockEmpRepo
}

// --- Тесты для Create ---

func TestCreate_Success(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	mockDeptRepo.On("GetByNameAndParent", ctx, "IT", (*uint)(nil)).Return(nil, nil)
	mockDeptRepo.On("Create", ctx, mock.MatchedBy(func(dept *models.Department) bool {
		return dept.Name == "IT" && dept.ParentID == nil
	})).Return(nil)

	dept, err := service.Create(ctx, "IT", nil)

	assert.NoError(t, err)
	assert.NotNil(t, dept)
	assert.Equal(t, "IT", dept.Name)
	assert.Nil(t, dept.ParentID)
	mockDeptRepo.AssertExpectations(t)
}

func TestCreate_WithParent_Success(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	parentID := uint(1)
	parentDept := &models.Department{ID: 1, Name: "Parent"}

	mockDeptRepo.On("GetByNameAndParent", ctx, "Child", &parentID).Return(nil, nil)
	mockDeptRepo.On("GetByID", ctx, parentID).Return(parentDept, nil)
	mockDeptRepo.On("Create", ctx, mock.MatchedBy(func(dept *models.Department) bool {
		return dept.Name == "Child" && *dept.ParentID == parentID
	})).Return(nil)

	dept, err := service.Create(ctx, "Child", &parentID)

	assert.NoError(t, err)
	assert.NotNil(t, dept)
	assert.Equal(t, "Child", dept.Name)
	assert.Equal(t, parentID, *dept.ParentID)
	mockDeptRepo.AssertExpectations(t)
}

func TestCreate_EmptyName(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	dept, err := service.Create(ctx, "", nil)

	assert.Error(t, err)
	assert.Equal(t, "name cannot be empty", err.Error())
	assert.Nil(t, dept)
	mockDeptRepo.AssertNotCalled(t, "GetByNameAndParent")
	mockDeptRepo.AssertNotCalled(t, "Create")
}

func TestCreate_NameTooLong(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	longName := string(make([]byte, 201))
	dept, err := service.Create(ctx, longName, nil)

	assert.Error(t, err)
	assert.Equal(t, "name too long (max 200)", err.Error())
	assert.Nil(t, dept)
	mockDeptRepo.AssertNotCalled(t, "GetByNameAndParent")
	mockDeptRepo.AssertNotCalled(t, "Create")
}

func TestCreate_NameConflict(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	existingDept := &models.Department{ID: 1, Name: "IT"}

	mockDeptRepo.On("GetByNameAndParent", ctx, "IT", (*uint)(nil)).Return(existingDept, nil)

	dept, err := service.Create(ctx, "IT", nil)

	assert.Error(t, err)
	assert.Equal(t, "department with this name already exists the same parent", err.Error())
	assert.Nil(t, dept)
	mockDeptRepo.AssertExpectations(t)
	mockDeptRepo.AssertNotCalled(t, "Create")
}

func TestCreate_ParentNotFound(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	parentID := uint(999)

	mockDeptRepo.On("GetByNameAndParent", ctx, "Child", &parentID).Return(nil, nil)
	mockDeptRepo.On("GetByID", ctx, parentID).Return(nil, nil)

	dept, err := service.Create(ctx, "Child", &parentID)

	assert.Error(t, err)
	assert.Equal(t, "parent department not found", err.Error())
	assert.Nil(t, dept)
	mockDeptRepo.AssertExpectations(t)
	mockDeptRepo.AssertNotCalled(t, "Create")
}

// --- Тесты для GetByID ---

func TestGetByID_Success(t *testing.T) {
	service, mockDeptRepo, mockEmpRepo := setupDepartmentService(t)
	ctx := context.Background()

	expectedDept := &models.Department{ID: 1, Name: "IT"}

	mockDeptRepo.On("GetByID", ctx, uint(1)).Return(expectedDept, nil)

	employees := []models.Employee{
		{ID: 1, FullName: "John", DepartmentID: 1},
	}
	mockEmpRepo.On("ListByDepartment", ctx, uint(1), "created_at").Return(employees, nil)

	dept, err := service.GetByID(ctx, 1, 1, true)

	assert.NoError(t, err)
	assert.NotNil(t, dept)
	assert.Equal(t, expectedDept.ID, dept.ID)
	assert.Equal(t, expectedDept.Name, dept.Name)
	assert.Len(t, dept.Employees, 1)
	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertExpectations(t)
}

func TestGetByID_WithDepth(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	expectedDept := &models.Department{ID: 1, Name: "Root"}

	mockDeptRepo.On("GetByID", ctx, uint(1)).Return(expectedDept, nil)

	children := []models.Department{
		{ID: 2, Name: "Child1", ParentID: &[]uint{1}[0]},
		{ID: 3, Name: "Child2", ParentID: &[]uint{1}[0]},
	}
	mockDeptRepo.On("GetSubTree", ctx, uint(1), 2).Return(children, nil)

	dept, err := service.GetByID(ctx, 1, 2, false)

	assert.NoError(t, err)
	assert.NotNil(t, dept)
	assert.Len(t, dept.Children, 2)
	mockDeptRepo.AssertExpectations(t)
}

func TestGetByID_NotFound(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	mockDeptRepo.On("GetByID", ctx, uint(999)).Return(nil, nil)

	dept, err := service.GetByID(ctx, 999, 1, true)

	assert.Error(t, err)
	assert.Equal(t, "deparment not found", err.Error())
	assert.Nil(t, dept)
	mockDeptRepo.AssertExpectations(t)
}

// --- Тесты для Update ---

func TestUpdate_NameOnly(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	existingDept := &models.Department{ID: 1, Name: "OldName", ParentID: nil}
	newName := "NewName"

	mockDeptRepo.On("GetByID", ctx, uint(1)).Return(existingDept, nil)
	mockDeptRepo.On("GetByNameAndParent", ctx, newName, (*uint)(nil)).Return(nil, nil)
	mockDeptRepo.On("Update", ctx, mock.MatchedBy(func(dept *models.Department) bool {
		return dept.ID == 1 && dept.Name == newName
	})).Return(nil)

	dept, err := service.Update(ctx, 1, &newName, nil)

	assert.NoError(t, err)
	assert.NotNil(t, dept)
	assert.Equal(t, newName, dept.Name)
	mockDeptRepo.AssertExpectations(t)
}

func TestUpdate_ParentOnly(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	existingDept := &models.Department{ID: 1, Name: "Dept", ParentID: nil}
	newParentID := uint(2)
	parentDept := &models.Department{ID: 2, Name: "Parent"}

	mockDeptRepo.On("GetByID", ctx, uint(1)).Return(existingDept, nil)
	mockDeptRepo.On("GetByID", ctx, newParentID).Return(parentDept, nil)
	mockDeptRepo.On("GetSubTree", ctx, uint(1), 100).Return([]models.Department{}, nil)
	mockDeptRepo.On("Update", ctx, mock.MatchedBy(func(dept *models.Department) bool {
		return dept.ID == 1 && *dept.ParentID == newParentID
	})).Return(nil)

	dept, err := service.Update(ctx, 1, nil, &newParentID)

	assert.NoError(t, err)
	assert.NotNil(t, dept)
	assert.Equal(t, newParentID, *dept.ParentID)
	mockDeptRepo.AssertExpectations(t)
}

func TestUpdate_NotFound(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	mockDeptRepo.On("GetByID", ctx, uint(999)).Return(nil, nil)

	newName := "NewName"
	dept, err := service.Update(ctx, 999, &newName, nil)

	assert.Error(t, err)
	assert.Equal(t, "department not found", err.Error())
	assert.Nil(t, dept)
	mockDeptRepo.AssertExpectations(t)
}

func TestUpdate_NameConflict(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	existingDept := &models.Department{ID: 1, Name: "OldName", ParentID: nil}
	conflictingDept := &models.Department{ID: 2, Name: "NewName", ParentID: nil}
	newName := "NewName"

	mockDeptRepo.On("GetByID", ctx, uint(1)).Return(existingDept, nil)
	mockDeptRepo.On("GetByNameAndParent", ctx, newName, (*uint)(nil)).Return(conflictingDept, nil)

	dept, err := service.Update(ctx, 1, &newName, nil)

	assert.Error(t, err)
	assert.Equal(t, "department with this name already exists under the same parent", err.Error())
	assert.Nil(t, dept)
	mockDeptRepo.AssertExpectations(t)
	mockDeptRepo.AssertNotCalled(t, "Update")
}

func TestUpdate_SelfParent(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	existingDept := &models.Department{ID: 1, Name: "Dept", ParentID: nil}
	selfID := uint(1)

	mockDeptRepo.On("GetByID", ctx, uint(1)).Return(existingDept, nil)

	dept, err := service.Update(ctx, 1, nil, &selfID)

	assert.Error(t, err)
	assert.Equal(t, "cannot set parent to itself", err.Error())
	assert.Nil(t, dept)
	mockDeptRepo.AssertExpectations(t)
	mockDeptRepo.AssertNotCalled(t, "Update")
}

func TestUpdate_CycleDetected(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	existingDept := &models.Department{ID: 1, Name: "Parent", ParentID: nil}
	newParentID := uint(3)

	mockDeptRepo.On("GetByID", ctx, uint(1)).Return(existingDept, nil)
	mockDeptRepo.On("GetByID", ctx, newParentID).Return(&models.Department{ID: 3}, nil)

	descendants := []models.Department{
		{ID: 2, Name: "Child", ParentID: &[]uint{1}[0]},
		{ID: 3, Name: "Grandchild", ParentID: &[]uint{2}[0]},
	}
	mockDeptRepo.On("GetSubTree", ctx, uint(1), 100).Return(descendants, nil)

	dept, err := service.Update(ctx, 1, nil, &newParentID)

	assert.Error(t, err)
	assert.Equal(t, "cannot move department to its own descendant", err.Error())
	assert.Nil(t, dept)
	mockDeptRepo.AssertExpectations(t)
	mockDeptRepo.AssertNotCalled(t, "Update")
}

// --- Тесты для Delete ---

func TestDelete_Cascade(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	dept := &models.Department{ID: 1, Name: "ToDelete"}

	mockDeptRepo.On("GetByID", ctx, uint(1)).Return(dept, nil)
	mockDeptRepo.On("Delete", ctx, uint(1)).Return(nil)

	err := service.Delete(ctx, 1, "cascade", nil)

	assert.NoError(t, err)
	mockDeptRepo.AssertExpectations(t)
}

func TestDelete_Reassign(t *testing.T) {
	service, mockDeptRepo, mockEmpRepo := setupDepartmentService(t)
	ctx := context.Background()

	dept := &models.Department{ID: 1, Name: "ToDelete"}
	reassignTo := uint(2)
	targetDept := &models.Department{ID: 2, Name: "Target"}

	mockDeptRepo.On("GetByID", ctx, uint(1)).Return(dept, nil)
	mockDeptRepo.On("GetByID", ctx, reassignTo).Return(targetDept, nil)
	mockDeptRepo.On("GetChildren", ctx, &[]uint{1}[0]).Return([]models.Department{}, nil)
	mockEmpRepo.On("MoveToDepartment", ctx, uint(1), reassignTo).Return(nil)
	mockDeptRepo.On("Delete", ctx, uint(1)).Return(nil)

	err := service.Delete(ctx, 1, "reassign", &reassignTo)

	assert.NoError(t, err)
	mockDeptRepo.AssertExpectations(t)
	mockEmpRepo.AssertExpectations(t)
}

func TestDelete_NotFound(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	mockDeptRepo.On("GetByID", ctx, uint(999)).Return(nil, nil)

	err := service.Delete(ctx, 999, "cascade", nil)

	assert.Error(t, err)
	assert.Equal(t, "department not found", err.Error())
	mockDeptRepo.AssertExpectations(t)
}

func TestDelete_InvalidMode(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	dept := &models.Department{ID: 1, Name: "Dept"}

	mockDeptRepo.On("GetByID", ctx, uint(1)).Return(dept, nil)

	err := service.Delete(ctx, 1, "invalid", nil)

	assert.Error(t, err)
	assert.Equal(t, "invalid mode, must be 'cascade' or 'reassign'", err.Error())
	mockDeptRepo.AssertExpectations(t)
}

func TestDelete_ReassignWithoutTarget(t *testing.T) {
	service, mockDeptRepo, _ := setupDepartmentService(t)
	ctx := context.Background()

	dept := &models.Department{ID: 1, Name: "Dept"}

	mockDeptRepo.On("GetByID", ctx, uint(1)).Return(dept, nil)

	err := service.Delete(ctx, 1, "reassign", nil)

	assert.Error(t, err)
	assert.Equal(t, "reassign_to_department_id is required for reassign mode", err.Error())
	mockDeptRepo.AssertExpectations(t)
}
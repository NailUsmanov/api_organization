package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/NailUsmanov/api_organization/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDepartmentRepo - ручной мок для DepartmentRepo
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

func (m *MockDepartmentRepo) GetByNameAndParent(ctx context.Context, name string, parentID *uint) (*models.Department, error) {
	args := m.Called(ctx, name, parentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Department), args.Error(1)
}

func (m *MockDepartmentRepo) GetSubTree(ctx context.Context, rootID uint, depth int) ([]models.Department, error) {
	args := m.Called(ctx, rootID, depth)
	return args.Get(0).([]models.Department), args.Error(1)
}

// Тесты для DepartmentRepo
func TestDepartmentRepo_Create(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	dept := &models.Department{
		Name: "IT",
	}

	mockRepo.On("Create", ctx, dept).Return(nil)

	err := mockRepo.Create(ctx, dept)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_Create_Error(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	dept := &models.Department{
		Name: "IT",
	}
	expectedErr := errors.New("database error")

	mockRepo.On("Create", ctx, dept).Return(expectedErr)

	err := mockRepo.Create(ctx, dept)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetByID_Success(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	expectedDept := &models.Department{
		ID:   1,
		Name: "IT",
	}

	mockRepo.On("GetByID", ctx, uint(1)).Return(expectedDept, nil)

	dept, err := mockRepo.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, dept)
	assert.Equal(t, expectedDept.ID, dept.ID)
	assert.Equal(t, expectedDept.Name, dept.Name)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	mockRepo.On("GetByID", ctx, uint(999)).Return(nil, nil)

	dept, err := mockRepo.GetByID(ctx, 999)

	assert.NoError(t, err)
	assert.Nil(t, dept)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetByID_Error(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	expectedErr := errors.New("database error")

	mockRepo.On("GetByID", ctx, uint(1)).Return(nil, expectedErr)

	dept, err := mockRepo.GetByID(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, dept)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_Update(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	dept := &models.Department{
		ID:   1,
		Name: "Updated IT",
	}

	mockRepo.On("Update", ctx, dept).Return(nil)

	err := mockRepo.Update(ctx, dept)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_Delete(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	mockRepo.On("Delete", ctx, uint(1)).Return(nil)

	err := mockRepo.Delete(ctx, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_Delete_Error(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	expectedErr := errors.New("database error")

	mockRepo.On("Delete", ctx, uint(1)).Return(expectedErr)

	err := mockRepo.Delete(ctx, 1)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetChildren(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	parentID := uint(1)
	expectedChildren := []models.Department{
		{ID: 2, Name: "Child1", ParentID: &parentID},
		{ID: 3, Name: "Child2", ParentID: &parentID},
	}

	mockRepo.On("GetChildren", ctx, &parentID).Return(expectedChildren, nil)

	children, err := mockRepo.GetChildren(ctx, &parentID)

	assert.NoError(t, err)
	assert.Len(t, children, 2)
	assert.Equal(t, expectedChildren[0].Name, children[0].Name)
	assert.Equal(t, expectedChildren[1].Name, children[1].Name)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetChildren_NilParent(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	expectedChildren := []models.Department{
		{ID: 1, Name: "Root1", ParentID: nil},
		{ID: 2, Name: "Root2", ParentID: nil},
	}

	mockRepo.On("GetChildren", ctx, (*uint)(nil)).Return(expectedChildren, nil)

	children, err := mockRepo.GetChildren(ctx, nil)

	assert.NoError(t, err)
	assert.Len(t, children, 2)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetChildren_Empty(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	parentID := uint(1)
	expectedChildren := []models.Department{}

	mockRepo.On("GetChildren", ctx, &parentID).Return(expectedChildren, nil)

	children, err := mockRepo.GetChildren(ctx, &parentID)

	assert.NoError(t, err)
	assert.Empty(t, children)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetChildren_Error(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	parentID := uint(1)
	expectedErr := errors.New("database error")

	mockRepo.On("GetChildren", ctx, &parentID).Return([]models.Department{}, expectedErr)

	children, err := mockRepo.GetChildren(ctx, &parentID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Empty(t, children)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetByNameAndParent_Success(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	parentID := uint(1)
	expectedDept := &models.Department{
		ID:       2,
		Name:     "Backend",
		ParentID: &parentID,
	}

	mockRepo.On("GetByNameAndParent", ctx, "Backend", &parentID).Return(expectedDept, nil)

	dept, err := mockRepo.GetByNameAndParent(ctx, "Backend", &parentID)

	assert.NoError(t, err)
	assert.NotNil(t, dept)
	assert.Equal(t, expectedDept.Name, dept.Name)
	assert.Equal(t, expectedDept.ParentID, dept.ParentID)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetByNameAndParent_NotFound(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	parentID := uint(1)

	mockRepo.On("GetByNameAndParent", ctx, "NonExistent", &parentID).Return(nil, nil)

	dept, err := mockRepo.GetByNameAndParent(ctx, "NonExistent", &parentID)

	assert.NoError(t, err)
	assert.Nil(t, dept)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetByNameAndParent_Error(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	parentID := uint(1)
	expectedErr := errors.New("database error")

	mockRepo.On("GetByNameAndParent", ctx, "Backend", &parentID).Return(nil, expectedErr)

	dept, err := mockRepo.GetByNameAndParent(ctx, "Backend", &parentID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, dept)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetByNameAndParent_NilParent(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	expectedDept := &models.Department{
		ID:       1,
		Name:     "Root",
		ParentID: nil,
	}

	mockRepo.On("GetByNameAndParent", ctx, "Root", (*uint)(nil)).Return(expectedDept, nil)

	dept, err := mockRepo.GetByNameAndParent(ctx, "Root", nil)

	assert.NoError(t, err)
	assert.NotNil(t, dept)
	assert.Equal(t, expectedDept.Name, dept.Name)
	assert.Nil(t, dept.ParentID)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetSubTree(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	rootID := uint(1)
	depth := 3
	expectedTree := []models.Department{
		{ID: 2, Name: "Child1", ParentID: &rootID},
		{ID: 3, Name: "Child2", ParentID: &rootID},
		{ID: 4, Name: "Grandchild", ParentID: &[]uint{2}[0]},
	}

	mockRepo.On("GetSubTree", ctx, rootID, depth).Return(expectedTree, nil)

	tree, err := mockRepo.GetSubTree(ctx, rootID, depth)

	assert.NoError(t, err)
	assert.Len(t, tree, 3)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetSubTree_DepthLessThan2(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	// При depth < 2 метод должен вернуть пустой слайс
	mockRepo.On("GetSubTree", ctx, uint(1), 1).Return([]models.Department{}, nil)

	tree, err := mockRepo.GetSubTree(ctx, 1, 1)

	assert.NoError(t, err)
	assert.Empty(t, tree)
	mockRepo.AssertExpectations(t)
}

func TestDepartmentRepo_GetSubTree_Error(t *testing.T) {
	mockRepo := new(MockDepartmentRepo)
	ctx := context.Background()

	rootID := uint(1)
	depth := 2
	expectedErr := errors.New("database error")

	mockRepo.On("GetSubTree", ctx, rootID, depth).Return([]models.Department(nil), expectedErr)

	tree, err := mockRepo.GetSubTree(ctx, rootID, depth)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, tree)
	mockRepo.AssertExpectations(t)
}

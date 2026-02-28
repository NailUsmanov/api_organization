// Package service содержит бизнес-логику приложения.
// Сервисы реализуют основные операции с подразделениями и сотрудниками,
// обеспечивают валидацию данных, проверку бизнес-правил и координацию между различными репозиториями.
package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/NailUsmanov/api_organization/internal/models"
)

// DepartmentService определяет интерфейс для работы с подразделениями,
// который будет реализован сервисом и использован обработчиками HTTP
type DepartmentService interface {
	Create(ctx context.Context, name string, parentID *uint) (*models.Department, error)
	GetByID(ctx context.Context, id uint, depth int, includeEmployees bool) (*models.Department, error)
	Update(ctx context.Context, id uint, name *string, parentID *uint) (*models.Department, error)
	Delete(ctx context.Context, id uint, mode string, reassignTo *uint) error
}

// DepartmentRepository определяет интерфейс репозитория подразделений, необходимый для работы сервиса.
type DepartmentRepository interface {
	Create(ctx context.Context, dept *models.Department) error
	GetByID(ctx context.Context, id uint) (*models.Department, error)
	Update(ctx context.Context, dept *models.Department) error
	Delete(ctx context.Context, id uint) error
	GetChildren(ctx context.Context, parentID *uint) ([]models.Department, error)
	GetSubTree(ctx context.Context, rootID uint, depth int) ([]models.Department, error)
	GetByNameAndParent(ctx context.Context, name string, parentID *uint) (*models.Department, error)
}

// DepService реализует бизнес-логику для работы с подразделениями.
type DepService struct {
	deptRepo DepartmentRepository
	empRepo  EmployeeRepository
}

// NewDepartmentService создаёт новый экземпляр сервиса подразделений.
func NewDepartmentService(deptRepo DepartmentRepository, empRepo EmployeeRepository) *DepService {
	return &DepService{deptRepo: deptRepo, empRepo: empRepo}
}

// ValidateName проверяет и очищает название подразделения.
func ValidateName(name string) (string, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "", errors.New("name cannot be empty")
	}

	if len(trimmed) > 200 {
		return "", errors.New("name too long (max 200)")
	}

	return trimmed, nil
}

// Create реализует бизнес-логику создания нового подразделения.
func (s *DepService) Create(ctx context.Context, name string, parentID *uint) (*models.Department, error) {
	clearName, err := ValidateName(name)
	if err != nil {
		return nil, err
	}

	existing, err := s.deptRepo.GetByNameAndParent(ctx, clearName, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to check uniqueness: %w", err)
	}
	if existing != nil {
		return nil, errors.New("department with this name already exists the same parent")
	}

	if parentID != nil {
		parent, err := s.deptRepo.GetByID(ctx, *parentID)
		if err != nil {
			return nil, fmt.Errorf("failed to check parent existence: %w", err)
		}
		if parent == nil {
			return nil, errors.New("parent department not found")
		}
	}
	dept := &models.Department{
		Name:     clearName,
		ParentID: parentID,
	}

	if err := s.deptRepo.Create(ctx, dept); err != nil {
		return nil, fmt.Errorf("failed to create department: %w", err)
	}

	return dept, nil
}

// GetByID реализует бизнес-логику получения подразделения с поддеревом.
func (s *DepService) GetByID(ctx context.Context, id uint, depth int, includeEmployees bool) (*models.Department, error) {
	root, err := s.deptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if root == nil {
		return nil, errors.New("deparment not found")
	}

	if includeEmployees {
		employees, err := s.empRepo.ListByDepartment(ctx, id, "created_at")
		if err != nil {
			return nil, err
		}
		root.Employees = employees
	}
	if depth > 1 {
		children, err := s.deptRepo.GetSubTree(ctx, id, depth)
		if err != nil {
			return nil, err
		}
		root.Children = s.buildChildrenTree(children, id)
	} else {
		root.Children = []models.Department{}
	}

	return root, nil
}

// buildChildrenTree преобразует плоский список подразделений в иерархическую структуру.
func (s *DepService) buildChildrenTree(flat []models.Department, rootID uint) []models.Department {
	if len(flat) == 0 {
		return nil
	}
	childrenMap := make(map[uint][]*models.Department)
	for i := range flat {
		dept := &flat[i]
		if dept.ParentID != nil {
			pid := *dept.ParentID
			childrenMap[pid] = append(childrenMap[pid], dept)
		}
	}

	var fillChildren func(*models.Department)
	fillChildren = func(d *models.Department) {
		if kids, ok := childrenMap[uint(d.ID)]; ok {
			d.Children = make([]models.Department, len(kids))
			for i, kid := range kids {
				d.Children[i] = *kid
				fillChildren(&d.Children[i])
			}
		}
	}

	rootChildren := childrenMap[rootID]
	result := make([]models.Department, len(rootChildren))
	for i, kid := range rootChildren {
		result[i] = *kid
		fillChildren(&result[i])
	}
	return result
}

// Update реализует бизнес-логику обновления подразделения.
func (s *DepService) Update(ctx context.Context, id uint, name *string, parentID *uint) (*models.Department, error) {
	dept, err := s.deptRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if dept == nil {
		return nil, errors.New("department not found")
	}

	if name != nil {
		cleanName, err := ValidateName(*name)
		if err != nil {
			return nil, err
		}
		// Определяем родителя для проверки уникальности
		checkParent := dept.ParentID
		if parentID != nil {
			checkParent = parentID
		}
		existing, err := s.deptRepo.GetByNameAndParent(ctx, cleanName, checkParent)
		if err != nil {
			return nil, err
		}
		if existing != nil && uint(existing.ID) != id {
			return nil, errors.New("department with this name already exists under the same parent")
		}
		dept.Name = cleanName
	}

	if parentID != nil {
		if *parentID == id {
			return nil, errors.New("cannot set parent to itself")
		}
		if *parentID != 0 {
			parent, err := s.deptRepo.GetByID(ctx, *parentID)
			if err != nil {
				return nil, err
			}
			if parent == nil {
				return nil, errors.New("parent department not found")
			}
		}
		// Проверка на цикл
		descendants, err := s.deptRepo.GetSubTree(ctx, id, 100) // большая глубина для получения всех потомков
		if err != nil {
			return nil, err
		}
		for _, d := range descendants {
			if uint(d.ID) == *parentID {
				return nil, errors.New("cannot move department to its own descendant")
			}
		}
		dept.ParentID = parentID
	}

	if err := s.deptRepo.Update(ctx, dept); err != nil {
		return nil, err
	}
	return dept, nil
}

// Delete реализует бизнес-логику удаления подразделения с учётом режима.
func (s *DepService) Delete(ctx context.Context, id uint, mode string, reassignTo *uint) error {
	dept, err := s.deptRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if dept == nil {
		return errors.New("department not found")
	}

	switch mode {
	case "cascade":
		return s.deptRepo.Delete(ctx, id)
	case "reassign":
		if reassignTo == nil {
			return errors.New("reassign_to_department_id is required for reassign mode")
		}
		target, err := s.deptRepo.GetByID(ctx, *reassignTo)
		if err != nil {
			return err
		}
		if target == nil {
			return errors.New("target department not found")
		}
		// Проверяем наличие дочерних подразделений
		children, err := s.deptRepo.GetChildren(ctx, &id)
		if err != nil {
			return err
		}
		if len(children) > 0 {
			return errors.New("cannot reassign department with children; delete children first or use cascade mode")
		}
		// Переназначаем сотрудников
		if err := s.empRepo.MoveToDepartment(ctx, id, *reassignTo); err != nil {
			return err
		}
		return s.deptRepo.Delete(ctx, id)
	default:
		return errors.New("invalid mode, must be 'cascade' or 'reassign'")
	}
}

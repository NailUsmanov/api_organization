// Package service содержит бизнес-логику приложения.
// Сервисы реализуют основные операции с подразделениями и сотрудниками,
// обеспечивают валидацию данных, проверку бизнес-правил и координацию между различными репозиториями.
package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/NailUsmanov/api_organization/internal/models"
)

// EmployeeRepository определяет интерфейс репозитория сотрудников, необходимый для работы сервиса.
type EmployeeRepository interface {
	Create(ctx context.Context, emp *models.Employee) error
	GetByID(ctx context.Context, id uint) (*models.Employee, error)
	Update(ctx context.Context, emp *models.Employee) error
	Delete(ctx context.Context, id uint) error
	ListByDepartment(ctx context.Context, departmenID uint, orderBy string) ([]models.Employee, error)
	MoveToDepartment(ctx context.Context, departmentID uint, targerDerpartmentID uint) error
}

// EmployeeService определяет интерфейс для работы с сотрудниками,
// который будет реализован сервисом и использован обработчиками HTTP.
type EmployeeService interface {
	Create(ctx context.Context, departmentID uint, fullName, position string, hiredAt *time.Time) (*models.Employee, error)
}

// EmpService реализует бизнес-логику для работы с сотрудниками.
type EmpService struct {
	empRepo  EmployeeRepository
	deptRepo DepartmentRepository
}

// NewEmpService создаёт новый экземпляр сервиса сотрудников.
func NewEmpService(epmRepo EmployeeRepository, deptRepo DepartmentRepository) *EmpService {
	return &EmpService{empRepo: epmRepo, deptRepo: deptRepo}
}

// Create реализует бизнес-логику создания нового сотрудника.
func (e *EmpService) Create(ctx context.Context, departmentID uint, fullName, position string, hiredAt *time.Time) (*models.Employee, error) {
	dept, err := e.deptRepo.GetByID(ctx, departmentID)
	if err != nil {
		return nil, err
	}

	if dept == nil {
		return nil, errors.New("department not found")
	}

	cleanFullName := strings.TrimSpace(fullName)
	if cleanFullName == "" {
		return nil, errors.New("full_name must be not empty")
	}

	if len(cleanFullName) > 200 {
		return nil, errors.New("full_name can't be longer than 200")
	}

	cleanPosition := strings.TrimSpace(position)
	if cleanPosition == "" || len(cleanPosition) > 200 {
		return nil, errors.New("position must be non-empty and max 200 characters")
	}

	emp := &models.Employee{
		DepartmentID: int(departmentID),
		FullName:     cleanFullName,
		Position:     cleanPosition,
		HiredAt:      hiredAt,
	}

	if err := e.empRepo.Create(ctx, emp); err != nil {
		return nil, err
	}

	return emp, nil
}

// Package repository предоставляет реализацию доступа к данным для работы с базой данных.
package repository

import (
	"context"
	"errors"

	"github.com/NailUsmanov/api_organization/internal/models"
	"gorm.io/gorm"
)

// EmployeeRepo реализует репозиторий для работы с сотрудниками в базе данных.
type EmployeeRepo struct {
	db *gorm.DB
}

// NewEmployeeRepo создаёт новый экземпляр репозитория сотрудников.
func NewEmployeeRepo(db *gorm.DB) *EmployeeRepo {
	return &EmployeeRepo{db: db}
}

// Create сохраняет нового сотрудника в базе данных.
func (e *EmployeeRepo) Create(ctx context.Context, emp *models.Employee) error {
	return e.db.WithContext(ctx).Create(emp).Error
}

// GetByID возвращает сотрудника по его идентификатору.
func (e *EmployeeRepo) GetByID(ctx context.Context, id uint) (*models.Employee, error) {
	var emp models.Employee
	err := e.db.WithContext(ctx).First(&emp, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &emp, err
}

// Update обновляет существующего сотрудника в базе данных.
func (e *EmployeeRepo) Update(ctx context.Context, emp *models.Employee) error {
	return e.db.WithContext(ctx).Save(emp).Error
}

// Delete удаляет сотрудника по его идентификатору.
func (e *EmployeeRepo) Delete(ctx context.Context, id uint) error {
	return e.db.WithContext(ctx).Delete(&models.Employee{}, id).Error
}

// ListByDepartment возвращает список сотрудников указанного отдела.
func (e *EmployeeRepo) ListByDepartment(ctx context.Context, departmenID uint, orderBy string) ([]models.Employee, error) {
	var emps []models.Employee
	query := e.db.WithContext(ctx).Where("department_id = ?", departmenID)
	if orderBy == "full_name" {
		query = query.Order("full_name")
	} else {
		query = query.Order("created_at")
	}

	err := query.Find(&emps).Error

	return emps, err
}

// MoveToDepartment перемещает всех сотрудников из одного отдела в другой.
func (e *EmployeeRepo) MoveToDepartment(ctx context.Context, departmentID uint, targerDerpartmentID uint) error {
	return e.db.WithContext(ctx).Model(&models.Employee{}).Where("depatment_id").
		Update("department_id", targerDerpartmentID).Error
}

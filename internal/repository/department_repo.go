// Package repository предоставляет реализацию доступа к данным для работы с базой данных.
package repository

import (
	"context"
	"errors"

	"github.com/NailUsmanov/api_organization/internal/models"
	"gorm.io/gorm"
)

// DepartmentRepo реализует репозиторий для работы с подразделениями в базе данных.
type DepartmentRepo struct {
	db *gorm.DB
}

// NewDepartmentRepository создаёт новый экземпляр репозитория подразделений.
func NewDepartmentRepository(db *gorm.DB) *DepartmentRepo {
	return &DepartmentRepo{db: db}
}

// Create сохраняет новое подразделение в базе данных.
func (d *DepartmentRepo) Create(ctx context.Context, dept *models.Department) error {
	return d.db.WithContext(ctx).Create(dept).Error
}

// GetByID возвращает подразделение по его идентификатору.
func (d *DepartmentRepo) GetByID(ctx context.Context, id uint) (*models.Department, error) {
	var dept models.Department
	err := d.db.WithContext(ctx).First(&dept, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dept, err
}

// Update обновляет существующее подразделение в базе данных.
func (d *DepartmentRepo) Update(ctx context.Context, dept *models.Department) error {
	return d.db.WithContext(ctx).Save(dept).Error
}

// Delete удаляет подразделение по его идентификатору.
func (d *DepartmentRepo) Delete(ctx context.Context, id uint) error {
	return d.db.WithContext(ctx).Delete(models.Department{}, id).Error
}

// GetChildren возвращает список прямых дочерних подразделений для указанного родителя.
func (d *DepartmentRepo) GetChildren(ctx context.Context, parentID *uint) ([]models.Department, error) {
	var depts []models.Department
	err := d.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&depts).Error
	return depts, err
}

// GetByNameAndParent возвращает подразделение по его имени и родителю.
func (d *DepartmentRepo) GetByNameAndParent(ctx context.Context, name string, parentID *uint) (*models.Department, error) {
	var dept models.Department
	err := d.db.WithContext(ctx).Where("name = ? AND parent_id = ?", name, parentID).First(&dept).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dept, err
}

// GetSubTree возвращает всех потомков указанного подразделения до заданной глубины.
func (d *DepartmentRepo) GetSubTree(ctx context.Context, rootID uint, depth int) ([]models.Department, error) {
	if depth < 2 {
		return nil, nil
	}
	maxLevel := depth - 1
	query := `
		WITH RECURSIVE dept_tree AS (
			SELECT id, name, parent_id, created_at, 1 AS level
			FROM departments
			WHERE parent_id = $1
			UNION ALL
			SELECT d.id, d.name, d.parent_id, d.created_at, dt.level + 1
			FROM departments d
			INNER JOIN dept_tree dt ON d.parent_id = dt.id
			WHERE dt.level < $2
		)
		SELECT id, name, parent_id, created_at FROM dept_tree ORDER BY id;
	`
	rows, err := d.db.WithContext(ctx).Raw(query, rootID, maxLevel).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var depts []models.Department
	for rows.Next() {
		var d models.Department
		if err := rows.Scan(&d.ID, &d.Name, &d.ParentID, &d.CreatedAt); err != nil {
			return nil, err
		}
		depts = append(depts, d)
	}
	return depts, nil
}

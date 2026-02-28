// Package models содержит GORM-модели данных, соответствующие таблицам в базе данных.
package models

import "time"

// Department представляет модель подразделения (отдела) в организационной структуре.
type Department struct {
	ID        int          `gorm:"primaryKey" json:"id"`
	Name      string       `gorm:"size:200;not null;uniqueIndex:idx_parent_name,priority:2" json:"name"`
	ParentID  *uint        `gorm:"index;uniqueIndex:idx_parent_name,priority:1" json:"parent_id"`
	CreatedAt time.Time    `json:"created_at"`
	Children  []Department `gorm:"foreignkey:ParentID" json:"children,omitempty"`
	Employees []Employee   `json:"employees,omitempty"`
}

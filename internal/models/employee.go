// Package models содержит GORM-модели данных, соответствующие таблицам в базе данных.
package models

import "time"

// Employee представляет модель сотрудника в организационной структуре.
type Employee struct {
	ID           int        `gorm:"primaryKey" json:"id"`
	DepartmentID int        `gorm:"not null;index" json:"department_id"`
	FullName     string     `gorm:"size:200;not null" json:"full_name"`
	Position     string     `gorm:"size:200;not null" json:"position"`
	HiredAt      *time.Time `json:"hired_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

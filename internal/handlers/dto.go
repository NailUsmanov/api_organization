package handlers

import "time"

// createDepartmentRequest представляет структуру JSON-запроса для создания нового подразделения.
type createDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

// updateDepartmentRequest представляет структуру JSON-запроса для обновления подразделения.
type updateDepartmentRequest struct {
	Name     *string `json:"name,omitempty"`
	ParentID *uint   `json:"parent_id,omitempty"`
}

// createEmployeeRequest представляет структуру JSON-запроса для создания нового сотрудника.
type createEmployeeRequest struct {
	FullName string     `json:"full_name"`
	Position string     `json:"position"`
	HiredAt  *time.Time `json:"hired_at"`
}

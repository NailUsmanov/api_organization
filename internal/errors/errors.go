// Package apperrors содержит централизованное определение всех бизнес-ошибок, используемых в приложении.
package apperrors

import "errors"

var (
	// Department errors
	ErrDepartmentNotFound       = errors.New("department not found")
	ErrDepartmentNameConflict   = errors.New("department with this name already exists under the same parent")
	ErrParentNotFound           = errors.New("parent department not found")
	ErrSelfParent               = errors.New("cannot set parent to itself")
	ErrCycleDetected            = errors.New("cannot move department to its own descendant")
	ErrReassignWithChildren     = errors.New("cannot reassign department with children; delete children first or use cascade mode")
	ErrInvalidMode              = errors.New("invalid mode, must be 'cascade' or 'reassign'")
	ErrReassignTargetRequired   = errors.New("reassign_to_department_id is required for reassign mode")
	ErrTargetDepartmentNotFound = errors.New("target department not found")

	// Employee errors
	ErrEmployeeDepartmentNotFound = errors.New("department not found")
	ErrInvalidFullName            = errors.New("full_name must be non-empty and max 200 characters")
	ErrInvalidPosition            = errors.New("position must be non-empty and max 200 characters")
)

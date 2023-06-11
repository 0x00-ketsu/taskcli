package repository

import (
	"time"

	"github.com/0x00-ketsu/taskcli/internal/model"
)

type Task interface {
	// GetTodayTodoCount returns todo items count in today.
	GetTodayTodoCount() int

	// Search returns match tasks with search fields.
	Search(title string, isCompleted, isDeleted bool) ([]model.Task, error)

	// GetByTitle returns a task with title.
	GetByTitle(title string) (model.Task, error)

	// GetAllCompleted returns all completed tasks.
	GetAllCompleted() ([]model.Task, error)

	// GetAllDeleted returns all deleted tasks.
	GetAllDeleted() ([]model.Task, error)

	// GetAllExpired returns all expired tasks.
	GetAllExpired() ([]model.Task, error)

	// GetAllByDate returns all tasks in specific date.
	GetAllByDate(datetime time.Time) ([]model.Task, error)

	// GetAllByDateRange returns all tasks in specific date range.
	GetAllByDateRange(from, to time.Time) ([]model.Task, error)

	// Create creates a new task.
	Create(title, content string) (*model.Task, error)

	// Update updates a task.
	Update(t *model.Task) error

	// UpdateField updates a task with specific field.
	UpdateField(t *model.Task, fieldName string, value interface{}) error

	// Delete deletes a task.
	Delete(t *model.Task) error

	// IsTaskExist accroding task title check is alreay in Today's task list.
	IsTaskExist(title string) bool
}

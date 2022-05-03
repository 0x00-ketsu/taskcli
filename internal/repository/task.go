package repository

import (
	"time"

	"github.com/0x00-ketsu/taskcli/internal/model"
)

type Task interface {
	GetTodayTodoCount() int
	IsTaskExist(title string) bool

	Search(title string, isCompleted, isDeleted bool) ([]model.Task, error)

	GetByID(ID string) (model.Task, error)
	GetByTitle(title string) (model.Task, error)

	GetAllCompleted() ([]model.Task, error)
	GetAllDeleted() ([]model.Task, error)
	GetAllExpired() ([]model.Task, error)
	GetAllByDate(datetime time.Time) ([]model.Task, error)
	GetAllByDateRange(from, to time.Time) ([]model.Task, error)

	Create(title, content string) (*model.Task, error)

	Update(t *model.Task) error
	UpdateField(t *model.Task, fieldName string, value interface{}) error

	Delete(t *model.Task) error
}

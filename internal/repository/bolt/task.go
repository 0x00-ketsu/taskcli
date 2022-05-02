package bolt

import (
	"errors"
	"fmt"
	"time"

	"github.com/0x00-ketsu/taskcli/internal/model"
	"github.com/0x00-ketsu/taskcli/internal/repository"
	"github.com/0x00-ketsu/taskcli/internal/utils"
	"github.com/asdine/storm/q"
	"github.com/asdine/storm/v3"
)

type TaskBolt struct {
	DB *storm.DB
}

// NewTask creates an object
func NewTask(db *storm.DB) repository.Task {
	return &TaskBolt{db}
}

// GetTodayTodoCount returns todo items count in today
func (t *TaskBolt) GetTodayTodoCount() int {
	var tasks []model.Task

	today := utils.ToDate(time.Now())
	err := t.DB.
		Select(q.Eq("IsDeleted", false), q.Eq("DueDate", today), q.Eq("IsCompleted", false)).
		Find(&tasks)
	if err != nil {
		return 0
	} else {
		return len(tasks)
	}
}

// GetByID returns a task with ID
func (t *TaskBolt) GetByID(ID string) (model.Task, error) {
	panic("Not implemented error")
}

// GetByTitle returns a task with title
func (t *TaskBolt) GetByTitle(title string) (model.Task, error) {
	return t.getOneByField("Title", title)
}

// GetAllCompleted returns all completed tasks
func (t *TaskBolt) GetAllCompleted() ([]model.Task, error) {
	var tasks []model.Task

	err := t.DB.
		Select(q.Eq("IsCompleted", true), q.Eq("IsDeleted", false)).
		OrderBy("CreatedAt").
		Reverse().
		Find(&tasks)

	return tasks, err
}

// GetAllExpired returns all expired tasks
func (t *TaskBolt) GetAllExpired() ([]model.Task, error) {
	var tasks []model.Task

	today := utils.ToDate(time.Now())
	err := t.DB.
		Select(q.Lt("DueDate", today), q.Eq("IsCompleted", false)).
		OrderBy("CreatedAt").
		Reverse().
		Find(&tasks)

	return tasks, err
}

// GetAllDeleted returns all deleted tasks
func (t *TaskBolt) GetAllDeleted() ([]model.Task, error) {
	var tasks []model.Task

	err := t.DB.
		Select(q.Eq("IsDeleted", true)).
		OrderBy("CreatedAt").
		Reverse().
		Find(&tasks)

	return tasks, err
}

// GetAllByDate returns all tasks in specific date
func (t *TaskBolt) GetAllByDate(datetime time.Time) ([]model.Task, error) {
	start := utils.GetDatetimeStart(datetime)
	end := utils.GetDatetimeEnd(datetime)

	return t.GetAllByDateRange(start, end)
}

// GetAllByDateRange returns all tasks in specific date range
func (t *TaskBolt) GetAllByDateRange(from, to time.Time) ([]model.Task, error) {
	var tasks []model.Task

	err := t.DB.
		Select(q.Gte("DueDate", from), q.Lte("DueDate", to), q.Eq("IsDeleted", false)).
		OrderBy("CreatedAt").
		Reverse().
		Find(&tasks)

	return tasks, err
}

// Create creates a new task
func (t *TaskBolt) Create(title, content string) (*model.Task, error) {
	if t.IsTaskExist(title) {
		msg := fmt.Sprintf("task title: %v is alreay exist", title)
		return nil, errors.New(msg)

	}

	// Create
	today := utils.ToDate(time.Now())
	task := model.Task{
		Title:     title,
		Content:   content,
		DueDate:   today,
		CreatedAt: time.Now(),
	}

	err := t.DB.Save(&task)
	return &task, err
}

// Update updates a task
func (t *TaskBolt) Update(task *model.Task) error {
	return t.DB.Update(task)
}

// UpdateField updates a task with specific field
func (t *TaskBolt) UpdateField(task *model.Task, fieldName string, value interface{}) error {
	return t.DB.UpdateField(task, fieldName, value)
}

// Delete deletes a task
func (t *TaskBolt) Delete(task *model.Task) error {
	task.IsDeleted = true
	return t.Update(task)
}

// IsTaskExist accroding task title check is alreay in Today's task list
func (t *TaskBolt) IsTaskExist(title string) bool {
	var tasks []model.Task
	var err error

	today := utils.ToDate(time.Now())
	err = t.DB.
		Select(q.Eq("IsDeleted", false), q.Eq("Title", title), q.Eq("DueDate", today)).
		Find(&tasks)

	if err == nil && len(tasks) > 0 {
		return true
	} else {
		return false
	}
}

func (t *TaskBolt) getOneByField(fieldName string, value interface{}) (model.Task, error) {
	var task model.Task
	err := t.DB.Select(q.Eq(fieldName, value), q.Eq("IsDeleted", false)).First(&task)

	return task, err
}

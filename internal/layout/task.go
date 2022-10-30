package layout

import (
	"fmt"
	"time"

	"github.com/0x00-ketsu/taskcli/internal/model"
	"github.com/asdine/storm/v3"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	TODO_HEADER      = "[::b]> Todo items[-::]"
	COMPLETED_HEADER = "[green::b]> Completed[-::]"
	EXPIRED_HEADER   = "[yellow::b]> Expired[-::]"
	DELETED_HEADER   = "[red::b]> Deleted[-::]"
)

type TaskView struct {
	*tview.Flex
	list *tview.List
	hint *tview.TextView

	filter     string // Filter view select type
	newTask    *tview.InputField
	activeTask *model.Task
	tasks      []model.Task
}

func NewTaskView() *TaskView {
	view := TaskView{
		Flex:    tview.NewFlex().SetDirection(tview.FlexRow),
		list:    tview.NewList().ShowSecondaryText(false),
		hint:    tview.NewTextView().SetTextColor(tcell.ColorYellow).SetTextAlign(tview.AlignCenter),
		filter:  "",
		newTask: taskInputFiled("+ Add a Task to Today's List"),
	}

	view.AddItem(view.newTask, 1, 0, false).
		AddItem(view.list, 0, 1, true).
		AddItem(view.hint, 0, 1, false)

	view.SetBorder(true).SetTitle(" Tasks ")
	view.setHintMessage()

	// Select task
	view.list.SetBorderPadding(1, 0, 0, 0)

	// Create new task
	view.newTask.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			name := view.newTask.GetText()
			if !validateTaskName(name) {
				return
			}

			// Save new task to db & update Tasks
			_, err := taskRepo.Create(name, "")
			if err != nil {
				msg := fmt.Sprintf("[red]Create task failed, error: %v", err.Error())
				statusView.showForSeconds(msg, 5)
				return
			}
			statusView.showForSeconds("[green]Task create successful", 5)

			// Reload tasks
			view.reloadTasks()

			// Reset
			view.newTask.SetText("")
			app.SetFocus(view)
			return
		case tcell.KeyEsc:
			view.newTask.SetText("")
			app.SetFocus(view)
		}
	})

	return &view
}

// Reload tasks based on selected Filter
func (p *TaskView) reloadTasks() {
	filter := p.filter
	tasks, _ := p.getFiterTasks(filter)
	p.renderTaskList(tasks)

	// update Today view
	todayView.updateTodoCount()
}

// Loads tasks based on selected item in Filter view
// choices: today, tomorrow, last 7 days
func (p *TaskView) loadFilterTasks(filter string) {
	p.clearTaskList()
	p.RemoveItem(p.hint)
	p.filter = filter

	tasks, err := p.getFiterTasks(filter)
	if err == storm.ErrIdxNotFound {
		p.AddItem(p.hint, 0, 1, false)
		statusView.showForSeconds("[yellow]No task in list - "+filter, 5)
	} else if err != nil {
		p.AddItem(p.hint, 0, 1, false)
		statusView.showForSeconds("[red]Error: "+err.Error(), 5)
	} else {
		p.RemoveItem(p.hint)
		p.renderTaskList(tasks)
		statusView.showForSeconds("[yellow]Displaying tasks of "+filter, 3)
	}

	app.SetFocus(taskView)
	removeTaskDetailView()
}

// Get tasks based on selected item in Filter view
func (p *TaskView) getFiterTasks(filter string) ([]model.Task, error) {
	var tasks []model.Task
	var err error

	switch filter {
	case "today":
		tasks, err = taskRepo.GetAllByDate(today)

	case "tomorrow":
		tasks, err = taskRepo.GetAllByDate(tomorrow)

	case "last 7 days":
		week := today.Add(time.Hour * 7 * 24)
		tasks, err = taskRepo.GetAllByDateRange(today, week)

	case "completed":
		tasks, err = taskRepo.GetAllCompleted()

	case "expired":
		tasks, err = taskRepo.GetAllExpired()

	case "trash":
		tasks, err = taskRepo.GetAllDeleted()
	}

	return tasks, err
}

// Render task list
// - todo
// - completed
func (p *TaskView) renderTaskList(tasks []model.Task) {
	p.clearTaskList()

	switch p.filter {
	case "today":
		p.classifyTasks(tasks)
	case "tomorrow":
		p.classifyTasks(tasks)
	case "last 7 days":
		p.classifyTasks(tasks)
	case "completed":
		p.unclassifyTasks(tasks)
	case "expired":
		p.unclassifyTasks(tasks)
	case "trash":
		p.unclassifyTasks(tasks)
	case "search":
		p.classifyTasks(tasks)
	}

	// Remove hint
	if p.list.GetItemCount() > 0 {
		p.RemoveItem(p.hint)
	}
}

// Classify tasks to: todo, completed
func (p *TaskView) classifyTasks(tasks []model.Task) {
	var todoTasks, completedTasks, expiredTasks, deletedTasks []model.Task

	for _, task := range tasks {
		if task.IsDeleted {
			deletedTasks = append(deletedTasks, task)
		} else if task.IsCompleted {
			completedTasks = append(completedTasks, task)
		} else {
			if task.DueDate.Before(today) {
				expiredTasks = append(expiredTasks, task)
			} else {
				todoTasks = append(todoTasks, task)
			}
		}
	}

	var text string
	emptyTask := model.Task{}
	// Todo tasks
	todoCount := len(todoTasks)
	if todoCount > 0 {
		p.list.AddItem(fmt.Sprintf("%s [green::](%d)[-::]", TODO_HEADER, todoCount), "", 0, nil)
		p.addTaskToList(emptyTask)

		for _, task := range todoTasks {
			p.addTaskToList(task)

			dueDate := task.DueDate.Format("2006-01-02")
			if p.filter == "search" {
				text = fmt.Sprintf("%v    [lime::i]-- Due: %v", renderTaskTitle(task), dueDate)
			} else {
				text = renderTaskTitle(task)
			}
			p.list.AddItem(text, "", 0, func() func() {
				return func() { p.activateTask(task) }
			}())
		}

		// Select first todo task item
		p.list.SetCurrentItem(1)
	}

	// Completed tasks
	if len(completedTasks) > 0 {
		if len(todoTasks) > 0 {
			p.list.AddItem("", "", 0, nil)
			p.addTaskToList(emptyTask)
		}

		p.list.AddItem(COMPLETED_HEADER, "", 0, nil)
		p.addTaskToList(emptyTask)

		for _, task := range completedTasks {
			p.addTaskToList(task)

			dueDate := task.DueDate.Format("2006-01-02")
			if p.filter == "search" {
				text = fmt.Sprintf("%v    [lime::i]-- Due: %v", renderTaskTitle(task), dueDate)
			} else {
				text = renderTaskTitle(task)
			}
			p.list.AddItem(text, "", 0, func() func() {
				return func() { p.activateTask(task) }
			}())
		}
	}

	// Expired tasks
	if len(expiredTasks) > 0 {
		p.list.AddItem(EXPIRED_HEADER, "", 0, nil)
		p.addTaskToList(emptyTask)

		for _, task := range expiredTasks {
			p.addTaskToList(task)

			dueDate := task.DueDate.Format("2006-01-02")
			if p.filter == "search" {
				text = fmt.Sprintf("%v    [lime::i]-- Due: %v", renderTaskTitle(task), dueDate)
			} else {
				text = renderTaskTitle(task)
			}
			p.list.AddItem(text, "", 0, func() func() {
				return func() { p.activateTask(task) }
			}())
		}
	}

	// Deleted tasks
	if len(deletedTasks) > 0 {
		p.list.AddItem(DELETED_HEADER, "", 0, nil)
		p.addTaskToList(emptyTask)

		for _, task := range deletedTasks {
			p.addTaskToList(task)

			dueDate := task.DueDate.Format("2006-01-02")
			if p.filter == "search" {
				text = fmt.Sprintf("%v    [lime::i]-- Due: %v", renderTaskTitle(task), dueDate)
			} else {
				text = renderTaskTitle(task)
			}
			p.list.AddItem(text, "", 0, func() func() {
				return func() { p.activateTask(task) }
			}())
		}
	}
}

// Unclassify tasks, display all tasks
// Extra add due date
func (p *TaskView) unclassifyTasks(tasks []model.Task) {
	for _, task := range tasks {
		p.addTaskToList(task)

		dueDate := task.DueDate.Format("2006-01-02")
		text := fmt.Sprintf("%v    [lime::i]-- Due: %v", renderTaskTitle(task), dueDate)
		p.list.AddItem(text, "", 0, func() func() {
			return func() { p.activateTask(task) }
		}())
	}
}

// Add a task to task list
func (p *TaskView) addTaskToList(task model.Task) {
	p.tasks = append(p.tasks, task)
}

// Remove all task items from Task view
func (p *TaskView) clearTaskList() {
	p.list.Clear()
	p.tasks = nil
	p.activeTask = nil
}

// Marks a task is actived & loads detail in Task Detail view
func (p *TaskView) activateTask(task model.Task) {
	removeTaskDetailView()

	focusTask := p.getFocusTask()
	taskDetailView.setTask(focusTask)

	main.AddItem(taskDetailView, 0, 4, false)
	app.SetFocus(taskDetailView)
}

func (p *TaskView) setHintMessage() {
	if len(p.tasks) == 0 {
		p.hint.SetText("Welcome to Taskcli!\n------------------------------\n Create Task at the top of Tasks view.\n (Press n)")
	} else {
		p.hint.SetText("Select a Task (Press Enter) to load task detail.\nOr create a new Task (Press n).")
	}
}

// Move to next item
// Skip sepreation line
func (p *TaskView) lineDown() {
	curItemIndex := p.list.GetCurrentItem()
	itemCount := p.list.GetItemCount()

	if curItemIndex >= 0 && curItemIndex < itemCount-1 {
		nextItemIndex := curItemIndex + 1
		nextItemText, _ := p.list.GetItemText(nextItemIndex)
		if nextItemText == "" {
			nextItemIndex += 1
		}
		p.list.SetCurrentItem(nextItemIndex)
	}
}

// Move to previous item
// Skip sepreation line
func (p *TaskView) lineUp() {
	curItemIndex := p.list.GetCurrentItem()
	itemCount := p.list.GetItemCount()

	if curItemIndex < itemCount && curItemIndex > 0 {
		prevItemIndex := curItemIndex - 1
		prevItemText, _ := p.list.GetItemText(prevItemIndex)
		if prevItemText == "" {
			prevItemIndex -= 1
		}
		p.list.SetCurrentItem(prevItemIndex)
	}
}

// Get current focus task(pointer) item in Task view
func (p *TaskView) getFocusTask() *model.Task {
	curItemIndex := p.list.GetCurrentItem()
	return &p.tasks[curItemIndex]
}

// Toggle task status (completed / uncompleted)
func (p *TaskView) toggleTaskStatus(task *model.Task) {
	status := !task.IsCompleted
	if taskRepo.UpdateField(task, "IsCompleted", status) == nil {
		task.IsCompleted = status
		// reload
		p.reloadTasks()
		// update Today view
		todayView.updateTodoCount()
	}
}

// Rename current focused task title
func (p *TaskView) renameCurrentTask(task *model.Task, newTitle string) {
	if !validateTaskName(newTitle) {
		return
	}

	if taskRepo.IsTaskExist(newTitle) {
		msg := fmt.Sprintf("[red]Task title: %v is already exist", newTitle)
		statusView.showForSeconds(msg, 5)
	} else {
		originalTitle := task.Title
		task.Title = newTitle
		if err := taskRepo.Update(task); err != nil {
			msg := fmt.Sprintf("[red]Update task title: %v failed, error: %v", newTitle, err.Error())
			statusView.showForSeconds(msg, 5)
		} else {
			msg := fmt.Sprintf("[green]Update task title[%s] -> %s successful", originalTitle, newTitle)
			statusView.showForSeconds(msg, 5)
			// reload
			p.reloadTasks()
		}
	}
}

// Delete current focused task
func (p *TaskView) deleteCurrentTask() {
	task := p.getFocusTask()
	if err := taskRepo.Delete(task); err != nil {
		msg := fmt.Sprintf("[red]Delete task: %v failed, error: %v", task.Title, err.Error())
		statusView.showForSeconds(msg, 5)
	} else {
		msg := fmt.Sprintf("[green]Delete task: %v successful", task.Title)
		statusView.showForSeconds(msg, 5)
		// reload
		p.reloadTasks()
		// update Today view
		todayView.updateTodoCount()
	}
}

// Restore current focused task
// For deleted task
func (p *TaskView) restoreCurrentTask() {
	task := p.getFocusTask()
	if err := taskRepo.UpdateField(task, "IsDeleted", false); err != nil {
		msg := fmt.Sprintf("[red]Restore task: %v failed, error: %v", task.Title, err.Error())
		statusView.showForSeconds(msg, 5)
	} else {
		msg := fmt.Sprintf("[green]Restore task: %v successful", task.Title)
		statusView.showForSeconds(msg, 5)
		// reload
		p.reloadTasks()
		// update Today view
		todayView.updateTodoCount()
	}
}

// Render task title
func renderTaskTitle(task model.Task) string {
	checkbox := "[ []"

	if task.IsCompleted {
		checkbox = "[x[]"
	}

	return fmt.Sprintf(" [%s]%s %s", getTaskTitleColor(task), checkbox, task.Title)
}

func taskInputFiled(placeholder string) *tview.InputField {
	return tview.NewInputField().
		SetPlaceholder(placeholder).
		SetPlaceholderTextColor(tcell.ColorWhite).
		SetFieldTextColor(tcell.ColorBlack).
		SetFieldBackgroundColor(tcell.ColorLightBlue)
}

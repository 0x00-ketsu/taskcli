package layout

import (
	"fmt"
	"time"

	"github.com/0x00-ketsu/taskcli/internal/global"
	"github.com/0x00-ketsu/taskcli/internal/model"
	"github.com/asdine/storm/v3"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	TODO_HEADER      = "[::b]> Todo items"
	COMPLETE_DHEADER = "[green::b]> Completed"
)

type TaskPanel struct {
	*tview.Flex
	list *tview.List
	hint *tview.TextView

	filter     string // Filter panel select type
	newTask    *tview.InputField
	activeTask *model.Task
	tasks      []model.Task
}

func NewTaskPanel() *TaskPanel {
	app := global.App

	panel := TaskPanel{
		Flex:    tview.NewFlex().SetDirection(tview.FlexRow),
		list:    tview.NewList().ShowSecondaryText(false),
		hint:    tview.NewTextView().SetTextColor(tcell.ColorYellow).SetTextAlign(tview.AlignCenter),
		filter:  "",
		newTask: taskInputFiled("+ Add a Task to Today's List"),
	}

	panel.AddItem(panel.newTask, 1, 0, false).
		AddItem(panel.list, 0, 1, true).
		AddItem(panel.hint, 0, 1, false)

	panel.SetBorder(true).SetTitle(" Tasks ")
	panel.setHintMessage()

	// Select task
	panel.list.SetBorderPadding(1, 0, 0, 0)

	// Create new task
	panel.newTask.SetDoneFunc(func(key tcell.Key) {
		repo := global.TaskRepo

		switch key {
		case tcell.KeyEnter:
			name := panel.newTask.GetText()
			if !validateTaskName(name) {
				return
			}

			// Save new task to db & update Tasks
			_, err := repo.Create(name, "")
			if err != nil {
				msg := fmt.Sprintf("[red]Create task failed, error: %v", err.Error())
				statusPanel.showForSeconds(msg, 5)
				return
			}
			statusPanel.showForSeconds("[green]Task create successful", 5)

			// Reload tasks
			panel.reloadTasks()

			// Reset
			panel.newTask.SetText("")
			app.SetFocus(panel)
			return
		case tcell.KeyEsc:
			panel.newTask.SetText("")
			app.SetFocus(panel)
		}
	})

	return &panel
}

// Reload tasks based on selected Filter
func (p *TaskPanel) reloadTasks() {
	filter := p.filter
	tasks, _ := p.getFiterTasks(filter)
	p.renderTaskList(tasks)

	// update Today panel
	todayPanel.updateTodoCount()
}

// Loads tasks based on selected item in Filter panel
// choices: today, tomorrow, last 7 days
func (p *TaskPanel) loadFilterTasks(filter string) {
	app := global.App

	p.clearTaskList()
	p.RemoveItem(p.hint)
	p.filter = filter

	tasks, err := p.getFiterTasks(filter)
	if err == storm.ErrIdxNotFound {
		p.AddItem(p.hint, 0, 1, false)
		statusPanel.showForSeconds("[yellow]No task in list - "+filter, 5)
	} else if err != nil {
		p.AddItem(p.hint, 0, 1, false)
		statusPanel.showForSeconds("[red]Error: "+err.Error(), 5)
	} else {
		p.RemoveItem(p.hint)
		p.renderTaskList(tasks)
		statusPanel.showForSeconds("[yellow]Displaying tasks of "+filter, 3)
	}

	app.SetFocus(taskPanel)
	removeTaskDetailPanel()
}

// Get tasks based on selected item in Filter panel
func (p *TaskPanel) getFiterTasks(filter string) ([]model.Task, error) {
	var tasks []model.Task
	var err error

	repo := global.TaskRepo

	switch filter {
	case "today":
		tasks, err = repo.GetAllByDate(today)

	case "tomorrow":
		tasks, err = repo.GetAllByDate(tomorrow)

	case "last 7 days":
		week := today.Add(time.Hour * 7 * 24)
		tasks, err = repo.GetAllByDateRange(today, week)

	case "completed":
		tasks, err = repo.GetAllCompleted()

	case "expired":
		tasks, err = repo.GetAllExpired()

	case "trash":
		tasks, err = repo.GetAllDeleted()
	}

	return tasks, err
}

// Render task list
// - todo
// - completed
func (p *TaskPanel) renderTaskList(tasks []model.Task) {
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
	}

	// Remove hint
	if p.list.GetItemCount() > 0 {
		p.RemoveItem(p.hint)
	}
}

// Classify tasks to: todo, completed
func (p *TaskPanel) classifyTasks(tasks []model.Task) {
	var todoTasks, completedTasks []model.Task

	for _, task := range tasks {
		if task.IsCompleted {
			completedTasks = append(completedTasks, task)
		} else {
			todoTasks = append(todoTasks, task)
		}
	}

	emptyTask := model.Task{}

	// Todo tasks
	if len(todoTasks) > 0 {
		p.list.AddItem(TODO_HEADER, "", 0, nil)
		p.addTaskToList(emptyTask)

		for _, task := range todoTasks {
			p.addTaskToList(task)
			p.list.AddItem(renderTaskTitle(task), "", 0, func() func() {
				return func() { p.activateTask(task) }
			}())
		}
	}

	// Completed tasks
	if len(completedTasks) > 0 {
		if len(todoTasks) > 0 {
			p.list.AddItem("", "", 0, nil)
			p.addTaskToList(emptyTask)
		}

		p.list.AddItem(COMPLETE_DHEADER, "", 0, nil)
		p.addTaskToList(emptyTask)

		for _, task := range completedTasks {
			p.addTaskToList(task)
			p.list.AddItem(renderTaskTitle(task), "", 0, func() func() {
				return func() { p.activateTask(task) }
			}())
		}
	}
}

// Unclassify tasks, display all tasks
// Extra add due date
func (p *TaskPanel) unclassifyTasks(tasks []model.Task) {
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
func (p *TaskPanel) addTaskToList(task model.Task) {
	p.tasks = append(p.tasks, task)
}

// Remove all task items from Task panel
func (p *TaskPanel) clearTaskList() {
	p.list.Clear()
	p.tasks = nil
	p.activeTask = nil
}

// Marks a task is actived & loads detail in Task Detail panel
func (p *TaskPanel) activateTask(task model.Task) {
	app := global.App

	removeTaskDetailPanel()

	focusTask := p.getFocusTask()
	taskDetailPanel.setTask(focusTask)

	main.AddItem(taskDetailPanel, 0, 4, false)
	app.SetFocus(taskDetailPanel)
}

func (p *TaskPanel) setHintMessage() {
	if len(p.tasks) == 0 {
		p.hint.SetText("Welcome to Taskcli!\n------------------------------\n Create Task at the top of Tasks panel.\n (Press n)")
	} else {
		p.hint.SetText("Select a Task (Press Enter) to load task detail.\nOr create a new Task (Press n).")
	}
}

// Move to next item
// Skip sepreation line
func (p *TaskPanel) lineDown() {
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
func (p *TaskPanel) lineUp() {
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

// Get current focus task(pointer) item in Task panel
func (p *TaskPanel) getFocusTask() *model.Task {
	curItemIndex := p.list.GetCurrentItem()
	return &p.tasks[curItemIndex]
}

// Toggle task status (completed / uncompleted)
func (p *TaskPanel) toggleTaskStatus(task *model.Task) {
	repo := global.TaskRepo
	status := !task.IsCompleted
	if repo.UpdateField(task, "IsCompleted", status) == nil {
		task.IsCompleted = status
		// reload
		p.reloadTasks()
		// update Today panel
		todayPanel.updateTodoCount()
	}
}

// Rename current focused task title
func (p *TaskPanel) renameCurrentTask(task *model.Task, newTitle string) {
	repo := global.TaskRepo

	if !validateTaskName(newTitle) {
		return
	}

	if repo.IsTaskExist(newTitle) {
		msg := fmt.Sprintf("[red]Task title: %v is already exist", newTitle)
		statusPanel.showForSeconds(msg, 5)
	} else {
		originalTitle := task.Title
		task.Title = newTitle
		if err := repo.Update(task); err != nil {
			msg := fmt.Sprintf("[red]Update task title: %v failed, error: %v", newTitle, err.Error())
			statusPanel.showForSeconds(msg, 5)
		} else {
			msg := fmt.Sprintf("[green]Update task title[%s] -> %s successful", originalTitle, newTitle)
			statusPanel.showForSeconds(msg, 5)
			// reload
			p.reloadTasks()
		}
	}
}

// Delete current focused task
func (p *TaskPanel) deleteCurrentTask() {
	task := p.getFocusTask()
	if err := global.TaskRepo.Delete(task); err != nil {
		msg := fmt.Sprintf("[red]Delete task: %v failed, error: %v", task.Title, err.Error())
		statusPanel.showForSeconds(msg, 5)
	} else {
		msg := fmt.Sprintf("[green]Delete task: %v successful", task.Title)
		statusPanel.showForSeconds(msg, 5)
		// reload
		p.reloadTasks()
		// update Today panel
		todayPanel.updateTodoCount()
	}
}

// Restore current focused task
// For deleted task
func (p *TaskPanel) restoreCurrentTask() {
	task := p.getFocusTask()
	if err := global.TaskRepo.UpdateField(task, "IsDeleted", false); err != nil {
		msg := fmt.Sprintf("[red]Restore task: %v failed, error: %v", task.Title, err.Error())
		statusPanel.showForSeconds(msg, 5)
	} else {
		msg := fmt.Sprintf("[green]Restore task: %v successful", task.Title)
		statusPanel.showForSeconds(msg, 5)
		// reload
		p.reloadTasks()
		// update Today panel
		todayPanel.updateTodoCount()
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

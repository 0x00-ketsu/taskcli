package layout

import (
	"github.com/0x00-ketsu/taskcli/internal/global"
	"github.com/rivo/tview"
)

// Store form input values
var (
	titleVal        string
	isDeletedFlag   bool
	isCompletedFlag bool
)

type SearchPanel struct {
	*tview.Flex
	form *tview.Form
}

func NewSearchPanel() *SearchPanel {
	panel := SearchPanel{
		Flex: tview.NewFlex(),
		form: tview.NewForm(),
	}
	panel.init()

	panel.AddItem(panel.form, 0, 1, true)
	panel.SetBorder(true).SetTitle(" Search ")

	return &panel
}

func (p *SearchPanel) init() {
	p.form.
		AddInputField("Title:", "", 20, nil, getInputTitle).
		AddCheckbox("IsCompleted:", false, getChkIsCompleted).
		AddCheckbox("IsDeleted:", false, getChkIsDeleted).
		AddButton("Search", p.search).
		AddButton("Reset", p.reset)
}

// Search tasks
func (p *SearchPanel) search() {
	taskPanel.clearTaskList()
	taskPanel.RemoveItem(taskPanel.hint)

	tasks, err := global.TaskRepo.Search(titleVal, isCompletedFlag, isDeletedFlag)
	if err != nil {
		statusPanel.showForSeconds("[red]Search failed, error: "+err.Error(), 5)
		return
	} else {
		taskPanel.filter = "search"
		taskPanel.RemoveItem(taskPanel.hint)
		taskPanel.renderTaskList(tasks)
		statusPanel.showForSeconds("[yellow]Displaying tasks of search", 3)

		app := global.App
		app.SetFocus(taskPanel)
		removeTaskDetailPanel()
	}
}

// Reset form
func (p *SearchPanel) reset() {
	p.form.Clear(true)
	p.init()

	app := global.App
	app.SetFocus(p)
}

// Assign input title to `titleVal`
func getInputTitle(text string) {
	titleVal = text
}

// Assign checkbox IsCompleted to `isCompletedFlag`
func getChkIsCompleted(checked bool) {
	isCompletedFlag = checked
}

// Assign checkbox IsDeleted to `isDeletedFlag`
func getChkIsDeleted(checked bool) {
	isDeletedFlag = checked
}

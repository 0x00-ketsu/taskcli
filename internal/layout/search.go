package layout

import (
	"github.com/rivo/tview"
)

// Store form input values
var (
	titleVal        string
	isDeletedFlag   bool
	isCompletedFlag bool
)

type SearchView struct {
	*tview.Flex
	form *tview.Form
}

func NewSearchView() *SearchView {
	view := SearchView{
		Flex: tview.NewFlex(),
		form: tview.NewForm(),
	}
	view.init()

	view.AddItem(view.form, 0, 1, true)
	view.SetBorder(true).SetTitle(" Search ")

	return &view
}

func (p *SearchView) init() {
	p.form.
		AddInputField("Title:", "", 20, nil, getInputTitle).
		AddCheckbox("IsCompleted:", false, getChkIsCompleted).
		AddCheckbox("IsDeleted:", false, getChkIsDeleted).
		AddButton("Search", p.search).
		AddButton("Reset", p.reset)
}

// Search tasks
func (p *SearchView) search() {
	taskView.clearTaskList()
	taskView.RemoveItem(taskView.hint)

	tasks, err := taskRepo.Search(titleVal, isCompletedFlag, isDeletedFlag)
	if err != nil {
		statusView.showForSeconds("[red]Search failed, error: "+err.Error(), 5)
		return
	} else {
		taskView.filter = "search"
		taskView.RemoveItem(taskView.hint)
		taskView.renderTaskList(tasks)
		statusView.showForSeconds("[yellow]Displaying tasks of search", 3)

		app.SetFocus(taskView)
		removeTaskDetailView()
	}
}

// Reset form
func (p *SearchView) reset() {
	p.form.Clear(true)
	p.init()

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

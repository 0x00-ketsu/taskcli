package layout

import (
	"github.com/0x00-ketsu/taskcli/internal/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Display menu choices
type choice struct {
	*tview.Flex
	title *tview.TextView
	items *tview.List
}

type MenuView struct {
	*tview.Pages
	choice      *choice
	renameInput *tview.InputField
}

type menuItem struct {
	name     string
	shortcut rune
	handler  func() func()
}

const (
	MENU_PAGE   = "default"
	RENAME_PAGE = "renameInput"
)

// Common task view menu items
var commonTaskMenus = []menuItem{
	{
		name:     "reload tasks",
		shortcut: 'l',
		handler: func() func() {
			return func() {
				taskView.reloadTasks()

				removeTaskDetailView()
				closeMenuView()
			}
		},
	},
	{
		name:     "rename the current task title",
		shortcut: 'r',
		handler: func() func() {
			return func() {
				task := taskView.getFocusTask()
				currentTitle := task.Title

				menuView.Pages.SwitchToPage(RENAME_PAGE)
				menuView.renameInput.SetLabel("New task title: ").
					SetText(currentTitle).
					SetPlaceholder("Enter new task title").
					SetPlaceholderTextColor(tcell.ColorDarkSlateBlue).
					SetFieldTextColor(tcell.ColorBlack).
					SetFieldBackgroundColor(tcell.ColorLightBlue)

				menuView.renameInput.SetDoneFunc(func(key tcell.Key) {
					switch key {
					case tcell.KeyEnter:
						newTitle := menuView.renameInput.GetText()
						taskView.renameCurrentTask(task, newTitle)
						menuView.Pages.SwitchToPage(MENU_PAGE)

						removeTaskDetailView()
						closeMenuView()
						return
					case tcell.KeyEsc:
						menuView.Pages.SwitchToPage(MENU_PAGE)
					}
				})
			}
		},
	},
}

// Normal task view menu items
var normalTaskMenus = []menuItem{
	{
		name:     "toggle the current task status",
		shortcut: 't',
		handler: func() func() {
			return func() {
				task := taskView.getFocusTask()
				taskView.toggleTaskStatus(task)

				removeTaskDetailView()
				closeMenuView()
			}
		},
	},
	{
		name:     "[red]delete[white] the current task",
		shortcut: 'd',
		handler: func() func() {
			return func() {
				taskView.deleteCurrentTask()

				removeTaskDetailView()
				closeMenuView()
			}
		},
	},
}

// Deleted task view menu items
var deletedTaskMenus = []menuItem{
	{
		name:     "[green]restore[white] the current task",
		shortcut: 's',
		handler: func() func() {
			return func() {
				taskView.restoreCurrentTask()

				removeTaskDetailView()
				closeMenuView()
			}
		},
	},
}

// NewMenuView displays menu view
// NOTE: current only works for Task
func NewMenuView() *MenuView {
	view := MenuView{
		Pages:       tview.NewPages(),
		renameInput: tview.NewInputField(),
	}
	view.SetBorder(true)

	return &view
}

// Open Menu view and focus
func (p *MenuView) open() {
	focusTask := taskView.getFocusTask()
	p.loadMenus(*focusTask)

	curItemIndex := taskView.list.GetCurrentItem()
	curItemText, _ := taskView.list.GetItemText(curItemIndex)
	if curItemText != TODO_HEADER && curItemText != COMPLETED_HEADER && curItemText != "" {
		layout.
			RemoveItem(statusView).
			AddItem(p, p.getSize(), 1, true).
			AddItem(statusView, 1, 1, false)
		app.SetFocus(p)
	}
}

func (p *MenuView) loadMenus(task model.Task) {
	choice := &choice{
		Flex:  tview.NewFlex().SetDirection(tview.FlexRow),
		title: tview.NewTextView(),
		items: tview.NewList().ShowSecondaryText(false),
	}
	choice.title.SetText("Menu. Use j/k/enter, or the shortcuts indicated\n===============================================\n")
	choice.addMenuToList(task)

	choice.
		AddItem(choice.title, 0, 1, false).
		AddItem(choice.items, choice.items.GetItemCount(), 1, true)

	p.choice = choice

	p.Pages.AddPage(MENU_PAGE, p.choice, true, true)
	p.Pages.AddPage(RENAME_PAGE, p.renameInput, true, false)
}

// Close Menu view and focus Task view
func closeMenuView() {
	layout.RemoveItem(menuView)
	app.SetFocus(taskView)
}

// Returns the size of MenuView rows
func (p *MenuView) getSize() int {
	return p.choice.items.GetItemCount() + 4
}

// Add menu choices to Menu view
func (c *choice) addMenuToList(task model.Task) {
	for _, menu := range commonTaskMenus {
		c.items.AddItem(menu.name, "", menu.shortcut, menu.handler())
	}

	// for normal tasks
	if !task.IsDeleted {
		for _, menu := range normalTaskMenus {
			c.items.AddItem(menu.name, "", menu.shortcut, menu.handler())
		}
	}

	// for deleted tasks
	if task.IsDeleted {
		for _, menu := range deletedTaskMenus {
			c.items.AddItem(menu.name, "", menu.shortcut, menu.handler())
		}
	}
}

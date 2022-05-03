package layout

import (
	"github.com/0x00-ketsu/taskcli/internal/global"
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

type MenuPanel struct {
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

// Common task panel menu items
var commonTaskMenus = []menuItem{
	{
		name:     "reload tasks",
		shortcut: 'l',
		handler: func() func() {
			return func() {
				taskPanel.reloadTasks()

				removeTaskDetailPanel()
				closeMenuPanel()
			}
		},
	},
	{
		name:     "rename the current task title",
		shortcut: 'r',
		handler: func() func() {
			return func() {
				task := taskPanel.getFocusTask()
				currentTitle := task.Title

				menuPanel.Pages.SwitchToPage(RENAME_PAGE)
				menuPanel.renameInput.SetLabel("New task title: ").
					SetText(currentTitle).
					SetPlaceholder("Enter new task title").
					SetPlaceholderTextColor(tcell.ColorDarkSlateBlue).
					SetFieldTextColor(tcell.ColorBlack).
					SetFieldBackgroundColor(tcell.ColorLightBlue)

				menuPanel.renameInput.SetDoneFunc(func(key tcell.Key) {
					switch key {
					case tcell.KeyEnter:
						newTitle := menuPanel.renameInput.GetText()
						taskPanel.renameCurrentTask(task, newTitle)
						menuPanel.Pages.SwitchToPage(MENU_PAGE)

						removeTaskDetailPanel()
						closeMenuPanel()
						return
						// }
					case tcell.KeyEsc:
						menuPanel.Pages.SwitchToPage(MENU_PAGE)
					}
				})
			}
		},
	},
}

// Normal task panel menu items
var normalTaskMenus = []menuItem{
	{
		name:     "toggle the current task status",
		shortcut: 't',
		handler: func() func() {
			return func() {
				task := taskPanel.getFocusTask()
				taskPanel.toggleTaskStatus(task)

				removeTaskDetailPanel()
				closeMenuPanel()
			}
		},
	},
	{
		name:     "[red]delete[white] the current task",
		shortcut: 'd',
		handler: func() func() {
			return func() {
				taskPanel.deleteCurrentTask()

				removeTaskDetailPanel()
				closeMenuPanel()
			}
		},
	},
}

// Deleted task panel menu items
var deletedTaskMenus = []menuItem{
	{
		name:     "[green]restore[white] the current task",
		shortcut: 's',
		handler: func() func() {
			return func() {
				taskPanel.restoreCurrentTask()

				removeTaskDetailPanel()
				closeMenuPanel()
			}
		},
	},
}

// NewMenuPanel displays menu panel
// NOTE: current only works for Task
func NewMenuPanel() *MenuPanel {
	panel := MenuPanel{
		Pages:       tview.NewPages(),
		renameInput: tview.NewInputField(),
	}
	panel.SetBorder(true)

	return &panel
}

// Open Menu Panel and focus
func (p *MenuPanel) open() {
	focusTask := taskPanel.getFocusTask()
	p.loadMenus(*focusTask)

	app := global.App
	curItemIndex := taskPanel.list.GetCurrentItem()
	curItemText, _ := taskPanel.list.GetItemText(curItemIndex)
	if curItemText != TODO_HEADER && curItemText != COMPLETED_HEADER && curItemText != "" {
		layout.
			RemoveItem(statusPanel).
			AddItem(p, p.getSize(), 1, true).
			AddItem(statusPanel, 1, 1, false)
		app.SetFocus(p)
	}

}

func (p *MenuPanel) loadMenus(task model.Task) {
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

// Close Menu panel and focus Task panel
func closeMenuPanel() {
	app := global.App
	layout.RemoveItem(menuPanel)
	app.SetFocus(taskPanel)
}

// Returns the size of MenuPanel rows
func (p *MenuPanel) getSize() int {
	return p.choice.items.GetItemCount() + 4
}

// Add menu choices to Menu panel
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

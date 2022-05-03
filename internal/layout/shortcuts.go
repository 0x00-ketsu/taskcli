package layout

import (
	"reflect"

	"github.com/0x00-ketsu/taskcli/internal/global"
	"github.com/0x00-ketsu/taskcli/internal/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func setKeyboardShortcuts() *tview.Application {
	app := global.App
	return app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if ignoreKeyEvent(app) {
			return event
		}

		switch event.Rune() {
		case 'q':
			app.Stop()
			return nil
		case '/':
			app.SetFocus(searchPanel)
			return nil
		case '?':
			layout.Clear().AddItem(helpPanel, 0, 1, true)
			app.SetFocus(helpPanel)
			return nil
		}

		switch {
		case filterPanel.HasFocus():
			event = handleFilterPanelShortcuts(app, event)

		case searchPanel.HasFocus():
			event = handleSearchPanelShortcuts(app, event)

		case taskPanel.HasFocus():
			event = handleTaskPanelShortcuts(app, event)

		case taskDetailPanel.HasFocus():
			event = handleTaskDetailPanelShortcuts(app, event)

		case menuPanel.HasFocus():
			event = handleMenuPanelShortcuts(event)

		case helpPanel.HasFocus():
			event = handleHelpPanelShortcuts(app, event)
		}

		return event
	})
}

// Shortcuts for Filter panel
func handleFilterPanelShortcuts(app *tview.Application, event *tcell.EventKey) *tcell.EventKey {
	this := filterPanel

	switch event.Key() {
	case tcell.KeyRune:
		switch event.Rune() {
		case 'j':
			this.lineDown()
			return nil
		case 'k':
			this.lineUp()
			return nil
		case 'g':
			this.list.SetCurrentItem(0)
			return nil
		case 'G':
			this.list.SetCurrentItem(this.list.GetItemCount() - 1)
			return nil
		}
	}

	return event
}

// Shortcuts for Search panel
func handleSearchPanelShortcuts(app *tview.Application, event *tcell.EventKey) *tcell.EventKey {
	// this := searchPanel
	switch event.Key() {
	case tcell.KeyEsc:
		app.SetFocus(filterPanel)
		return nil
	}

	return event
}

// Shortcuts for Task panel
func handleTaskPanelShortcuts(app *tview.Application, event *tcell.EventKey) *tcell.EventKey {
	this := taskPanel

	switch event.Key() {
	case tcell.KeyEsc:
		app.SetFocus(filterPanel)
		return nil
	case tcell.KeyRune:
		switch event.Rune() {
		case 'n':
			app.SetFocus(this.newTask)
			return nil
		case 'j':
			this.lineDown()
			return nil
		case 'k':
			this.lineUp()
			return nil
		case 'g':
			this.list.SetCurrentItem(0)
			return nil
		case 'G':
			this.list.SetCurrentItem(this.list.GetItemCount() - 1)
			return nil
		case 'm':
			menuPanel.open()
			return nil
		}
	}

	return event
}

func handleTaskDetailPanelShortcuts(app *tview.Application, event *tcell.EventKey) *tcell.EventKey {
	this := taskDetailPanel

	switch event.Key() {
	case tcell.KeyEsc:
		removeTaskDetailPanel()
		app.SetFocus(taskPanel)
		return nil
	case tcell.KeyCtrlD:
		this.contentView.ScrollDown(1)
		return nil
	case tcell.KeyCtrlU:
		this.contentView.ScrollUp(1)
		return nil
	case tcell.KeyRune:
		switch event.Rune() {
		case 'h':
			this.contentView.CursorLeft()
			return nil
		case 'l':
			this.contentView.CursorRight()
			return nil
		case 'j':
			this.contentView.Cursor.DownN(1)
			return nil
		case 'k':
			this.contentView.Cursor.UpN(1)
			return nil
		case 'i':
			this.activateEditor()
			return nil
		case 'c':
			this.copyTaskContent()
			return nil
		case 'v':
			this.editInExternalEditor()
			return nil
		case 'r':
			this.header.showRename()
			return nil
		case ' ':
			this.toggleTaskStatus()
			return nil
		case 't':
			this.todaySelector()
			return nil
		case '+':
			this.nextDaySelector()
			return nil
		case '-':
			this.prevDaySelector()
			return nil
		}
	}

	return event
}

// Shortcuts for Menu panel
func handleMenuPanelShortcuts(event *tcell.EventKey) *tcell.EventKey {
	this := menuPanel
	totalItemCount := this.choice.items.GetItemCount()
	currentItemIndex := this.choice.items.GetCurrentItem()
	prevItemIndex := currentItemIndex - 1
	nextItemIndex := currentItemIndex + 1

	switch event.Key() {
	case tcell.KeyEsc:
		closeMenuPanel()
		return nil
	case tcell.KeyRune:
		switch event.Rune() {
		case 'j':
			// Move Down
			if nextItemIndex >= totalItemCount {
				this.choice.items.SetCurrentItem(0)
			} else {
				this.choice.items.SetCurrentItem(nextItemIndex)
			}
			return nil
		case 'k':
			// Move Up
			if prevItemIndex >= totalItemCount {
				this.choice.items.SetCurrentItem(0)
			} else {
				this.choice.items.SetCurrentItem(prevItemIndex)
			}
			return nil
		}
	}

	return event
}

// Shortcuts for Help panel
func handleHelpPanelShortcuts(app *tview.Application, event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q':
		layout.Clear().
			AddItem(main, 0, 1, true).
			AddItem(statusPanel, 1, 1, false)
		app.SetFocus(main)
		return nil
	default:
		layout.Clear().
			AddItem(main, 0, 1, true).
			AddItem(statusPanel, 1, 1, false)
		app.SetFocus(main)
		return event
	}

}

func ignoreKeyEvent(app *tview.Application) bool {
	textInputs := []string{"*tview.InputField", "*femto.View"}

	return utils.InArray(reflect.TypeOf(app.GetFocus()).String(), textInputs)
}

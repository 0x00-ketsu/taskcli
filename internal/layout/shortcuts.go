package layout

import (
	"reflect"

	"github.com/0x00-ketsu/taskcli/internal/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func setKeyboardShortcuts() *tview.Application {
	return app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if ignoreKeyEvent(app) {
			return event
		}

		switch event.Rune() {
		case 'q':
			db.Close()
			app.Stop()
			return nil
		case '/':
			app.SetFocus(searchView)
			return nil
		case '?':
			layout.Clear().AddItem(helpView, 0, 1, true)
			app.SetFocus(helpView)
			return nil
		}

		switch {
		case filterView.HasFocus():
			event = handleFilterViewShortcuts(event)

		case searchView.HasFocus():
			event = handleSearchViewShortcuts(app, event)

		case taskView.HasFocus():
			event = handleTaskViewShortcuts(app, event)

		case taskDetailView.HasFocus():
			event = handleTaskDetailViewShortcuts(app, event)

		case menuView.HasFocus():
			event = handleMenuViewShortcuts(event)

		case helpView.HasFocus():
			event = handleHelpViewShortcuts(app, event)
		}
		return event
	})
}

// Shortcuts for Filter view
func handleFilterViewShortcuts(event *tcell.EventKey) *tcell.EventKey {
	this := filterView
	if event.Key() == tcell.KeyRune {
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

// Shortcuts for Search view
func handleSearchViewShortcuts(app *tview.Application, event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyEsc {
		app.SetFocus(filterView)
		return nil
	}
	return event
}

// Shortcuts for Task view
func handleTaskViewShortcuts(app *tview.Application, event *tcell.EventKey) *tcell.EventKey {
	this := taskView
	switch event.Key() {
	case tcell.KeyEsc:
		app.SetFocus(filterView)
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
			menuView.open()
			return nil
		}
	}
	return event
}

func handleTaskDetailViewShortcuts(app *tview.Application, event *tcell.EventKey) *tcell.EventKey {
	this := taskDetailView
	switch event.Key() {
	case tcell.KeyEsc:
		removeTaskDetailView()
		app.SetFocus(taskView)
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

// Shortcuts for Menu view
func handleMenuViewShortcuts(event *tcell.EventKey) *tcell.EventKey {
	this := menuView
	totalItemCount := this.choice.items.GetItemCount()
	currentItemIndex := this.choice.items.GetCurrentItem()
	prevItemIndex := currentItemIndex - 1
	nextItemIndex := currentItemIndex + 1
	switch event.Key() {
	case tcell.KeyEsc:
		closeMenuView()
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

// Shortcuts for Help view
func handleHelpViewShortcuts(app *tview.Application, event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case 'q':
		layout.Clear().
			AddItem(main, 0, 1, true).
			AddItem(statusView, 1, 1, false)
		app.SetFocus(main)
		return nil
	default:
		layout.Clear().
			AddItem(main, 0, 1, true).
			AddItem(statusView, 1, 1, false)
		app.SetFocus(main)
		return nil
	}
}

func ignoreKeyEvent(app *tview.Application) bool {
	textInputs := []string{"*tview.InputField", "*femto.View"}
	return utils.InArray(reflect.TypeOf(app.GetFocus()).String(), textInputs)
}

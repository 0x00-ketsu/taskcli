package layout

import (
	"time"

	"github.com/0x00-ketsu/taskcli/internal/global"
	"github.com/0x00-ketsu/taskcli/internal/utils"
	"github.com/rivo/tview"
)

// Declare all panels
var (
	layout, main *tview.Flex

	todayPanel      *TodayPanel
	filterPanel     *FilterPanel
	taskPanel       *TaskPanel
	taskDetailPanel *TaskDetailPanel
	menuPanel       *MenuPanel
	statusPanel     *StatusPanel
	helpPanel       *HelpPanel
)

// Declare package global variables
var (
	today    = utils.ToDate(time.Now())
	tomorrow = today.AddDate(0, 0, 1)
)

func Load() *tview.Flex {
	// GUI panels
	todayPanel = NewTodayPanel()
	filterPanel = NewFilterPanel()
	taskPanel = NewTaskPanel()
	taskDetailPanel = NewTaskDetailPanel()
	menuPanel = NewMenuPanel()
	statusPanel = NewStatusPanel(global.App)
	helpPanel = NewHelpPanel()

	// GUI main
	main = tview.NewFlex()
	main.AddItem(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(todayPanel, 3, 1, false).
			AddItem(filterPanel, 0, 1, true),
		35, 1, true)
	main.AddItem(taskPanel, 0, 2, false)

	// GUI root layout
	layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(main, 0, 1, true).
		AddItem(statusPanel, 1, 1, false)

	// Bind keyboard shortcuts
	setKeyboardShortcuts()

	return layout
}

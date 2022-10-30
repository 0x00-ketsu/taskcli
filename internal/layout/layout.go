package layout

import (
	"time"

	"github.com/0x00-ketsu/taskcli/internal/cmd/flags"
	"github.com/0x00-ketsu/taskcli/internal/database"
	"github.com/0x00-ketsu/taskcli/internal/repository"
	"github.com/0x00-ketsu/taskcli/internal/repository/bolt"
	"github.com/0x00-ketsu/taskcli/internal/utils"
	"github.com/asdine/storm/v3"
	"github.com/rivo/tview"
)

// Declare all views
var (
	db *storm.DB
	taskRepo repository.Task

	app          *tview.Application
	layout, main *tview.Flex

	todayView      *TodayView
	filterView     *FilterView
	searchView     *SearchView
	taskView       *TaskView
	taskDetailView *TaskDetailView
	menuView       *MenuView
	statusView     *StatusView
	helpView       *HelpView
)

// Declare package global variables
var (
	today    = utils.ToDate(time.Now())
	tomorrow = today.AddDate(0, 0, 1)
)

func Load(application *tview.Application) *tview.Flex {
	// Connect DB
	db = database.Connect(flags.Storage)
	taskRepo = bolt.NewTask(db)

	app = application

	// GUI views
	todayView = NewTodayView()
	filterView = NewFilterView()
	searchView = NewSearchView()
	taskView = NewTaskView()
	taskDetailView = NewTaskDetailView()
	menuView = NewMenuView()
	statusView = NewStatusView()
	helpView = NewHelpView()

	// GUI main
	main = tview.NewFlex()
	main.AddItem(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(todayView, 3, 1, false).
			AddItem(filterView, 0, 1, true).
			AddItem(searchView, 0, 1, false),
		35, 1, true)
	main.AddItem(taskView, 0, 2, false)

	// GUI root layout
	layout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(main, 0, 1, true).
		AddItem(statusView, 1, 1, false)

	// Bind keyboard shortcuts
	setKeyboardShortcuts()

	return layout
}

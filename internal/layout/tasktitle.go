package layout

import (
	"fmt"

	"github.com/0x00-ketsu/taskcli/internal/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// TaskDetailHeader displays Task title and relevant action in TaskDetail view
type TaskDetailHeader struct {
	*tview.Flex
	pages       *tview.Pages
	renameInput *tview.InputField
	task        *model.Task
	taskName    *tview.TextView
}

func NewTaskDetailHeader() *TaskDetailHeader {
	header := TaskDetailHeader{
		Flex:        tview.NewFlex().SetDirection(tview.FlexRow),
		pages:       tview.NewPages(),
		renameInput: makeLightTextInput("Task title"),
		taskName:    tview.NewTextView().SetDynamicColors(true),
	}
	header.pages.AddPage("title", header.taskName, true, true)
	header.pages.AddPage("rename", header.renameInput, true, false)
	header.bindRenameEvent()
	renameView := tview.NewTextView().SetTextColor(tcell.ColorDimGray).SetText("r = Rename task title")
	tips := tview.NewFlex().AddItem(renameView, 0, 1, false)
	header.
		AddItem(header.pages, 1, 1, true).
		AddItem(blankCell, 1, 1, false).
		AddItem(tips, 1, 1, false).
		AddItem(makeHorizontalLine(tcell.RuneS3, tcell.ColorGray), 1, 1, false)
	return &header
}

func (p *TaskDetailHeader) bindRenameEvent() *tview.InputField {
	return p.renameInput.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			name := p.renameInput.GetText()
			if !validateTaskName(name) {
				return
			}
			taskView.renameCurrentTask(p.task, name)
			p.setTitle(p.task)
			p.pages.SwitchToPage("title")
			taskView.reloadTasks()
		case tcell.KeyEsc:
			p.pages.SwitchToPage("title")
		}
	})
}

func (p *TaskDetailHeader) setTitle(task *model.Task) {
	p.task = task
	p.taskName.SetText(fmt.Sprintf("[%s::b]# %s", getTaskTitleColor(*task), task.Title))
}

// Activate edit option of task title
func (p *TaskDetailHeader) showRename() {
	p.pages.SwitchToPage("rename")
	p.renameInput.SetText(p.task.Title)

	app.SetFocus(p)
}

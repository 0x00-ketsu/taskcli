package layout

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/0x00-ketsu/taskcli/internal/cmd/flags"
	"github.com/0x00-ketsu/taskcli/internal/model"
	"github.com/0x00-ketsu/taskcli/internal/utils"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/pgavlin/femto"
	"github.com/pgavlin/femto/runtime"
	"github.com/rivo/tview"
)

var (
	dateLayout      = "2006-01-02"
	dateHumanLayout = "02 Jan, Monday"
)

type TaskDetailView struct {
	*tview.Flex
	header      *TaskDetailHeader
	contentHint *tview.TextView
	statusHint  *tview.TextView

	task               *model.Task
	contentView        *femto.View
	colorScheme        femto.Colorscheme
	taskDueDate        *tview.InputField
	taskDueDateDisplay *tview.TextView
	taskToggleStatus   *tview.Button
}

func NewTaskDetailView() *TaskDetailView {
	view := TaskDetailView{
		Flex:               tview.NewFlex().SetDirection(tview.FlexRow),
		header:             NewTaskDetailHeader(),
		contentHint:        tview.NewTextView(),
		statusHint:         tview.NewTextView(),
		taskDueDateDisplay: tview.NewTextView().SetDynamicColors(true),
		taskToggleStatus:   makeButton("Uncompleted", nil).SetLabelColor(tcell.ColorGray),
	}

	view.loadEditor()

	view.taskToggleStatus.SetSelectedFunc(func() {
		view.toggleTaskStatus()
	})
	view.contentHint = tview.NewTextView().
		SetText(" i = insert, c = copy, h/j/k/l = move cursor, v = external editor").
		SetTextColor(tcell.ColorDimGray)
	view.statusHint = tview.NewTextView().
		SetTextColor(tcell.ColorDimGray).
		SetText(" <space> to toggle task status")

	editorLabel := tview.NewFlex().
		AddItem(tview.NewTextView().SetText("[lime::b]Task Content").SetDynamicColors(true), 0, 1, false)
	editorHelp := tview.NewFlex().
		AddItem(view.contentHint, 0, 1, false).
		AddItem(tview.NewTextView().SetTextAlign(tview.AlignRight).
			SetText(fmt.Sprintf("syntax:markdown (%v)", "monokai")).
			SetTextColor(tcell.ColorDimGray), 0, 1, false)

	view.
		AddItem(view.header, 3, 1, true).
		AddItem(blankCell, 1, 1, false).
		AddItem(view.makeDateRow(), 1, 1, true).
		AddItem(blankCell, 1, 1, false).
		AddItem(editorLabel, 1, 1, false).
		AddItem(view.contentView, 0, 10, false).
		AddItem(editorHelp, 1, 1, false).
		AddItem(blankCell, 0, 1, false).
		AddItem(view.statusHint, 1, 1, false).
		AddItem(view.taskToggleStatus, 1, 1, false)

	view.SetBorder(true).SetTitle(" Task Detail ")

	return &view
}

// Loads and shows task detail
func (p *TaskDetailView) setTask(task *model.Task) {
	p.task = task

	p.header.setTitle(task)
	p.contentView.Buf = makeBufferFromString(p.task.Content)
	p.contentView.SetColorscheme(p.colorScheme)
	p.contentView.Start()
	p.setTaskDate(p.task.DueDate, false)
	p.updateTaskToggleDisplay()
	p.deactivateEditor()
}

func (p *TaskDetailView) activateEditor() {
	p.contentView.Readonly = false
	p.contentView.SetBorderColor(tcell.ColorDarkOrange)
	p.contentHint.SetText(" Esc to save changes")

	app.SetFocus(p.contentView)
}

func (p *TaskDetailView) deactivateEditor() {
	p.contentView.Readonly = true
	p.contentView.SetBorderColor(tcell.ColorLightSlateGray)
	p.contentHint.SetText(" i = insert, c = copy, h/j/k/l = move cursor, v = external editor")

	app.SetFocus(p)
}

func (p *TaskDetailView) editInExternalEditor() {
	tmpFileName, err := utils.WriteToTempFile(p.task.Content, "taskcli_note_*.md")
	if err != nil {
		statusView.showForSeconds("[red]Failed to create tmp file. Try in-app editing by pressing i", 5)
		return
	}

	var messageToShow, updatedContent string
	app.Suspend(func() {
		editor := flags.Editor
		cmd := exec.Command(editor, tmpFileName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			messageToShow = "[red]Failed to save content. Try in-app editing by pressing i"
			return
		}

		if content, readErr := ioutil.ReadFile(tmpFileName); readErr == nil {
			updatedContent = string(content)
		} else {
			messageToShow = "[red]Failed to load external editing. Try in-app editing by pressing i"
		}
	})

	if messageToShow != "" {
		statusView.showForSeconds(messageToShow, 10)
	}

	if updatedContent != "" {
		p.updateTaskContent(updatedContent)
		p.setTask(p.task)
	}

	app.EnableMouse(true)

	_ = os.Remove(tmpFileName)
}

// Load detail edit view
func (p *TaskDetailView) loadEditor() {
	p.contentView = femto.NewView(makeBufferFromString(""))
	p.contentView.SetRuntimeFiles(runtime.Files)

	if monokai := runtime.Files.FindFile(femto.RTColorscheme, "monokai"); monokai != nil {
		if data, err := monokai.Data(); err == nil {
			p.colorScheme = femto.ParseColorscheme(string(data))
		}
	}

	p.contentView.SetColorscheme(p.colorScheme)
	p.contentView.SetBorder(true)
	p.contentView.SetBorderColor(tcell.ColorLightSlateGray)

	p.contentView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			p.updateTaskContent(p.contentView.Buf.String())
			p.deactivateEditor()
			return nil
		}
		return event
	})
}

func (p *TaskDetailView) updateTaskContent(content string) {
	p.task.Content = content
	if err := taskRepo.UpdateField(p.task, "Content", content); err != nil {
		msg := fmt.Sprintf("[red]Save content failed, error: %v", err.Error())
		statusView.showForSeconds(msg, 5)
	} else {
		statusView.showForSeconds("[green]Save content successful", 5)
	}
}

func (p *TaskDetailView) makeDateRow() *tview.Flex {
	p.taskDueDate = makeLightTextInput("YYYY-mm-dd").
		SetLabel("Set: ").
		SetLabelColor(tcell.ColorWhite).
		SetFieldWidth(14).
		SetDoneFunc(func(key tcell.Key) {
			switch key {
			case tcell.KeyEnter:
				if date, err := utils.ParseStrToDate(p.taskDueDate.GetText(), dateLayout); err != nil {
					p.taskDueDate.SetBorderColor(tcell.ColorRed)
					statusView.showForSeconds("[red]Input new due date is invalid", 5)
				} else {
					if date.Before(today) {
						p.taskDueDate.SetBorderColor(tcell.ColorRed)
						statusView.showForSeconds("[red]Input new due date should greater than or equal to today", 5)
					} else {
						p.taskDueDate.SetBorderColor(tcell.ColorDefault)
						p.setTaskDate(date, true)
					}
				}
			case tcell.KeyEsc:
				p.setTaskDate(p.task.DueDate, false)
				p.taskDueDate.SetBorderColor(tcell.ColorDefault)
			}
			app.SetFocus(p)
		})

	return tview.NewFlex().
		AddItem(p.taskDueDateDisplay, 0, 2, true).
		AddItem(p.taskDueDate, 16, 0, true).
		AddItem(blankCell, 1, 0, false).
		AddItem(makeButton("[::u]t[::-]oday", p.todaySelector), 7, 1, false).
		AddItem(blankCell, 1, 0, false).
		AddItem(makeButton("[::u]+[::-]1", p.nextDaySelector), 4, 1, false).
		AddItem(blankCell, 1, 0, false).
		AddItem(makeButton("[::u]-[::-]1", p.prevDaySelector), 4, 1, false)
}

// Update task date if `update` is true
func (p *TaskDetailView) setTaskDate(date time.Time, update bool) {
	if update {
		if err := taskRepo.UpdateField(p.task, "DueDate", date); err != nil {
			msg := fmt.Sprintf("[red]Update task due date failed, error: %v", err.Error())
			statusView.showForSeconds(msg, 5)
			return
		}
	}

	color := "whilte"
	humanDate := date.Format(dateHumanLayout)
	if date.Before(today) {
		color = "yellow"
	}
	p.taskDueDate.SetText(date.Format(dateLayout))
	p.taskDueDateDisplay.SetText(fmt.Sprintf("Due: [%s]%s", color, humanDate))
}

// If task is deleted, hide tips: <space> to toggle task status
func (p *TaskDetailView) updateTaskToggleDisplay() {
	if p.task.IsDeleted {
		p.RemoveItem(p.statusHint)
		p.taskToggleStatus.SetLabel("Deleted").SetBackgroundColor(tcell.ColorRed)
	} else {
		p.RemoveItem(p.statusHint)
		p.RemoveItem(p.taskToggleStatus)

		p.AddItem(p.statusHint, 1, 1, false).
			AddItem(p.taskToggleStatus, 1, 1, false)

		if p.task.IsCompleted {
			p.taskToggleStatus.SetLabel("Resume").SetBackgroundColor(tcell.ColorYellow)
		} else {
			p.taskToggleStatus.SetLabel("Completed").SetBackgroundColor(tcell.ColorDarkGreen)
		}
	}
}

func (p *TaskDetailView) toggleTaskStatus() {
	if !p.task.IsDeleted {
		taskView.toggleTaskStatus(p.task)
		p.updateTaskToggleDisplay()
	}
}

func (p *TaskDetailView) todaySelector() {
	p.setTaskDate(p.task.DueDate, true)
}

func (p *TaskDetailView) nextDaySelector() {
	if date, err := utils.ParseStrToDate(p.taskDueDate.GetText(), dateLayout); err == nil {
		p.setTaskDate(date.AddDate(0, 0, 1), true)
	}
}

func (p *TaskDetailView) prevDaySelector() {
	if date, err := utils.ParseStrToDate(p.taskDueDate.GetText(), dateLayout); err == nil {
		p.setTaskDate(date.AddDate(0, 0, -1), true)
	}
}

// Copy task detail to clipboard
func (p *TaskDetailView) copyTaskContent() {
	var content bytes.Buffer
	content.WriteString(p.task.Content)
	_ = clipboard.WriteAll(content.String())

	app.SetFocus(p)
	statusView.showForSeconds("[green]Task content copyed. Try Pasting anywhere", 5)
}

func makeBufferFromString(content string) *femto.Buffer {
	buff := femto.NewBufferFromString(content, "")
	buff.Settings["filetype"] = "markdown"
	buff.Settings["keepautoindent"] = true
	buff.Settings["statusline"] = false
	buff.Settings["softwrap"] = true
	buff.Settings["scrollbar"] = true

	return buff
}

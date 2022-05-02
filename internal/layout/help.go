package layout

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HelpPanel struct {
	*tview.Flex
}

func NewHelpPanel() *HelpPanel {
	panel := HelpPanel{
		Flex: tview.NewFlex().SetDirection(tview.FlexRow),
	}

	tips := tview.NewTextView().SetText("\tPress any key to return.")
	panel.
		AddItem(drawBanner(), 9, 1, false).
		AddItem(drawGlobal(), 4, 1, false).
		AddItem(drawFilter(), 6, 1, false).
		AddItem(drawTask(), 7, 1, false).
		AddItem(drawTaskDetail(), 15, 1, false).
		AddItem(tips, 1, 1, false).
		AddItem(tview.NewTextView(), 0, 1, false)

	panel.SetBackgroundColor(tcell.ColorBlack)
	return &panel
}

// drawBanner ...
func drawBanner() *tview.Flex {
	text :=
		`
    ____________________
    < Welcome to Taskcli >
    --------------------
            \   ^__^
             \  (oo)\_______
                (__)\       )\/\
                    ||----w |
                    ||     ||
    
 
	`
	banner := tview.NewTextView().SetText(text)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(banner, 0, 9, false)
}

func drawGlobal() *tview.Flex {
	global := tview.NewTextView().SetText("\tGlobal")
	globalDesc := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\tq: quit"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t Esc: step back"), 0, 1, false)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(global, 1, 1, false).
		AddItem(globalDesc, 3, 1, false)
}

func drawFilter() *tview.Flex {
	filter := tview.NewTextView().SetText("\tFilter Panel")
	filterDesc := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t  j/k: down/up"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\tg: go to first item"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\tG: go to last item"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\tEnter: Activate task"), 0, 1, false)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(filter, 1, 1, false).
		AddItem(filterDesc, 5, 1, false)
}

func drawTask() *tview.Flex {
	task := tview.NewTextView().SetText("\tTask Panel")
	taskDesc := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\tn: create a new task"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t  j/k: down/up"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\tg: go to first item"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\tG: go to last item"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\tm: show menu"), 0, 1, false)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(task, 1, 1, false).
		AddItem(taskDesc, 6, 1, false)
}

func drawTaskDetail() *tview.Flex {
	detail := tview.NewTextView().SetText("\tTask Detail Panel")
	detailDesc := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\tr: rename task title"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\tt: set to today"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\t+: set to next day"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\t-: set to previous day"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\tspace: toggle task status"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText(""), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\ti: edit task content view"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\tc: copy task content"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t  j/k: move down/up"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t  h/l: move left/right"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("   Ctrl-d: scroll down"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("   Ctrl-u: scroll up"), 0, 1, false).
		AddItem(tview.NewTextView().SetTextColor(tcell.ColorLimeGreen).SetText("\t\tv: edit task content with external editor"), 0, 1, false)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(detail, 1, 1, false).
		AddItem(detailDesc, 14, 1, false)
}

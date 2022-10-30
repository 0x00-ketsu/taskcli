package layout

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TodayView struct {
	*tview.Flex
	hint *tview.TextView
}

func NewTodayView() *TodayView {
	view := TodayView{
		Flex: tview.NewFlex(),
		hint: tview.NewTextView().SetTextColor(tcell.ColorGreen),
	}

	hint := renderHint()
	view.hint.SetText(hint)
	view.AddItem(view.hint, 0, 1, false)
	view.SetBorder(true).SetTitle("Today").SetTitleAlign(tview.AlignLeft)

	return &view
}

func (p *TodayView) updateTodoCount() {
	hint := renderHint()
	p.hint.SetText(hint)
}

func renderHint() string {
	var hint string
	todayTodoCount := taskRepo.GetTodayTodoCount()
	if todayTodoCount == 0 {
		hint = " have a nice day"
	} else {
		hint = fmt.Sprintf(" todo items: %d", todayTodoCount)
	}

	return hint
}

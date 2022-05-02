package layout

import (
	"fmt"

	"github.com/0x00-ketsu/taskcli/internal/global"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TodayPanel struct {
	*tview.Flex
	hint *tview.TextView
}

func NewTodayPanel() *TodayPanel {
	panel := TodayPanel{
		Flex: tview.NewFlex(),
		hint: tview.NewTextView().SetTextColor(tcell.ColorGreen),
	}

	hint := renderHint()
	panel.hint.SetText(hint)
	panel.AddItem(panel.hint, 0, 1, false)
	panel.SetBorder(true).SetTitle("Today").SetTitleAlign(tview.AlignLeft)

	return &panel
}

func (p *TodayPanel) updateTodoCount() {
	hint := renderHint()
	p.hint.SetText(hint)
}

func renderHint() string {
	var hint string
	todayTodoCount := global.TaskRepo.GetTodayTodoCount()
	if todayTodoCount == 0 {
		hint = " have a nice day"
	} else {
		hint = fmt.Sprintf(" todo items: %d", todayTodoCount)
	}

	return hint
}

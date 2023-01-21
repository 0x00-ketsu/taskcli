package layout

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var SEP = "[::d]" + strings.Repeat(string(tcell.RuneS3), 35)

// Filter panel
type FilterView struct {
	*tview.Flex
	list *tview.List
}

func NewFilterView() *FilterView {
	view := FilterView{
		Flex: tview.NewFlex(),
		list: tview.NewList().ShowSecondaryText(false),
	}

	view.AddItem(view.list, 0, 1, true)
	view.SetBorder(true).
		SetTitle(" Filter ")

	view.list.AddItem(" 📅 Today", "", 0, func() { taskView.loadFilterTasks("today") }).
		AddItem(" 📅 Tomorrow", "", 0, func() { taskView.loadFilterTasks("tomorrow") }).
		AddItem(" 📅 Last 7 days", "", 0, func() { taskView.loadFilterTasks("last 7 days") }).
		AddItem(SEP, "", 0, nil).
		AddItem(" ✅ [green]Completed[white]", "", 0, func() { taskView.loadFilterTasks("completed") }).
		AddItem(" 💁 [yellow]Expired[white]", "", 0, func() { taskView.loadFilterTasks("expired") }).
		AddItem(" 🚮 [red]Trash[white]", "", 0, func() { taskView.loadFilterTasks("trash") })

	return &view
}

// Move to next item
// Skip sepreation line
func (p *FilterView) lineDown() {
	curItemIndex := p.list.GetCurrentItem()
	itemCount := p.list.GetItemCount()

	if curItemIndex >= 0 && curItemIndex < itemCount-1 {
		nextItemIndex := curItemIndex + 1
		nextItemText, _ := p.list.GetItemText(nextItemIndex)
		if nextItemText == SEP || nextItemText == "" {
			nextItemIndex++
		}
		p.list.SetCurrentItem(nextItemIndex)
	}
}

// Move to previous item
// Skip sepreation line
func (p *FilterView) lineUp() {
	curItemIndex := p.list.GetCurrentItem()
	itemCount := p.list.GetItemCount()

	if curItemIndex < itemCount && curItemIndex > 0 {
		prevItemIndex := curItemIndex - 1
		prevItemText, _ := p.list.GetItemText(prevItemIndex)
		if prevItemText == SEP || prevItemText == "" {
			prevItemIndex--
		}
		p.list.SetCurrentItem(prevItemIndex)
	}
}

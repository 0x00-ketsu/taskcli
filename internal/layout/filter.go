package layout

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var SEP = "[::d]" + strings.Repeat(string(tcell.RuneS3), 35)

// Filter panel
type FilterPanel struct {
	*tview.Flex
	list *tview.List
}

func NewFilterPanel() *FilterPanel {
	panel := FilterPanel{
		Flex: tview.NewFlex(),
		list: tview.NewList().ShowSecondaryText(false),
	}

	panel.AddItem(panel.list, 0, 1, true)
	panel.SetBorder(true).
		SetTitle(" Filter ")

	panel.list.AddItem(" 📅 Today", "", 0, func() { taskPanel.loadFilterTasks("today") }).
		AddItem(" 📅 Tomorrow", "", 0, func() { taskPanel.loadFilterTasks("tomorrow") }).
		AddItem(" 📅 Last 7 days", "", 0, func() { taskPanel.loadFilterTasks("last 7 days") }).
		AddItem(SEP, "", 0, nil).
		AddItem(" ✅ [green]Completed[white]", "", 0, func() { taskPanel.loadFilterTasks("completed") }).
		AddItem(" 💁 [yellow]Expired[white]", "", 0, func() { taskPanel.loadFilterTasks("expired") }).
		AddItem(" 🚮 [red]Trash[white]", "", 0, func() { taskPanel.loadFilterTasks("trash") })

	return &panel
}

// Move to next item
// Skip sepreation line
func (p *FilterPanel) lineDown() {
	curItemIndex := p.list.GetCurrentItem()
	itemCount := p.list.GetItemCount()

	if curItemIndex >= 0 && curItemIndex < itemCount-1 {
		nextItemIndex := curItemIndex + 1
		nextItemText, _ := p.list.GetItemText(nextItemIndex)
		if nextItemText == SEP || nextItemText == "" {
			nextItemIndex += 1
		}
		p.list.SetCurrentItem(nextItemIndex)
	}
}

// Move to previous item
// Skip sepreation line
func (p *FilterPanel) lineUp() {
	curItemIndex := p.list.GetCurrentItem()
	itemCount := p.list.GetItemCount()

	if curItemIndex < itemCount && curItemIndex > 0 {
		prevItemIndex := curItemIndex - 1
		prevItemText, _ := p.list.GetItemText(prevItemIndex)
		if prevItemText == SEP || prevItemText == "" {
			prevItemIndex -= 1
		}
		p.list.SetCurrentItem(prevItemIndex)
	}
}

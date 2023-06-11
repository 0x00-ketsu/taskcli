package layout

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StatusView struct {
	*tview.Flex
	container *tview.Application
	hint      *tview.TextView
}

func NewStatusView() *StatusView {
	view := &StatusView{
		Flex:      tview.NewFlex(),
		container: app,
		hint:      tview.NewTextView().SetDynamicColors(true),
	}
	view.setDefaultHint()
	view.AddItem(view.hint, 0, 1, false)
	return view
}

func (p *StatusView) setDefaultHint() {
	p.hint.
		SetTextColor(tcell.ColorBlue).
		SetTextAlign(tview.AlignLeft).
		SetText("q: quit, ?: help, /: search")
}

func (p *StatusView) restore() {
	p.container.QueueUpdateDraw(func() {
		p.setDefaultHint()
	})
}

// Used to skip queued restore of Status view
// in case of new showForSeconds within waiting period
var restorInQ = 0

// Show message in Status view for seconds
func (p *StatusView) showForSeconds(message string, timeout int) {
	if p.container == nil {
		return
	}
	p.hint.SetText(message)
	restorInQ++

	go func() {
		time.Sleep(time.Second * time.Duration(timeout))
		// Apply restore only if this is the last pending restore
		if restorInQ == 1 {
			p.restore()
		}
		restorInQ--
	}()
}

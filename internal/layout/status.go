package layout

import (
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type StatusPanel struct {
	*tview.Flex
	container *tview.Application
	hint      *tview.TextView
}

func NewStatusPanel(app *tview.Application) *StatusPanel {
	panel := &StatusPanel{
		Flex:      tview.NewFlex(),
		container: app,
		hint:      tview.NewTextView().SetDynamicColors(true),
	}

	panel.setDefaultHint()
	panel.AddItem(panel.hint, 0, 1, false)

	return panel
}

func (p *StatusPanel) setDefaultHint() {
	p.hint.
		SetTextColor(tcell.ColorBlue).
		SetTextAlign(tview.AlignLeft).
		SetText("q: quit, ?: help")
}

func (p *StatusPanel) restore() {
	p.container.QueueUpdateDraw(func() {
		p.setDefaultHint()
	})
}

// Used to skip queued restore of Status panel
// in case of new showForSeconds within waiting period
var restorInQ = 0

// Show message in Status panel for seconds
func (p *StatusPanel) showForSeconds(message string, timeout int) {
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

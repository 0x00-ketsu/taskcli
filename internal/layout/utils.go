package layout

import (
	"github.com/0x00-ketsu/taskcli/internal/model"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var blankCell = tview.NewTextView()

func removeTaskDetailView() {
	main.RemoveItem(taskDetailView)
}

func getTaskTitleColor(task model.Task) string {
	var color string

	switch {
	case task.IsCompleted:
		color = "gray"
	case task.IsDeleted:
		color = "red"
	case task.DueDate.Before(today):
		color = "yellow"
	default:
		color = "white"
	}

	return color
}

func makeLightTextInput(placeholder string) *tview.InputField {
	return tview.NewInputField().
		SetPlaceholder(placeholder).
		SetPlaceholderTextColor(tcell.ColorDarkSlateBlue).
		SetFieldTextColor(tcell.ColorBlack).
		SetFieldBackgroundColor(tcell.ColorLightBlue)
}

func makeHorizontalLine(lineChar rune, color tcell.Color) *tview.TextView {
	hr := tview.NewTextView()
	hr.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		style := tcell.StyleDefault.Foreground(color).Background(tcell.ColorBlack)
		centerY := y + height/2
		for cx := x; cx < x+width; cx++ {
			screen.SetContent(cx, centerY, lineChar, nil, style)
		}

		return x + 1, centerY + 1, width - 2, height - (centerY + 1 - y)
	})

	return hr
}

func makeButton(label string, handler func()) *tview.Button {
	btn := tview.NewButton(label).SetSelectedFunc(handler).
		SetLabelColor(tcell.ColorWhite)

	btn.SetBackgroundColor(tcell.ColorCornflowerBlue)

	return btn
}

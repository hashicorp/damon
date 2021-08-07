package primitives

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Modal struct {
	primitive *tview.Modal
	container *tview.Flex
}

func NewModal(title string, buttons []string, c tcell.Color) *Modal {
	m := tview.NewModal()
	m.SetTitle(title)
	m.SetTitleAlign(tview.AlignCenter)
	m.SetBackgroundColor(c)
	m.SetTextColor(tcell.ColorBlack)
	m.AddButtons(buttons)

	f := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(m, 10, 1, true).
			AddItem(nil, 0, 1, false), 80, 1, false).
		AddItem(nil, 0, 1, false)

	return &Modal{
		primitive: m,
		container: f,
	}
}

func (m *Modal) SetDoneFunc(handler func(buttonIndex int, buttonLabel string)) {
	m.primitive.SetDoneFunc(handler)
}

func (m *Modal) SetText(text string) {
	m.primitive.SetText(text)
}

func (m *Modal) SetFocus(index int) {
	m.primitive.SetFocus(index)
}

func (m *Modal) Container() tview.Primitive {
	return m.container
}

func (m *Modal) Primitive() tview.Primitive {
	return m.primitive
}

package primitives

import (
	"github.com/rivo/tview"
)

type TextView struct {
	primitive *tview.TextView
}

func NewTextView(align int) *TextView {
	t := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(align)

	return &TextView{primitive: t}
}

func (t *TextView) Primitive() tview.Primitive {
	return t.primitive
}

func (t *TextView) GetText() string {
	return t.primitive.GetText(true)
}

func (t *TextView) SetText(text string) {
	t.primitive.SetText(text)
}

func (t *TextView) Clear() {
	t.primitive.Clear()
}

func (t *TextView) ModifyPrimitive(f func(t *tview.TextView)) {
	f(t.primitive)
}

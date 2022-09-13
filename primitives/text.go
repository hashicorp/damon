package primitives

import (
	"github.com/rivo/tview"
)

type TextView struct {
	*tview.TextView
}

func NewTextView(align int) *TextView {
	t := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(align)

	return &TextView{TextView: t}
}

func (t *TextView) Primitive() tview.Primitive {
	return t.TextView
}

func (t *TextView) ModifyPrimitive(f func(t *tview.TextView)) {
	f(t.TextView)
}

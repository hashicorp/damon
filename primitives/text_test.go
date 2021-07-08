package primitives_test

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/primitives"
)

func TestTextView(t *testing.T) {
	r := require.New(t)

	tv := primitives.NewTextView(tview.AlignRight)
	p := tv.Primitive().(*tview.TextView)

	tv.SetText("test")
	r.Equal(tv.GetText(), "test")

	tv.Clear()
	r.Equal(tv.GetText(), "")

	tv.ModifyPrimitive(func(v *tview.TextView) {
		v.SetBackgroundColor(tcell.ColorWhite)
	})

	r.Equal(p.GetBackgroundColor(), tcell.ColorWhite)
}

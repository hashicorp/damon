// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

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
	r.Equal(tv.GetText(true), "test")

	tv.Clear()
	r.Equal(tv.GetText(true), "")

	tv.ModifyPrimitive(func(v *tview.TextView) {
		v.SetBackgroundColor(tcell.ColorWhite)
	})

	r.Equal(p.GetBackgroundColor(), tcell.ColorWhite)
}

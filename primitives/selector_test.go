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

func TestSelectionModal(t *testing.T) {
	r := require.New(t)

	m := primitives.NewSelectionModal()

	tb := m.Primitive().(*tview.Table)
	c := m.Container().(*tview.Flex)
	table := m.GetTable()

	table.RenderRow([]string{"item1", "item2"}, 0, tcell.ColorWhite)
	table.RenderRow([]string{"item3", "item4"}, 1, tcell.ColorWhite)

	r.NotNil(tb)
	r.NotNil(c)
	r.Equal(tb, table.Primitive())

	item1 := table.GetCellContent(0, 0)
	item2 := table.GetCellContent(0, 1)
	item3 := table.GetCellContent(1, 0)
	item4 := table.GetCellContent(1, 1)

	r.Equal(item1, "item1")
	r.Equal(item2, "item2")
	r.Equal(item3, "item3")
	r.Equal(item4, "item4")
}

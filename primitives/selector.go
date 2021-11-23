package primitives

import (
	"github.com/rivo/tview"
)

type SelectionModal struct {
	Table     *Table
	container *tview.Flex
}

func NewSelectionModal() *SelectionModal {
	t := NewTable()
	f := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(t.primitive, 10, 1, false).
			AddItem(nil, 0, 1, false), 80, 1, false).
		AddItem(nil, 0, 1, false)

	return &SelectionModal{
		Table:     t,
		container: f,
	}
}

func (s *SelectionModal) Container() tview.Primitive {
	return s.container
}

func (s *SelectionModal) Primitive() tview.Primitive {
	return s.Table.primitive
}

func (s *SelectionModal) GetTable() *Table {
	return s.Table
}

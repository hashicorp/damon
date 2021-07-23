package component

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/models"
	primitive "github.com/hcjulz/damon/primitives"
	"github.com/hcjulz/damon/styles"
)

const TableTitleNamespaces = "Namespaces"

var (
	TableHeaderNamespaces = []string{
		LabelName,
		LabelDescription,
	}
)

type NamespaceTable struct {
	Table Table
	Props *NamespacesProps

	slot *tview.Flex
}

type NamespacesProps struct {
	HandleNoResources models.HandlerFunc
	Data              []*models.Namespace
}

func NewNamespaceTable() *NamespaceTable {
	t := primitive.NewTable()

	return &NamespaceTable{
		Table: t,
		Props: &NamespacesProps{},
	}
}

func (n *NamespaceTable) Bind(slot *tview.Flex) {
	n.slot = slot
}

func (n *NamespaceTable) Render() error {
	if n.slot == nil {
		return ErrComponentNotBound
	}

	if n.Props.HandleNoResources == nil {
		return ErrComponentPropsNotSet
	}

	n.reset()

	if len(n.Props.Data) == 0 {
		n.Props.HandleNoResources(
			"%sno namespaces available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)

		return nil
	}

	n.Table.SetTitle(TableTitleNamespaces)

	n.Table.RenderHeader(TableHeaderNamespaces)
	n.renderRows()

	n.slot.AddItem(n.Table.Primitive(), 0, 1, false)
	return nil
}

func (n *NamespaceTable) reset() {
	n.Table.Clear()
	n.slot.Clear()
}

func (n *NamespaceTable) renderRows() {
	for i, ns := range n.Props.Data {
		row := []string{
			ns.Name,
			ns.Description,
		}

		index := i + 1
		n.Table.RenderRow(row, index, tcell.ColorWhite)
	}
}

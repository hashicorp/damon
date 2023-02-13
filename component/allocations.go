// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package component

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/models"
	primitive "github.com/hcjulz/damon/primitives"
	"github.com/hcjulz/damon/styles"
)

const (
	TableTitleAllocations = "Allocations"
)

var (
	TableHeaderAllocations = []string{
		LabelID,
		LabelTaskGroup,
		LabelJobID,
		LabelType,
		LabelNamespace,
		LabelAddresses,
		LabelNodeID,
		LabelNodeName,
		LabelDesiredStatus,
	}
)

type SelectAllocationFunc func(allocID string)

type AllocationTable struct {
	Table Table
	Props *AllocationTableProps

	slot *tview.Flex
}

type AllocationTableProps struct {
	SelectAllocation  SelectAllocationFunc
	HandleNoResources models.HandlerFunc

	JobID string

	Data []*models.Alloc
}

func NewAllocationTable() *AllocationTable {
	t := primitive.NewTable()

	return &AllocationTable{
		Table: t,
		Props: &AllocationTableProps{},
	}
}

func (t *AllocationTable) Bind(slot *tview.Flex) {
	t.slot = slot
}

func (t *AllocationTable) Render() error {
	if t.Props.SelectAllocation == nil {
		return ErrComponentPropsNotSet
	}

	if t.Props.HandleNoResources == nil {
		return ErrComponentPropsNotSet
	}

	if t.slot == nil {
		return ErrComponentNotBound
	}

	t.reset()

	if len(t.Props.Data) == 0 {
		t.Props.HandleNoResources(
			"%sno allocations available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)

		return nil
	}

	t.Table.SetSelectedFunc(t.allocationSelected)

	t.Table.RenderHeader(TableHeaderAllocations)
	t.renderRows()

	t.Table.SetTitle(fmt.Sprintf("%s (Job: %s)", TableTitleAllocations, t.Props.JobID))
	t.slot.AddItem(t.Table.Primitive(), 0, 1, false)

	return nil
}

func (t *AllocationTable) reset() {
	t.slot.Clear()
	t.Table.Clear()
}

func (t *AllocationTable) renderRows() {
	for i, a := range t.Props.Data {
		hostAddr := fmt.Sprintf("%v", a.HostAddresses)
		row := []string{
			a.ID,
			a.TaskGroup,
			a.JobID,
			a.JobType,
			a.Namespace,
			hostAddr,
			a.NodeID,
			a.NodeName,
			a.DesiredStatus,
		}

		index := i + 1

		c := t.getCellColor(a.DesiredStatus)
		t.Table.RenderRow(row, index, c)
	}
}

func (t *AllocationTable) getCellColor(status string) tcell.Color {
	c := tcell.ColorWhite

	switch status {
	case models.DesiredStatusStop:
		c = tcell.ColorDarkGray
	}

	return c
}

func (t *AllocationTable) allocationSelected(row, column int) {
	allocID := t.Table.GetCellContent(row, 0)
	t.Props.SelectAllocation(allocID)
}

func (t *AllocationTable) GetIDForSelection() string {
	row, _ := t.Table.GetSelection()
	return t.Table.GetCellContent(row, 0)
}

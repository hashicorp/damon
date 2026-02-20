// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component_test

import (
	"testing"

	"github.com/rivo/tview"
	"github.com/stretchr/testify/require"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/component/componentfakes"
	"github.com/hcjulz/damon/primitives"
)

func TestSelectorModal_Happy(t *testing.T) {
	r := require.New(t)

	fakeSelectorModal := &componentfakes.FakeSelector{}
	table := primitives.NewTable()

	pages := tview.NewPages()

	m := component.NewSelectorModal()
	m.Modal = fakeSelectorModal
	m.Bind(pages)
	m.Props.Items = []string{
		"task-1",
		"task-2",
	}

	fakeSelectorModal.GetTableReturns(table)
	err := m.Render()
	r.NoError(err)

	r.Equal(fakeSelectorModal.GetTableCallCount(), 2)

	item1 := table.GetCellContent(0, 0)
	item2 := table.GetCellContent(1, 0)

	r.Equal(item1, "task-1")
	r.Equal(item2, "task-2")

	r.Equal(pages.GetPageCount(), 1)
}

func TestSelectorModal_Sad(t *testing.T) {
	t.Run("When the component is not bound", func(t *testing.T) {
		r := require.New(t)

		fakeSelectorModal := &componentfakes.FakeSelector{}
		table := primitives.NewTable()

		m := component.NewSelectorModal()
		m.Modal = fakeSelectorModal
		m.Props.Items = []string{
			"task-1",
			"task-2",
		}

		fakeSelectorModal.GetTableReturns(table)
		err := m.Render()
		r.Error(err)

		r.ErrorIs(err, component.ErrComponentNotBound)
	})

	t.Run("When component properites are not set", func(t *testing.T) {
		r := require.New(t)

		fakeSelectorModal := &componentfakes.FakeSelector{}
		table := primitives.NewTable()

		pages := tview.NewPages()

		m := component.NewSelectorModal()
		m.Modal = fakeSelectorModal

		m.Bind(pages)

		fakeSelectorModal.GetTableReturns(table)
		err := m.Render()
		r.Error(err)

		r.ErrorIs(err, component.ErrComponentPropsNotSet)
	})
}

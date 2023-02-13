// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"regexp"

	"github.com/rivo/tview"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/models"
)

func (v *View) Namespaces() {
	v.viewSwitch()

	v.Layout.Body.SetTitle(titleNamespaces)

	v.state.Elements.TableMain = v.components.NamespaceTable.Table.Primitive().(*tview.Table)
	v.components.Commands.Update(component.NoViewCommands)
	v.Layout.Container.SetInputCapture(v.InputNamespaces)

	update := func() {
		v.components.NamespaceTable.Props.Data = v.filterNamespaces(v.state.Namespaces)
		v.components.NamespaceTable.Render()
		v.Draw()
	}

	v.components.Search.Props.ChangedFunc = func(text string) {
		v.state.Filter.Namespaces = text
		update()
	}

	v.Watcher.SubscribeToNamespaces(update)

	update()

	v.components.Selections.Namespace.SetSelectedFunc(func(text string, index int) {
		v.state.SelectedNamespace = text
		v.Namespaces()
	})

	v.addToHistory(v.state.SelectedNamespace, models.TopicNamespace, v.Namespaces)
	v.Layout.Container.SetFocus(v.components.NamespaceTable.Table.Primitive())
}

func getNamespaceNameIndex(name string, ns []*models.Namespace) int {
	var index int
	for i, n := range ns {
		if n.Name == name {
			index = i
		}
	}

	return index
}

func (v *View) filterNamespaces(data []*models.Namespace) []*models.Namespace {
	filter := v.state.Filter.Namespaces
	if filter != "" {
		rx, _ := regexp.Compile(filter)
		result := []*models.Namespace{}
		for _, ns := range v.state.Namespaces {
			switch true {
			case rx.MatchString(ns.Name),
				rx.MatchString(ns.Description):
				result = append(result, ns)
			}
		}

		return result
	}

	return data
}

func (v *View) resetSearch() {
	if v.state.Toggle.Search {
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		v.Layout.Footer.RemoveItem(v.components.Search.InputField.Primitive())
		v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 0)
		v.state.Toggle.Search = false
	}
}

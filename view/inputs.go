// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"github.com/gdamore/tcell/v2"
)

func (v *View) InputJobs(event *tcell.EventKey) *tcell.EventKey {
	event = v.InputMainCommands(event)
	return v.inputJobs(event)
}

func (v *View) InputDeployments(event *tcell.EventKey) *tcell.EventKey {
	return v.InputMainCommands(event)
}

func (v *View) InputNamespaces(event *tcell.EventKey) *tcell.EventKey {
	return v.InputMainCommands(event)
}

func (v *View) InputTaskGroups(event *tcell.EventKey) *tcell.EventKey {
	return v.InputMainCommands(event)
}

func (v *View) InputAllocations(event *tcell.EventKey) *tcell.EventKey {
	v.InputMainCommands(event)
	return v.inputAllocs(event)
}

func (v *View) InputMainCommands(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	switch event.Key() {
	case tcell.KeyCtrlJ:
		v.Jobs()

	case tcell.KeyCtrlN:
		v.Namespaces()

	case tcell.KeyCtrlD:
		v.Deployments()

	case tcell.KeyCtrlO, tcell.KeyEsc:
		v.GoBack()

	case tcell.KeyCtrlP:
		if !v.Layout.Footer.HasFocus() {
			v.Layout.Container.SetFocus(v.components.LogSearch.InputField.Primitive())
			if !v.state.Toggle.JumpToJob {
				v.viewSwitch()
				v.JumpToJob()
				v.state.Toggle.JumpToJob = true
			} else {
				v.Layout.Container.SetFocus(v.components.JumpToJob.InputField.Primitive())
			}
		}
	case tcell.KeyRune:
		switch event.Rune() {

		case 's':
			if !v.Layout.Footer.HasFocus() {
				v.Layout.Container.SetFocus(v.state.Elements.DropDownNamespace)
			}
		}
	}

	return event
}

func (v *View) inputAllocs(event *tcell.EventKey) *tcell.EventKey {
	return event
}

func (v *View) InputLogs(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyEsc, tcell.KeyCtrlO, tcell.KeyEnter:
		if v.components.LogStream.TextView.Primitive().HasFocus() {
			v.GoBack()
			return nil
		}
	case tcell.KeyRune:
		switch event.Rune() {
		case '/':
			if !v.Layout.Footer.HasFocus() {
				if !v.state.Toggle.LogSearch {
					v.state.Toggle.LogSearch = true
					v.LogSearch()
					return nil
				} else {
					v.Layout.Container.SetFocus(v.components.LogSearch.InputField.Primitive())
				}

			}
		case 'h':
			if !v.Layout.Footer.HasFocus() {
				if !v.state.Toggle.LogHighlight {
					v.state.Toggle.LogHighlight = true
					v.LogHighlight()
					return nil
				} else {
					v.Layout.Container.SetFocus(v.components.LogHighlight.InputField.Primitive())
				}
			}
		case 's':
			if !v.Layout.Footer.HasFocus() {
				v.Watcher.Unsubscribe()
			}
		case 'r':
			if !v.Layout.Footer.HasFocus() {
				v.Watcher.ResumeLogs()

			}
		}
	}

	return event
}

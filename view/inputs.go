// Copyright IBM Corp. 2021, 2023
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
	event = v.InputMainCommands(event)
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

	case tcell.KeyLeft:
		// Same as Esc, but only when not editing a footer input field,
		// where the left arrow needs to keep moving the text cursor.
		if !v.Layout.Footer.HasFocus() {
			v.GoBack()
		}

	case tcell.KeyRight:
		// Same as pressing Enter on the currently selected row, but only
		// when not editing a footer input field, where the right arrow
		// needs to keep moving the text cursor.
		if !v.Layout.Footer.HasFocus() {
			return tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
		}

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

		case 'n':
			if !v.Layout.Footer.HasFocus() {
				v.Layout.Container.SetFocus(v.state.Elements.DropDownNamespace)
			}

		case 'h':
			if !v.Layout.Footer.HasFocus() {
				v.GoBack()
			}

		case 'j':
			if !v.Layout.Footer.HasFocus() {
				return tcell.NewEventKey(tcell.KeyDown, 0, tcell.ModNone)
			}

		case 'k':
			if !v.Layout.Footer.HasFocus() {
				return tcell.NewEventKey(tcell.KeyUp, 0, tcell.ModNone)
			}

		case 'l':
			if !v.Layout.Footer.HasFocus() {
				return tcell.NewEventKey(tcell.KeyEnter, 0, tcell.ModNone)
			}

		case ':':
			if !v.Layout.Footer.HasFocus() {
				if !v.state.Toggle.GotoLine {
					v.state.Toggle.GotoLine = true
					v.GotoLine()
				} else {
					v.Layout.Container.SetFocus(v.components.GotoLine.InputField.Primitive())
				}
				return nil
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
	case tcell.KeyEsc, tcell.KeyCtrlO, tcell.KeyEnter, tcell.KeyLeft:
		if v.components.LogStream.TextView.Primitive().HasFocus() {
			v.GoBack()
			return nil
		}
	case tcell.KeyRight:
		// There's nothing further to "enter" from the log view - flash
		// instead of leaving, like a terminal's visual bell.
		if v.components.LogStream.TextView.Primitive().HasFocus() {
			v.FlashLogBorder()
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
			if v.components.LogStream.TextView.Primitive().HasFocus() {
				v.GoBack()
				return nil
			}

		case 'l':
			// There's nothing further to "enter" from the log view - flash
			// instead of leaving, like a terminal's visual bell.
			if v.components.LogStream.TextView.Primitive().HasFocus() {
				v.FlashLogBorder()
				return nil
			}
		case 'H':
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
		case 'p':
			if !v.Layout.Footer.HasFocus() {
				v.state.Toggle.PrettyJSON = !v.state.Toggle.PrettyJSON
				v.components.LogStream.Props.PrettyJSON = v.state.Toggle.PrettyJSON
				v.components.LogStream.Render()
				v.Draw()
			}
		}
	}

	return event
}

// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package view

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gdamore/tcell/v2"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/styles"
)

func (v *View) Init(version string, clusterStatsInterval time.Duration) {
	// ClusterInfo
	v.components.ClusterInfo.Props.Info = fmt.Sprintf(
		"%sAddress%s: %s\n%sVersion:%s %s",
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		v.state.NomadAddress,
		styles.HighlightSecondaryTag,
		styles.StandardColorTag,
		version,
	)

	v.components.ClusterInfo.Bind(v.Layout.Elements.ClusterInfo)
	v.components.ClusterInfo.Render()

	// JumpToJob
	v.components.JumpToJob.Bind(v.Layout.Footer)
	v.components.JumpToJob.Props.DoneFunc = func(key tcell.Key) {
		v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 0)
		v.Layout.Footer.RemoveItem(v.components.JumpToJob.InputField.Primitive())
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)

		id := v.components.JumpToJob.InputField.GetText()
		if id != "" {
			jobID := v.components.JumpToJob.InputField.GetText()
			v.Allocations(jobID)
		}

		v.components.JumpToJob.InputField.SetText("")
		v.state.Toggle.JumpToJob = false
	}

	// GotoLine
	v.components.GotoLine.Bind(v.Layout.Footer)
	v.components.GotoLine.Props.DoneFunc = func(key tcell.Key) {
		v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 0)
		v.Layout.Footer.RemoveItem(v.components.GotoLine.InputField.Primitive())
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)

		if key == tcell.KeyEnter {
			text := v.components.GotoLine.InputField.GetText()
			if line, err := strconv.Atoi(text); err == nil && v.state.Elements.TableMain != nil {
				v.state.Elements.TableMain.Select(line, 0)
			}
		}

		v.components.GotoLine.InputField.SetText("")
		v.state.Toggle.GotoLine = false
	}

	// LogSearchField
	v.components.LogSearch.Bind(v.Layout.Footer)
	v.components.LogSearch.Props.ChangedFunc = func(text string) {
		v.state.Filter.Logs = text
		v.components.LogStream.Props.Filter = text
		v.components.LogStream.Render()
		v.Draw()
	}

	v.components.LogSearch.Props.DoneFunc = func(key tcell.Key) {
		v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 0)
		v.Layout.Footer.RemoveItem(v.components.LogSearch.InputField.Primitive())
		v.Layout.Container.SetFocus(v.components.LogStream.TextView.Primitive())
		v.state.Toggle.LogSearch = false

		v.components.LogStream.Render()
		v.Draw()
	}

	// LogHighlightfield
	v.components.LogHighlight.Bind(v.Layout.Footer)
	v.components.LogHighlight.Props.ChangedFunc = func(text string) {
		v.components.LogStream.Props.Highlight = text
		v.components.LogStream.Render()
		v.Draw()
	}

	v.components.LogHighlight.Props.DoneFunc = func(key tcell.Key) {
		v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 0)
		v.Layout.Footer.RemoveItem(v.components.LogHighlight.InputField.Primitive())
		v.Layout.Container.SetFocus(v.components.LogStream.TextView.Primitive())
		v.state.Toggle.LogHighlight = false

		v.components.LogStream.Render()
		v.Draw()
	}

	// SearchField
	v.components.Search.Bind(v.Layout.Footer)
	v.components.Search.Props.DoneFunc = func(key tcell.Key) {
		// v.components.Search.InputField.SetText("")
		v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 0)
		v.Layout.Footer.RemoveItem(v.components.Search.InputField.Primitive())
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		v.state.Toggle.Search = false
	}

	// JobTable
	v.components.JobTable.Bind(v.Layout.Body)
	v.components.JobTable.Props.HandleNoResources = v.handleNoResources

	// JobStatus
	v.components.JobStatus.Bind(v.Layout.Body)

	// DeploymentTable
	v.components.DeploymentTable.Bind(v.Layout.Body)
	v.components.DeploymentTable.Props.SelectDeployment = func(jobID string) {
		//TODO
	}
	v.components.DeploymentTable.Props.HandleNoResources = v.handleNoResources

	// NamespaceTable
	v.components.NamespaceTable.Bind(v.Layout.Body)
	v.components.NamespaceTable.Props.HandleNoResources = v.handleNoResources

	// Alllocations
	v.components.AllocationTable.Bind(v.Layout.Body)
	v.components.AllocationTable.Props.HandleNoResources = v.handleNoResources
	v.components.AllocationTable.Props.SelectAllocation = func(allocID string) {
		alloc, ok := v.getAllocation(allocID)
		if !ok {
			return
		}

		v.Tasks(alloc)
	}

	v.components.TaskTable.Bind(v.Layout.Body)
	v.components.TaskTable.Props.HandleNoResources = v.handleNoResources
	v.components.TaskTable.Props.SelectTask = func(taskName, allocID string) {
		v.components.LogSearch.InputField.SetText("")
		v.Logs(taskName, allocID, "stdout")
	}
	v.components.TaskTable.BindKey(tcell.KeyCtrlE, func(event *tcell.EventKey) {
		taskName := v.components.TaskTable.GetNameForSelection()
		allocID := v.components.TaskTable.Props.AllocationID

		v.Logs(taskName, allocID, "stderr")
	})

	v.components.TaskTable.BindKey(tcell.KeyRune, func(event *tcell.EventKey) {
		switch event.Rune() {
		case 'e':
			allocID := v.components.TaskTable.Props.AllocationID
			taskName := v.components.TaskTable.GetNameForSelection()
			v.TaskEvents(allocID, taskName)
		case 'x':
			allocID := v.components.TaskTable.Props.AllocationID
			taskName := v.components.TaskTable.GetNameForSelection()
			v.Exec(taskName, allocID)
		}
	})

	// TaskGroupTable
	v.components.TaskGroupTable.Bind(v.Layout.Body)
	v.components.TaskGroupTable.Props.SelectTaskGroup = func(taskGroupID string) {
		//TODO
		v.handleInfo("You selected TaskGroup: %s\n Sorry, selecting task groups isn't implemented yet!", taskGroupID)
	}

	// TaskEventsTable
	v.components.TaskEventsTable.Bind(v.Layout.Body)

	// Logs
	v.components.LogStream.Bind(v.Layout.Body)
	v.components.LogStream.Props.HandleNoResources = v.handleNoResources
	v.components.LogStream.Props.App = v.Layout.Container

	// ClusterStats
	v.components.ClusterStats.Bind(v.Layout.Header.SlotLogo)
	v.components.ClusterStats.Render()

	updateClusterStats := func() {
		v.components.ClusterStats.Props.Stats = v.state.ClusterStats
		v.components.ClusterStats.Props.Err = v.state.ClusterStatsErr
		v.components.ClusterStats.Props.Loaded = v.state.ClusterStats != nil

		v.components.ClusterStats.Render()
		v.Draw()
	}

	v.Watcher.WatchClusterStats(clusterStatsInterval, updateClusterStats)

	// Commands
	v.components.Commands.Bind(v.Layout.Header.SlotCmd)
	v.components.Commands.BindView(v.Layout.Header.SlotViewCmd)
	v.components.Commands.Render()

	// Selections
	v.components.Selections.Bind(v.Layout.Elements.Dropdowns)

	v.components.Selections.Render()

	// Error
	v.components.Error.Bind(v.Layout.Pages)
	v.components.Error.Props.Done = func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Quit" {
			v.Layout.Container.Stop()
			return
		}

		v.Layout.Pages.RemovePage(component.PageNameError)
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		v.GoBack()
	}

	// Info
	v.components.Info.Bind(v.Layout.Pages)
	v.components.Info.Props.Done = func(buttonIndex int, buttonLabel string) {
		v.Layout.Pages.RemovePage(component.PageNameInfo)
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		v.GoBack()
	}

	// Warn
	v.components.Failure.Bind(v.Layout.Pages)
	v.components.Failure.Props.Done = func(buttonIndex int, buttonLabel string) {
		v.Layout.Pages.RemovePage(component.PageNameInfo)
		v.Layout.Container.SetFocus(v.state.Elements.TableMain)
		v.GoBack()
	}

	v.components.Confirm.Bind(v.Layout.Pages)
	selectorModal := v.components.SelectorModal
	selectorModal.Bind(v.Layout.Pages)
	selectorModal.BindKey(tcell.KeyEsc, func() {
		selectorModal.Close()
	})

	v.Watcher.SubscribeHandler(models.HandleError, v.handleError)
	v.Watcher.SubscribeHandler(models.HandleFatal, v.handleFatal)

	stop := make(chan struct{})

	go v.DrawLoop(stop)

	// Set initial view to jobs
	v.Jobs()
}

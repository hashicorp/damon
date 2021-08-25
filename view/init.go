package view

import (
	"fmt"

	"github.com/gdamore/tcell/v2"

	"github.com/hcjulz/damon/component"
	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/styles"
)

func (v *View) Init(version string) {
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

	// LogSearchField
	v.components.LogSearch.Bind(v.Layout.Footer)
	v.components.LogSearch.Props.ChangedFunc = func(text string) {
		v.state.Filter.Logs = text
		v.components.LogStream.Props.Data = filterLogs(v.state.Logs, text)
		v.components.LogStream.Render()
		v.Draw()
	}
	v.components.LogSearch.Props.DoneFunc = func(key tcell.Key) {
		v.Layout.MainPage.ResizeItem(v.Layout.Footer, 0, 0)
		v.Layout.Footer.RemoveItem(v.components.LogSearch.InputField.Primitive())
		v.Layout.Container.SetFocus(v.components.LogStream.TextView.Primitive())
		v.state.Toggle.LogSearch = false
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
		// Reset the LogSearch field in case there was a search string
		// entered previously.
		v.components.LogSearch.InputField.SetText("")
		v.Logs(allocID, "stdout")
	}

	// TaskGroupTable
	v.components.TaskGroupTable.Bind(v.Layout.Body)
	v.components.TaskGroupTable.Props.SelectTaskGroup = func(taskGroupID string) {
		//TODO
		v.handleInfo("You selected TaskGroup: %s\n Sorry, selecting task groups isn't implemented yet!", taskGroupID)
	}

	// Logs
	v.components.LogStream.Bind(v.Layout.Body)
	v.components.LogStream.Props.HandleNoResources = v.handleNoResources

	// Logo
	v.components.Logo.Bind(v.Layout.Header.SlotLogo)
	v.components.Logo.Render()

	// Commands
	v.components.Commands.Bind(v.Layout.Header.SlotCmd)
	v.components.Commands.Render()

	// Selections
	v.components.Selections.Bind(v.Layout.Elements.Dropdowns)
	v.components.Selections.Props.Rerender = func() {
		v.Watcher.ForceUpdate()
	}

	v.components.Selections.Render()

	// Error
	v.components.Error.Bind(v.Layout.Pages)
	v.components.Error.Props.Done = func(buttonIndex int, buttonLabel string) {
		v.Layout.Container.Stop()
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

	v.Watcher.SubscribeHandler(models.HandleError, v.handleError)
	v.Watcher.SubscribeHandler(models.HandleFatal, v.handleFatal)

	stop := make(chan struct{})

	go v.DrawLoop(stop)

	// Set initial view to jobs
	v.Jobs()
}

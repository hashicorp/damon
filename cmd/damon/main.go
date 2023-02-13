// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jessevdk/go-flags"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/nomad"
	"github.com/hcjulz/damon/state"
	"github.com/hcjulz/damon/styles"
	"github.com/hcjulz/damon/version"
	"github.com/hcjulz/damon/view"
	"github.com/hcjulz/damon/watcher"

	"github.com/hcjulz/damon/component"
)

var refreshIntervalDefault = time.Second * 2

type options struct {
	Version bool `short:"v" long:"version" description:"Show Damon version"`
}

func main() {
	// globally overwrite the background color
	tview.Styles.PrimitiveBackgroundColor = tcell.NewRGBColor(40, 44, 48)

	var opts options
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Println("Damon", version.GetHumanVersion())
		os.Exit(0)
	}

	nomadClient, _ := nomad.New(nomad.Default)

	state := initializeState(nomadClient)

	clusterInfo := component.NewClusterInfo()
	selections := component.NewSelections(state)
	selectorModal := component.NewSelectorModal()
	commands := component.NewCommands()
	logo := component.NewLogo()
	jobs := component.NewJobsTable()
	jobStatus := component.NewJobStatus()
	depl := component.NewDeploymentTable()
	namespaces := component.NewNamespaceTable()
	allocations := component.NewAllocationTable()
	taskGroups := component.NewTaskGroupTable()
	taskEvents := component.NewTaskEventsTable()
	taskTable := component.NewTaskTable()
	logs := component.NewLogger()
	jumpToJob := component.NewJumpToJob()
	logSearch := component.NewSearchField("/")
	logHighlight := component.NewSearchField("highlight")
	errorComp := component.NewError()
	info := component.NewInfo()
	failure := component.NewInfo()
	confirm := component.NewModal(
		"confirm",
		"confirm",
		[]string{"cancel", "confirm"},
		styles.TcellColorAttention,
	)

	components := &view.Components{
		ClusterInfo:     clusterInfo,
		Selections:      selections,
		SelectorModal:   selectorModal,
		Commands:        commands,
		Logo:            logo,
		JobTable:        jobs,
		JobStatus:       jobStatus,
		DeploymentTable: depl,
		NamespaceTable:  namespaces,
		AllocationTable: allocations,
		TaskGroupTable:  taskGroups,
		TaskEventsTable: taskEvents,
		TaskTable:       taskTable,
		LogStream:       logs,
		LogHighlight:    logHighlight,
		JumpToJob:       jumpToJob,
		Error:           errorComp,
		Info:            info,
		Failure:         failure,
		LogSearch:       logSearch,
		Confirm:         confirm,
	}

	watcher := watcher.NewWatcher(state, nomadClient, refreshIntervalDefault)
	go watcher.Watch()

	view := view.New(components, watcher, nomadClient, state)
	view.Init(version.GetHumanVersion())

	err = view.Layout.Container.Run()
	if err != nil {
		log.Fatal("cannot initialize view.")
	}
}

func initializeState(client *nomad.Nomad) *state.State {
	state := state.New()
	namespaces, err := client.Namespaces(nil)
	if err != nil {
		log.Fatal("cannot initialize view. Is Nomad running?")
	}

	state.NomadAddress = client.Address()
	state.Namespaces = namespaces

	return state
}

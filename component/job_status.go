// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component

import (
	"fmt"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/models"
	primitive "github.com/hcjulz/damon/primitives"
)

const (
	TitleJobStatus = "Status"
)

type JobStatus struct {
	TextView TextView
	Props    *JobStatusProps
	slot     *tview.Flex
}

type JobStatusProps struct {
	Data *models.JobStatus
}

func NewJobStatus() *JobStatus {
	textView := primitive.NewTextView(tview.AlignLeft)
	return &JobStatus{
		TextView: textView,
		Props:    &JobStatusProps{},
	}
}

func (jobStatus *JobStatus) Bind(slot *tview.Flex) {
	jobStatus.slot = slot
}

func (jobStatus *JobStatus) Render() error {
	if jobStatus.slot == nil {
		return ErrComponentNotBound
	}

	jobStatus.slot.Clear()
	jobStatus.slot.SetDirection(tview.FlexRow)

	if jobStatus.Props.Data.ID == "" {
		jobStatus.TextView.SetText("Status not available.")
		return nil
	}

	jobStatus.TextView.ModifyPrimitive(func(textView *tview.TextView) {
		applyModifiers(textView, jobStatus.Props.Data.ID)
	})

	infoData, err := jobStatus.renderInfoData()
	if err != nil {
		return err
	}

	summary, err := jobStatus.renderSummary()
	if err != nil {
		return err
	}

	latestDeployment, err := jobStatus.renderLatestDeployment()
	if err != nil {
		return err
	}

	deployments, err := jobStatus.renderDeployments()
	if err != nil {
		return err
	}

	allocs, err := jobStatus.renderAllocs()
	if err != nil {
		return err
	}

	var jobStatusText []string
	jobStatusText = append(jobStatusText,
		"\n",
		infoData,
		fmt.Sprintf("\n  %s\n", LabelStatusSummary),
		summary,
		fmt.Sprintf("\n  %s\n", LabelLatestDeployment),
		latestDeployment,
		fmt.Sprintf("\n  %s\n", LabelDeployed),
		deployments,
		fmt.Sprintf("\n  %s\n", LabelAllocations),
		allocs)

	jobStatus.TextView.SetText(strings.Join(jobStatusText, ""))
	jobStatus.slot.AddItem(jobStatus.TextView.Primitive(), 0, 1, true)
	return nil
}

func (jobStatus *JobStatus) renderInfoData() (string, error) {
	j := jobStatus.Props.Data

	tableString := &strings.Builder{}
	tableWriter := newTable(tableString)

	frmt := func(value interface{}) string {
		return fmt.Sprintf("= %s", value)
	}

	infoData := [][]string{
		{LabelID, frmt(j.ID)},
		{LabelName, frmt(j.Name)},
		{LabelSubmitTime, frmt(j.SubmitDate.Format("2006-01-02 15:04:05"))},
		{LabelType, frmt(j.Type)},
		{LabelPriority, frmt(fmt.Sprint(j.Priority))},
		{LabelDatacenters, frmt(j.Datacenters)},
		{LabelNamespace, frmt(j.Namespace)},
		{LabelStatus, frmt(j.Status)},
		{LabelPeriodic, frmt(fmt.Sprint(j.Periodic))},
		{LabelParameterized, frmt(fmt.Sprint(j.Parameterized))},
	}
	if err := tableWriter.Bulk(infoData); err != nil {
		return "", err
	}

	if err := tableWriter.Render(); err != nil {
		return "", err
	}

	return tableString.String(), nil
}

func (jobStatus *JobStatus) renderSummary() (string, error) {
	tasks := jobStatus.Props.Data.TaskGroups

	tableString := &strings.Builder{}
	tableWriter := newTable(tableString)
	tableWriter.Header([]string{LabelTaskGroup, LabelQueued, LabelStarting, LabelRunning, LabelFailed, LabelComplete, LabelLost})

	for _, tg := range tasks {
		row := []string{
			tg.Name,
			fmt.Sprint(tg.Queued),
			fmt.Sprint(tg.Starting),
			fmt.Sprint(tg.Running),
			fmt.Sprint(tg.Failed),
			fmt.Sprint(tg.Complete),
			fmt.Sprint(tg.Lost),
		}
		if err := tableWriter.Append(row); err != nil {
			return "", err
		}
	}

	if err := tableWriter.Render(); err != nil {
		return "", err
	}

	return tableString.String(), nil
}

func (jobStatus *JobStatus) renderLatestDeployment() (string, error) {
	tgs := jobStatus.Props.Data.TaskGroupStatus

	if len(tgs) <= 0 {
		return "", nil
	}

	ts := tgs[0]
	tableString := &strings.Builder{}
	tableWriter := newTable(tableString)

	frmt := func(value interface{}) string {
		return fmt.Sprintf("= %s", value)
	}

	infoData := [][]string{
		{LabelID, frmt(ts.ID)},
		{LabelStatus, frmt(ts.Status)},
		{LabelStatusDescriptionLong, frmt(ts.StatusDescription)},
	}
	if err := tableWriter.Bulk(infoData); err != nil {
		return "", err
	}

	if err := tableWriter.Render(); err != nil {
		return "", err
	}

	return tableString.String(), nil
}

func (jobStatus *JobStatus) renderDeployments() (string, error) {
	tgs := jobStatus.Props.Data.TaskGroupStatus

	tableString := &strings.Builder{}
	tableWriter := newTable(tableString)

	tableWriter.Header([]string{LabelTaskGroup, LabelDesired, LabelPlaced, LabelHealthy, LabelUnhealthy, LabelProgressDeadline})

	for _, t := range tgs {
		row := []string{
			fmt.Sprint(t.ID),
			fmt.Sprint(t.Desired),
			fmt.Sprint(t.Placed),
			fmt.Sprint(t.Healthy),
			fmt.Sprint(t.Unhealthy),
			fmt.Sprint(t.ProgressDeadline),
		}
		if err := tableWriter.Append(row); err != nil {
			return "", err
		}
	}

	if err := tableWriter.Render(); err != nil {
		return "", err
	}

	return tableString.String(), nil
}

func (jobStatus *JobStatus) renderAllocs() (string, error) {
	tgs := jobStatus.Props.Data.Allocations

	tableString := &strings.Builder{}
	tableWriter := newTable(tableString)

	tableWriter.Header([]string{LabelID, LabelNodeID, LabelTaskGroup, LabelVersion, LabelDesired, LabelStatus, LabelCreated, LabelModified})

	for _, t := range tgs {
		row := []string{
			t.ID[0:8],
			t.NodeID[0:8],
			t.TaskGroup,
			fmt.Sprint(t.Version),
			t.DesiredStatus,
			t.Status,
			fmt.Sprintf("%s ago", time.Since(t.Created).Round(time.Second)),
			fmt.Sprintf("%s ago", time.Since(t.Modified).Round(time.Second)),
		}
		if err := tableWriter.Append(row); err != nil {
			return "", err
		}
	}

	if err := tableWriter.Render(); err != nil {
		return "", err
	}

	return tableString.String(), nil
}

// newTable returns a borderless, left-aligned table without separators,
// headers kept verbatim and content never wrapped.
func newTable(w *strings.Builder) *tablewriter.Table {
	cells := tw.CellConfig{
		Formatting: tw.CellFormatting{AutoFormat: tw.Off, AutoWrap: tw.WrapNone},
		Alignment:  tw.CellAlignment{Global: tw.AlignLeft},
	}

	return tablewriter.NewTable(w,
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Borders:  tw.BorderNone,
			Settings: tw.Settings{Lines: tw.LinesNone, Separators: tw.SeparatorsNone},
		})),
		tablewriter.WithConfig(tablewriter.Config{Header: cells, Row: cells}),
	)
}

func applyModifiers(t *tview.TextView, jobID string) {
	t.SetScrollable(true)
	t.SetBorder(true)
	t.SetTitle(fmt.Sprintf("%s (Job: %s)", TitleJobStatus, jobID))
}

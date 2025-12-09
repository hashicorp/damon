// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component

import (
	"fmt"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
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

	var jobStatusText []string
	jobStatusText = append(jobStatusText,
		"\n",
		jobStatus.renderInfoData(),
		fmt.Sprintf("\n  %s\n", LabelStatusSummary),
		jobStatus.renderSummary(),
		fmt.Sprintf("\n  %s\n", LabelLatestDeployment),
		jobStatus.renderLatestDeployment(),
		fmt.Sprintf("\n  %s\n", LabelDeployed),
		jobStatus.renderDeployments(),
		fmt.Sprintf("\n  %s\n", LabelAllocations),
		jobStatus.renderAllocs())

	jobStatus.TextView.SetText(strings.Join(jobStatusText, ""))
	jobStatus.slot.AddItem(jobStatus.TextView.Primitive(), 0, 1, true)
	return nil
}

func (jobStatus *JobStatus) renderInfoData() string {
	j := jobStatus.Props.Data

	tableString := &strings.Builder{}
	tableWriter := tablewriter.NewWriter(tableString)
	format(tableWriter)

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
	tableWriter.AppendBulk(infoData)
	tableWriter.Render()
	return tableString.String()
}

func (jobStatus *JobStatus) renderSummary() string {
	tasks := jobStatus.Props.Data.TaskGroups

	tableString := &strings.Builder{}
	tableWriter := tablewriter.NewWriter(tableString)
	format(tableWriter)
	tableWriter.SetHeader([]string{LabelTaskGroup, LabelQueued, LabelStarting, LabelRunning, LabelFailed, LabelComplete, LabelLost})

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
		tableWriter.Append(row)
	}
	tableWriter.Render()
	return tableString.String()
}

func (jobStatus *JobStatus) renderLatestDeployment() string {
	tgs := jobStatus.Props.Data.TaskGroupStatus

	if len(tgs) <= 0 {
		return ""
	}

	ts := tgs[0]
	tableString := &strings.Builder{}
	tableWriter := tablewriter.NewWriter(tableString)
	format(tableWriter)

	frmt := func(value interface{}) string {
		return fmt.Sprintf("= %s", value)
	}

	infoData := [][]string{
		{LabelID, frmt(ts.ID)},
		{LabelStatus, frmt(ts.Status)},
		{LabelStatusDescriptionLong, frmt(ts.StatusDescription)},
	}
	tableWriter.AppendBulk(infoData)
	tableWriter.Render()
	return tableString.String()
}

func (jobStatus *JobStatus) renderDeployments() string {
	tgs := jobStatus.Props.Data.TaskGroupStatus

	tableString := &strings.Builder{}
	tableWriter := tablewriter.NewWriter(tableString)
	format(tableWriter)

	tableWriter.SetHeader([]string{LabelTaskGroup, LabelDesired, LabelPlaced, LabelHealthy, LabelUnhealthy, LabelProgressDeadline})

	for _, t := range tgs {
		row := []string{
			fmt.Sprint(t.ID),
			fmt.Sprint(t.Desired),
			fmt.Sprint(t.Placed),
			fmt.Sprint(t.Healthy),
			fmt.Sprint(t.Unhealthy),
			fmt.Sprint(t.ProgressDeadline),
		}
		tableWriter.Append(row)
	}
	tableWriter.Render()
	return tableString.String()
}

func (jobStatus *JobStatus) renderAllocs() string {
	tgs := jobStatus.Props.Data.Allocations

	tableString := &strings.Builder{}
	tableWriter := tablewriter.NewWriter(tableString)
	format(tableWriter)

	tableWriter.SetHeader([]string{LabelID, LabelNodeID, LabelTaskGroup, LabelVersion, LabelDesired, LabelStatus, LabelCreated, LabelModified})

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
		tableWriter.Append(row)
	}
	tableWriter.Render()
	return tableString.String()
}

func format(tw *tablewriter.Table) {
	tw.SetBorder(false)
	tw.SetAutoWrapText(false)
	tw.SetAutoFormatHeaders(false)
	tw.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	tw.SetAlignment(tablewriter.ALIGN_LEFT)
	tw.SetCenterSeparator("")
	tw.SetColumnSeparator("")
	tw.SetRowSeparator("")
	tw.SetHeaderLine(false)
	tw.SetTablePadding("\t")
}

func applyModifiers(t *tview.TextView, jobID string) {
	t.SetScrollable(true)
	t.SetBorder(true)
	t.SetTitle(fmt.Sprintf("%s (Job: %s)", TitleJobStatus, jobID))
}

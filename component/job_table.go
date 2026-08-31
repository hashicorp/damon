// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/models"
	primitive "github.com/hcjulz/damon/primitives"
	"github.com/hcjulz/damon/styles"
)

const (
	TableTitleJobs = "Jobs"

	IndicatorWarning = "⚠️"
	IndicatorError   = "❌"
	IndicatorSuccess = "✅"
	IndicatorWaiting = "⌛"
	IndicatorEmpty   = "---"

	// SortBy* are the columns the Jobs table can be sorted by. Ready isn't
	// sortable on its own since it's derived from Status/deployment health.
	SortByID         = "id"
	SortByName       = "name"
	SortByType       = "type"
	SortByNamespace  = "namespace"
	SortByStatus     = "status"
	SortBySubmitTime = "submit_time"
	SortByUptime     = "uptime"
)

var (
	TableHeaderJobs = []string{
		LabelLineNumber,
		LabelID,
		LabelName,
		LabelType,
		LabelNamespace,
		LabelStatus,
		LabelReady,
		LabelSubmitTime,
		LabelUptime,
	}

	// sortCycleOrder is the order CycleSort() advances through.
	sortCycleOrder = []string{
		SortByID,
		SortByName,
		SortByType,
		SortByNamespace,
		SortByStatus,
		SortBySubmitTime,
		SortByUptime,
	}

	sortColumnLabels = map[string]string{
		SortByID:         LabelID,
		SortByName:       LabelName,
		SortByType:       LabelType,
		SortByNamespace:  LabelNamespace,
		SortByStatus:     LabelStatus,
		SortBySubmitTime: LabelSubmitTime,
		SortByUptime:     LabelUptime,
	}
)

//go:generate counterfeiter . SelectJobFunc
type SelectJobFunc func(jobID string)

type JobTable struct {
	Table Table
	Props *JobTableProps

	slot *tview.Flex
}

type JobTableProps struct {
	SelectJob         SelectJobFunc
	HandleNoResources models.HandlerFunc

	Data      []*models.Job
	Namespace string

	// SortColumn is one of the SortBy* constants, or "" for no active sort
	// (natural order). SortAscending is only meaningful when SortColumn is set.
	SortColumn    string
	SortAscending bool
}

func NewJobsTable() *JobTable {
	t := primitive.NewTable()

	jt := &JobTable{
		Table: t,
		Props: &JobTableProps{},
	}

	return jt
}

func (j *JobTable) Bind(slot *tview.Flex) {
	j.slot = slot
}

func (j *JobTable) Render() error {
	if err := j.validate(); err != nil {
		return err
	}

	j.reset()

	j.Table.SetTitle("%s (%s)", TableTitleJobs, j.Props.Namespace)

	if len(j.Props.Data) == 0 {
		j.Props.HandleNoResources(
			"%sno jobs available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)

		return nil
	}

	j.Table.SetSelectedFunc(j.jobSelected)
	j.Table.RenderHeader(j.header())
	j.renderRows()

	j.slot.AddItem(j.Table.Primitive(), 0, 1, false)
	return nil
}

// header returns TableHeaderJobs with an ascending/descending arrow appended
// to the currently active sort column's label, if any.
func (j *JobTable) header() []string {
	header := make([]string, len(TableHeaderJobs))
	copy(header, TableHeaderJobs)

	label, ok := sortColumnLabels[j.Props.SortColumn]
	if !ok {
		return header
	}

	arrow := "▲"
	if !j.Props.SortAscending {
		arrow = "▼"
	}

	for i, h := range header {
		if h == label {
			header[i] = fmt.Sprintf("%s %s", h, arrow)
			break
		}
	}

	return header
}

// CycleSort advances to the next sortable column (wrapping around), always
// resetting to ascending order.
func (j *JobTable) CycleSort() {
	next := sortCycleOrder[0]
	for i, col := range sortCycleOrder {
		if col == j.Props.SortColumn {
			next = sortCycleOrder[(i+1)%len(sortCycleOrder)]
			break
		}
	}

	j.Props.SortColumn = next
	j.Props.SortAscending = true
}

// FlipSortDirection toggles ascending/descending on the current sort column.
// It's a no-op if no column is currently sorted.
func (j *JobTable) FlipSortDirection() {
	if j.Props.SortColumn == "" {
		return
	}

	j.Props.SortAscending = !j.Props.SortAscending
}

func (j *JobTable) GetIDForSelection() string {
	row, _ := j.Table.GetSelection()
	return j.Table.GetCellContent(row, 1)
}

func (j *JobTable) validate() error {
	if j.Props.SelectJob == nil || j.Props.HandleNoResources == nil {
		return ErrComponentPropsNotSet
	}

	if j.slot == nil {
		return ErrComponentNotBound
	}

	return nil
}

func (j *JobTable) reset() {
	j.slot.Clear()
	j.Table.Clear()
}

func (j *JobTable) jobSelected(row, _ int) {
	jobID := j.Table.GetCellContent(row, 1)
	j.Props.SelectJob(jobID)
}

func (j *JobTable) renderRows() {
	data := make([]*models.Job, len(j.Props.Data))
	copy(data, j.Props.Data)
	sortJobs(data, j.Props.SortColumn, j.Props.SortAscending)

	for i, job := range data {
		index := i + 1

		ready, rowColor := readyStatus(job.ReadyStatus, job.Status, job.DeploymentStatus)

		row := []string{
			strconv.Itoa(index),
			job.ID,
			job.Name,
			job.Type,
			job.Namespace,
			job.Status,
			ready,
			job.SubmitTime.Format(time.RFC3339),
			formatTimeSince(time.Since(job.SubmitTime)),
		}

		j.Table.RenderRow(row, index, rowColor)
	}
}

// sortJobs sorts jobs in place by column, ascending or descending. It's a
// no-op if column is empty or unrecognized.
func sortJobs(jobs []*models.Job, column string, ascending bool) {
	compare := func(a, b *models.Job) int {
		switch column {
		case SortByID:
			return strings.Compare(a.ID, b.ID)
		case SortByName:
			return strings.Compare(a.Name, b.Name)
		case SortByType:
			return strings.Compare(a.Type, b.Type)
		case SortByNamespace:
			return strings.Compare(a.Namespace, b.Namespace)
		case SortByStatus:
			return strings.Compare(a.Status, b.Status)
		case SortBySubmitTime:
			return compareTime(a.SubmitTime, b.SubmitTime)
		case SortByUptime:
			// Uptime is derived from SubmitTime, but as its own metric:
			// smaller uptime (most recently submitted) sorts first when
			// ascending, independent of SubmitTime's own sort direction.
			return compareDuration(time.Since(a.SubmitTime), time.Since(b.SubmitTime))
		default:
			return 0
		}
	}

	if _, ok := sortColumnLabels[column]; !ok {
		return
	}

	sort.SliceStable(jobs, func(i, k int) bool {
		c := compare(jobs[i], jobs[k])
		if !ascending {
			c = -c
		}
		return c < 0
	})
}

func compareTime(a, b time.Time) int {
	switch {
	case a.Before(b):
		return -1
	case a.After(b):
		return 1
	default:
		return 0
	}
}

func compareDuration(a, b time.Duration) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// readyStatus renders a job's healthy/desired allocation count along with
// a status indicator based on the job's latest deployment.
func readyStatus(status models.ReadyStatus, jobStatus string, deploymentStatus string) (string, tcell.Color) {
	if jobStatus != models.StatusRunning {
		return IndicatorEmpty, tcell.ColorDarkGrey
	}

	statusIndicator := IndicatorWarning
	color := tcell.ColorWhite

	if status.Unhealthy > 0 {
		statusIndicator = IndicatorError
		color = tcell.ColorRed
	} else if status.Desired == status.Healthy {
		statusIndicator = IndicatorSuccess
	} else if deploymentStatus == models.StatusRunning {
		statusIndicator = IndicatorWaiting
		color = tcell.ColorOrange
	}

	return fmt.Sprintf("%d/%d %s", status.Healthy, status.Desired, statusIndicator), color
}

func formatTimeSince(since time.Duration) string {
	if since.Seconds() < 60 {
		return fmt.Sprintf("%.0fs", since.Seconds())
	}

	if since.Minutes() < 60 {
		return fmt.Sprintf("%.0fm", since.Minutes())
	}

	if since.Hours() < 60 {
		return fmt.Sprintf("%.0fh", since.Hours())
	}

	return fmt.Sprintf("%.0fd", (since.Hours() / 24))
}

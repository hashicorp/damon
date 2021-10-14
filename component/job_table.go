package component

import (
	"fmt"
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
)

var (
	TableHeaderJobs = []string{
		LabelID,
		LabelName,
		LabelType,
		LabelNamespace,
		LabelStatus,
		LabelReady,
		LabelSubmitTime,
		LabelUptime,
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
	j.Table.RenderHeader(TableHeaderJobs)
	j.renderRows()

	j.slot.AddItem(j.Table.Primitive(), 0, 1, false)
	return nil
}

func (j *JobTable) GetIDForSelection() string {
	row, _ := j.Table.GetSelection()
	return j.Table.GetCellContent(row, 0)
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
	jobID := j.Table.GetCellContent(row, 0)
	j.Props.SelectJob(jobID)
}

func (j *JobTable) renderRows() {
	for i, job := range j.Props.Data {
		reaadyStatus, rowColor := readyStatus(job.ReadyStatus, job.Status, job.DeploymentStatus)
		row := []string{
			job.ID,
			job.Name,
			job.Type,
			job.Namespace,
			job.Status,
			reaadyStatus,
			job.SubmitTime.Format(time.RFC3339),
			formatTimeSince(time.Since(job.SubmitTime)),
		}

		index := i + 1

		j.Table.RenderRow(row, index, rowColor)
	}
}

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

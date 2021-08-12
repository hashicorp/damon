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
)

var (
	TableHeaderJobs = []string{
		LabelID,
		LabelName,
		LabelType,
		LabelNamespace,
		LabelStatus,
		LabelStatusSummary,
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
		row := []string{
			job.ID,
			job.Name,
			job.Type,
			job.Namespace,
			job.Status,
			fmt.Sprintf("%d/%d", job.StatusSummary.Running, job.StatusSummary.Total),
			job.SubmitTime.Format(time.RFC3339),
			formatTimeSince(time.Since(job.SubmitTime)),
		}

		index := i + 1

		c := j.cellColor(job.Status, job.Type, job.StatusSummary)

		j.Table.RenderRow(row, index, c)
	}
}

func (j *JobTable) cellColor(status, typ string, summary models.Summary) tcell.Color {
	c := tcell.ColorWhite

	switch status {
	case models.StatusRunning:
		if summary.Total != summary.Running &&
			typ == models.TypeService {
			c = styles.TcellColorAttention
		}
	case models.StatusPending:
		c = tcell.ColorYellow
	case models.StatusDead, models.StatusFailed:
		c = tcell.ColorRed

		if typ == models.TypeBatch {
			c = tcell.ColorDarkGrey
		}
	}

	return c
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

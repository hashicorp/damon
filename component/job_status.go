package component

import (
	"fmt"
	"strings"
	"time"

	"github.com/hcjulz/damon/models"
	primitive "github.com/hcjulz/damon/primitives"
	"github.com/olekukonko/tablewriter"
	"github.com/rivo/tview"
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
	textView.ModifyPrimitive(applyModifiers)
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

	var jobStatusText []string
	jobStatusText = append(jobStatusText,
		"\n",
		jobStatus.renderInfoData(),
		"\n  Summary\n",
		jobStatus.renderSummary(),
		"\n  Latest Deployment\n",
		jobStatus.renderLatestDeployment(),
		"\n  Deployed\n",
		jobStatus.renderDeployments(),
		"\n  Allocations\n",
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
		{"ID", frmt(j.ID)},
		{"Name", frmt(j.Name)},
		{"Submit Date", frmt(j.SubmitDate.Format("2006-01-02 15:04:05"))},
		{"Type", frmt(j.Type)},
		{"Priority", frmt(fmt.Sprint(j.Priority))},
		{"Datacenters", frmt(j.Datacenters)},
		{"Namespace", frmt(j.Namespace)},
		{"Status", frmt(j.Status)},
		{"Periodic", frmt(fmt.Sprint(j.Periodic))},
		{"Parameterized", frmt(fmt.Sprint(j.Parameterized))},
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
	tableWriter.SetHeader([]string{"Task Group", "Queued", "Starting", "Running", "Failed", "Complete", "Lost"})

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
		{"ID", frmt(ts.ID)},
		{"Status", frmt(ts.Status)},
		{"Status Description", frmt(ts.StatusDescription)},
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

	tableWriter.SetHeader([]string{"Task Group", "Desired", "Placed", "Healthy", "Unhealthy", "Progress Deadline"})

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

	tableWriter.SetHeader([]string{"ID", "Node ID", "Task Group", "Version", "Desired", "Status", "Created", "Modified"})

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

func applyModifiers(t *tview.TextView) {
	t.SetScrollable(true)
	t.SetBorder(true)
	t.SetTitle("Job Status")
}

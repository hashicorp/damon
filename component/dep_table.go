package component

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/models"
	primitive "github.com/hcjulz/damon/primitives"
	"github.com/hcjulz/damon/styles"
)

const (
	TableTitleDeployments = "Deployments"
)

var (
	TableHeaderDeployments = []string{
		LabelID,
		LabelJobID,
		LabelNamespace,
		LabelStatus,
		LabelStatusDescription,
	}
)

//go:generate counterfeiter . SelectJobFunc
type SelectFunc func(id string)

type DeploymentTable struct {
	Table Table
	Props *DeploymentTableProps

	slot *tview.Flex
}

type DeploymentTableProps struct {
	SelectDeployment  SelectFunc
	HandleNoResources models.HandlerFunc

	Data      []*models.Deployment
	Namespace string
}

func NewDeploymentTable() *DeploymentTable {
	t := primitive.NewTable()

	dt := &DeploymentTable{
		Table: t,
		Props: &DeploymentTableProps{},
	}

	dt.Table.SetSelectedFunc(dt.deploymentSelected)

	return dt
}

func (d *DeploymentTable) Bind(slot *tview.Flex) {
	d.slot = slot

}

func (d *DeploymentTable) Render() error {
	if d.Props.SelectDeployment == nil || d.Props.HandleNoResources == nil {
		return ErrComponentPropsNotSet
	}

	if d.slot == nil {
		return ErrComponentNotBound
	}

	d.reset()

	if len(d.Props.Data) == 0 {
		d.Props.HandleNoResources(
			"%sno deployments available\n¯%s\\_( ͡• ͜ʖ ͡•)_/¯",
			styles.HighlightPrimaryTag,
			styles.HighlightSecondaryTag,
		)

		return nil
	}

	d.Table.SetTitle(fmt.Sprintf("%s (%s)", TableTitleDeployments, d.Props.Namespace))

	d.Table.RenderHeader(TableHeaderDeployments)
	d.renderRows()

	d.slot.AddItem(d.Table.Primitive(), 0, 1, false)
	return nil
}

func (d *DeploymentTable) reset() {
	d.slot.Clear()
	d.Table.Clear()
}

func (d *DeploymentTable) deploymentSelected(row, column int) {
	deplID := d.Table.GetCellContent(row, 0)
	d.Props.SelectDeployment(deplID)
}

func (d *DeploymentTable) renderRows() {
	for i, dep := range d.Props.Data {
		row := []string{
			dep.ID,
			dep.JobID,
			dep.Namespace,
			dep.Status,
			dep.StatusDescription,
		}

		index := i + 1

		c := d.getCellColor(dep.Status)
		d.Table.RenderRow(row, index, c)
	}
}

func (d *DeploymentTable) getCellColor(status string) tcell.Color {
	c := tcell.ColorWhite

	switch status {
	case models.StatusRunning:
		c = styles.TcellColorHighlighPrimary
	case models.StatusPending:
		c = tcell.ColorYellow
	case models.StatusFailed:
		c = tcell.ColorRed
	}

	return c
}

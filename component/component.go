package component

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/models"
)

const (
	LabelID                = "ID"
	LabelJobID             = "JobID"
	LabelType              = "Type"
	LabelName              = "Name"
	LabelNamespace         = "Namespace"
	LabelStatus            = "Status"
	LabelStatusDescription = "Description"
	LabelStatusSummary     = "Summary"
	LabelDescription       = "Description"
	LabelCount             = "Count"
	LabelSubmitTime        = "SubmitTime"
	LabelUptime            = "Uptime"
	LabelDesiredStatus     = "DesiredStatus"
	LabelTaskGroup         = "TaskGroup"
	LabelTime              = "Time"
	LabelMessage           = "Message"

	LabelRunning  = "Running"
	LabelStarting = "Starting"
	LabelComplete = "Complete"
	LabelQueued   = "Queued"
	LabelLost     = "Lost"
	LabelFailed   = "Failed"

	LabelDesired               = "Desired"
	LabelHealthy               = "Healthy"
	LabelUnhealthy             = "Unhealthy"
	LabelPlaced                = "Placed"
	LabelProgressDeadline      = "Progress Deadline"
	LabelVersion               = "Version"
	LabelStatusDescriptionLong = "Status Description"
	LabelPriority              = "Priority"
	LabelDatacenters           = "Datacenters"
	LabelPeriodic              = "Periodic"
	LabelParameterized         = "Parameterized"

	LabelLatestDeployment = "Latest Deployment"
	LabelDeployed         = "Deployed"
	LabelAllocations      = "Allocations"

	LabelCreated  = "Created"
	LabelModified = "Modified"

	LabelNodeID   = "NodeID"
	LabelNodeName = "NodeName"

	ErrComponentNotBound    = models.Sentinel("component not bound")
	ErrComponentPropsNotSet = models.Sentinel("component properties not set")
)

//go:generate counterfeiter . DoneModalFunc
type DoneModalFunc func(buttonIndex int, buttonLabel string)

type Primitive interface {
	Primitive() tview.Primitive
}

//go:generate counterfeiter . Table
type Table interface {
	Primitive
	SetTitle(format string, args ...interface{})
	GetCellContent(row, column int) string
	GetSelection() (row, column int)
	Clear()
	RenderHeader(data []string)
	RenderRow(data []string, index int, c tcell.Color)
	SetSelectedFunc(fn func(row, column int))
}

//go:generate counterfeiter . TextView
type TextView interface {
	Primitive
	GetText() string
	SetText(text string)
	Clear()
	ModifyPrimitive(f func(t *tview.TextView))
}

//go:generate counterfeiter . Modal
type Modal interface {
	Primitive
	SetDoneFunc(handler func(buttonIndex int, buttonLabel string))
	SetText(text string)
	SetFocus(index int)
	Container() tview.Primitive
}

//go:generate counterfeiter . InputField
type InputField interface {
	Primitive
	SetDoneFunc(handler func(k tcell.Key))
	SetChangedFunc(handler func(text string))
	SetAutocompleteFunc(callback func(currentText string) (entries []string))
	SetText(text string)
	GetText() string
}

//go:generate counterfeiter . DropDown
type DropDown interface {
	Primitive
	SetOptions(options []string, selected func(text string, index int))
	SetCurrentOption(index int)
	SetSelectedFunc(selected func(text string, index int))
}

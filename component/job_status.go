package component

import (
	"github.com/rivo/tview"
	primitive "github.com/hcjulz/damon/primitives"
)

type JobStatus struct {
	TextView TextView
	Status    string
	slot     *tview.Flex
}

func NewJobStatus() *JobStatus {
	textView := primitive.NewTextView(tview.AlignLeft)
	textView.ModifyPrimitive(applyModifiers)
	return &JobStatus{
		TextView: textView,
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
	jobStatus.TextView.Clear()

	if len(jobStatus.Status) > 0 {
		jobStatus.TextView.SetText(jobStatus.Status)
	} else {
		jobStatus.TextView.SetText("Status not available.")
	}
	jobStatus.slot.AddItem(jobStatus.TextView.Primitive(), 0, 1, true)
	return nil
}

func applyModifiers(t *tview.TextView) {
	t.SetScrollable(true)
	t.SetBorder(true)
	t.ScrollToEnd()
	t.SetTitle("Job Status")
}

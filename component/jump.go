package component

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/hcjulz/damon/models"
	"github.com/hcjulz/damon/primitives"
)

const jumpToJobPlaceholder = "(hit enter or esc to leave)"

type SetDoneFunc func(key tcell.Key)

type JumpToJob struct {
	InputField InputField
	Props      *JumpToJobProps
	slot       *tview.Flex
}

type JumpToJobProps struct {
	DoneFunc SetDoneFunc
	Jobs     []*models.Job
}

func NewJumpToJob() *JumpToJob {
	jj := &JumpToJob{}
	jj.Props = &JumpToJobProps{}

	in := primitives.NewInputField("jump: ", jumpToJobPlaceholder)

	in.SetAutocompleteFunc(func(currentText string) (entries []string) {
		return jj.find(currentText)
	})

	jj.InputField = in
	return jj
}

func (jj *JumpToJob) Render() error {
	if err := jj.validate(); err != nil {
		return err
	}

	jj.InputField.SetDoneFunc(jj.Props.DoneFunc)
	jj.slot.AddItem(jj.InputField.Primitive(), 0, 2, false)
	return nil
}

func (jj *JumpToJob) validate() error {
	if jj.Props.DoneFunc == nil {
		return ErrComponentPropsNotSet
	}

	if jj.slot == nil {
		return ErrComponentNotBound
	}

	return nil
}

func (jj *JumpToJob) Bind(slot *tview.Flex) {
	jj.slot = slot
}

func (jj *JumpToJob) find(text string) []string {
	result := []string{}
	if text == "" {
		return result
	}

	for _, j := range jj.Props.Jobs {
		ok := strings.Contains(j.ID, text)
		if ok {
			result = append(result, j.ID)
		}
	}

	return result
}

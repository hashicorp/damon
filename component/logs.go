package component

import (
	"strings"

	"github.com/rivo/tview"

	"github.com/hcjulz/damon/models"
	primitive "github.com/hcjulz/damon/primitives"
	"github.com/hcjulz/damon/styles"
)

type Logger struct {
	TextView TextView
	Props    *LogStreamProps
	slot     *tview.Flex
	buf      strings.Builder
}

type LogStreamProps struct {
	HandleNoResources models.HandlerFunc
	Data              []byte
}

func NewLogger() *Logger {
	t := primitive.NewTextView(tview.AlignLeft)
	t.ModifyPrimitive(applyLogModifiers)

	return &Logger{
		TextView: t,
		Props:    &LogStreamProps{},
		buf:      strings.Builder{},
	}
}

func (l *Logger) Bind(slot *tview.Flex) {
	l.slot = slot
}

func (l *Logger) Render() error {
	if l.slot == nil {
		return ErrComponentNotBound
	}

	if l.Props.HandleNoResources == nil {
		return ErrComponentPropsNotSet
	}

	l.slot.Clear()
	l.TextView.Clear()

	if len(l.Props.Data) == 0 {
		l.Props.HandleNoResources(
			"%sWHOOOPS, no Logs found",
			styles.HighlightSecondaryTag,
		)
		return nil
	}

	l.TextView.SetText(string(l.Props.Data))
	l.slot.AddItem(l.TextView.Primitive(), 0, 1, true)
	return nil
}

func applyLogModifiers(t *tview.TextView) {
	t.SetScrollable(true)
	t.SetBorder(true)
	t.ScrollToEnd()
	t.SetTitle("Logs")
}

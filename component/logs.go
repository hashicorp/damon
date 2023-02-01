// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package component

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/rivo/tview"

	"github.com/hcjulz/damon/models"
	primitive "github.com/hcjulz/damon/primitives"
	"github.com/hcjulz/damon/styles"
)

type renderFunc func() error

type Logger struct {
	TextView TextView
	Props    *LogStreamProps
	slot     *tview.Flex
	buf      strings.Builder
	mutex    sync.Mutex
}

type LogStreamProps struct {
	HandleNoResources models.HandlerFunc
	Filter            string
	Highlight         string
	Data              []byte
	ChangedFunc       func()
	TaskName          string
	App               *tview.Application
}

func NewLogger() *Logger {
	t := primitive.NewTextView(tview.AlignLeft)

	l := &Logger{
		TextView: t,
		Props:    &LogStreamProps{},
		buf:      strings.Builder{},
	}

	t.ModifyPrimitive(l.applyLogModifiers)
	return l

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

	l.ClearDisplay()

	if len(l.Props.Data) == 0 {
		l.Props.HandleNoResources(
			"%sWHOOOPS, no Logs found",
			styles.HighlightSecondaryTag,
		)
		return nil
	}

	lines := bytes.Split(l.Props.Data, []byte("\n"))
	if len(lines) > 1000 {
		rem := len(lines) % 1000
		lines = lines[rem:]
		l.Props.Data = bytes.Join(lines, []byte("\n"))
	}

	display := filter(l.Props.Data, l.Props.Filter)

	if l.Props.Filter == "" {
		display = highlight(display, l.Props.Highlight)
	}

	l.Clear()

	l.SetText(string(display))

	l.Display()

	return nil
}

func (l *Logger) ClearDisplay() {
	l.slot.Clear()
}

func (l *Logger) Display() {
	l.slot.AddItem(l.TextView.Primitive(), 0, 1, true)
}

func (l *Logger) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// l.Props.App.QueueUpdate(func() {
	l.TextView.Clear()
	// })
}

func (l *Logger) SetText(log string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	// l.Props.App.QueueUpdateDraw(func() {
	l.TextView.SetText(log)
	// })
}

func filter(logs []byte, filter string) []byte {
	if filter == "" {
		return logs
	}

	buf := bytes.Buffer{}
	defer buf.Reset()

	rx := regexp.MustCompile(filter)

	logLines := bytes.Split(logs, []byte("\n"))
	var result []byte
	for _, log := range logLines {
		if rx.Match(log) {
			idx := rx.FindIndex([]byte(log))
			fmt.Fprintf(
				&buf,
				"%s%s%s%s%s%s\n",
				[]byte(styles.ColorLighGreyTag),
				log[:idx[0]],
				[]byte(styles.HighlightSecondaryTag),
				log[idx[0]:idx[1]],
				[]byte(styles.ColorLighGreyTag),
				log[idx[1]:],
			)

			result = append(result, buf.Bytes()...)
			buf.Reset()
		}
	}
	return result
}

func highlight(logs []byte, highlight string) []byte {
	if highlight == "" {
		return logs
	}

	buf := bytes.Buffer{}
	defer buf.Reset()

	rx, _ := regexp.Compile(highlight)
	logLines := bytes.Split(logs, []byte("\n"))
	var result []byte
	for _, log := range logLines {
		if rx.Match(log) {
			fmt.Fprintf(&buf, "%s%s%s\n",
				[]byte(styles.HighlightSecondaryTag),
				log,
				[]byte(styles.ColorWhiteTag),
			)
		} else {
			fmt.Fprintf(&buf, "%s\n", log)
		}

		result = append(result, buf.Bytes()...)
		buf.Reset()
	}

	return result
}

func (l *Logger) applyLogModifiers(t *tview.TextView) {
	t.SetScrollable(true)
	t.SetBorder(true)
	t.ScrollToEnd()
	t.SetTitle("Logs")
	// t.SetMaxLines(1000)
}

// Copyright IBM Corp. 2021, 2023
// SPDX-License-Identifier: MPL-2.0

package component

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"

	primitive "github.com/hcjulz/damon/primitives"
	"github.com/hcjulz/damon/styles"
)

var (
	MainCommands = []string{
		fmt.Sprintf("%sCommands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<ctrl-j>%s to display Jobs", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<ctrl-d>%s to display Deployments", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<ctrl-n>%s to display Namespaces", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<ctrl-p>%s to jump to a Job", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<ctrl-c>%s to Quit", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<n>%s to focus the Namespace dropdown", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<h/j/k/l>%s left/down/up/right (same as esc/down/up/enter)", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<:NUM>%s to jump to line NUM", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}

	JobCommands = []string{
		fmt.Sprintf("%sJob Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<Enter>%s to display allocations", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<t>%s to display TaskGroups for the selected Job", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<i>%s to display information for the selected Job", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<ctrl-s>%s start/stop the selected Job", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s</>%s apply filter", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<s>%s cycle sort column", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<S>%s flip sort direction", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}

	AllocCommands = []string{
		// fmt.Sprintf("\n%sAlloc Commands:", styles.HighlightSecondaryTag),
	}

	TaskCommands = []string{
		fmt.Sprintf("%s<e>%s to display events for a Task", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<ctrl-e>%s to display STDERR logs", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<Enter>%s to display STDOUT logs", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<x>%s to exec into a Task", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}

	LogCommands = []string{
		fmt.Sprintf("%sLog Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<Enter> | <ESC> | <h/l> | <left/right>%s to leave", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s</>%s apply filter", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<H>%s highlight", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<s>%s stop log stream", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<r>%s resume log stream", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<p>%s pretty-print JSON lines", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}

	DeploymentCommands = []string{}

	NoViewCommands = []string{}
)

type Commands struct {
	TextView     TextView
	ViewTextView TextView
	Props        *CommandsProps
	slot         *tview.Flex
	slotView     *tview.Flex
}

type CommandsProps struct {
	MainCommands []string
	ViewCommands []string
}

func NewCommands() *Commands {
	textView := primitive.NewTextView(tview.AlignLeft)
	textView.ModifyPrimitive(disableWrap)

	viewTextView := primitive.NewTextView(tview.AlignLeft)
	viewTextView.ModifyPrimitive(disableWrap)

	return &Commands{
		TextView:     textView,
		ViewTextView: viewTextView,
		Props: &CommandsProps{
			MainCommands: MainCommands,
			ViewCommands: JobCommands,
		},
	}
}

// disableWrap keeps each command on a single row: on a narrow terminal a
// wrapped line pushes every entry below it down, which can shove later
// commands past the header's visible height entirely. Truncating a long
// line is preferable to losing whole entries off-screen.
func disableWrap(t *tview.TextView) {
	t.SetWrap(false)
}

func (c *Commands) Update(commands []string) {
	c.Props.ViewCommands = commands

	c.updateText()
}

func (c *Commands) Render() error {
	if c.slot == nil {
		return ErrComponentNotBound
	}

	c.updateText()

	c.slot.AddItem(c.TextView.Primitive(), 0, 1, false)

	if c.slotView != nil {
		c.slotView.AddItem(c.ViewTextView.Primitive(), 0, 1, false)
	}

	return nil
}

func (c *Commands) updateText() {
	c.TextView.SetText(strings.Join(c.Props.MainCommands, "\n"))

	if c.slotView != nil {
		c.ViewTextView.SetText(strings.Join(c.Props.ViewCommands, "\n"))
	}
}

func (c *Commands) Bind(slot *tview.Flex) {
	c.slot = slot
}

// BindView binds the slot that displays the current view's context-specific
// commands, rendered to the right of the main Commands box.
func (c *Commands) BindView(slot *tview.Flex) {
	c.slotView = slot
}

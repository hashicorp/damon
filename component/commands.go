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
	}

	JobCommands = []string{
		fmt.Sprintf("\n%sJob Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<Enter>%s to display allocations", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<t>%s to display TaskGroups for a Job", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<i>%s to display information for Job", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<ctrl-s>%s start/stop a Job", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s</>%s apply filter", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}

	AllocCommands = []string{
		fmt.Sprintf("\n%sAlloc Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<Enter>%s to display STDOUT logs", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<ctrl-e>%s to display STDERR logs", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s<e>%s to display events for an allocation", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}

	LogCommands = []string{
		fmt.Sprintf("\n%sLog Commands:", styles.HighlightSecondaryTag),
		fmt.Sprintf("%s<Enter> | <ESC>%s to leave", styles.HighlightPrimaryTag, styles.StandardColorTag),
		fmt.Sprintf("%s</>%s apply filter", styles.HighlightPrimaryTag, styles.StandardColorTag),
	}

	DeploymentCommands = []string{}

	NoViewCommands = []string{}
)

type Commands struct {
	TextView TextView
	Props    *CommandsProps
	slot     *tview.Flex
}

type CommandsProps struct {
	MainCommands []string
	ViewCommands []string
}

func NewCommands() *Commands {
	return &Commands{
		TextView: primitive.NewTextView(tview.AlignLeft),
		Props: &CommandsProps{
			MainCommands: MainCommands,
			ViewCommands: JobCommands,
		},
	}
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
	return nil
}

func (c *Commands) updateText() {
	commands := append(c.Props.MainCommands, c.Props.ViewCommands...)
	cmds := strings.Join(commands, "\n")
	c.TextView.SetText(cmds)
}

func (c *Commands) Bind(slot *tview.Flex) {
	c.slot = slot
}

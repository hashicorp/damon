package styles

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

func GetBackgroundColor() tcell.Color {
	return tcell.NewRGBColor(40, 44, 48)
}

var (
	TcellBackgroundColor = tcell.NewRGBColor(40, 44, 48)

	HighlightPrimaryHex   = "#26ffe6"
	HighlightSecondaryHex = "#baff26"
	StandardColorHex      = "#00b57c"
	ColorActiveHex        = "#b3f1ff"
	ColorWhiteHex         = "#ffffff"
	ColorLightGreyHex     = "#cccccc"
	ColorModalInfoHex     = "#61877f"
	ColorAttentionHex     = "#d98b6a"

	StandardColorTag      = fmt.Sprintf("[%s]", StandardColorHex)
	HighlightPrimaryTag   = fmt.Sprintf("[%s]", HighlightPrimaryHex)
	HighlightSecondaryTag = fmt.Sprintf("[%s]", HighlightSecondaryHex)
	ColorActiveTag        = fmt.Sprintf("[%s]", ColorActiveHex)
	ColorWhiteTag         = fmt.Sprintf("[%s]", ColorWhiteHex)
	ColorLighGreyTag      = fmt.Sprintf("[%s]", ColorLightGreyHex)

	TcellColorHighlighPrimary   = tcell.GetColor(HighlightPrimaryHex)
	TcellColorHighlighSecondary = tcell.GetColor(HighlightSecondaryHex)
	TcellColorStandard          = tcell.GetColor(StandardColorHex)
	TcellColorActive            = tcell.GetColor(ColorActiveHex)
	TcellColorModalInfo         = tcell.GetColor(ColorModalInfoHex)
	TcellColorAttention         = tcell.GetColor(ColorAttentionHex)
)

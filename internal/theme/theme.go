package theme

import (
	"fmt"
	"os"

	"github.com/robbiew/artvu-ansi-gallery/internal/ansi"
	"github.com/robbiew/artvu-ansi-gallery/pkg/showrune"
	escapes "github.com/snugfox/ansi-escapes"
)

var (
	Block = " "
)

func ShowHeader(w int) {
	fmt.Fprintf(os.Stdout, escapes.CursorHide)

	count := 0
	for i := 1; i <= w; i++ {
		fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 0))
		fmt.Fprintf(os.Stdout, ansi.BgCyan+Block+ansi.Reset)
		count += i
	}
	fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 0))

	s := "ArVu Galery v1.0"
	hw := w
	fmt.Fprintf(os.Stdout, ansi.BgCyan+ansi.BrightWhite+(fmt.Sprintf("%[1]*s", -hw, fmt.Sprintf("%[1]*s", (hw+len(s))/2, s)))+ansi.Reset)
}

func ShowFooter(w int, h int) {
	fmt.Fprintf(os.Stdout, escapes.CursorHide)

	fmt.Fprintf(os.Stdout, escapes.CursorPos(0, h))

	sChars := 21
	s := ansi.BrightCyan + "made with " + ansi.BrightRed + showrune.Heart + ansi.BrightCyan + " by alpha"
	fw := w - 2 + sChars

	fmt.Fprintf(os.Stdout, (fmt.Sprintf(" "+ansi.BgCyan+"%[1]*s", -fw, fmt.Sprintf("%[1]*s", (fw+len(s))/2, s)))+ansi.Reset)
}

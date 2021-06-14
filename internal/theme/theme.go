package theme

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/robbiew/artvu-ansi-gallery/internal/ansi"
	"github.com/robbiew/artvu-ansi-gallery/pkg/showrune"
	escapes "github.com/snugfox/ansi-escapes"
)

//go:embed splash.80.ans
var s80 string

//go:embed splash.132.ans
var s132 string

var (
	Block = " "
)

func ShowHeader(w int, headerH int, path string, rootDir string, themeDir string) {
	// Dir item label
	var res string
	str := path
	res = strings.Join(strings.Split(str, rootDir+"/"), "")
	if res == rootDir {
		res = "dir: root"

	} else {
		res = "dir: " + res + "/"
	}
	r := "Art Vu v.01"
	s := res
	ls := len(r)
	hws := w + ls
	t := "[Q] Quit"
	lt := len(t)
	hwt := w - lt - 1

	// Print backgound bar + centered logo
	fmt.Fprintf(os.Stdout, escapes.CursorPos(0, headerH-2))
	fmt.Fprintf(os.Stdout, ansi.BgBlue+ansi.BrightYellow+(fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (hws)/2, r)))+ansi.Reset)

	// Print Menu commands
	fmt.Fprintf(os.Stdout, escapes.CursorPos(hwt, headerH-2))
	fmt.Fprintf(os.Stdout, ansi.BgBlue+ansi.BrightCyan+t+ansi.Reset)

	// Print current directory
	fmt.Fprintf(os.Stdout, escapes.CursorPos(1, headerH-2))
	fmt.Fprintf(os.Stdout, ansi.BgBlue+ansi.BrightWhite+showrune.ArrowDown+" "+ansi.BrightCyan+s+ansi.Reset)

}

func ShowFooter(w int, h int) {
	fmt.Fprintf(os.Stdout, escapes.CursorHide)
	fmt.Fprintf(os.Stdout, escapes.CursorPos(0, h))

	sChars := 21
	s := ansi.BrightBlack + "made with " + ansi.BrightRed + showrune.Heart + ansi.BrightBlack + " by alpha"
	fw := w - 2 + sChars

	fmt.Fprintf(os.Stdout, (fmt.Sprintf(" "+"%[1]*s", -fw, fmt.Sprintf("%[1]*s", (fw+len(s))/2, s)))+ansi.Reset)
}

func ShowSplash(w int) {

	fmt.Fprintf(os.Stdout, escapes.CursorHide)
	fmt.Fprintf(os.Stdout, escapes.ClearScreen)
	fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 0))

	if w < 132 {
		ansi.ShowArt(s80)
	} else {
		ansi.ShowArt(s132)
	}

}

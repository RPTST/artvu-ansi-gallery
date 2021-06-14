package showfile

import (
	"fmt"
	"os"

	"github.com/robbiew/artvu-ansi-gallery/internal/ansi"
	"github.com/robbiew/artvu-ansi-gallery/internal/file"
	"github.com/robbiew/artvu-ansi-gallery/pkg/showrune"
	"github.com/robbiew/artvu-ansi-gallery/pkg/textutil"
)

func Show(dirList []string, fileList []string, visibleDirIdx int, currentDir int, rootDir string, headerH int, h int, w int) (int, string, int) {
	var action int
	var selected string

	var wPad int
	// var wBtwn int
	var wFile int

	// var tpos string

	// var wTitle int

	// Set X Position for each column

	wPad = 4 // padding from left side of screen
	// wBtwn = 1 // padding between columns
	wFile = 25
	// tpos = escapes.CursorPosX(wFile + wPad + wBtwn)
	// wTitle = 40
	// wAuthor := 50
	// wGroup := 60
	// wDimen := 67
	// wDate := 73

	startList := "\033[" + fmt.Sprint(headerH+1) + ";2H"
	fmt.Fprintf(os.Stdout, startList)
	for i, v := range append(dirList, fileList...) {
		up := "\033[1A"
		xpos := "\033[" + fmt.Sprint(wPad) + "G"
		// n := selectedAbs + "/" + v
		// record := sauce.GetSauce(n)
		// t := record.Sauceinf.Title
		// a := record.Sauceinf.Author
		// g := record.Sauceinf.Group
		// d := record.Sauceinf.Date "+"
		// t1 := strconv.Itoa(int(record.Sauceinf.Tinfo1))
		// t2 := strconv.Itoa(int(record.Sauceinf.Tinfo2))
		// fmt.Fprintf(os.Stdout, "%s", a)
		// fmt.Fprintf(os.Stdout, "%s", g)
		// fmt.Fprintf(os.Stdout, "%s", d)
		// fmt.Fprintf(os.Stdout, "%s", t1)
		// fmt.Fprintf(os.Stdout, "%s\r\n", t2)

		for i >= visibleDirIdx && i < visibleDirIdx+(h-(headerH+1)) {
			switch {
			case i == currentDir:
				switch {
				case v == "./":
					action = 1 // up to root
					selected = v
					fmt.Println(xpos + ansi.BgCyan + " " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + " " + ansi.BgCyan + ansi.BrightWhite + " " + textutil.TruncateText("root", wFile-2) + " " + ansi.Reset)
				case v == "../":
					action = 2 // back one dir
					selected = "root"
					fmt.Println(xpos + ansi.BgCyan + " " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + " " + ansi.BgCyan + ansi.BrightWhite + " " + textutil.TruncateText(showrune.ArrowUp, wFile-2) + " " + ansi.Reset)
				case file.IsDirectory(v) == true && v != rootDir:
					action = 0 // enter directory
					selected = v
					fmt.Println(xpos + ansi.BgCyan + " " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + " " + ansi.BgCyan + ansi.BrightWhite + " " + showrune.ArrowRight + " " + textutil.TruncateText(v, wFile-2) + " " + ansi.Reset)
				case file.IsDirectory(v) == false:
					fmt.Println(xpos + ansi.BgCyan + " " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + " " + ansi.BgCyan + ansi.BrightWhite + " " + textutil.TruncateText(v, wFile-2) + " " + ansi.Reset)

					// fmt.Fprintf(os.Stdout, escapes.CursorPosX(wFile+wPad+wBtwn))
					// fmt.Fprintf(os.Stdout, tpos+up+"%s", t)"root"
				default:
					fmt.Println("this is a default")
				}
			case i != currentDir:
				switch {
				case v == "./":
					fmt.Println(xpos + ansi.Cyan + " " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + " " + ansi.Magenta + " " + textutil.TruncateText("root", wFile-2) + ansi.Reset)
				case v == "../":
					fmt.Println(xpos + ansi.Cyan + " " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + " " + ansi.Magenta + " " + textutil.TruncateText(showrune.ArrowUp, wFile-2) + ansi.Reset)
				case file.IsDirectory(v) == true && v != rootDir:
					fmt.Println(xpos + ansi.Cyan + " " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + " " + ansi.Cyan + " " + ansi.Green + showrune.ArrowRight + " " + ansi.Cyan + textutil.TruncateText(v, wFile-2) + ansi.Reset)
				case file.IsDirectory(v) == false:
					fmt.Fprintf(os.Stdout, xpos+ansi.Cyan+" "+textutil.PadLeft(">", " ", wFile)+"\n"+ansi.Reset)
					fmt.Fprintf(os.Stdout, xpos+up+ansi.Reset+" "+ansi.White+" "+textutil.TruncateText(v, wFile-2)+"\n"+ansi.Reset)

				}
			default:
				fmt.Println("this is a default")
			}
			break
		}
	}
	// scrollArrows(headerH)
	return currentDir, selected, action
}

func scrollArrows(headerH int) {

	// up arrow
	scrollUpRight := "\033[" + fmt.Sprint(headerH) + ";2H"
	fmt.Fprintf(os.Stdout, scrollUpRight)
	fmt.Fprintf(os.Stdout, string([]rune{'\u0018'}))

	//down arrow
	scrollDownRight := "\033[" + fmt.Sprint(headerH+1) + ";2H"
	fmt.Fprintf(os.Stdout, scrollDownRight)
	fmt.Fprintf(os.Stdout, string([]rune{'\u0019'}))

}

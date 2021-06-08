package menu

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/eiannone/keyboard"
	"github.com/robbiew/artvu-ansi-gallery/internal/ansi"
	"github.com/robbiew/artvu-ansi-gallery/internal/file"
	"github.com/robbiew/artvu-ansi-gallery/internal/theme"
	"github.com/robbiew/artvu-ansi-gallery/pkg/textutil"
	escapes "github.com/snugfox/ansi-escapes"
)

var (
	debug bool
)

type CurrentFile struct {
	CurrentDir    int
	Selected      string
	Action        int
	VisibleDirIdx int
}

type DirsFiles struct {
	FilesSlices []string
	DirSlices   []string
	CurrentPath string
}

func ShowDirs(dirList []string, fileList []string, visibleDirIdx int, currentDir int, rootDir string, headerH int, h int, w int) (int, string, int) {
	var action int
	var selected string

	var wPad int
	// var wBtwn int
	var wFile int

	// var tpos string

	// var wTitle int

	// Set X Position for each column

	wPad = 3 // padding from left side of screen
	// wBtwn = 1 // padding between columns
	wFile = 20
	// tpos = escapes.CursorPosX(wFile + wPad + wBtwn)
	// wTitle = 40
	// wAuthor := 50
	// wGroup := 60
	// wDimen := 67
	// wDate := 73

	fmt.Fprintf(os.Stdout, escapes.CursorPos(0, headerH))
	for i, v := range append(dirList, fileList...) {
		up := "\033[1A"
		xpos := "\033[" + fmt.Sprint(wPad) + "G"
		// n := selectedAbs + "/" + v
		// record := sauce.GetSauce(n)
		// t := record.Sauceinf.Title
		// a := record.Sauceinf.Author
		// g := record.Sauceinf.Group
		// d := record.Sauceinf.Date
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
				case v == rootDir:
					action = 1 // up to root
					selected = v
					fmt.Println(xpos + ansi.BgCyan + "  " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + "  " + ansi.BgCyan + ansi.BrightWhite + " " + textutil.TruncateText("^ Root", wFile-2) + " " + ansi.Reset)
				case v == "../":
					action = 2 // back one dir
					selected = v
					fmt.Println(xpos + ansi.BgCyan + "  " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + "  " + ansi.BgCyan + ansi.BrightWhite + " " + textutil.TruncateText("< Back", wFile-2) + " " + ansi.Reset)
				case file.IsDirectory(v) == true && v != rootDir:
					action = 0 // enter directory
					selected = v
					fmt.Println(xpos + ansi.BgCyan + "  " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + "  " + ansi.BgCyan + ansi.BrightWhite + " " + textutil.TruncateText("<DIR> "+v, wFile-2) + " " + ansi.Reset)
				case file.IsDirectory(v) == false:
					fmt.Println(xpos + ansi.BgCyan + "  " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + "  " + ansi.BgCyan + ansi.BrightWhite + " " + textutil.TruncateText(v, wFile-2) + " " + ansi.Reset)

					// fmt.Fprintf(os.Stdout, escapes.CursorPosX(wFile+wPad+wBtwn))
					// fmt.Fprintf(os.Stdout, tpos+up+"%s", t)

				default:
					fmt.Println("this is a default")
				}
			case i != currentDir:
				switch {
				case v == rootDir:
					fmt.Println(xpos + ansi.Cyan + "  " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + "  " + ansi.Magenta + " " + textutil.TruncateText("^ Root", wFile-2) + ansi.Reset)
				case v == "../":
					fmt.Println(xpos + ansi.Cyan + "  " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + "  " + ansi.Magenta + " " + textutil.TruncateText("< Back", wFile-2) + ansi.Reset)
				case file.IsDirectory(v) == true && v != rootDir:
					fmt.Println(xpos + ansi.Cyan + "  " + textutil.PadLeft(">", " ", wFile) + ansi.Reset)
					fmt.Println(xpos + up + ansi.Reset + "  " + ansi.Cyan + " " + textutil.TruncateText("<DIR> "+v, wFile-2) + ansi.Reset)
				case file.IsDirectory(v) == false:
					fmt.Fprintf(os.Stdout, xpos+ansi.Cyan+"  "+textutil.PadLeft(">", " ", wFile)+"\n"+ansi.Reset)
					fmt.Fprintf(os.Stdout, xpos+up+ansi.Reset+"  "+ansi.Green+" "+textutil.TruncateText(v, wFile-2)+"\n"+ansi.Reset)

				}
			default:
				fmt.Println("this is a default")
			}
			break
		}
	}
	return currentDir, selected, action
}

func MenuAction(rootDir string, h int, w int, headerH int) {

	debug = true

	fmt.Println(ansi.Clear)
	fmt.Println(ansi.Home)

	f := CurrentFile{}
	p := &f

	s := DirsFiles{}
	p1 := &s

	p.VisibleDirIdx = 0
	p.CurrentDir = 0
	p.Selected = rootDir

	p1.CurrentPath = rootDir

	p1.DirSlices, p1.FilesSlices, p1.CurrentPath = file.GetDirInfo(p.Selected, rootDir, p1.CurrentPath)
	p.CurrentDir, p.Selected, p.Action = ShowDirs(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w)

	theme.ShowHeader(w)
	theme.ShowFooter(w, h)

	if debug == true {
		debugMenu(p.CurrentDir, h, p.Selected, p.VisibleDirIdx, p.Action, p1.CurrentPath)
	}

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {

		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == keyboard.KeyEnter {

			if p.Action == 0 { // go to sub directory name
				fmt.Println(ansi.Clear)
				fmt.Println(ansi.Home)

				p.VisibleDirIdx = 0

				p1.DirSlices, p1.FilesSlices, p1.CurrentPath = file.GetDirInfo(p.Selected, rootDir, p1.CurrentPath)
				p.CurrentDir, p.Selected, p.Action = ShowDirs(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w)

				theme.ShowHeader(w)
				theme.ShowFooter(w, h)

				if debug == true {
					debugMenu(p.CurrentDir, h, p.Selected, p.VisibleDirIdx, p.Action, p1.CurrentPath)
				}

			}

			if p.Action == 1 { // return to main root dir
				fmt.Println(ansi.Clear)
				fmt.Println(ansi.Home)

				p1.DirSlices, p1.FilesSlices, p1.CurrentPath = file.GetDirInfo(rootDir, rootDir, p1.CurrentPath)
				p.CurrentDir, p.Selected, p.Action = ShowDirs(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w)

				theme.ShowHeader(w)
				theme.ShowFooter(w, h)

			}

			if p.Action == 2 { // back one dir
				fmt.Println(ansi.Clear)
				fmt.Println(ansi.Home)
				p1.DirSlices, p1.FilesSlices, p1.CurrentPath = file.GetDirInfo("../", rootDir, p1.CurrentPath)
				p.CurrentDir, p.Selected, p.Action = ShowDirs(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w)

				theme.ShowHeader(w)
				theme.ShowFooter(w, h)
			}

			if p.Action == 3 { // display selected ansi
				fmt.Println(ansi.Clear)
				fmt.Println(ansi.Home)
				// ansi.WriteAnsi(p.Selected)
			}
		}

		if key == keyboard.KeyEsc || string(char) == "q" || string(char) == "Q" {
			fmt.Println(ansi.Clear)
			fmt.Println(ansi.Home)
			fmt.Println(escapes.CursorShow)

			cmd := exec.Command("reset") //Linux only
			cmd.Stdout = os.Stdout
			fmt.Println("Thanks for using ArtVu!!")
			cmd.Run()
			os.Exit(0)
		}

		if key == keyboard.KeyArrowDown { //down arrow

			p1.DirSlices, p1.FilesSlices, p1.CurrentPath = file.GetDirInfo(".", rootDir, p1.CurrentPath)
			count := len(p1.DirSlices) + len(p1.FilesSlices)

			if p.VisibleDirIdx <= count-2 && p.CurrentDir <= count-2 {
				p.CurrentDir++
				if p.CurrentDir > p.VisibleDirIdx+(h-(headerH+2)) {
					p.VisibleDirIdx++
				}

				p.CurrentDir, p.Selected, p.Action = ShowDirs(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w)

				if debug == true {
					debugMenu(p.CurrentDir, h, p.Selected, p.VisibleDirIdx, p.Action, p1.CurrentPath)
				}
			}
		}

		if key == keyboard.KeyArrowUp { //down arrow

			p1.DirSlices, p1.FilesSlices, p1.CurrentPath = file.GetDirInfo(".", rootDir, p1.CurrentPath)
			count := len(p1.DirSlices) + len(p1.FilesSlices)

			if p.VisibleDirIdx >= 0 && p.CurrentDir > 0 && p.CurrentDir <= count {
				p.CurrentDir--
				if p.CurrentDir < p.VisibleDirIdx {
					p.VisibleDirIdx--
				}

				p.CurrentDir, p.Selected, p.Action = ShowDirs(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w)

				if debug == true {
					debugMenu(p.CurrentDir, h, p.Selected, p.VisibleDirIdx, p.Action, p1.CurrentPath)
				}
			}
		}
	}
}

func debugMenu(currentDir int, h int, selected string, visibleDirIdx int, action int, currentPath string) {

	fmt.Fprintf(os.Stdout, escapes.CursorPos(0, h))
	fmt.Fprintf(os.Stdout, escapes.TextDeleteLine)
	fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))
	fmt.Fprintf(os.Stdout, "p.Action: %v, p.CurrentDir: %v, p1.CurrentPath: %v, p.Selected: %v, p.VisibleDirIdx: %v", action, currentDir, currentPath, selected, visibleDirIdx)

}

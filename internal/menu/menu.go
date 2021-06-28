package menu

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	"github.com/eiannone/keyboard"
	"github.com/robbiew/artvu-ansi-gallery/internal/debugr"
	"github.com/robbiew/artvu-ansi-gallery/internal/file"
	"github.com/robbiew/artvu-ansi-gallery/internal/show"
	"github.com/robbiew/artvu-ansi-gallery/internal/theme"
	escapes "github.com/snugfox/ansi-escapes"
)

var (
	debug bool
)

//go:embed ansiFooter.80.ans
var f80 string

//go:embed ansiFooter.132.ans
var f132 string

//go:embed quit.80.ans
var q80 string

//go:embed quit.132.ans
var q132 string

type CurrentFile struct {
	CurrentDir    int
	Selected      string
	Action        int
	VisibleDirIdx int
	CurrentPath   string
	ViewAnsi      bool
}

type DirsFiles struct {
	FilesSlices []string
	DirSlices   []string
	Count       int
}

func MenuAction(rootDir string, h int, w int, headerH int) {

	debug = false

	var navOn bool

	fmt.Println(theme.Clear)
	fmt.Println(theme.Home)

	f := CurrentFile{}
	p := &f

	s := DirsFiles{}
	p1 := &s

	p.Action = 0
	p.VisibleDirIdx = 0
	p.CurrentDir = 0
	p.Selected = rootDir
	p.CurrentPath = rootDir
	p.ViewAnsi = true

	theme.ShowHeader(w, headerH, p.CurrentPath, rootDir)
	theme.ShowFooter(w, h, p.ViewAnsi)

	p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(p.Selected, rootDir, p.CurrentPath)
	p.CurrentDir, p.Selected, p.Action, p.ViewAnsi = show.Gallery(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w, p.CurrentPath)

	if debug == true {
		debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
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

			p.VisibleDirIdx = 0
			switch p.Action {

			case 0: // open directory

				p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(p.Selected, rootDir, p.CurrentPath)
				p.CurrentDir, p.Selected, p.Action, p.ViewAnsi = show.Gallery(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w, p.CurrentPath)

				fmt.Println(theme.Clear)
				fmt.Println(theme.Home)

				theme.ShowHeader(w, headerH, p.CurrentPath, rootDir)

				if p1.Count <= 1 {
					p.CurrentDir--
				}

				p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(".", rootDir, p.CurrentPath)
				p.CurrentDir, p.Selected, p.Action, p.ViewAnsi = show.Gallery(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w, p.CurrentPath)
				theme.ShowFooter(w, h, p.ViewAnsi)

				if debug == true {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}

			case 1: // view ansi

				fmt.Println(theme.Clear)
				fmt.Println(theme.Home)

				theme.WriteAnsi(p.CurrentPath+"/"+p.Selected, h, w, headerH, p.CurrentPath, rootDir)

				if debug == true {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}

				fmt.Fprintf(os.Stdout, "\n"+escapes.CursorPos(0, h))

				if w <= 80 {
					theme.ShowArt(f80)
				} else {
					theme.ShowArt(f132)
				}

				for {

					char, key, err := keyboard.GetKey()
					if err != nil {
						panic(err)
					}

					if string(char) == "q" || string(char) == "Q" || key == keyboard.KeyEsc {
						fmt.Fprintf(os.Stdout, "\033[2J") //clear screen
						theme.ShowHeader(w, headerH, p.CurrentPath, rootDir)
						theme.ShowFooter(w, h, p.ViewAnsi)
						break
					}

					if key == keyboard.KeyArrowUp {

					}

					if key == keyboard.KeyArrowDown {

					}

					if string(char) == "r" || string(char) == "R" {
						fmt.Println(theme.Clear)
						fmt.Println(theme.Home)

						theme.WriteAnsi(p.CurrentPath+"/"+p.Selected, h, w, headerH, p.CurrentPath, rootDir)

						if debug == true {
							debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
						}

						fmt.Fprintf(os.Stdout, "\n"+escapes.CursorPos(0, h))

						if w <= 80 {
							theme.ShowArt(f80)
						} else {
							theme.ShowArt(f132)
						}

					}
				}

				fmt.Println(theme.Clear)
				fmt.Println(theme.Home)

				theme.ShowHeader(w, headerH, p.CurrentPath, rootDir)

				p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(".", rootDir, p.CurrentPath)
				p.CurrentDir, p.Selected, p.Action, p.ViewAnsi = show.Gallery(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w, p.CurrentPath)
				theme.ShowFooter(w, h, p.ViewAnsi)

			case 2: // up one dir

				p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo("../", rootDir, p.CurrentPath)
				p.CurrentDir, p.Selected, p.Action, p.ViewAnsi = show.Gallery(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w, p.CurrentPath)

				fmt.Println(theme.Clear)
				fmt.Println(theme.Home)

				theme.ShowHeader(w, headerH, p.CurrentPath, rootDir)

				p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(".", rootDir, p.CurrentPath)
				p.CurrentDir, p.Selected, p.Action, p.ViewAnsi = show.Gallery(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w, p.CurrentPath)

				theme.ShowFooter(w, h, p.ViewAnsi)
				if debug == true {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}

			case 3: // don't anything, art is too wide

			}
		}

		if key == keyboard.KeyEsc || string(char) == "q" || string(char) == "Q" {
			fmt.Println(theme.Clear)
			fmt.Println(theme.Home)

			fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 0))

			if w <= 80 {
				theme.ShowArt(q80)
			} else {
				theme.ShowArt(q132)
			}

			char, key, err := keyboard.GetKey()
			if err != nil {
				panic(err)
			}

			if string(char) == "y" || string(char) == "Y" || key == keyboard.KeyEsc {
				fmt.Println(theme.Clear)
				fmt.Println(theme.Home)
				fmt.Println(escapes.CursorShow)

				cmd := exec.Command("reset") //Linux only
				cmd.Stdout = os.Stdout
				cmd.Run()
				os.Exit(0)

			} else {

				fmt.Println(theme.Clear)
				fmt.Println(theme.Home)

				theme.ShowHeader(w, headerH, p.CurrentPath, rootDir)

				p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(".", rootDir, p.CurrentPath)
				p.CurrentDir, p.Selected, p.Action, p.ViewAnsi = show.Gallery(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w, p.CurrentPath)

				theme.ShowFooter(w, h, p.ViewAnsi)
				if debug == true {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}

			}

		}

		if key == keyboard.KeyArrowDown { //down arrow

			if p.VisibleDirIdx <= p1.Count && p.CurrentDir <= p1.Count-2 {
				p.CurrentDir++
				if p.CurrentDir > p.VisibleDirIdx+(h-(headerH+2)) {
					p.VisibleDirIdx++
				}

				fmt.Println(theme.Clear)
				fmt.Println(theme.Home)

				theme.ShowHeader(w, headerH, p.CurrentPath, rootDir)

				p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(".", rootDir, p.CurrentPath)
				p.CurrentDir, p.Selected, p.Action, p.ViewAnsi = show.Gallery(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w, p.CurrentPath)
				theme.ShowFooter(w, h, p.ViewAnsi)
				if debug == true {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}
			}
		}

		if key == keyboard.KeyArrowUp { //down arrow

			if p.VisibleDirIdx >= 0 && p.CurrentDir > 0 && p.CurrentDir <= p1.Count {
				p.CurrentDir--
				if p.CurrentDir < p.VisibleDirIdx {
					p.VisibleDirIdx--
				}

				fmt.Println(theme.Clear)
				fmt.Println(theme.Home)

				theme.ShowHeader(w, headerH, p.CurrentPath, rootDir)

				p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(".", rootDir, p.CurrentPath)
				p.CurrentDir, p.Selected, p.Action, p.ViewAnsi = show.Gallery(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w, p.CurrentPath)

				theme.ShowFooter(w, h, p.ViewAnsi)

				if debug == true {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}
			}
		}
	}
}

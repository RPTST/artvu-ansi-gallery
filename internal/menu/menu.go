package menu

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/eiannone/keyboard"
	"github.com/robbiew/artvu-ansi-gallery/internal/ansi"
	"github.com/robbiew/artvu-ansi-gallery/internal/debugr"
	"github.com/robbiew/artvu-ansi-gallery/internal/file"
	"github.com/robbiew/artvu-ansi-gallery/internal/showfile"
	"github.com/robbiew/artvu-ansi-gallery/internal/theme"
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
	CurrentPath   string
}

type DirsFiles struct {
	FilesSlices []string
	DirSlices   []string
	Count       int
}

func MenuAction(rootDir string, h int, w int, headerH int, themeDir string) {

	debug = false

	var navOn bool

	fmt.Println(ansi.Clear)
	fmt.Println(ansi.Home)

	f := CurrentFile{}
	p := &f

	s := DirsFiles{}
	p1 := &s

	p.Action = 0
	p.VisibleDirIdx = 0
	p.CurrentDir = 0
	p.Selected = rootDir
	p.CurrentPath = rootDir

	p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(p.Selected, rootDir, p.CurrentPath)
	p.CurrentDir, p.Selected, p.Action = showfile.Show(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w)

	theme.ShowHeader(w, headerH, p.CurrentPath, rootDir, themeDir)
	theme.ShowFooter(w, h)

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
			if p.Action == 0 { // go to sub directory name
				if p.VisibleDirIdx <= p1.Count-3 && p.CurrentDir <= p1.Count-3 {

					p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(p.Selected, rootDir, p.CurrentPath)
					p.CurrentDir, p.Selected, p.Action = showfile.Show(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w)
				}
				theme.ShowHeader(w, headerH, p.CurrentPath, rootDir, themeDir)
				theme.ShowFooter(w, h)

				if debug == true {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}
			}

			if p.Action == 1 { // return to main root dir
				fmt.Println(ansi.Clear)
				fmt.Println(ansi.Home)
				if p.VisibleDirIdx <= p1.Count-3 && p.CurrentDir <= p1.Count-3 {

					p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(rootDir, rootDir, p.CurrentPath)
					p.CurrentDir, p.Selected, p.Action = showfile.Show(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w)
				}
				theme.ShowHeader(w, headerH, p.CurrentPath, rootDir, themeDir)
				theme.ShowFooter(w, h)

				if debug == true {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}
			}

			if p.Action == 2 { // back one dir
				fmt.Println(ansi.Clear)
				fmt.Println(ansi.Home)
				if p.VisibleDirIdx <= p1.Count-3 && p.CurrentDir <= p1.Count-3 {

					p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo("../", rootDir, p.CurrentPath)
					p.CurrentDir, p.Selected, p.Action = showfile.Show(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w)
				}
				theme.ShowHeader(w, headerH, p.CurrentPath, rootDir, themeDir)
				theme.ShowFooter(w, h)

				if debug == true {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}

			}

			if p.Action == 3 { // display selected ansi
				fmt.Println(ansi.Clear)
				fmt.Println(ansi.Home)

				if p.VisibleDirIdx <= p1.Count-3 && p.CurrentDir <= p1.Count-3 {
					// ansi.WriteAnsi(p.Selected)
				}

				if debug == true {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}

			} else {
				fmt.Println(ansi.Clear)
				fmt.Println(ansi.Home)

				if p.VisibleDirIdx <= p1.Count-3 && p.CurrentDir <= p1.Count-3 {

					p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(".", rootDir, p.CurrentPath)
					p.CurrentDir, p.Selected, p.Action = showfile.Show(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w)

				}
				theme.ShowHeader(w, headerH, p.CurrentPath, rootDir, themeDir)
				theme.ShowFooter(w, h)

				if debug == true {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}
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

			if p.VisibleDirIdx <= p1.Count-3 && p.CurrentDir <= p1.Count-2 {
				p.CurrentDir++
				if p.CurrentDir > p.VisibleDirIdx+(h-(headerH+2)) {
					p.VisibleDirIdx++
				}

				p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(".", rootDir, p.CurrentPath)
				p.CurrentDir, p.Selected, p.Action = showfile.Show(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w)

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

				p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(".", rootDir, p.CurrentPath)
				p.CurrentDir, p.Selected, p.Action = showfile.Show(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w)

				if debug == true {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}
			}
		}
	}
}

package menu

import (
	"bufio"
	"bytes"
	_ "embed"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/robbiew/artvu-ansi-gallery/internal/debugr"
	"github.com/robbiew/artvu-ansi-gallery/internal/file"
	"github.com/robbiew/artvu-ansi-gallery/internal/show"
	"github.com/robbiew/artvu-ansi-gallery/internal/theme"
	escapes "github.com/snugfox/ansi-escapes"
)

var (
	debug    bool
	rowCount int
	currLoc  int
)

//go:embed ansiFooter.80.ans
var f80 string

//go:embed ansiFooter.132.ans
var f132 string

//go:embed quit.80.ans
var q80 string

//go:embed quit.132.ans
var q132 string

//go:embed logoff.80.ans
var l80 string

//go:embed logoff.132.ans
var l132 string

//go:embed block.ans
var block string

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

	if debug {
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

		}

		if key == keyboard.KeyEnter {

			p.VisibleDirIdx = 0
			switch p.Action {

			case 0: // open directory

				p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(p.Selected, rootDir, p.CurrentPath)
				p.CurrentDir, p.Selected, p.Action, p.ViewAnsi = show.Gallery(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w, p.CurrentPath)

				fmt.Println(theme.Clear)
				fmt.Println(theme.Home)

				if p1.Count <= 1 {
					p.CurrentDir = 1
					p.VisibleDirIdx = 0
				} else {
					p.CurrentDir = 1
					p.VisibleDirIdx = 0

				}

				theme.ShowHeader(w, headerH, p.CurrentPath, rootDir)
				theme.ShowFooter(w, h, p.ViewAnsi)

				p1.DirSlices, p1.FilesSlices, p.CurrentPath, p1.Count = file.GetDirInfo(".", rootDir, p.CurrentPath)
				p.CurrentDir, p.Selected, p.Action, p.ViewAnsi = show.Gallery(p1.DirSlices, p1.FilesSlices, p.VisibleDirIdx, p.CurrentDir, rootDir, headerH, h, w, p.CurrentPath)

				if debug {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}

			case 1: // view ansi

				fmt.Println(theme.Clear)
				fmt.Println(theme.Home)

				rowCount = 0

				WriteAnsi(p.CurrentPath+"/"+p.Selected, h, w)
				fmt.Println(" ")

				if debug {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}

				fmt.Fprintf(os.Stdout, "\n"+escapes.CursorPos(0, h))

				if w <= 80 {
					theme.ShowArt(f80)
				} else {
					theme.ShowArt(f132)
				}

				fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))
				fmt.Fprintf(os.Stdout, theme.BgCyan+theme.BrightYellow+"[          ]"+theme.Reset)

				fmt.Fprintf(os.Stdout, escapes.CursorPos(4, h))
				theme.ShowArt(block)
				theme.ShowArt(block)
				theme.ShowArt(block)
				theme.ShowArt(block)
				theme.ShowArt(block)
				theme.ShowArt(block)
				theme.ShowArt(block)
				theme.ShowArt(block)
				theme.ShowArt(block)
				theme.ShowArt(block)

				fmt.Fprintf(os.Stdout, theme.Reset)

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

					if string(char) == "r" || string(char) == "R" {
						fmt.Println(theme.Clear)
						fmt.Println(theme.Home)

						WriteAnsi(p.CurrentPath+"/"+p.Selected, h, w)

						if debug {
							debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
						}

						fmt.Fprintf(os.Stdout, "\n"+escapes.CursorPos(0, h))

						if w <= 80 {
							theme.ShowArt(f80)
						} else {
							theme.ShowArt(f132)
						}

					}

					if key == keyboard.KeyArrowUp {
						ScrollAnsi(p.CurrentPath+"/"+p.Selected, h, w, "up")
						fmt.Fprintf(os.Stdout, "\n"+escapes.CursorPos(0, h))

						if w <= 80 {
							theme.ShowArt(f80)
						} else {
							theme.ShowArt(f132)
						}

						StatusBar(h)
					}

					if key == keyboard.KeyArrowDown {
						ScrollAnsi(p.CurrentPath+"/"+p.Selected, h, w, "down")
						fmt.Fprintf(os.Stdout, "\n"+escapes.CursorPos(0, h))

						if w <= 80 {
							theme.ShowArt(f80)
						} else {
							theme.ShowArt(f132)
						}
						StatusBar(h)
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
				if debug {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}

			case 3: // don't anything, art is too wide

			}
		}

		if string(char) == "q" || string(char) == "Q" || key == keyboard.KeyEsc {
			fmt.Println(theme.Clear)
			fmt.Println(theme.Home)

			// fmt.Println(escapes.CursorPos(0, 0))

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
				if w <= 80 {
					theme.ShowArt(l80)
				} else {
					theme.ShowArt(l132)
				}

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
				if debug {
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
				if debug {
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

				if debug {
					debugr.DebugInf(p.CurrentDir, h, p.Selected, p.Action, p.CurrentPath, p1.Count, navOn)
				}
			}
		}
	}
}

// WriteAnsi(string) dislays a full CP437 ansi art file
func WriteAnsi(selected string, h int, w int) {

	file, err := os.Open(selected)
	if err != nil {
		//handle error
		log.Fatal(err)
	}

	defer file.Close()
	s := bufio.NewScanner(file)

	fmt.Fprintf(os.Stdout, escapes.EraseScreen)
	fmt.Println(escapes.CursorPos(0, 0))

	for s.Scan() {
		read_line := s.Text()
		// trim the text if it's after a SAUCE RECORD
		trimmed := theme.TrimStringFromSauce(read_line)
		var b bytes.Buffer
		for {
			// add delay between each line to throttle speed
			fmt.Println(escapes.CursorPos(0, rowCount))
			time.Sleep(time.Duration(30) * time.Millisecond)
			// fmt.Fprintf(os.Stdout, escapes.CursorNextLine)
			b.Write([]byte(trimmed + "\r"))
			b.WriteTo(os.Stdout)
			rowCount++
			break
		}
	}
	currLoc = rowCount
	return
}

// WriteAnsi(string) dislays a full CP437 ansi art file
func ScrollAnsi(selected string, h int, w int, scroll string) {

	row := 0

	if scroll == "up" {
		if currLoc > 0+h {
			currLoc--

		}
	}

	if scroll == "down" {
		if currLoc < rowCount {
			currLoc++
		}

	}

	fmt.Fprintf(os.Stdout, escapes.EraseScreen)
	fmt.Println(escapes.CursorPos(0, 0))

	f, _ := os.Open(selected)
	// Create new Scanner
	scanner := bufio.NewScanner(f)

	// Use Scan
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := theme.TrimStringFromSauce(line)
		row++
		if row <= currLoc && row >= currLoc-h {
			fmt.Println(trimmed)

		}
	}
}

func StatusBar(h int) {

	fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))

	incr := float64(currLoc) / float64(rowCount)
	r := math.Floor(incr*100) / 100
	fmt.Fprintf(os.Stdout, theme.BgCyan+theme.BrightYellow+"[          ]"+theme.Reset)

	switch {
	case r <= 1 && r > .9:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(4, h))
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .95 && r > .89:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(4, h))
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .9 && r > .8:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(4, h))
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .8 && r > .7:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(4, h))
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .7 && r > .6:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(4, h))
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .6 && r > .5:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(4, h))
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .5 && r > .4:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(4, h))
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .4 && r > .3:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(4, h))
		theme.ShowArt(block)
		theme.ShowArt(block)
		theme.ShowArt(block)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .3 && r > .2:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(4, h))
		theme.ShowArt(block)
		theme.ShowArt(block)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .2 && r > .1:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(4, h))
		theme.ShowArt(block)
		fmt.Fprintf(os.Stdout, theme.Reset)
	}
}

// fmt.Fprintf(os.Stdout, theme.BgCyan+theme.BrightYellow+"%v "+theme.BrightCyan+"/ "+theme.BrightYellow+"%v"+theme.Reset, currLoc, rowCount)

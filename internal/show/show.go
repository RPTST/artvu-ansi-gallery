package show

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/robbiew/artvu-ansi-gallery/internal/file"
	"github.com/robbiew/artvu-ansi-gallery/internal/theme"
	"github.com/robbiew/artvu-ansi-gallery/pkg/sauce"
	"github.com/robbiew/artvu-ansi-gallery/pkg/showrune"
	"github.com/robbiew/artvu-ansi-gallery/pkg/textutil"
)

type SauceData struct {
	Title  string
	Author string
	Group  string
	Date   string
	Tinfo1 string
	Tinfo2 string
}

func Gallery(dirList []string, fileList []string, visibleDirIdx int, currentDir int, rootDir string, headerH int, h int, w int, currentPath string) (int, string, int, bool) {

	var action int      // directory, file or up
	var selected string // file name of active item
	var viewAnsi bool

	f := SauceData{}
	s := &f

	wName := 37   // width of Name field
	wAuthor := 26 // width of Author field
	wType := 8    // width of Type field
	wSize := 8    // width of Tinfo (size) fields

	authorX := 37 // x position of Author field
	typeX := 62   // x position of Type field
	sizeX := 70   // x position of Tinfo (size) fields

	up := "\033[1A"                             // ansi code to move cusror up one row
	apos := "\033[" + fmt.Sprint(authorX) + "C" // Author cursor pos
	tpos := "\033[" + fmt.Sprint(typeX) + "C"   // Type cursor pos
	spos := "\033[" + fmt.Sprint(sizeX) + "C"   // Tinfo (size) cursor pos

	startList := "\033[" + fmt.Sprint(headerH+1) + ";H" // Cursor position to start display
	fmt.Fprintf(os.Stdout, startList)

	for i, v := range append(dirList, fileList...) {
		for i >= visibleDirIdx && i < visibleDirIdx+(h-(headerH+1)) {
			fmt.Fprintf(os.Stdout, theme.Reset)
			switch {
			case i == currentDir:
				switch {
				case v == "../": // back one dir

					viewAnsi = true
					action = 2
					selected = "root"

					fmt.Println(theme.BgCyan + textutil.PadLeft(">", " ", wName))
					fmt.Println(up + theme.BgCyan + theme.BrightWhite + " " + textutil.TruncateText(showrune.ArrowUp, wName-2))

					fmt.Println(up + tpos + theme.BgCyan + textutil.PadLeft(">", " ", wType))
					fmt.Println(up + tpos + theme.White + " " + textutil.TruncateText("Folder", wType))

					fmt.Println(up + apos + theme.BgCyan + textutil.PadLeft(">", " ", wAuthor))
					fmt.Println(up + apos + theme.White + " " + textutil.TruncateText("-", wAuthor-2))

					fmt.Println(up + spos + theme.BgCyan + textutil.PadLeft(">", " ", wSize+1))
					fmt.Println(up + spos + theme.White + " " + textutil.TruncateText("-", wSize))

				case file.IsDirectory(v) && v != rootDir: // we've entered a directory

					viewAnsi = true
					selected = v
					action = 0

					fmt.Println(theme.BgCyan + textutil.PadLeft(">", " ", wName))
					fmt.Println(up + theme.BgCyan + theme.BrightWhite + " " + showrune.ArrowRight + " " + textutil.TruncateText(v, wName-2))

					fmt.Println(up + tpos + theme.BgCyan + textutil.PadLeft(">", " ", wType))
					fmt.Println(up + tpos + theme.White + " " + textutil.TruncateText("Folder", wType))

					fmt.Println(up + apos + theme.BgCyan + textutil.PadLeft(">", " ", wAuthor))
					fmt.Println(up + apos + theme.White + " " + textutil.TruncateText("-", wAuthor-2))

					fmt.Println(up + spos + theme.BgCyan + textutil.PadLeft(">", " ", wSize+1))
					fmt.Println(up + spos + theme.White + " " + textutil.TruncateText("-", wSize))

				case !file.IsDirectory(v) && v != "../":

					selected = v
					action = 1 // viewing ansi file

					fmt.Println(theme.BgCyan + textutil.PadLeft(">", " ", wName))
					fmt.Println(up + theme.BrightWhite + " " + textutil.TruncateText(v, wName-2))

					fmt.Println(up + tpos + theme.BgCyan + textutil.PadLeft(">", " ", wType))
					fmt.Println(up + tpos + theme.White + " " + textutil.TruncateText("File", wType))

					fmt.Println(up + apos + theme.BgCyan + textutil.PadLeft(">", " ", wAuthor))
					fmt.Println(up + apos + theme.White + " " + textutil.TruncateText("-", wAuthor-2))

					fmt.Println(up + spos + theme.BgCyan + textutil.PadLeft(">", " ", wSize+1))
					fmt.Println(up + spos + theme.White + " " + textutil.TruncateText("-", wSize))

					// SUACE records
					n := currentPath + "/" + v
					if sauce.CheckSauce(n) {

						record := sauce.GetSauce(n)
						s.Author = strings.TrimSpace(string(fmt.Sprintf("%s", record.Sauceinf.Author)[:]))
						s.Tinfo1 = strings.TrimSpace(string(fmt.Sprintf("%s", strconv.Itoa(int(record.Sauceinf.Tinfo1))[:])))
						s.Tinfo2 = strings.TrimSpace(string(fmt.Sprintf("%s", strconv.Itoa(int(record.Sauceinf.Tinfo2))[:])))

						sInt, err := strconv.Atoi(s.Tinfo1)
						if err != nil {
							log.Fatal(err)
						}

						var widthColor string
						if sInt > w {
							widthColor = theme.BrightRed
							viewAnsi = false
							action = 3
						} else {
							widthColor = theme.White
							viewAnsi = true
						}

						fmt.Println(up + apos + theme.BgCyan + textutil.PadLeft(">", " ", wAuthor))
						fmt.Println(up + apos + theme.White + " " + textutil.TruncateText(s.Author, wAuthor-2))

						if len(strings.TrimSpace(s.Tinfo1)) > 1 {

							fmt.Println(up + spos + theme.BgCyan + textutil.PadLeft(">", " ", wSize+1))
							fmt.Println(up + spos + widthColor + " " + textutil.TruncateText(s.Tinfo1+"x"+s.Tinfo2, wSize))
						}

					} else {

						viewAnsi = true

						fmt.Println(up + apos + textutil.PadLeft(">", " ", wAuthor))
						fmt.Println(up + apos + theme.White + " -")

						fmt.Println(up + spos + textutil.PadLeft(">", " ", wSize+1))
						fmt.Println(up + spos + " -")

						fmt.Println(up + spos + textutil.PadLeft(">", " ", wSize+1))
						fmt.Println(up + spos + " -")

					}
					fmt.Fprintf(os.Stdout, theme.Reset)

				}

			case i != currentDir:

				switch {
				case v == "../":

					fmt.Println(theme.Cyan + textutil.PadLeft(">", " ", wName))
					fmt.Println(up + theme.Magenta + " " + textutil.TruncateText(showrune.ArrowUp, wName-2))

					fmt.Println(up + tpos + textutil.PadLeft(">", " ", wType))
					fmt.Println(up + tpos + theme.White + " " + textutil.TruncateText("Folder", 7))

					fmt.Println(up + apos + textutil.PadLeft(">", " ", wAuthor))
					fmt.Println(up + apos + theme.White + " " + textutil.TruncateText("-", wAuthor-2))

					fmt.Println(up + spos + textutil.PadLeft(">", " ", wSize+1))
					fmt.Println(up + spos + theme.White + " " + textutil.TruncateText("-", wSize))

				case file.IsDirectory(v) && v != rootDir:

					fmt.Println(theme.Cyan + textutil.PadLeft(">", " ", wName))
					fmt.Println(up + theme.Cyan + " " + theme.Green + showrune.ArrowRight + " " + theme.Cyan + textutil.TruncateText(v, wName-2))

					fmt.Println(up + tpos + textutil.PadLeft(">", " ", wType))
					fmt.Println(up + tpos + theme.White + " " + textutil.TruncateText("Folder", 7))

					fmt.Println(up + apos + textutil.PadLeft(">", " ", wAuthor))
					fmt.Println(up + apos + theme.White + " " + textutil.TruncateText("-", wAuthor-2))

					fmt.Println(up + spos + textutil.PadLeft(">", " ", wSize+1))
					fmt.Println(up + spos + theme.White + " " + textutil.TruncateText("-", wSize))

				case !file.IsDirectory(v):

					fmt.Println(theme.Cyan + textutil.PadLeft(">", " ", wName))
					fmt.Println(up + theme.White + " " + textutil.TruncateText(v, wName-2))

					fmt.Println(up + tpos + textutil.PadLeft(">", " ", wType))
					fmt.Println(up + tpos + theme.White + " " + textutil.TruncateText("File", 7))

					// SUACE records

					n := currentPath + "/" + v
					if sauce.CheckSauce(n) {
						record := sauce.GetSauce(n)
						s.Author = strings.TrimSpace(string(fmt.Sprintf("%s", record.Sauceinf.Author)[:]))
						s.Tinfo1 = strings.TrimSpace(string(strconv.Itoa(int(record.Sauceinf.Tinfo1))[:]))
						s.Tinfo2 = strings.TrimSpace(string(strconv.Itoa(int(record.Sauceinf.Tinfo2))[:]))

						sInt, err := strconv.Atoi(s.Tinfo1)
						if err != nil {
							log.Fatal(err)
						}
						var widthColor string
						if sInt > w {
							widthColor = theme.BrightRed
						} else {
							widthColor = theme.White
						}

						fmt.Println(up + apos + textutil.PadLeft(">", " ", wAuthor))
						fmt.Println(up + apos + theme.White + " " + textutil.TruncateText(s.Author, wAuthor-2))

						if len(strings.TrimSpace(s.Tinfo1)) < 2 {
							fmt.Println(up + spos + textutil.PadLeft(">", " ", wSize+1))
							fmt.Println(up + spos + widthColor + " -")
						} else {
							fmt.Println(up + spos + textutil.PadLeft(">", " ", wSize+1))
							fmt.Println(up + spos + widthColor + " " + textutil.TruncateText(s.Tinfo1+"x"+s.Tinfo2, wSize))
						}
					} else {

						fmt.Println(up + apos + textutil.PadLeft(">", " ", wAuthor))
						fmt.Println(up + apos + theme.White + " -")

						fmt.Println(up + spos + textutil.PadLeft(">", " ", wSize+1))
						fmt.Println(up + spos + " -")

						fmt.Println(up + spos + textutil.PadLeft(">", " ", wSize+1))
						fmt.Println(up + spos + " -")

					}
				}
				fmt.Fprintf(os.Stdout, theme.Reset)
			}
			break
		}
	}

	fmt.Fprintf(os.Stdout, theme.Reset)

	return currentDir, selected, action, viewAnsi
}

func hasExts(path string, exts []string) bool {
	pathExt := strings.ToLower(filepath.Ext(path))
	for _, ext := range exts {
		if pathExt == strings.ToLower(ext) {
			return true
		}
	}
	return false
}

func hasDiz(path string) bool {

	var sourceExts []string

	file := path

	dir := file
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	sourceExts = []string{".diz", ".DIZ", "Diz"}
	for _, fi := range files {
		if hasExts(fi.Name(), sourceExts) {
			return true
		}
	}
	return false

}

func showDiz(root string, dirName string, xloc int, headerH int, h int) {

	var diz string
	var sourceExts []string

	file := root + "/" + dirName

	dir := file
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	sourceExts = []string{".diz", ".DIZ", "Diz"}
	for _, fi := range files {
		if hasExts(fi.Name(), sourceExts) {
			diz = fi.Name()
			// fmt.Println(fi.Name())
		}
	}

	dizFile, err := ioutil.ReadFile(root + "/" + dirName + "/" + diz)

	if err != nil {
		log.Fatal(err)
	}
	noSauce := theme.TrimStringFromSauce(string(dizFile))
	var r io.Reader = strings.NewReader(noSauce)

	if err != nil {
		log.Fatal(err)
	}
	rowCount := headerH
	s := bufio.NewScanner(r)
	for s.Scan() {
		read_line := s.Text()
		loc := "\033[" + fmt.Sprint(rowCount) + ";" + fmt.Sprint(xloc) + "f"
		for rowCount <= 25 {
			fmt.Fprintf(os.Stdout, loc+read_line+"\r\n")
			rowCount++
			break
		}
	}

}

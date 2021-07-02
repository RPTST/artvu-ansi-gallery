package theme

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/robbiew/artvu-ansi-gallery/pkg/sauce"
	"github.com/robbiew/artvu-ansi-gallery/pkg/showrune"
	escapes "github.com/snugfox/ansi-escapes"
)

//go:embed splash.80.ans
var s80 string

//go:embed splash.132.ans
var s132 string

//go:embed header.80.ans
var h80 string

//go:embed header.132.ans
var h132 string

var (
	Block = " "
)

const (
	// Cursor codes
	Clear = "\033[2J"
	Home  = "\033[H"

	// Ansi color codes
	Reset = "\u001b[0m"

	Black   = "\u001b[30m"
	Red     = "\u001b[31m"
	Green   = "\u001b[32m"
	Yellow  = "\u001b[33m"
	Blue    = "\u001b[34m"
	Magenta = "\u001b[35m"
	Cyan    = "\u001b[36m"
	White   = "\u001b[37m"

	BrightBlack   = "\u001b[30;1m"
	BrightRed     = "\u001b[31;1m"
	BrightGreen   = "\u001b[32;1m"
	BrightYellow  = "\u001b[33;1m"
	BrightBlue    = "\u001b[34;1m"
	BrightMagenta = "\u001b[35;1m"
	BrightCyan    = "\u001b[36;1m"
	BrightWhite   = "\u001b[37;1m"

	BgBlack   = "\u001b[40m"
	BgRed     = "\u001b[41m"
	BgGreen   = "\u001b[42m"
	BgYellow  = "\u001b[43m"
	BgBlue    = "\u001b[44m"
	BgMagenta = "\u001b[45m"
	BgCyan    = "\u001b[46m"
	BgWhite   = "\u001b[47m"
)

func CheckSauce(filename string) bool {

	// let's check the file for a valid SAUCE record
	record := sauce.GetSauce(filename)

	// if we find a SAUCE record, return bool flag
	if string(record.Sauceinf.ID[:]) == sauce.SauceID {
		return true
	} else {
		return false
	}
}

func GetSauce(filename string) ([35]byte, [20]byte, [20]byte, [8]byte, uint16, uint16) {

	var tinfo1 uint16
	var tinfo2 uint16

	record := sauce.GetSauce(filename)

	title := record.Sauceinf.Title
	author := record.Sauceinf.Author
	group := record.Sauceinf.Group
	date := record.Sauceinf.Date

	if record.Sauceinf.Tinfo1 != 0 {
		tinfo1 = record.Sauceinf.Tinfo1
	}
	if record.Sauceinf.Tinfo2 != 0 {
		tinfo2 = record.Sauceinf.Tinfo2
	}
	return title, author, group, date, tinfo1, tinfo2
}

func ShowArt(text string) {

	// trim the text if it's after a SAUCE RECORD

	noSauce := TrimStringFromSauce(string(text))
	var r io.Reader = strings.NewReader(noSauce)

	io.Copy(os.Stdout, r)
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func TrimStringFromSauce(s string) string {

	if idx := strings.Index(s, "COMNT"); idx != -1 {
		string := s
		delimiter := "COMNT"
		leftOfDelimiter := strings.Split(string, delimiter)[0]
		trim := TrimLastChar(leftOfDelimiter)
		return trim
		// rightOfDelimiter := strings.Join(strings.Split(string, delimiter)[1:], delimiter)
	}
	if idx := strings.Index(s, "SAUCE00"); idx != -1 {
		string := s
		delimiter := "SAUCE00"
		leftOfDelimiter := strings.Split(string, delimiter)[0]
		trim := TrimLastChar(leftOfDelimiter)
		return trim
		// rightOfDelimiter := strings.Join(strings.Split(string, delimiter)[1:], delimiter)
	}

	return s
}

func ShowHeader(w int, headerH int, path string, rootDir string) {

	fmt.Fprintf(os.Stdout, escapes.CursorHide)
	fmt.Fprintf(os.Stdout, escapes.ClearScreen)
	fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 0))

	if w < 132 {
		ShowArt(h80)
	} else {
		ShowArt(h132)
	}
	// Dir item label
	var res string
	str := path
	res = strings.Join(strings.Split(str, rootDir+"/"), "")
	if res == rootDir {
		res = "root"

	} else {
		res = res + "/"
	}

	s := res
	t := "[Q] Quit"
	lt := len(t)
	hwt := w - lt - 2

	// Print Menu commands
	fmt.Fprintf(os.Stdout, escapes.CursorPos(hwt, 1))
	fmt.Fprintf(os.Stdout, BgCyan+BrightCyan+t+Reset)

	// Print current directory
	fmt.Fprintf(os.Stdout, escapes.CursorPos(12, headerH-4))
	fmt.Fprintf(os.Stdout, BrightWhite+showrune.ArrowDown+" "+BrightCyan+s+Reset)

}

func ShowFooter(w int, h int, v bool) {
	var s string
	var sColor string
	fmt.Fprintf(os.Stdout, escapes.CursorHide)

	if v == false {
		s = "Art is too wide for your screen mode..."
		sColor = Red
	} else {
		s = "ArtVu v.01 by Alpha"
		sColor = BrightBlack
	}
	fmt.Fprintf(os.Stdout, escapes.CursorPos(0, h))
	fmt.Fprintf(os.Stdout, sColor+(fmt.Sprintf("%*s", -w+1, fmt.Sprintf("%*s", (w+1+len(s))/2, s)))+Reset)
}

func ShowSplash(w int) {

	fmt.Fprintf(os.Stdout, escapes.CursorHide)
	fmt.Fprintf(os.Stdout, escapes.ClearScreen)
	fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 0))

	if w < 132 {
		ShowArt(s80)
	} else {
		ShowArt(s132)
	}

}

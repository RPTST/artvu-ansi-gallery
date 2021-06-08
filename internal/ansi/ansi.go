package ansi

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/eiannone/keyboard"
	"github.com/robbiew/artvu-ansi-gallery/pkg/sauce"
	escapes "github.com/snugfox/ansi-escapes"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
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

// WriteAnsi(string) dislays a full CP437 ansi art file
func WriteAnsi(selected string) {
	content, err := ioutil.ReadFile(selected)
	if err != nil {
		log.Fatal(err)
	}
	noSauce := TrimStringFromSauce(string(content))
	var r io.Reader = strings.NewReader(noSauce)
	r = transform.NewReader(r, charmap.CodePage437.NewDecoder())

	var ready string
	if b, err := io.ReadAll(r); err == nil {
		ready = string(b)
	}
	rowCount := 0
	scanner := bufio.NewScanner(strings.NewReader(ready))
	for scanner.Scan() {
		read_line := scanner.Text()
		// trim the text if it's after a SAUCE RECORD
		trimmed := TrimStringFromSauce(read_line)
		for {
			fmt.Println(escapes.CursorPos(40, rowCount))
			time.Sleep(time.Duration(40) * time.Millisecond)
			fmt.Fprintf(os.Stdout, trimmed)
			rowCount++
			break
		}
	}
	char, key, err := keyboard.GetKey()
	if err != nil {
		panic(err)
	}

	if string(char) == "q" || string(char) == "Q" || key == keyboard.KeyEsc {
		fmt.Fprintf(os.Stdout, "\033[2J") //clear screen
		// Top(currentDirPath)
		// ShowDirs(subDirs, subFiles)
		// Bottom(selected, h)
	}
}

func Top(dp string, abs string) {
	fmt.Fprintf(os.Stdout, escapes.CursorPos(0, 0))
	fmt.Fprintf(os.Stdout, BrightWhite+"ArtVu Local Ansi Browser\n"+Reset)
	fmt.Fprintf(os.Stdout, Yellow+"Location: %v                             "+Reset, abs)
}

func Bottom(selected string, action int, h int, currentDir int) {
	fmt.Fprintf(os.Stdout, escapes.CursorPos(0, h))
	fmt.Fprintf(os.Stdout, "Footer - Selected: %v, Action: %v, CurrentDir: %v", selected, action, currentDir)
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

func TrimStringFromSauce(s string) string {
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

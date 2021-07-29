package art

import (
	"bufio"
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"runtime"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/robbiew/artvu-ansi-gallery/internal/ansi"
	"github.com/robbiew/artvu-ansi-gallery/internal/theme"
	escapes "github.com/snugfox/ansi-escapes"
)

var (

	//go:embed bar.ans
	bar string

	//go:embed bar10.ans
	bar10 string

	//go:embed bar20.ans
	bar20 string

	//go:embed bar30.ans
	bar30 string

	//go:embed bar40.ans
	bar40 string

	//go:embed bar50.ans
	bar50 string

	//go:embed bar60.ans
	bar60 string

	//go:embed bar70.ans
	bar70 string

	//go:embed bar80.ans
	bar80 string

	//go:embed bar90.ans
	bar90 string

	//go:embed bar100.ans
	bar100 string
)

const (
	CHR_ESCAPE = 0x1B
	CHR_CR     = 0x0D
	CHR_LF     = 0x0A
)

func Ansiart2utf8(file string, artw int, w int, h int) {
	var (
		oErr  error
		pFile *os.File
	)

	// ERROR LOGGING
	pLogErr := log.New(os.Stderr, "", log.Lshortfile)
	fnErrExit := func(oErr error) {

		if oErr != nil {

			pLogErr.Output(2, oErr.Error())
			os.Exit(1)
		}
	}

	// TRANSLATION ARRAY
	Array437 := [256]rune{
		'\x00', '☺', '☻', '♥', '♦', '♣', '♠', '•', '\b', '\t', '\n', '♂', '♀', '\r', '♫', '☼',
		'►', '◄', '↕', '‼', '¶', '§', '▬', '↨', '↑', '↓', '→', '\x1b', '∟', '↔', '▲', '▼',
		' ', '!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?',
		'@', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
		'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '[', '\\', ']', '^', '_',
		'`', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
		'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', '{', '|', '}', '~', '⌂',
		'\u0080', '\u0081', 'é', 'â', 'ä', 'à', 'å', 'ç', 'ê', 'ë', 'è', 'ï', 'î', 'ì', 'Ä', 'Å',
		'É', 'æ', 'Æ', 'ô', 'ö', 'ò', 'û', 'ù', 'ÿ', 'Ö', 'Ü', '¢', '£', '¥', '₧', 'ƒ',
		'á', 'í', 'ó', 'ú', 'ñ', 'Ñ', 'ª', 'º', '¿', '⌐', '¬', '½', '¼', '¡', '«', '»',
		'░', '▒', '▓', '│', '┤', '╡', '╢', '╖', '╕', '╣', '║', '╗', '╝', '╜', '╛', '┐',
		'└', '┴', '┬', '├', '─', '┼', '╞', '╟', '╚', '╔', '╩', '╦', '╠', '═', '╬', '╧',
		'╨', '╤', '╥', '╙', '╘', '╒', '╓', '╫', '╪', '┘', '┌', '█', '▄', '▌', '▐', '▀',
		'α', 'ß', 'Γ', 'π', 'Σ', 'σ', 'µ', 'τ', 'Φ', 'Θ', 'Ω', 'δ', '∞', 'φ', 'ε', '∩',
		'≡', '±', '≥', '≤', '⌠', '⌡', '÷', '≈', '°', '∙', '·', '√', 'ⁿ', '²', '■', '\u00a0',
	}

	runtime.GOMAXPROCS(1)
	pFile = nil

	// COMMAND PARAMETERS
	puiWidth := artw
	pszInput := file
	pbDebug := false
	pnRowBytes := 0
	// termWidth := w
	// termHeight := h

	// DEBUG LOGGING
	pLogDebug := log.New(os.Stdout, "", 0)
	fnDebug := func(v ...interface{}) {

		if pbDebug {
			pLogDebug.Output(2, fmt.Sprint(v...))
		}
	}

	// GET FILE HANDLE
	if strings.Compare(pszInput, "-") == 0 {

		pFile = os.Stdin

	} else {

		pFile, oErr = os.Open(pszInput)
		fnErrExit(oErr)

		// fnDebug("FILE: ", pszInput)
	}

	bsInput, oErr := ioutil.ReadAll(pFile)
	fnErrExit(oErr)

	bEsc := false

	bsSGR := ansi.SGR{}
	bsSGR.Reset()
	pGrid := ansi.GridNew(ansi.GridDim(puiWidth))

	// BUFFER OUTPUT
	pWriter := bufio.NewWriter(os.Stdout)
	curCode := ansi.ECode{}
	curPos := ansi.NewPos()
	curSaved := ansi.NewPos()

	trimmed := theme.TrimStringFromSauce(string(bsInput))

	b := []byte(trimmed)

	// ITERATE BYTES IN INPUT
	for _, chr := range b {

		// DROP \r
		if chr == CHR_CR {

			continue

			// BEGIN ESCAPE CODE
		} else if chr == CHR_ESCAPE {

			bEsc = true
			curCode.Reset()

			// HANDLE ESCAPE CODE SEQUENCE
		} else if bEsc {

			// NOPS

			/*
				UNHANDLED CODE:   ESC[J; [1]
				UNHANDLED CODE:   ESC[K; [1]
				UNHANDLED CODE:   ESCc;
				INVALID CODE:     ESC[MF;
				INVALID CODE:     ESC[m;
				INVALID CODE:     ESC[P;
				INVALID CODE:     ESC[T;
				INVALID CODE:     ESC[S;
				UNHANDLED CODE:   ESC[@K; [1]
				INVALID CODE:     ESC[@l;
				INVALID CODE:     ESC[@S;
				INVALID CODE:     ESC[@N;
				INVALID CODE:     ESC[@u;
				INVALID CODE:     ESC[@s;
				INVALID CODE:     ESC[Mo3egc;
			*/
			// TODO: ESC[?7h; - possibly "wrap" mode

			// ESCAPE CODE TERMINATING CHARS:
			// EXIT ESCAPE CODE FSM SUCCESSFULLY ON TERMINATING 'm' CHARACTER
			if strings.IndexByte(ansi.CodeTerminators(), chr) != -1 {

				bEsc = false
				curCode.Code = rune(chr)

				if curCode.Validate() {

					// ONLY RESTORE SGR ESCAPE CODES
					switch curCode.Code {

					case 'm':

						oErr = bsSGR.Merge(curCode.SubParams)

						if oErr != nil {
							fnDebug(oErr)
						}

					// UP
					case 'A':

						pGrid.IncClamp(&curPos, 0, -int(curCode.SubParams[0]))

					// DOWN
					case 'B':

						pGrid.IncClamp(&curPos, 0, int(curCode.SubParams[0]))

					// FORWARD
					case 'C':

						pGrid.IncClamp(&curPos, int(curCode.SubParams[0]), 0)

					// BACK
					case 'D':

						pGrid.IncClamp(&curPos, -int(curCode.SubParams[0]), 0)

					// TO X,Y
					case 'H', 'f':

						curPos.Y = ansi.GridDim(curCode.SubParams[0])
						curPos.X = ansi.GridDim(curCode.SubParams[1])

					// SAVE CURSOR POS
					case 's':

						curSaved = curPos

					// RESTORE CURSOR POS
					case 'u':

						curPos = curSaved

					default:

						fnDebug("UNHANDLED CODE: ", curCode.Debug())
						continue
					}

					// fnDebug("SUCCESS: ", curCode.Debug())

				} else {

					fnDebug("INVALID CODE: ", curCode.Debug())
				}

				continue

				// SKIP + IGNORE CONTROL CHARS DURING ESCAPE CODE
			} else if (chr > 0) && (chr <= 31) {

				continue

				// WRITE OUT COMPONENT OF ESCAPE SEQUENCE
			} else {

				curCode.Params += string(chr)
			}
		}

		// HANDLE WRITABLE CHARACTERS OUTSIDE OF ESCAPE MODE
		if !bEsc {

			if chr == CHR_LF {

				oErr = pGrid.Put(curPos, ' ', bsSGR)

				if oErr != nil {
					fnDebug(oErr)
				}

				curPos.Y += 1
				curPos.X = 1

			} else if chr != '\b' {

				oErr = pGrid.Put(curPos, Array437[chr], bsSGR)

				if oErr != nil {
					fnDebug(oErr)
				}

				pGrid.Inc(&curPos)
			}
		}
	}

	pGrid.Print(pWriter, int(pnRowBytes), pbDebug, w, artw)
	pWriter.WriteByte(CHR_LF)
	pWriter.Flush()

	return

}

// for art 80 cols wide or less
func RenderArt(file string, h int, w int) int {

	rowCount := 0

	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	noSauce := TrimStringFromSauce(string(content)) // strip off the SAUCE metadata
	s := bufio.NewScanner(strings.NewReader(string(noSauce)))

	for s.Scan() {

		fmt.Println(s.Text())
		time.Sleep(70 * time.Millisecond) // wait for a bit between lines
		rowCount++
	}

	return rowCount
}

func FontControl(file string) {

	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	s := bufio.NewScanner(strings.NewReader(string(content)))

	fmt.Println(s)

}

func ScrollAnsi(selected string, h int, w int, currLoc int) {

	row := 0

	fmt.Fprintf(os.Stdout, escapes.EraseScreen)
	fmt.Println(escapes.CursorPos(0, 0))

	f, _ := os.Open(selected)
	// Create new Scanner
	scanner := bufio.NewScanner(f)

	// Use Scan
	for scanner.Scan() {
		if w > 80 {

		} else {

			line := scanner.Text()
			trimmed := theme.TrimStringFromSauce(line)
			row++

			if row <= currLoc && row >= currLoc-h {
				fmt.Println(trimmed)

			}
		}
	}

}

func StatusBar(h int, rowCount int, currLoc int) {

	fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))

	incr := float64(currLoc) / float64(rowCount)
	r := math.Floor(incr*100) / 100
	fmt.Fprintf(os.Stdout, theme.BgCyan+theme.BrightYellow+"[          ]"+theme.Reset)

	switch {
	case r <= 1 && r > .9:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))
		theme.ShowArt(bar100)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .9 && r > .8:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))
		theme.ShowArt(bar90)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .8 && r > .7:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))
		theme.ShowArt(bar80)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .7 && r > .6:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))
		theme.ShowArt(bar70)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .6 && r > .5:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))
		theme.ShowArt(bar60)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .5 && r > .4:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))
		theme.ShowArt(bar50)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .4 && r > .3:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))
		theme.ShowArt(bar40)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .3 && r > .2:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))
		theme.ShowArt(bar30)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .2 && r > .1:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))
		theme.ShowArt(bar20)
		fmt.Fprintf(os.Stdout, theme.Reset)
	case r <= .1:
		fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))
		theme.ShowArt(bar10)
		fmt.Fprintf(os.Stdout, theme.Reset)
	}
}

func TrimStringFromSauce(s string) string {

	if idx := strings.Index(s, "COMNT"); idx != -1 {
		string := s
		delimiter := "COMNT"
		leftOfDelimiter := strings.Split(string, delimiter)[0]
		trim := TrimLastChar(leftOfDelimiter)
		return trim
	}
	if idx := strings.Index(s, "SAUCE00"); idx != -1 {
		string := s
		delimiter := "SAUCE00"
		leftOfDelimiter := strings.Split(string, delimiter)[0]
		trim := TrimLastChar(leftOfDelimiter)
		return trim
	}
	return s
}

func TrimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}

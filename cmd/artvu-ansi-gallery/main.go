package main

import (
	"fmt"
	"log"
	"os"

	"github.com/robbiew/artvu-ansi-gallery/internal/ansi"
	"github.com/robbiew/artvu-ansi-gallery/internal/menu"
	"github.com/robbiew/artvu-ansi-gallery/pkg/term"
)

func main() {

	var headerH int

	var h int
	var w int
	var rootDir string

	headerH = 2

	if len(os.Args) == 1 {
		log.Fatal("Please specify path, or .")
		return
	}
	if rootDir = os.Args[1]; rootDir == "" {
		log.Fatal("Please specify path, or .")
		return
	}
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		// path/to/whatever does not exist
		log.Fatal("Path does not exist!")
		return
	}
	if _, err := os.Stat(rootDir); !os.IsNotExist(err) {
		// path/to/whatever exists, continue
	}

	// Try and detect the user's term size
	h, w = term.GetTermSize()

	fmt.Println(ansi.Clear)
	fmt.Println(ansi.Home)

	menu.MenuAction(rootDir, h, w, headerH)
}

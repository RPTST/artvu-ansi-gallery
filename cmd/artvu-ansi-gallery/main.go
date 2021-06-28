package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/robbiew/artvu-ansi-gallery/internal/menu"
	"github.com/robbiew/artvu-ansi-gallery/internal/theme"
	"github.com/robbiew/artvu-ansi-gallery/pkg/term"
)

func main() {
	rootPtr := flag.String("root", "/foo/bar", "path to art folder")

	required := []string{"root"}
	flag.Parse()

	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			// or possibly use `log.Fatalf` instead of:
			fmt.Fprintf(os.Stderr, "missing required -%s argument/flag\n", req)
			os.Exit(2) // the same exit code flag.Parse uses
		}
	}

	flag.Parse()

	headerH := 5
	rootDir := *rootPtr

	// Try and detect the user's term size
	h, w := term.GetTermSize()

	fmt.Println(theme.Clear)
	fmt.Println(theme.Home)

	theme.ShowSplash(w)
	time.Sleep(2 * time.Second)

	menu.MenuAction(rootDir, h, w, headerH)
}

package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/robbiew/artvu-ansi-gallery/internal/ansi"
	"github.com/robbiew/artvu-ansi-gallery/internal/menu"
	"github.com/robbiew/artvu-ansi-gallery/internal/theme"
	"github.com/robbiew/artvu-ansi-gallery/pkg/term"
)

type Config struct {
	HeaderHeight int
	ArtDir       string
}

func main() {

	var h int
	var w int

	var conf Config
	if _, err := toml.DecodeFile("/robbiew/artvu-ansi-gallery/cfg/config.toml", &conf); err != nil {
		// handle error
		if err != nil {
			panic(err)
		}
	}

	headerH := conf.HeaderHeight
	rootDir := conf.ArtDir

	// Try and detect the user's term size
	h, w = term.GetTermSize()

	fmt.Println(ansi.Clear)
	fmt.Println(ansi.Home)

	theme.ShowSplash(w)
	time.Sleep(2 * time.Second)

	menu.MenuAction(rootDir, h, w, headerH)
}

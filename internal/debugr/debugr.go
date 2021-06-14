package debugr

import (
	"fmt"
	"os"

	escapes "github.com/snugfox/ansi-escapes"
)

func DebugInf(currentDir int, h int, selected string, action int, currentPath string, count int, navOn bool) {

	fmt.Fprintf(os.Stdout, escapes.CursorPos(0, h))
	fmt.Fprintf(os.Stdout, escapes.TextDeleteLine)
	fmt.Fprintf(os.Stdout, escapes.CursorPos(3, h))
	fmt.Fprintf(os.Stdout, "p.Action: %v, p.CurrentDir: %v, p1.CurrentPath: %v, p.Selected: %v, p.Count: %v, p.navOn: %v", action, currentDir, currentPath, selected, count, navOn)

}

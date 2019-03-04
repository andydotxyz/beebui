package beebui

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne"
)

func (b *beeb) LIST() {
	lines := strings.Split(b.program, "\n")
	for i, line := range lines {
		if i == len(lines)-1 {
			break
		}
		b.appendLine(line)
	}
}

func (b *beeb) QUIT(app fyne.App) {
	time.Sleep(lineDelay)
	app.Quit()
}

func (b *beeb) RESTART() {
	b.program = ""
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	b.appendLine(fmt.Sprintf("BBC Computer %dK", int(m.HeapSys/1024)))
	b.appendLine("")
	b.appendLine(strings.Title(runtime.GOOS) + " DFS")
	b.appendLine("")
	b.appendLine("BASIC")
	b.appendLine("")
	b.appendLine(">")
}
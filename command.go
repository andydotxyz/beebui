package beebui

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne"
)

func (b *beeb) QUIT(app fyne.App) {
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
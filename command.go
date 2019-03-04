package beebui

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"github.com/andydotxyz/gobasic/builtin"
	"github.com/andydotxyz/gobasic/object"
)

func (b *beeb) CLS(env builtin.Environment, args []object.Object) object.Object {
	for i := 0; i < len(b.content); i++ {
		text := b.content[i].(*canvas.Text)
		text.Text = ""
		canvas.Refresh(text)
	}
	b.current = -1

	return &object.NumberObject{Value: 0}
}

func (b *beeb) LIST() {
	lines := strings.Split(b.program, "\n")
	for i, line := range lines {
		if i == len(lines)-1 {
			break
		}
		b.appendLine(line)
	}
}

func (b *beeb) NEW() {
	b.program = ""
}

func (b *beeb) QUIT(app fyne.App) {
	time.Sleep(lineDelay)
	app.Quit()
}

func (b *beeb) RESTART() {
	b.NEW()

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

func (b *beeb) RUN() {
	b.runProg(b.program)
}
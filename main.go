// Package beebui emulates a BBC Micro Computer
package beebui

import (
	"bufio"
	"fmt"
	"image/color"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"

	"github.com/andydotxyz/gobasic/builtin"
	"github.com/andydotxyz/gobasic/eval"
	"github.com/andydotxyz/gobasic/object"
	"github.com/andydotxyz/gobasic/tokenizer"
)

const (
	screenInsetX = 130
	screenInsetY = 62

	screenLines = 25
	screenCols  = 40
)

var screenSize = fyne.Size{800, 600}

type beeb struct {
	content []fyne.CanvasObject
	overlay *canvas.Image
	current int

	program string
}

func (b *beeb) Write(p []byte) (n int, err error) {
	str := string(p)
	b.appendLine(str[:len(str)-1])
	return len(p), nil
}

func (b *beeb) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return screenSize
}

func (b *beeb) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	b.overlay.Resize(size)

	y := screenInsetY
	for i := 0; i < screenLines; i++ {
		b.content[i].Move(fyne.NewPos(screenInsetX, y))
		b.content[i].Resize(fyne.NewSize(size.Width-screenInsetX*2, 18))
		y += 19
	}
}

func (b *beeb) loadUI() fyne.CanvasObject {
	b.content = make([]fyne.CanvasObject, screenLines)

	for i := 0; i < screenLines; i++ {
		b.content[i] = canvas.NewText("", color.RGBA{0xbb, 0xbb, 0xbb, 0xff})
		b.content[i].(*canvas.Text).TextSize = 15
		b.content[i].(*canvas.Text).TextStyle.Monospace = true
	}

	b.overlay = canvas.NewImageFromResource(monitor)
	return fyne.NewContainerWithLayout(b, append(b.content, b.overlay)...)
}

func (b *beeb) appendLine(line string) {
	if b.current >= 0 {
		text := b.content[b.current].(*canvas.Text)

		if len(text.Text) > 0 && text.Text[len(text.Text)-1] == '_' {
			text.Text = text.Text[:len(text.Text)-1]
			canvas.Refresh(text)
		}
	}

	if b.current == screenLines-1 {
		b.scroll()
	}
	b.current++
	text := b.content[b.current].(*canvas.Text)
	text.Text = line

	canvas.Refresh(text)
}

func (b *beeb) append(line string) {
	text := b.content[b.current].(*canvas.Text)
	if len(text.Text) > 0 && text.Text[len(text.Text)-1] == '_' {
		text.Text = text.Text[:len(text.Text)-1] + line + "_"
	} else {
		text.Text = text.Text + line
	}

	canvas.Refresh(text)
}

func (b *beeb) blink() {
	for {
		time.Sleep(time.Second / 2)
		line := b.content[b.current].(*canvas.Text)

		if line.Text == "" {
			continue
		}
		if line.Text[len(line.Text)-1] == '_' {
			line.Text = line.Text[:len(line.Text)-1]
		} else {
			line.Text = line.Text + "_"
		}
		canvas.Refresh(b.content[b.current])
	}
}

func (b *beeb) CLS(env builtin.Environment, args []object.Object) object.Object {
	for i := 0; i < len(b.content); i++ {
		text := b.content[i].(*canvas.Text)
		text.Text = ""
		canvas.Refresh(text)
	}
	b.current = -1

	return &object.NumberObject{Value: 0}
}

func (b *beeb) scroll() {
	for i := 0; i < len(b.content)-1; i++ {
		text1 := b.content[i].(*canvas.Text)
		text2 := b.content[i+1].(*canvas.Text)
		text1.Text = text2.Text

		canvas.Refresh(text1)
	}

	text := b.content[len(b.content)-1].(*canvas.Text)
	text.Text = ""
	canvas.Refresh(text)

	b.current -= 1
}

func (b *beeb) onRune(r rune) {
	b.append(string(r))
}

func (b *beeb) onKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyReturn:
		text := b.content[b.current].(*canvas.Text).Text[1:]
		if len(text) > 0 && text[len(text)-1] == '_' {
			text = text[:len(text)-1]
		}
		prog := strings.TrimSpace(text) + "\n"

		first := prog[0]
		if first >= '0' && first <= '9' {
			b.program += prog
		} else if prog == "RUN\n" {
			b.runProg(b.program)
		} else if prog == "NEW\n" {
			b.program = ""
		} else {
			b.runProg(prog)
		}
		b.appendLine(">")
	case fyne.KeyBackspace:
		line := b.content[b.current].(*canvas.Text)
		text := line.Text[1:]
		if len(text) > 0 && text[len(text)-1] == '_' {
			text = text[:len(text)-1]
		}
		if len(text) > 0 {
			line.Text = ">" + text[:len(text)-1]
			canvas.Refresh(line)
		}
	}
}

func (b *beeb) runProg(prog string) {
	t := tokenizer.New(prog)
	e, err := eval.New(t)
	e.STDOUT = bufio.NewWriterSize(b, screenCols)
	e.RegisterBuiltin("CLS", 0, b.CLS)
	if err != nil {
		fmt.Println("Error parsing program", err)
	} else {
		err = e.Run()
		if err != nil {
			fmt.Println("Error running program", err)
		}
	}
}

// Show starts a new beeb computer simulator
func Show(app fyne.App) {
	b := beeb{}
	b.current = -1
	app.Settings().SetTheme(&beebTheme{})

	window := app.NewWindow("BBC Emulator")
	window.SetContent(b.loadUI())
	window.SetPadded(false)
	window.SetFixedSize(true)
	window.Resize(screenSize)

	window.Canvas().SetOnTypedRune(b.onRune)
	window.Canvas().SetOnTypedKey(b.onKey)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	b.appendLine(fmt.Sprintf("BBC Computer %dK", int(m.HeapSys/1024)))
	b.appendLine("")
	b.appendLine(strings.Title(runtime.GOOS) + " DFS")
	b.appendLine("")
	b.appendLine("BASIC")
	b.appendLine("")
	b.appendLine(">")

	go b.blink()

	window.Show()
}

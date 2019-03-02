// Package beebui emulates a BBC Micro Computer
package beebui

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"

	"github.com/skx/gobasic/tokenizer"
	"github.com/skx/gobasic/eval"
)

const (
	screenInsetX = 130
	screenInsetY = 62
)

var screenSize = fyne.Size{800, 600}

type beeb struct {
	content []fyne.CanvasObject
	overlay *canvas.Image
	current int
}

func (b *beeb) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return screenSize
}

func (b *beeb) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	b.overlay.Resize(size)

	y := screenInsetY
	for i := 0; i < 25; i++ {
		b.content[i].Move(fyne.NewPos(screenInsetX, y))
		b.content[i].Resize(fyne.NewSize(size.Width-screenInsetX*2, 18))
		y += 19
	}
}

func (b *beeb) loadUI() fyne.CanvasObject {
	b.content = make([]fyne.CanvasObject, 25)

	for i := 0; i < 25; i++ {
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
		time.Sleep(time.Second/2)
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

func (b *beeb) onRune(r rune) {
	b.append(string(r))
}

func (b *beeb) onKey(ev *fyne.KeyEvent) {
	switch ev.Name {
	case fyne.KeyReturn:
		prog := b.content[b.current].(*canvas.Text).Text[1:]+"\n"
		t := tokenizer.New(prog)
		e, err := eval.New(t)
//		e.RegisterBuiltin("CLS", 0, clear)
		if err != nil {
			fmt.Println("Error parsing program", err)
		} else {
			err = e.Run()
			if err != nil {
				fmt.Println("Error running program", err)
			}
		}
		b.appendLine(">")
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

	b.appendLine("BBC Computer 16K")
	b.appendLine("")
	b.appendLine("Acorn DFS")
	b.appendLine("")
	b.appendLine("BASIC")
	b.appendLine("")
	b.appendLine(">")

	go b.blink()

	window.Show()
}

// Package beebui emulates a BBC Micro Computer
package beebui

import (
	"bufio"
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"

	"github.com/skx/gobasic/eval"
	"github.com/skx/gobasic/tokenizer"
)

const (
	screenInsetX = 130
	screenInsetY = 62

	screenLines = 25
	screenCols  = 40
)

var (
	screenSize = fyne.Size{800, 600}
	lineDelay = time.Second / 10

	Icon = icon
)

type beeb struct {
	content []fyne.CanvasObject
	overlay *canvas.Image
	current int

	program  string
	bufInput []byte
	endInput bool
	nextAuto int
}

func (b *beeb) Read(p []byte) (n int, err error) {
	if b.endInput {
		b.endInput = false
		b.bufInput = nil
		p[0] = '\n'
		return 1, nil
	}
	b.bufInput = p

	if p[0] != 0 {
		return 1, nil
	}

	time.Sleep(lineDelay)
	return 0, nil
}

func (b *beeb) Write(p []byte) (n int, err error) {
	str := string(p)
	if str[len(str)-1] == '\n' {
		b.appendLine(str[:len(str)-1])
	} else {
		b.append(str)
	}
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
	b.append(line)

	b.newLine()
}

func (b *beeb) newLine() {
	time.Sleep(lineDelay)
	text := b.content[b.current].(*canvas.Text)

	if len(text.Text) > 0 && text.Text[len(text.Text)-1] == '_' {
		text.Text = text.Text[:len(text.Text)-1]
		canvas.Refresh(text)
	}

	if b.current == screenLines-1 {
		b.scroll()
	}
	b.current++
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
	if b.bufInput != nil {
		b.bufInput[0] = byte(r) // TODO could we have typed another?
	}
	b.append(string(r))
}

func (b *beeb) onKey(ev *fyne.KeyEvent) {
	if b.bufInput != nil {
		if ev.Name == fyne.KeyReturn {
			b.endInput = true
		}
		return
	}
	switch ev.Name {
	case fyne.KeyReturn:
		prog := ">"
		text := b.content[b.current].(*canvas.Text).Text
		if len(text) > 1 {
			text = text[1:]
			if len(text) > 0 && text[len(text)-1] == '_' {
				text = text[:len(text)-1]
			}
			prog = strings.TrimSpace(text) + "\n"
		}
		b.appendLine("")
		first := prog[0]
		if first >= '0' && first <= '9' {
			b.program += prog
		} else {
			// commands that can't be called from within a program
			cmd := strings.ToUpper(prog[:len(prog)-1])
			if cmd == "AUTO" {
				b.nextAuto = 10
			} else if cmd == "RUN" {
				b.RUN()
			} else if cmd == "NEW" {
				b.NEW()
			} else if cmd == "LIST" {
				b.LIST()
			} else if cmd == "QUIT" || cmd == "EXIT" {
				b.QUIT(fyne.CurrentApp())
			} else {
				b.runProg(prog)
			}
		}
		b.append(">")
		if b.nextAuto > 0 {
			b.append(fmt.Sprintf("%d ", b.nextAuto))
			b.nextAuto += 10
		}
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
	case fyne.KeyEscape:
		text := b.content[b.current].(*canvas.Text)
		if len(text.Text) == 0 || text.Text[0] != '>' {
			break
		}

		b.nextAuto = 0
		text.Text = ">"
		canvas.Refresh(text)
	}
}

func (b *beeb) runProg(prog string) {
	t := tokenizer.New(prog)
	e, err := eval.New(t)
	e.STDIN = bufio.NewReaderSize(b, screenCols)
	e.STDOUT = bufio.NewWriterSize(b, screenCols)
	e.STDERR = e.STDOUT
	e.LINEEND = "\n"
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
	app.Settings().SetTheme(&beebTheme{})

	window := app.NewWindow("BBC Emulator")
	window.SetContent(b.loadUI())
	window.SetPadded(false)
	window.SetFixedSize(true)
	window.Resize(screenSize)

	window.Canvas().SetOnTypedRune(b.onRune)
	window.Canvas().SetOnTypedKey(b.onKey)
	window.Canvas().AddShortcut(&desktop.CustomShortcut{
		Modifier: desktop.ControlModifier,
		KeyName: fyne.KeyD,
	}, func(fyne.Shortcut) {
		b.append("QUIT")
		b.appendLine("")
		go func() {
			b.QUIT(app)
		}()
	})

	b.RESTART()

	go b.blink()

	window.Show()
}

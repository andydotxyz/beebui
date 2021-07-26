// Package main launches the beeb computer emulator directly
package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/mobile"
	"github.com/andydotxyz/beebui"
)

func main() {
	app := app.New()
	app.SetIcon(beebui.Icon)

	beebui.Show(app)

	if mob, ok := fyne.CurrentDevice().(mobile.Device); ok {
		go func() {
			time.Sleep(100 * time.Millisecond)
			mob.ShowVirtualKeyboard()
		}()
	}

	app.Run()
}

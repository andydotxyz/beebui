// Package main launches the beeb computer emulator directly
package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/andydotxyz/beebui"
)

func main() {
	app := app.New()
	app.SetIcon(beebui.Icon)

	beebui.Show(app)
	app.Run()
}

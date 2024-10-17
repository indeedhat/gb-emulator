package ui

import (
	"image"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	fynecanvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/indeedhat/gb-emulator/internal/emu/config"
	"github.com/indeedhat/gb-emulator/internal/emu/enum"
	"github.com/indeedhat/gb-emulator/internal/emu/types"
)

func NewFyneRenderer() (fyne.App, fyne.Window) {
	runner := fyneapp.NewWithID("dev.indeedhat.gb-emu")

	win := runner.NewWindow("Emulator")
	win.Resize(fyne.NewSize(640, 480))

	im := fynecanvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, 0, 0)))
	im.FillMode = fynecanvas.ImageFillContain
	win.SetContent(im)

	canvas := win.Canvas()
	app := &App{
		window: win,
		runner: runner,
		frame:  im,
	}
	app.menu = NewMenu(runner, app)

	dc := canvas.(desktop.Canvas)
	dc.SetOnKeyDown(app.handleKeyDown)
	dc.SetOnKeyUp(app.handleKeyUp)

	win.SetMainMenu(app.menu.Root)

	return runner, win
}

func generateImage(data []types.Pixel) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, config.PpuXRes, config.PpuYRes))

	for i, px := range data {
		img.Pix[i*4] = px.R
		img.Pix[i*4+1] = px.G
		img.Pix[i*4+2] = px.B
		img.Pix[i*4+3] = 0xFF
	}

	return img
}

func mapKeyCode(p fyne.Preferences, e *fyne.KeyEvent) enum.KeyCode {
	switch e.Name {
	case fyne.KeyName(p.StringWithFallback(PrefControlsUp, PrefControlsUpFallback)):
		return enum.KeyUp
	case fyne.KeyName(p.StringWithFallback(PrefControlsRight, PrefControlsRightFallback)):
		return enum.KeyRight
	case fyne.KeyName(p.StringWithFallback(PrefControlsDown, PrefControlsDownFallback)):
		return enum.KeyDown
	case fyne.KeyName(p.StringWithFallback(PrefControlsLeft, PrefControlsLeftFallback)):
		return enum.KeyLeft
	case fyne.KeyName(p.StringWithFallback(PrefControlsA, PrefControlsAFallback)):
		return enum.KeyA
	case fyne.KeyName(p.StringWithFallback(PrefControlsB, PrefControlsBFallback)):
		return enum.KeyB
	case fyne.KeyName(p.StringWithFallback(PrefControlsStart, PrefControlsStartFallback)):
		return enum.KeyStart
	case fyne.KeyName(p.StringWithFallback(PrefControlsSelect, PrefControlsSelectFallback)):
		return enum.KeySelect
	}

	return enum.KeyUnknown
}

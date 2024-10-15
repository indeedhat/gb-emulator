package ui

import (
	"image"

	"fyne.io/fyne/v2"
	fyneapp "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/indeedhat/gb-emulator/internal/emu/config"
	"github.com/indeedhat/gb-emulator/internal/emu/enum"
	"github.com/indeedhat/gb-emulator/internal/emu/types"
)

func NewFyneRenderer() (fyne.App, fyne.Window) {
	runner := fyneapp.New()

	win := runner.NewWindow("Emulator")
	win.Resize(fyne.NewSize(640, 480))

	canvas := win.Canvas()
	app := &App{window: win}

	dc := canvas.(desktop.Canvas)
	dc.SetOnKeyDown(app.handleKeyDown)
	dc.SetOnKeyUp(app.handleKeyUp)

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.FolderOpenIcon(), app.handleLoadRom),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.MediaPlayIcon(), app.handleUnPauseEmulation),
		widget.NewToolbarAction(theme.MediaPauseIcon(), app.handlePauseEmulation),
		widget.NewToolbarAction(theme.MediaStopIcon(), app.handleStopEmulation),

		widget.NewToolbarAction(theme.DocumentSaveIcon(), app.handleSaveState),
		widget.NewToolbarAction(theme.LoginIcon(), app.handleLoadState),
	)

	app.container = container.NewBorder(toolbar, nil, nil, nil, container.NewWithoutLayout())
	win.SetContent(app.container)

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

func mapKeyCode(e *fyne.KeyEvent) enum.KeyCode {
	switch e.Name {
	case fyne.KeyComma:
		return enum.KeyUp
	case fyne.KeyE:
		return enum.KeyRight
	case fyne.KeyO:
		return enum.KeyDown
	case fyne.KeyA:
		return enum.KeyLeft
	case fyne.KeyEnter, fyne.KeyReturn:
		return enum.KeyA
	case fyne.KeyJ:
		return enum.KeyB
	case fyne.KeyPeriod:
		return enum.KeyStart
	case fyne.KeyApostrophe:
		return enum.KeySelect
	}

	return enum.KeyUnknown
}

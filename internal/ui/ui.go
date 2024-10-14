package ui

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/indeedhat/gb-emulator/internal/emu/config"
	"github.com/indeedhat/gb-emulator/internal/emu/context"
	"github.com/indeedhat/gb-emulator/internal/emu/enum"
	"github.com/indeedhat/gb-emulator/internal/emu/types"
)

func NewFyneRenderer(ctx *context.Context) (fyne.App, fyne.Window) {
	a := app.New()

	w := a.NewWindow("Emulator")
	w.Resize(fyne.NewSize(640, 480))

	c := w.Canvas()

	dc := c.(desktop.Canvas)
	dc.SetOnKeyDown(func(e *fyne.KeyEvent) {
		code := mapKeyCode(e)
		if code == enum.KeyUnknown {
			return
		}

		ctx.JoypadCh <- types.KeyEvent{
			Key:  code,
			Down: true,
		}
	})
	dc.SetOnKeyUp(func(e *fyne.KeyEvent) {
		code := mapKeyCode(e)
		if code == enum.KeyUnknown {
			return
		}

		ctx.JoypadCh <- types.KeyEvent{
			Key:  code,
			Down: false,
		}
	})

	go func() {
		for img := range ctx.FrameCh {
			frame := canvas.NewImageFromImage(generateImage(img))
			frame.FillMode = canvas.ImageFillContain
			c.SetContent(frame)
		}
	}()

	return a, w
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

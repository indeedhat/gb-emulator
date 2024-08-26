package ui

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/indeedhat/gb-emulator/internal/emu"
)

func NewFyneRenderer(ctx *emu.Context) (fyne.App, fyne.Window) {
	a := app.New()

	w := a.NewWindow("Emulator")
	w.Resize(fyne.NewSize(640, 480))

	c := w.Canvas()

	dc := c.(desktop.Canvas)
	dc.SetOnKeyDown(func(e *fyne.KeyEvent) {
		code := mapKeyCode(e)
		if code == emu.KeyUnknown {
			return
		}

		ctx.JoypadCh <- emu.KeyEvent{
			Key:  code,
			Down: true,
		}
	})
	dc.SetOnKeyUp(func(e *fyne.KeyEvent) {
		code := mapKeyCode(e)
		if code == emu.KeyUnknown {
			return
		}

		ctx.JoypadCh <- emu.KeyEvent{
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

func generateImage(data []emu.Pixel) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, emu.PpuXRes, emu.PpuYRes))

	for i, px := range data {
		img.Pix[i*4] = px.R
		img.Pix[i*4+1] = px.G
		img.Pix[i*4+2] = px.B
		img.Pix[i*4+3] = 0xFF
	}

	return img
}

func mapKeyCode(e *fyne.KeyEvent) emu.KeyCode {
	switch e.Name {
	case fyne.KeyComma:
		return emu.KeyUp
	case fyne.KeyE:
		return emu.KeyRight
	case fyne.KeyO:
		return emu.KeyDown
	case fyne.KeyA:
		return emu.KeyLeft
	case fyne.KeyEnter, fyne.KeyReturn:
		return emu.KeyA
	case fyne.KeyJ:
		return emu.KeyB
	case fyne.KeyPeriod:
		return emu.KeyStart
	case fyne.KeyApostrophe:
		return emu.KeySelect
	}

	return emu.KeyUnknown
}

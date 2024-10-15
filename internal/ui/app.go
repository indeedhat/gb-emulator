package ui

import (
	"os"

	"fyne.io/fyne/v2"
	fynecanvas "fyne.io/fyne/v2/canvas"
	fynedialog "fyne.io/fyne/v2/dialog"

	"github.com/sqweek/dialog"

	"github.com/indeedhat/gb-emulator/internal/emu"
	"github.com/indeedhat/gb-emulator/internal/emu/context"
	"github.com/indeedhat/gb-emulator/internal/emu/enum"
	"github.com/indeedhat/gb-emulator/internal/emu/types"
)

type App struct {
	emu  *emu.Emulator
	ctx  *context.Context
	done chan struct{}

	window    fyne.Window
	container *fyne.Container
}

func (a *App) renderLoop() {
	for {
		select {
		case <-a.done:
			break
		case img := <-a.ctx.FrameCh:
			frame := fynecanvas.NewImageFromImage(generateImage(img))
			frame.FillMode = fynecanvas.ImageFillContain
			a.container.Objects[0] = frame
			a.container.Refresh()
		}
	}
}

func (a *App) handleLoadRom() {
	if a.emu != nil {
		a.handleStopEmulation()
	}

	dir, err := os.Getwd()
	if err != nil {
		fynedialog.ShowError(err, a.window)
		return
	}

	filename, err := dialog.File().SetStartDir(dir).Filter(".gb files", "gb").Load()
	if err != nil {
		fynedialog.ShowError(err, a.window)
		return
	}

	a.emu, a.ctx, err = emu.NewEmulator(filename, false)
	if err != nil {
		fynedialog.ShowError(err, a.window)
		a.handleStopEmulation()
		return
	}

	a.done = make(chan struct{})

	go a.emu.Run()
	go a.renderLoop()
}

func (a *App) handlePauseEmulation() {
	if a.emu != nil {
		a.emu.Pause()
	}
}

func (a *App) handleUnPauseEmulation() {
	if a.emu != nil {
		a.emu.Play()
	}
}

func (a *App) handleStopEmulation() {
	if a.emu != nil {
		a.done <- struct{}{}
		close(a.done)

		e := a.emu
		defer e.Stop()

		a.emu = nil
		a.container.Objects[0] = fyne.NewContainerWithoutLayout()
		a.container.Refresh()
	}
}

func (a *App) handleKeyUp(e *fyne.KeyEvent) {
	if a.emu == nil || !a.emu.IsRunning() || a.emu.IsPaused() {
		return
	}

	code := mapKeyCode(e)
	if code == enum.KeyUnknown {
		return
	}

	a.ctx.JoypadCh <- types.KeyEvent{
		Key:  code,
		Down: false,
	}
}

func (a *App) handleKeyDown(e *fyne.KeyEvent) {
	if a.emu == nil || !a.emu.IsRunning() || a.emu.IsPaused() {
		return
	}

	code := mapKeyCode(e)
	if code == enum.KeyUnknown {
		return
	}

	a.ctx.JoypadCh <- types.KeyEvent{
		Key:  code,
		Down: true,
	}
}

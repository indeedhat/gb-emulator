package ui

import (
	"image"
	"os"
	"time"

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

	runner fyne.App
	window fyne.Window
	frame  *fynecanvas.Image

	menu *Menu

	stateSlot       int
	stateSlotRotate bool
	stateAutoSave   bool
}

func (a *App) renderLoop() {
	for {
		select {
		case <-a.done:
			break
		case img := <-a.ctx.FrameCh:
			a.frame.Image = generateImage(img)
			a.frame.Refresh()
		}
	}
}

func (a *App) autosaveLoop() {
	interval := a.runner.Preferences().IntWithFallback(PrefAutoSaveInterval, PrefAutoSaveIntervalFallback)
	ticker := time.NewTicker(time.Duration(interval) * time.Second)

	for {
		select {
		case <-a.done:
			break
		case <-ticker.C:
			if a.emu != nil && a.runner.Preferences().Bool(PrefAutoSaveState) {
				a.handleSaveState(a.menu.statePath(10))()
			}
		}
	}
}

func (a *App) handleLoadRom(filename string) func() {
	return func() {
		var err error
		if filename == "" {
			if a.emu != nil {
				a.handleStopEmulation()
			}

			dir, err := os.Getwd()
			if err != nil {
				fynedialog.ShowError(err, a.window)
				return
			}

			filename, err = dialog.File().SetStartDir(dir).Filter(".gb files", "gb").Load()
			if err != nil {
				fynedialog.ShowError(err, a.window)
				return
			}
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
		go a.autosaveLoop()

		a.menu.TriggerEmuRunnung()
		a.menu.TriggerRecentReload(filename)
		a.menu.TriggerStateReload()
	}
}

func (a *App) handleAutosaveToggle() bool {
	current := a.runner.Preferences().Bool(PrefAutoSaveState)
	a.runner.Preferences().SetBool(PrefAutoSaveState, !current)
	a.menu.TriggerStateReload()
	return !current
}

func (a *App) handlePauseEmulation() {
	if a.emu != nil {
		a.emu.Pause()
		a.menu.TriggerEmuPause()
	}
}

func (a *App) handleUnPauseEmulation() {
	if a.emu != nil {
		a.emu.Play()
		a.menu.TriggerEmuRunnung()
	}
}

func (a *App) handleStopEmulation() {
	if a.emu != nil {
		// TODO: this is a nasty hack to close both the render and autosave loops but its midnight
		//       and i can't be bothered to do this properly
		a.done <- struct{}{}
		a.done <- struct{}{}
		close(a.done)

		e := a.emu
		defer e.Stop()
		defer a.menu.TriggerEmuStop()

		a.emu = nil
		a.frame.Image = image.NewRGBA(image.Rect(0, 0, 0, 0))
		a.frame.Refresh()
	}
}

func (a *App) handleSaveState(path string) func() {
	return func() {
		a.emu.SaveState(path)
		a.menu.TriggerStateReload()
	}
}

func (a *App) handleLoadState(path string) func() {
	return func() {
		a.emu.LoadState(path)
	}
}

func (a *App) handleKeyUp(e *fyne.KeyEvent) {
	if a.emu == nil || !a.emu.IsRunning() || a.emu.IsPaused() {
		return
	}

	code := mapKeyCode(a.runner.Preferences(), e)
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

	code := mapKeyCode(a.runner.Preferences(), e)
	if code == enum.KeyUnknown {
		return
	}

	a.ctx.JoypadCh <- types.KeyEvent{
		Key:  code,
		Down: true,
	}
}

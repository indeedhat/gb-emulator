package ui

import (
	"fmt"
	"os"
	"path"
	"slices"
	"time"

	"fyne.io/fyne/v2"
	"github.com/indeedhat/gb-emulator/internal/emu/cart"
)

type Menu struct {
	Root *fyne.MainMenu

	File struct {
		Root   *fyne.Menu
		Open   *fyne.MenuItem
		Recent *fyne.MenuItem
	}

	Emulator struct {
		Root  *fyne.Menu
		Play  *fyne.MenuItem
		Pause *fyne.MenuItem
		Stop  *fyne.MenuItem
	}

	State struct {
		Root     *fyne.Menu
		Save     *fyne.MenuItem
		Load     *fyne.MenuItem
		AutoSave *fyne.MenuItem
	}

	Window struct {
		Root        *fyne.Menu
		Fullscreen  *fyne.MenuItem
		Preferences *fyne.MenuItem
	}

	app    *App
	runner fyne.App
}

func NewMenu(runner fyne.App, app *App) *Menu {
	m := &Menu{
		app:    app,
		runner: runner,
	}

	m.initStateMenu()
	m.initEmulatorMenu()
	m.initFileMenu()
	m.initWindowMenu()

	m.Root = fyne.NewMainMenu(
		m.File.Root,
		m.Emulator.Root,
		m.State.Root,
		m.Window.Root,
	)

	return m
}

func (m *Menu) TriggerEmuRunnung() {
	m.Emulator.Play.Disabled = true
	m.Emulator.Pause.Disabled = false
	m.Emulator.Stop.Disabled = false

	for i := range 10 {
		m.State.Save.ChildMenu.Items[i].Disabled = false
	}

	m.State.Root.Refresh()
	m.Emulator.Root.Refresh()
}

func (m *Menu) TriggerEmuPause() {
	m.Emulator.Play.Disabled = false
	m.Emulator.Pause.Disabled = true
	m.Emulator.Stop.Disabled = false

	m.Emulator.Root.Refresh()
}

func (m *Menu) TriggerEmuStop() {
	m.Emulator.Play.Disabled = true
	m.Emulator.Pause.Disabled = true
	m.Emulator.Stop.Disabled = true

	for i := range 10 {
		m.State.Save.ChildMenu.Items[i].Disabled = true
	}

	m.Emulator.Root.Refresh()
	m.State.Root.Refresh()
}

func (m *Menu) TriggerRecentReload(current ...string) {
	if len(current) != 0 {
		m.updateRecentList(current[0])
	}

	items := m.runner.Preferences().StringList(PrefRecentRomsList)
	m.File.Recent.ChildMenu.Items = nil

	if len(items) == 0 {
		m.File.Recent.Disabled = true
		return
	}

	m.File.Recent.Disabled = false

	for i, item := range items {
		filename := fmt.Sprintf("%d: %s", i+1, path.Base(item))
		m.File.Recent.ChildMenu.Items = append(m.File.Recent.ChildMenu.Items,
			fyne.NewMenuItem(filename, m.app.handleLoadRom(item)),
		)
	}

	m.File.Root.Refresh()
}

func (m *Menu) TriggerStateReload() {
	for i := range 10 {
		path := m.statePath(i)

		if stat, err := os.Stat(path); err != nil {
			m.State.Load.ChildMenu.Items[i].Disabled = true
			m.State.Load.ChildMenu.Items[i].Label = fmt.Sprintf("Load Slot %d", i+1)
			m.State.Save.ChildMenu.Items[i].Label = fmt.Sprintf("Save Slot %d", i+1)
		} else {
			m.State.Load.ChildMenu.Items[i].Disabled = false
			m.State.Load.ChildMenu.Items[i].Label = fmt.Sprintf("Load Slot %d (%s)",
				i+1,
				stat.ModTime().Format(time.DateTime),
			)

			m.State.Save.ChildMenu.Items[i].Label = fmt.Sprintf("Save Slot %d (%s)",
				i+1,
				stat.ModTime().Format(time.DateTime),
			)
		}
	}

	path := m.statePath(10)
	if stat, err := os.Stat(path); err != nil {
		m.State.Load.ChildMenu.Items[10].Disabled = true
		m.State.Load.ChildMenu.Items[10].Label = "Load Autosave"
	} else {
		m.State.Load.ChildMenu.Items[10].Disabled = false
		m.State.Load.ChildMenu.Items[10].Label = fmt.Sprintf("Load Autosave (%s)",
			stat.ModTime().Format(time.DateTime),
		)
	}

	m.State.Root.Refresh()
}

func (m *Menu) initFileMenu() {
	m.File.Root = fyne.NewMenu("File")

	m.File.Open = fyne.NewMenuItem("Open", m.app.handleLoadRom(""))

	m.File.Recent = fyne.NewMenuItem("Recent", nil)
	m.File.Recent.ChildMenu = fyne.NewMenu("")
	m.TriggerRecentReload()

	m.File.Root.Items = append(m.File.Root.Items,
		m.File.Open,
		m.File.Recent,
	)
}

func (m *Menu) initEmulatorMenu() {
	m.Emulator.Root = fyne.NewMenu("Emulator")

	m.Emulator.Play = fyne.NewMenuItem("Play", m.app.handleUnPauseEmulation)
	m.Emulator.Play.Disabled = true
	m.Emulator.Pause = fyne.NewMenuItem("Pause", m.app.handlePauseEmulation)
	m.Emulator.Pause.Disabled = true
	m.Emulator.Stop = fyne.NewMenuItem("Stop", m.app.handleStopEmulation)
	m.Emulator.Stop.Disabled = true

	m.Emulator.Root.Items = append(m.Emulator.Root.Items,
		m.Emulator.Play,
		m.Emulator.Pause,
		m.Emulator.Stop,
	)
}

func (m *Menu) initStateMenu() {
	m.State.Root = fyne.NewMenu("State")

	// save
	m.State.Save = fyne.NewMenuItem("Save", nil)
	m.State.Save.ChildMenu = fyne.NewMenu("")
	for i := range 10 {
		item := fyne.NewMenuItem(fmt.Sprintf("Save Slot %d", i+1), func() {
			m.app.handleSaveState(m.statePath(i))()
		})
		item.Disabled = true
		m.State.Save.ChildMenu.Items = append(m.State.Save.ChildMenu.Items, item)
	}

	// load
	m.State.Load = fyne.NewMenuItem("Load", nil)
	m.State.Load.ChildMenu = fyne.NewMenu("")
	for i := range 10 {
		item := fyne.NewMenuItem(fmt.Sprintf("Load Slot %d", i+1), func() {
			m.app.handleLoadState(m.statePath(i))()
		})
		m.State.Load.ChildMenu.Items = append(m.State.Load.ChildMenu.Items, item)
	}

	autoSaveSlot := fyne.NewMenuItem("Load Autosave", func() {
		m.app.handleLoadState(m.statePath(10))()
	})
	m.State.Load.ChildMenu.Items = append(m.State.Load.ChildMenu.Items, autoSaveSlot)

	// autosave
	m.State.AutoSave = fyne.NewMenuItem("Autosave", func() {
		m.State.AutoSave.Checked = m.app.handleAutosaveToggle()
		m.TriggerStateReload()
	})
	m.State.AutoSave.Checked = m.runner.Preferences().Bool(PrefAutoSaveState)

	m.State.Root.Items = append(m.State.Root.Items,
		m.State.Load,
		m.State.Save,
		m.State.AutoSave,
	)

	m.TriggerStateReload()
}

func (m *Menu) initWindowMenu() {
	m.Window.Root = fyne.NewMenu("Window")

	m.Window.Fullscreen = fyne.NewMenuItem("Fullscreen", func() {
		current := m.app.window.FullScreen()
		if current {
			m.Window.Fullscreen.Label = "Fullscreen"
		} else {
			m.Window.Fullscreen.Label = "Exit Fullscreen"
		}

		m.app.window.SetFullScreen(!current)

		m.Window.Root.Refresh()
	})

	m.Window.Preferences = fyne.NewMenuItem("Preferences", func() {
		p := NewPreferencesWindow(m.runner)
		p.window.Show()
	})

	m.Window.Root.Items = append(m.Window.Root.Items,
		m.Window.Fullscreen,
		m.Window.Preferences,
	)
}

func (m *Menu) statePath(i int) string {
	if m.app.ctx == nil {
		return ""
	}

	filepath := m.app.ctx.Cart.(*cart.Cartridge).Filepath()
	base := path.Base(filepath)
	dir := path.Dir(filepath)

	return fmt.Sprintf("%s/.gbstate/%s.%d.gbstate", dir, base, i)
}

func (m *Menu) updateRecentList(path string) {
	existing := m.runner.Preferences().StringList(PrefRecentRomsList)
	existing = slices.DeleteFunc(existing, func(val string) bool {
		return val == path
	})

	listSize := m.runner.Preferences().IntWithFallback(PrefRecentRomsCount, PrefRecentRomsCountFallback)
	if len(existing) > listSize {
		existing = existing[:listSize]
	}

	m.runner.Preferences().SetStringList(PrefRecentRomsList, append([]string{path}, existing...))
}

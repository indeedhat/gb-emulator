package ui

import (
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const (
	PrefAutoSaveState            = "auto-save.state"
	PrefAutoSaveInterval         = "auto-save.interval"
	PrefAutoSaveIntervalFallback = 30

	PrefRecentRomsList          = "recent-roms.list"
	PrefRecentRomsCount         = "recent-roms.count"
	PrefRecentRomsCountFallback = 5

	PrefControlsA              = "controls.a"
	PrefControlsAFallback      = "X"
	PrefControlsB              = "controls.b"
	PrefControlsBFallback      = "C"
	PrefControlsUp             = "controls.up"
	PrefControlsUpFallback     = "W"
	PrefControlsDown           = "controls.down"
	PrefControlsDownFallback   = "S"
	PrefControlsLeft           = "controls.left"
	PrefControlsLeftFallback   = "A"
	PrefControlsRight          = "controls.right"
	PrefControlsRightFallback  = "D"
	PrefControlsStart          = "controls.start"
	PrefControlsStartFallback  = "Return"
	PrefControlsSelect         = "controls.select"
	PrefControlsSelectFallback = "Space"
)

type Preferences struct {
	runner fyne.App
	window fyne.Window
}

func NewPreferencesWindow(runner fyne.App) *Preferences {
	window := runner.NewWindow("Preferences")

	p := &Preferences{
		runner: runner,
		window: window,
	}

	window.SetContent(container.NewVBox(
		p.initAutosaveSection(),
		p.initRecentSection(),
		p.initControlsSection(),
	))

	return p
}

func (p *Preferences) initAutosaveSection() *fyne.Container {
	title := widget.NewLabel("Autosave")
	title.TextStyle.Bold = true
	title.TextStyle.Underline = true

	label := widget.NewLabel("Interval")
	interval := widget.NewEntry()
	interval.PlaceHolder = strconv.Itoa(PrefAutoSaveIntervalFallback)
	interval.SetText(
		strconv.Itoa(
			p.runner.Preferences().IntWithFallback(PrefAutoSaveInterval, PrefAutoSaveIntervalFallback),
		),
	)
	interval.Validator = validation.NewRegexp(`\d+`, "Must be a number")
	interval.OnChanged = func(s string) {
		val, err := strconv.Atoi(s)
		if err != nil {
			return
		}
		p.runner.Preferences().SetInt(PrefAutoSaveInterval, val)
	}

	spacer := canvas.NewLine(color.White)

	cont := container.NewVBox(
		title,
		label,
		interval,
		spacer,
	)
	return cont
}

func (p *Preferences) initRecentSection() *fyne.Container {
	title := widget.NewLabel("Recent Items")
	title.TextStyle.Bold = true
	title.TextStyle.Underline = true

	label := widget.NewLabel("List Size")
	interval := widget.NewEntry()
	interval.PlaceHolder = strconv.Itoa(PrefRecentRomsCountFallback)
	interval.SetText(
		strconv.Itoa(
			p.runner.Preferences().IntWithFallback(PrefRecentRomsCount, PrefRecentRomsCountFallback),
		),
	)
	interval.Validator = validation.NewRegexp(`\d+`, "Must be a number")
	interval.OnChanged = func(s string) {
		val, err := strconv.Atoi(s)
		if err != nil {
			return
		}
		p.runner.Preferences().SetInt(PrefRecentRomsCount, val)
	}

	spacer := canvas.NewLine(color.White)

	return container.NewVBox(
		title,
		label,
		interval,
		spacer,
	)
}

func (p *Preferences) initControlsSection() *fyne.Container {
	title := widget.NewLabel("Controls")
	title.TextStyle.Bold = true
	title.TextStyle.Underline = true

	aLabel, aButton := p.initControlButton("A Button", PrefControlsA, PrefControlsAFallback)
	bLabel, bButton := p.initControlButton("B Button", PrefControlsB, PrefControlsBFallback)
	upLabel, upButton := p.initControlButton("Up Button", PrefControlsUp, PrefControlsUpFallback)
	downLabel, downButton := p.initControlButton("Down Button", PrefControlsDown, PrefControlsDownFallback)
	leftLabel, leftButton := p.initControlButton("Left Button", PrefControlsLeft, PrefControlsLeftFallback)
	rightLabel, rightButton := p.initControlButton("Right Button", PrefControlsRight, PrefControlsRightFallback)
	startLabel, startButton := p.initControlButton("Start Button", PrefControlsStart, PrefControlsStartFallback)
	selectLabel, selectButton := p.initControlButton("Select Button", PrefControlsSelect, PrefControlsSelectFallback)

	spacer := canvas.NewLine(color.White)

	return container.NewVBox(
		title,
		aLabel,
		aButton,
		bLabel,
		bButton,
		upLabel,
		upButton,
		downLabel,
		downButton,
		leftLabel,
		leftButton,
		rightLabel,
		rightButton,
		startLabel,
		startButton,
		selectLabel,
		selectButton,
		spacer,
	)
}

func (p *Preferences) initControlButton(label, pref, fallback string) (*widget.Label, *widget.Button) {
	l := widget.NewLabel(label)

	var b *widget.Button
	b = widget.NewButton(
		p.runner.Preferences().StringWithFallback(pref, fallback),
		func() {
			picker := dialog.NewCustom(
				"Control Key Picker",
				"Cancel",
				canvas.NewText("Press a key", color.Gray{}),
				p.window,
			)

			c := p.window.Canvas().(desktop.Canvas)
			c.SetOnKeyUp(func(ke *fyne.KeyEvent) {
				defer c.SetOnKeyUp(nil)

				p.runner.Preferences().SetString(pref, string(ke.Name))
				b.SetText(string(ke.Name))

				picker.Hide()
			})

			picker.Show()
		},
	)

	return l, b
}

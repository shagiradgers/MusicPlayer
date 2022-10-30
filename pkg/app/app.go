package app

import (
	"MusicPlayer/pkg/player"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var Player player.Player
var PlayPath string

type App struct {
	fyne.App
	fyne.Window
	fyne.Size
}

func NewApp() (*App, error) {
	a := &App{
		App: app.New(),
	}
	a.Window = a.App.NewWindow("Music player")
	a.Size = fyne.Size{
		Width:  350,
		Height: 200,
	}

	a.Window.Resize(a.Size)
	a.Window.CenterOnScreen()

	err := a.initInterface()
	return a, err
}

func (a *App) initInterface() error {
	if a.Window == nil {
		return fmt.Errorf("nil Window")
	}

	// volume slider
	volumeSlider := widget.NewSlider(0, 100)
	volumeSlider.SetValue(100)
	volumeSlider.OnChanged = func(volume float64) {
		if Player != nil {
			// у setVolume диапазон от 0 до 1, а у слайдер от 0 до 100
			Player.SetVolume(volumeSlider.Value / 100)
		}
	}

	// volume label
	volumeLabel := widget.NewLabel("Volume: ")

	// dialog window to choose path
	dialogToChoosePath := dialog.NewFileOpen(
		func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowInformation("Error", err.Error(), a.Window)
			}

			if reader != nil {
				PlayPath = reader.URI().Path()
			}
			if reader == nil {
				dialog.ShowInformation("Info", "canceled", a.Window)
			}
			a.Window.Resize(a.Size)
		},
		a.Window,
	)

	// button, that call dialog window to choose path
	chooseFilePathBtn := widget.NewButton("Chose Music",
		func() {
			var resizeCoefficient float32 = 200
			a.Window.Resize(fyne.Size{
				Width:  a.Size.Width + resizeCoefficient,
				Height: a.Size.Height + resizeCoefficient,
			})
			dialogToChoosePath.Show()
		},
	)

	// play button
	playBtn := widget.NewButton("Play",
		func() {
			var err error

			if PlayPath == "" {
				return
			}

			// проверка на то, что музыка может еще играть
			if Player != nil {
				err = Player.Close()
				if err != nil {
					dialog.ShowInformation("Error", err.Error(), a.Window)
					return
				}
			}

			Player, err = player.Play(PlayPath)
			if err != nil {
				dialog.ShowInformation("Error", err.Error(), a.Window)
				return
			}
			// у setVolume диапазон от 0 до 1, а у слайдер от 0 до 100
			Player.SetVolume(volumeSlider.Value / 100)

			Player.Play()
		},
	)

	// pause button
	pauseResumeBtn := widget.NewButton("Pause/Resume",
		func() {
			if PlayPath == "" {
				return
			}

			if Player == nil {
				return
			}

			if Player.IsPlaying() {
				Player.Pause()
			} else {
				Player.Play()
			}
		},
	)

	// stop music button
	stopBtn := widget.NewButton("Stop",
		func() {
			if Player != nil {
				err := Player.Close()
				if err != nil {
					dialog.ShowInformation("Error", err.Error(), a.Window)
				}
			}
		},
	)

	volumeContainer := container.NewGridWithColumns(2, volumeLabel, volumeSlider)
	btnContainer := container.NewGridWithColumns(3, playBtn, pauseResumeBtn, stopBtn)
	contentContainer := container.NewVBox(chooseFilePathBtn, btnContainer, volumeContainer)

	a.Window.SetContent(contentContainer)

	return nil
}

func (a *App) Run() {
	a.Window.ShowAndRun()
}

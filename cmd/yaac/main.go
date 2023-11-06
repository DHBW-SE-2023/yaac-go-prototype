package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	yaac_io "github.com/DHBW-SE-2023/yaac-go-prototype/pkg"
)

const APP_NAME = "YAAC-Go-Prototype"

var App fyne.App
var mainWindow fyne.Window

func main() {
	App = app.NewWithID(APP_NAME)

	// setuping window
	mainWindow = App.NewWindow(APP_NAME)

	// set icon
	r, _ := yaac_io.LoadResourceFromPath("./Icon.png")
	mainWindow.SetIcon(r)

	// setup systray
	if desk, ok := App.(desktop.App); ok {
		m := fyne.NewMenu(APP_NAME,
			fyne.NewMenuItem("Show", func() {
				mainWindow.Show()
			}))
		desk.SetSystemTrayMenu(m)
		desk.SetSystemTrayIcon(r)
	}
	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})

	// handle main window
	mainWindow.SetContent(makeUI_w1())
	mainWindow.Show()

	App.Run()
}

func makeUI_w1() *widget.Label {
	clock := widget.NewLabel("Hi!")
	return clock
}

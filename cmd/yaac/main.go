package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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
	mainWindow.SetMaster()

	// handle window 2
	w2 := App.NewWindow("Larger")
	w2.SetContent(makeUI_w2())
	w2.Resize(fyne.NewSize(250, 100))
	w2.Show()

	App.Run()
}

func updateTime(clock *widget.Label) {
	formatted := time.Now().Format("Current system time: 15:04:05")
	clock.SetText(formatted)
}

func makeUI_w1() *widget.Label {
	clock := widget.NewLabel("")
	updateTime(clock)
	return clock
}

func makeUI_w2() *widget.Button {
	return widget.NewButton("Open new", func() {
		w3 := App.NewWindow("Third")
		w3.SetContent(container.NewVBox(makeUI_w3()))
		w3.Resize(fyne.NewSize(200, 50))
		w3.Show()
	})
}

func makeUI_w3() (*widget.Label, *widget.Entry) {
	out_label := widget.NewLabel("Hello World!")
	in_entry := widget.NewEntry()

	in_entry.OnChanged = func(s string) {
		out_label.SetText(s + " to you as well!")
	}

	return out_label, in_entry
}

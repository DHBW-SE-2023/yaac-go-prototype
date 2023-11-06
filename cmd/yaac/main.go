package main

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const APP_NAME = "YAAC-Go-Prototype"

var App fyne.App
var mainWindow fyne.Window

func main() {
	// setuping window
	mainWindow = App.NewWindow(APP_NAME)

	// set icon
	r, _ := LoadResourceFromPath("./Icon.png")
	mainWindow.SetIcon(r)
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

func makeUI_w2(a fyne.App) *widget.Button {
	return widget.NewButton("Open new", func() {
		w3 := a.NewWindow("Third")
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

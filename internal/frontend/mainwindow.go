package yaac_frontend

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	yaac_consts "github.com/DHBW-SE-2023/yaac-go-prototype/internal/consts"
	resource "github.com/DHBW-SE-2023/yaac-go-prototype/pkg/resource_manager"
)

var App fyne.App
var mainWindow fyne.Window

func (f *Frontend) OpenMainWindow() {
	App = app.NewWithID(yaac_consts.APP_NAME)

	// setuping window
	mainWindow = App.NewWindow(yaac_consts.APP_NAME)

	// set icon
	r, _ := resource.LoadResourceFromPath("./Icon.png")
	mainWindow.SetIcon(r)

	// setup systray
	if desk, ok := App.(desktop.App); ok {
		m := fyne.NewMenu(yaac_consts.APP_NAME,
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
	mainWindow.SetContent(makeMainWindow(f))
	mainWindow.Show()

	App.Run()
}

func makeMainWindow(f *Frontend) *fyne.Container {
	header := widget.NewLabel("Select an action:")
	mail_button := widget.NewButton("Open Mail Window", f.OpenMailWindow)

	return container.NewVBox(
		header,
		mail_button,
	)
}

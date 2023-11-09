package yaac_frontend

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	resource "github.com/DHBW-SE-2023/yaac-go-prototype/pkg/resource_manager"
)

var opencvDemoWindow fyne.Window

func (f *Frontend) OpenOpencvDemoWindow() {
	// setuping window
	opencvDemoWindow = App.NewWindow("OpenCV Demo")

	// set icon
	r, _ := resource.LoadResourceFromPath("./Icon.png")
	opencvDemoWindow.SetIcon(r)

	// handle main window
	opencvDemoWindow.SetContent(makeOpencvDemoWindow(f))
	opencvDemoWindow.Show()

	App.Run()
}

func makeOpencvDemoWindow(f *Frontend) *fyne.Container {
	header := widget.NewLabel("Select an action:")
	mail_button := widget.NewButton("Open Mail Window", f.OpenMailWindow)

	return container.NewVBox(
		header,
		mail_button,
	)
}

package yaac_frontend

import (
	"io"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	resource "github.com/DHBW-SE-2023/yaac-go-prototype/pkg/resource_manager"
)

type OpencvDemoWindow struct {
	Window fyne.Window
	Image  *canvas.Image
}

var opencvDemoWindow OpencvDemoWindow

func (f *Frontend) OpenOpencvDemoWindow() {
	opencvDemoWindow = OpencvDemoWindow{}

	// setuping window
	opencvDemoWindow.Window = App.NewWindow("OpenCV Demo")

	// set icon
	r, _ := resource.LoadResourceFromPath("./Icon.png")
	opencvDemoWindow.Window.SetIcon(r)

	// handle main window
	opencvDemoWindow.Window.SetContent(makeOpencvDemoWindow(f))
	opencvDemoWindow.Window.Resize(fyne.NewSize(800, 600))
	opencvDemoWindow.Window.Show()

	App.Run()
}

func makeOpencvDemoWindow(f *Frontend) *fyne.Container {
	header := widget.NewLabel("Please select an Input image:")
	//openFile := widget.NewButton("File Open", openOpencvDemoWindowFileDialog)
	openFile := widget.NewButton("File Open With Filter (.jpg or .png)", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, opencvDemoWindow.Window)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}

			showImage(reader)
		}, opencvDemoWindow.Window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})
	opencvDemoWindow.Image = canvas.NewImageFromResource(theme.FyneLogo())
	input_image := container.NewScroll(opencvDemoWindow.Image)
	//input_image.Content
	//input_image.Resize(input_image.Size())
	//input_image.Resize(fyne.NewSize(800, 800))
	input_image_container := container.NewCenter(input_image)
	input_image_container.Resize(input_image.Size())

	return container.NewVBox(
		header,
		openFile,
		input_image_container,
	)
}

func loadImage(f fyne.URIReadCloser) *canvas.Image {
	data, err := io.ReadAll(f)
	if err != nil {
		fyne.LogError("Failed to load image data", err)
		return nil
	}
	res := fyne.NewStaticResource(f.URI().Name(), data)

	return canvas.NewImageFromResource(res)
}

func showImage(f fyne.URIReadCloser) {
	if f == nil {
		log.Println("Cancelled")
		return
	}
	defer f.Close()
	img := loadImage(f)
	if img == nil {
		return
	}
	img.FillMode = canvas.ImageFillOriginal

	//w := fyne.CurrentApp().NewWindow(f.URI().Name())
	//w.SetContent(container.NewScroll(img))
	//w.Resize(fyne.NewSize(320, 240))
	//w.Show()

}

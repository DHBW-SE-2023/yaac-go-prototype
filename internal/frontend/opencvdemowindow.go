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
	Window              fyne.Window
	Image               *canvas.Image
	InputImageContainer *fyne.Container
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

			showImage(reader, opencvDemoWindow.InputImageContainer)
		}, opencvDemoWindow.Window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})
	opencvDemoWindow.Image = canvas.NewImageFromResource(theme.FyneLogo())
	inputImage := container.NewScroll(opencvDemoWindow.Image) //previously: container.NewScroll(opencvDemoWindow.Image) - somehow works both (NewCenter)
	opencvDemoWindow.InputImageContainer = container.NewCenter(inputImage)
	opencvDemoWindow.InputImageContainer.Resize(inputImage.Size())

	showFile := widget.NewButton("Test 12345 Lebkuchen", func() {
		log.Println("tapped")
	})

	return container.NewVBox(
		header,
		openFile,
		showFile,
		opencvDemoWindow.InputImageContainer,
	)
}

func loadImage(f fyne.URIReadCloser) *canvas.Image {
	data, err := io.ReadAll(f)
	if err != nil {
		fyne.LogError("Error at loading file", err)
		return nil
	}
	res := fyne.NewStaticResource(f.URI().Name(), data)

	img := canvas.NewImageFromResource(res)
	if img == nil {
		fyne.LogError("Error at creating file object", err)
		return nil
	}

	// Calculate the desired size for the loaded image based on a fixed width
	fixedWidth := 400
	aspectRatio := float32(img.MinSize().Height) / float32(img.MinSize().Width)
	fixedHeight := int(float32(fixedWidth) * aspectRatio)

	// Set the calculated size directly on the image
	img.Resize(fyne.NewSize(float32(fixedWidth), float32(fixedHeight)))

	return img
}

func showImage(f fyne.URIReadCloser, imgContainer *fyne.Container) {
	if f == nil {
		log.Println("Cancelled")
		return
	}
	defer f.Close()
	img := loadImage(f)
	if img == nil {
		log.Println("Error at loading image")
		return
	}

	img.FillMode = canvas.ImageFillContain //.ImageFillOriginal

	// Create a container with dynamic sizing
	containerWithDynamicSizing := fyne.NewContainer(img)

	// Set the content of the main container to the new container with dynamic sizing
	imgContainer.Objects = []fyne.CanvasObject{containerWithDynamicSizing}

	// Actualize and show window
	opencvDemoWindow.Window.Content().Refresh()
	opencvDemoWindow.Window.RequestFocus()
	opencvDemoWindow.Window.Show()
}

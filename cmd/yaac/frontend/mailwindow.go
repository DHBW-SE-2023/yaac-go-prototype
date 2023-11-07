package yaac_frontend

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	yaac_consts "github.com/DHBW-SE-2023/yaac-go-prototype/cmd/yaac/consts"
	yaac_mvvm "github.com/DHBW-SE-2023/yaac-go-prototype/cmd/yaac/mvvm"
	resource "github.com/DHBW-SE-2023/yaac-go-prototype/pkg/resource_manager"
)

var mailWindow fyne.Window
var result_label *widget.Label

type EmailData struct {
	Name, Email string
}

func OpenMailWindow() {
	mailWindow = App.NewWindow(yaac_consts.APP_NAME)

	// set icon
	r, _ := resource.LoadResourceFromPath("./Icon.png")
	mailWindow.SetIcon(r)

	mailWindow.SetContent(makeMailWindow())
	mailWindow.Show()
}

func makeMailWindow() *fyne.Container {
	top_label := widget.NewLabel("Please enter your credentials:")

	formStruct := yaac_mvvm.EmailData{}

	formData := binding.BindStruct(&formStruct)
	form := newFormWithData(formData)

	form.OnSubmit = func() {
		yaac_mvvm.MailFormUpdated(formStruct)
	}

	result_label = widget.NewLabel("")

	return container.NewVBox(
		top_label,
		form,
		result_label,
	)
}

func UpdateResultLabel(content string) {
	result_label.Text = content
}

func newFormWithData(data binding.DataMap) *widget.Form {
	keys := data.Keys()
	items := make([]*widget.FormItem, len(keys))
	for i, k := range keys {
		data, err := data.GetItem(k)
		if err != nil {
			items[i] = widget.NewFormItem(k, widget.NewLabel(err.Error()))
		}
		items[i] = widget.NewFormItem(k, createBoundItem(data))
	}

	return widget.NewForm(items...)
}

func createBoundItem(v binding.DataItem) fyne.CanvasObject {
	switch val := v.(type) {
	case binding.Bool:
		return widget.NewCheckWithData("", val)
	case binding.Float:
		s := widget.NewSliderWithData(0, 1, val)
		s.Step = 0.01
		return s
	case binding.Int:
		return widget.NewEntryWithData(binding.IntToString(val))
	case binding.String:
		return widget.NewEntryWithData(val)
	default:
		return widget.NewLabel("")
	}
}

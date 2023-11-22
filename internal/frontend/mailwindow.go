// yaac_frontend is the package containing all frontend functionality for the yaac protoype
package yaac_frontend

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	yaac_shared "github.com/DHBW-SE-2023/yaac-go-prototype/internal/shared"
	resource "github.com/DHBW-SE-2023/yaac-go-prototype/pkg/resource_manager"
)

// mailWindow defines the window shown to the user
var mailWindow fyne.Window

// result_label defines the text field for the output
var result_label *widget.Label

// OpenMailWindow is the public function used to create the mail window
func (f *Frontend) OpenMailWindow() {
	mailWindow = App.NewWindow("Mail Demo")

	// set icon
	r, _ := resource.LoadResourceFromPath("./Icon.png")
	mailWindow.SetIcon(r)

	mailWindow.SetContent(makeFormTab(mailWindow, f))
	mailWindow.Show()
}

// UpdateResultLabel is the public function to update the output shown in the mail window
func (f *Frontend) UpdateResultLabel(content string) {
	result_label.SetText(content)
	fyne.CurrentApp().SendNotification(&fyne.Notification{
		Title: content,
	})
}

// makeFormTab is the private helper function defining the mail windows form
func makeFormTab(_ fyne.Window, f *Frontend) fyne.CanvasObject {
	// mailServer is used to collect the mail servers address as user input
	mailServer := widget.NewEntry()
	mailServer.SetPlaceHolder("John Smith")

	// email is used to collect the user accounts email as input
	email := widget.NewEntry()
	email.SetPlaceHolder("test@example.com")
	email.Validator = validation.NewRegexp(`\w{1,}@\w{1,}\.\w{1,4}`, "not a valid email")

	// password is used to collect the user accounts password as input
	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	result_label = widget.NewLabel("")

	// form adds the dynamic functionality to the mail input form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "MailServer", Widget: mailServer, HintText: "Specify the address of your mail server with ':port'"},
			{Text: "Email", Widget: email, HintText: "Your email address"},
		},
		OnCancel: func() {
			fmt.Println("Cancelled")
		},
		OnSubmit: func() {
			// formStruct holds the collected credentials to log in to the server
			formStruct := yaac_shared.EmailData{
				MailServer: mailServer.Text,
				Email:      email.Text,
				Password:   password.Text,
			}
			f.MVVM.MailFormUpdated(formStruct)
		},
	}
	form.Append("Password", password)
	form.Append("Your first unread message:", result_label)

	return form
}

package yaac_mvvm

import (
	yaac_backend "github.com/DHBW-SE-2023/yaac-go-prototype/cmd/yaac/backend"
)

type EmailData struct {
	Email, Password string
}

func MailFormUpdated(data EmailData) {
	response := yaac_backend.GetResponse(data.Email)
}

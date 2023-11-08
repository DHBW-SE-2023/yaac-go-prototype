package yaac_mvvm

import (
	yaac_backend "github.com/DHBW-SE-2023/yaac-go-prototype/internal/backend"
	yaac_frontend "github.com/DHBW-SE-2023/yaac-go-prototype/internal/frontend"
)

func (m *MVVM) MailFormUpdated(data yaac_frontend.EmailData) {
	var frontend = yaac_frontend.New(m)
	var backend = yaac_backend.New(m)

	resp := backend.GetResponse(data.Email)
	frontend.UpdateResultLabel(resp)
}

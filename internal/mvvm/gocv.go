package yaac_mvvm

import yaac_backend "github.com/DHBW-SE-2023/yaac-go-prototype/internal/backend"

func (m *MVVM) StartGoCV() {
	backend := yaac_backend.New(m)
	backend.StartGoCV()
}

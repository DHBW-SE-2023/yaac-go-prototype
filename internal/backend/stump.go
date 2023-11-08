package yaac_backend

import (
	"fmt"

	yaac_shared "github.com/DHBW-SE-2023/yaac-go-prototype/internal/shared"
)

func (b *Backend) GetResponse(input yaac_shared.EmailData) string {
	return fmt.Sprintf("Hello %s!", input.Email)
}

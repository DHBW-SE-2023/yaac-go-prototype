package yaac_backend

import "fmt"

func (b *Backend) GetResponse(input string) string {
	return fmt.Sprintf("Hello %s!", input)
}

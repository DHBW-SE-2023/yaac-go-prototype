package yaac_backend

import "fmt"

func GetResponse(input string) string {
	return fmt.Sprintf("Hello %s!", input)
}

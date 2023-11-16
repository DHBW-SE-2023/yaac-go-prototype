package yaac_mvvm

import (
	"fmt"

	yaac_backend "github.com/DHBW-SE-2023/yaac-go-prototype/internal/backend"
)

func (m *MVVM) StartGoCV() {
	backend := yaac_backend.New(m)
	var msg string
	var suc bool
	ch := make(chan int)
	go func() {
		msg, suc = backend.StartGoCV("./assets/list.jpg", ch)
	}()

	for elem := range ch {
		fmt.Printf("Progress: %d\n", elem)
	}
	fmt.Printf("Done!\n\tSuccess: %s\n\tMessage: %s\n", fmt.Sprint(suc), msg)
}

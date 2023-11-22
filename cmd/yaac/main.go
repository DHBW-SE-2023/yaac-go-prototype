// Main package for the yaac go prototype
package main

import yaac_mvvm "github.com/DHBW-SE-2023/yaac-go-prototype/internal/mvvm"

// Main function for the yaac go prototype
func main() {
	mvvm := yaac_mvvm.New()
	mvvm.OpenMainWindow()
}

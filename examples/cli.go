package examples

import (
	"fmt"
	"os"
)

func GetAddress() string {
	// Get address of device
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run about.go <address>")
		os.Exit(1)
	}

	return os.Args[1]
}

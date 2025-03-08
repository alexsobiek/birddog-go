package main

import (
	"fmt"
	"time"

	"github.com/alexsobiek/birddog-go"
	"github.com/alexsobiek/birddog-go/examples"
)

func main() {
	dev := birddog.NewDevice(birddog.NewAPI(examples.GetAddress()))

	dev.OnAvailable = func() {
		fmt.Println("Device is available!")

	}

	dev.OnUnavailable = func() {
		fmt.Println("Device is no longer available!")
	}

	dev.OnError = func(err error) {
		fmt.Println("Error:", err)
	}

	dev.Query(5 * time.Second) // query every 5 seconds, indefinitely
}

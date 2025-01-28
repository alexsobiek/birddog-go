package main

import (
	"fmt"

	"github.com/alexsobiek/birddog-go"
	"github.com/alexsobiek/birddog-go/examples"
)

func main() {
	api := birddog.NewAPI(examples.GetAddress())

	about, err := api.About()

	if err != nil {
		panic(err)
	}

	// Print the about struct
	fmt.Println(about)
}

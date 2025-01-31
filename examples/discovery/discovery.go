package main

import (
	"fmt"

	"github.com/alexsobiek/birddog-go"
)

func main() {
	discovery := birddog.NewBirddogDiscovery(&birddog.DiscoveryNetOptions{})

	discovery.OnDiscover(func(api *birddog.API) {
		about, err := api.About()
		if err != nil {
			panic(err)
		}

		// Print the about struct
		fmt.Println("Found Birddog device:", about.HostName, "at", about.IPAddress)
	})

	discovery.OnLost(func(api *birddog.API) {
		fmt.Println("Birddog device no longer available:", api.Host)
	})

	discovery.Find()
}

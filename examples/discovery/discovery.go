package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/alexsobiek/birddog-go"
)

func main() {
	discovery := birddog.NewBirddogDiscovery(&birddog.DiscoveryNetOptions{
		Interval: time.Duration(5) * time.Second, // Search every 5 seconds (devices inactive for 2*duration are considered "lost")
		Logger:   log.New(ioutil.Discard, "", 0), // discard logs
	})

	discovery.OnDiscover(func(api *birddog.API) {
		fmt.Println("Birddog device discovered:", api.Host)
		about, err := api.About()
		if err != nil {
			fmt.Println("Error getting about struct (device may not be ready):", err)
			return
		}

		// Print the about struct
		fmt.Println("About device:", about.HostName, "at", about.IPAddress)

	})

	discovery.OnLost(func(api *birddog.API) {
		fmt.Println("Birddog device no longer available: ", api.Host)
	})

	discovery.Find()
}

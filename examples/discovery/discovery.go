package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/alexsobiek/birddog-go"
)

func main() {
	// query interfaces

	ifaces, err := net.Interfaces()

	if err != nil {
		fmt.Println("Error getting interfaces:", err)
		return
	}

	targetInterfaceName := "eno2"

	var iface *net.Interface

	for _, i := range ifaces {
		if i.Name == targetInterfaceName {
			fmt.Println("Found interface:", i.Name)
			iface = &i
		}
	}

	logger := log.New(os.Stdout, "", 0)

	discovery := birddog.NewBirddogDiscovery(&birddog.DiscoveryNetOptions{
		Interval:  time.Duration(5) * time.Second, // Search every 5 seconds (devices inactive for 2*duration are considered "lost")
		Logger:    logger,
		Interface: iface, // Use default interface
		Domain: "",
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

	discovery.OnTimeout(func(api *birddog.API) {
		fmt.Println("Birddog device no longer available: ", api.Host)
	})

	discovery.Find()
}

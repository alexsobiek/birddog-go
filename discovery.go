package birddog

import (
	"log"
	"net"
	"time"

	"github.com/hashicorp/mdns"
)

type DiscoveryNetOptions struct {
	Interface *net.Interface
	Logger    *log.Logger
	Interval  time.Duration
	Domain    string
}

type BirddogDiscovery struct {
	opts       *DiscoveryNetOptions
	discovered map[string]*Device
	close      chan struct{}
	active     bool

	onDiscover func(*API)
	onTimeout  func(*API)
}

func NewBirddogDiscovery(opts *DiscoveryNetOptions) *BirddogDiscovery {
	if opts == nil {
		opts = &DiscoveryNetOptions{}
	}

	if opts.Interval == 0 {
		opts.Interval = 5 * time.Second
	}

	return &BirddogDiscovery{
		opts:       opts,
		discovered: make(map[string]*Device),
		active:     true,
		close:      make(chan struct{}),
	}
}

func (b *BirddogDiscovery) Find() {
	entriesCh := make(chan *mdns.ServiceEntry)
	params := mdns.QueryParam{
		Service:     "_birddog._tcp",
		Entries:     entriesCh,
		Interface:   b.opts.Interface,
		Domain:      b.opts.Domain,
		DisableIPv6: true,
		Logger:      b.opts.Logger,
	}

	// Producer goroutine

	go func() {
		for {
			if !b.active {
				return
			}
			err := mdns.Query(&params)
			if err != nil {
				log.Fatalf("mDNS Query error: %v", err)
			}
			time.Sleep(b.opts.Interval)
		}
	}()

	// Consumer loop
	for {
		select {
		case <-b.close:
			b.active = false
			return
		case next := <-entriesCh:
			if next != nil {
				addr := next.AddrV4.String()

				// attempt to ping device

				if _, ok := b.discovered[addr]; ok {
					continue
				}

				dev := NewDevice(NewAPI(next.AddrV4.String()))

				dev.OnAvailable = func() {
					b.onDiscover(dev.api)
				}

				dev.OnUnavailable = func() {
					b.onTimeout(dev.api)
				}

				dev.OnError = func(err error) {
					b.opts.Logger.Printf("Error: %v", err)
				}

				b.discovered[next.AddrV4.String()] = dev

				go dev.Query(b.opts.Interval)
			}
		}
	}
}

func (b *BirddogDiscovery) OnDiscover(f func(*API)) {
	b.onDiscover = f
}

func (b *BirddogDiscovery) OnTimeout(f func(*API)) {
	b.onTimeout = f
}

func (b *BirddogDiscovery) Close() {
	close(b.close)
	for _, dev := range b.discovered {
		dev.Close()
	}
}

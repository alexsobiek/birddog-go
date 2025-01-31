package birddog

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/hashicorp/mdns"
)

type DiscoveryNetOptions struct {
	Interface *net.Interface
	Logger    *log.Logger
	Interval  time.Duration
	Domain    string
}

type DiscoveredDevice struct {
	api      *API
	lastSeen time.Time
}

type BirddogDiscovery struct {
	opts       *DiscoveryNetOptions
	discovered map[string]*DiscoveredDevice
	close      chan struct{}
	active     bool
	mu         sync.Mutex

	onDiscover func(*API)
	onLost     func(*API)
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
		discovered: make(map[string]*DiscoveredDevice),
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

			// filter out devices that haven't been seen in a while
			b.mu.Lock()
			for host, dev := range b.discovered {
				if time.Since(dev.lastSeen) > 2*b.opts.Interval {
					delete(b.discovered, host)
					b.onLost(dev.api)
				}
			}
			b.mu.Unlock()
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
				b.mu.Lock()
				addr := next.AddrV4.String()

				if _, ok := b.discovered[addr]; ok {
					b.discovered[addr].lastSeen = time.Now()
					b.mu.Unlock()
					continue
				}
				dev := &DiscoveredDevice{
					api:      NewAPI(addr),
					lastSeen: time.Now(),
				}

				b.discovered[next.AddrV4.String()] = dev
				b.mu.Unlock()
				b.onDiscover(dev.api)
			}
		}
	}
}

func (b *BirddogDiscovery) OnDiscover(f func(*API)) {
	b.onDiscover = f
}

func (b *BirddogDiscovery) OnLost(f func(*API)) {
	b.onLost = f
}

func (b *BirddogDiscovery) Close() {
	close(b.close)
}

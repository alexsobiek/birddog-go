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
}

type BirddogDiscovery struct {
	opts       *DiscoveryNetOptions
	discovered map[string]*API
	mu         sync.RWMutex
	close      chan struct{}

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
		discovered: make(map[string]*API),
	}
}

func (b *BirddogDiscovery) Find() {
	entriesCh := make(chan *mdns.ServiceEntry)
	params := &mdns.QueryParam{
		Service:     "_birddog._tcp",
		Entries:     entriesCh,
		DisableIPv6: true,
	}

	for {
		select {
		case <-b.close:
			return
		default:
			err := mdns.Query(params)
			if err != nil {
				panic(err)
			}

			time.Sleep(b.opts.Interval)

			services := make(map[string]*API)

			for i := 0; i < len(entriesCh); i++ {
				entry := <-entriesCh
				services[entry.Host] = NewAPI(entry.Host)
			}

			b.mu.Lock()
			// Find new services
			for host, api := range services {
				if _, ok := b.discovered[host]; !ok {
					b.discovered[host] = api
					b.onDiscover(api)
				}
			}

			// Find lost services
			for host, api := range b.discovered {
				if _, ok := services[host]; !ok {
					delete(b.discovered, host)
					b.onLost(api)
				}
			}

			b.mu.Unlock()
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

package birddog

import (
	"time"
)

type Device struct {
	api           *API
	OnError       func(error)
	OnAvailable   func()
	OnUnavailable func()
	online        bool
	close         chan struct{}
}

func NewDevice(api *API) *Device {
	dev := &Device{
		api:   api,
		close: make(chan struct{}),
	}

	return dev
}

func (d *Device) Query(interval time.Duration) {
	for {
		select {
		case <-d.close:
			return

		default:
			a, err := d.api.About()

			if err != nil {
				d.OnError(err)

				if d.online {
					d.online = false
					d.OnUnavailable()
				}

				continue
			}

			if a.Status != StatusOffline {
				if !d.online {
					d.online = true
					d.OnAvailable()
				}
			} else {
				if d.online {
					d.online = false
					d.OnUnavailable()
				}
			}

			time.Sleep(interval)
		}
	}
}

func (d *Device) API() *API {
	return d.api
}

func (d *Device) IsOnline() bool {
	return d.online
}

func (d *Device) Close() {
	close(d.close)
	d.OnAvailable()
}

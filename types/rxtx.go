package types

type TransportMode string

const (
	TransportModeMulticast TransportMode = "Multicast"
	TransportModeMultiTCP  TransportMode = "Multi-TCP"
	TransportModeTCP       TransportMode = "TCP"
	TransportModeUDP       TransportMode = "UDP"
)

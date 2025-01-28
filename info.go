package birddog

type Status string

const (
	StatusActivating         Status = "activating"
	StatusActive             Status = "active"
	StatusOnline             Status = "online"              // untested
	StatusOffline            Status = "offline"             // untested
	StatusCameraInitializing Status = "camera initializing" // untested
	StatusNoVideo            Status = "no video"            // untested
)

type NetworkConfigMethod string

const (
	NetworkConfigMethodDHCP   NetworkConfigMethod = "dhcp"
	NetworkConfigMethodStatic NetworkConfigMethod = "static"
)

type About struct {
	FallbackIP             string              `json:"FallbackIP"`
	FirmwareVersion        string              `json:"FirmwareVersion"`
	Format                 string              `json:"Format"`
	GateWay                string              `json:"GateWay"`
	HardwareVersion        string              `json:"HardwareVersion"`
	HostName               string              `json:"HostName"`
	IPAddress              string              `json:"IPAddress"`
	MCUVersion             string              `json:"MCUVersion"`
	NetworkingConfigMethod NetworkConfigMethod `json:"NetworkingConfigMethod"`
	NetworkMask            string              `json:"NetworkMask"`
	SerialNumber           string              `json:"SerialNumber"`
	Status                 Status              `json:"Status"`
}

func (a *API) About() (*About, error) {
	var about About
	_, err := a.get("about", &about)
	if err != nil {
		return nil, err
	}
	return &about, nil
}

func (a *API) Restart() error {
	_, err := a.post("restart", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

func (a *API) Reboot() error {
	_, err :=a. post("reboot", nil, nil)
	if err != nil {
		return err
	}
	return nil
}

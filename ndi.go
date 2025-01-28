package birddog

import "github.com/alexsobiek/birddog-go/types"

func (a *API) NDIList() (map[string]string, error) {
	var ndiList map[string]string
	_, err := a.get("list", &ndiList)
	if err != nil {
		return nil, err
	}

	return ndiList, nil
}

type NDIDiscoveryServerSettings struct {
	Status types.NDIDiscoveryServerStatus `json:"NDIDisServ"`
	IP     string                         `json:"NDIDisServIP"`
}

// untested
func (a *API) GetNDIDiscoveryServerSettings() (*NDIDiscoveryServerSettings, error) {
	var settings NDIDiscoveryServerSettings
	_, err := a.get("NDIDisServer", &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

// untested
func (a *API) SetNDIDiscoveryServerSettings(settings NDIDiscoveryServerSettings) (*NDIDiscoveryServerSettings, error) {
	var res NDIDiscoveryServerSettings
	_, err := a.post("NDIDisServer", &settings, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// untested
func (a *API) GetNDIGroupName() (string, error) {
	var groupName string
	_, err := a.get("NDIGrpName", &groupName)
	if err != nil {
		return "", err
	}

	return groupName, nil
}

// untested
func (a *API) SetNDIGroupName(groupName string) (string, error) {
	var res string
	_, err := a.post("NDIGrpName", groupName, &res)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (a *API) RefreshNDISources() error {
	_, err := a.post("refresh", nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (a *API) ResetNDISources() error {
	_, err := a.post("reset", nil, nil)
	if err != nil {
		return err
	}

	return nil
}

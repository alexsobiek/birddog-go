package birddog

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/alexsobiek/birddog-go/types"
)

type CaptureOpts struct {
	types.ChNum
	Status types.OperationMode `json:"status"`
}

func (a *API) Capture(opts CaptureOpts) (string, error) {
	var res string

	_, err := a.get("capture?ChNum="+opts.ChNum.String()+"&status="+string(opts.Status), &res)

	if err != nil {
		return "", err
	}
	return res, nil
}

type Connect struct {
	types.ChNum
	SourceName string              `json:"sourceName"`
	Status     types.OperationMode `json:"status,omitempty"`
}

func (a *API) Connect(opts Connect) (*Connect, error) {
	var res Connect

	_, err := a.post("connectTo", &opts, &res)

	if err != nil {
		return nil, err
	}
	return &res, nil
}

type DecodeTransport struct {
	Rxpm types.TransportMode `json:"Rxpm"`
}

func (a *API) GetDecodeTransport() (*DecodeTransport, error) {
	var res DecodeTransport

	_, err := a.get("decodeTransport", &res)

	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (a *API) SetDecodeTransport(transport DecodeTransport) (*DecodeTransport, error) {
	var res DecodeTransport
	_, err := a.post("decodeTransport", &transport, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type DecodeSettings struct {
	types.ChNum
	ColorSpace      types.ColorSpace      `json:"ColorSpace"`
	TallyMode       types.TallyMode       `json:"TallyMode"`
	ScreenSaverMode types.ScreenSaverMode `json:"ScreenSaverMode"`
	NDIAudio        types.NDIAudio        `json:"NDIAudio"`
}

func (a *API) GetDecodeSettings() (*DecodeSettings, error) {
	var res DecodeSettings

	_, err := a.get("decodesetup", &res)

	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (a *API) SetDecodeSettings(settings DecodeSettings) (*DecodeSettings, error) {
	var res DecodeSettings
	_, err := a.post("decodesetup", &settings, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type DecodeStatus struct {
	VideoResolution string `json:"Videoresolution"`
	VideoFrameRate  string `json:"VideoFramerate"`
	VideoSampleRate string `json:"VideoSamplerate"`
	AudioChannels   int    `json:"Audiochannels"`
	AudioSampleRate int    `json:"Audiosamplerate"`
	AverateBitrate  int    `json:"AverateBitrate"`
}

func (d *DecodeStatus) UnmarshalJSON(data []byte) error {
	type Alias DecodeStatus
	aux := struct {
		AudioChannels   string `json:"Audiochannels"`
		AudioSampleRate string `json:"Audiosamplerate"`
		AverateBitrate  string `json:"AverateBitrate"`
		*Alias
	}{
		Alias: (*Alias)(d),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error

	if aux.AudioChannels == "" {
		aux.AudioChannels = "0"
	}

	if aux.AudioSampleRate == "" {
		aux.AudioSampleRate = "0"
	}

	if aux.AverateBitrate == "" {
		aux.AverateBitrate = "0"
	}

	if d.AudioChannels, err = strconv.Atoi(aux.AudioChannels); err != nil {
		return err
	}

	if d.AudioSampleRate, err = strconv.Atoi(aux.AudioSampleRate); err != nil {
		return err
	}

	if d.AverateBitrate, err = strconv.Atoi(aux.AverateBitrate); err != nil {
		return err
	}

	return nil
}

func (a *API) GetDecodeStatus(chnum types.ChNum) (*DecodeStatus, error) {
	var res DecodeStatus

	_, err := a.get(fmt.Sprintf("decodestatus?ChNum=%s", chnum.String()), &res)

	if err != nil {
		return nil, err
	}
	return &res, nil
}

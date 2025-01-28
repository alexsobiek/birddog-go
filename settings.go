package birddog

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/alexsobiek/birddog-go/types"
)

type AnalogAudio struct {
	AnalogAudioInGain       int                     `json:"AnalogAudioInGain"`
	AnalogAudioOutGain      int                     `json:"AnalogAudioOutGain"`
	AnalogAudiooutputselect types.AnalogAudioOutput `json:"AnalogAudiooutputselect"`
}

func (a *AnalogAudio) UnmarshalJSON(data []byte) error {
	type Alias AnalogAudio
	aux := struct {
		AnalogAudioInGain  string `json:"AnalogAudioInGain"`
		AnalogAudioOutGain string `json:"AnalogAudioOutGain"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var err error
	if a.AnalogAudioInGain, err = strconv.Atoi(aux.AnalogAudioInGain); err != nil {
		return fmt.Errorf("error parsing AnalogAudioInGain: %w", err)
	}
	if a.AnalogAudioOutGain, err = strconv.Atoi(aux.AnalogAudioOutGain); err != nil {
		return fmt.Errorf("error parsing AnalogAudioOutGain: %w", err)
	}

	return nil
}

func (a AnalogAudio) MarshalJSON() ([]byte, error) {
	type Alias AnalogAudio
	return json.Marshal(&struct {
		AnalogAudioInGain  string `json:"AnalogAudioInGain"`
		AnalogAudioOutGain string `json:"AnalogAudioOutGain"`
		*Alias
	}{
		AnalogAudioInGain:  strconv.Itoa(a.AnalogAudioInGain),
		AnalogAudioOutGain: strconv.Itoa(a.AnalogAudioOutGain),
		Alias:              (*Alias)(&a),
	})
}

func (a *API) GetAnalogAudio() (*AnalogAudio, error) {
	var analogAudio AnalogAudio
	_, err := a.get("analogaudiosetup", &analogAudio)
	if err != nil {
		return nil, err
	}
	return &analogAudio, nil
}

func (a *API) SetAnalogAudio(analogAudio *AnalogAudio) error {
	_, err := a.post("analogaudiosetup", analogAudio, nil)
	if err != nil {
		return err
	}
	return nil
}

// Returns the current operation mode of the device.
func (a *API) GetOperationMode() (types.OperationMode, error) {
	var operationMode types.OperationMode
	_, err := a.get("operationmode", &operationMode)
	if err != nil {
		return "", err
	}
	return operationMode, nil
}

// Attempts to set the operation mode of the device. Returns the new operation mode.
func (a *API) SetOperationMode(operationMode types.OperationMode) (types.OperationMode, error) {
	var newMode types.OperationMode
	_, err := a.post("operationmode", &struct {
		OperationMode types.OperationMode `json:"OperationMode"`
	}{
		OperationMode: operationMode,
	}, &newMode)
	if err != nil {
		return "", err
	}

	return newMode, nil
}

func (a *API) GetVideoOutput() (types.VideoOutput, error) {
	url := a.Host + "videooutputinterface"
	var videoOutput types.VideoOutput
	_, err := a.get(url, &videoOutput)
	if err != nil {
		return "", err
	}
	return videoOutput, nil
}

func (a *API) SetVideoOutput(videoOutput types.VideoOutput) (types.VideoOutput, error) {
	var newOutput types.VideoOutput
	_, err := a.post("videooutputinterface", &struct {
		VideoOutput types.VideoOutput `json:"VideoOutput"`
	}{
		VideoOutput: videoOutput,
	}, &newOutput)
	if err != nil {
		return "", err
	}

	return newOutput, nil
}

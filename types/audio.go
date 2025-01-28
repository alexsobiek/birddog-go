package types

type AnalogAudioOutput string

const (
	AnalogAudioOutputDecodeMain  AnalogAudioOutput = "DecodeMain"
	AnalogAudioOutputDecodeComms AnalogAudioOutput = "DecodeComms"
	AnalogAudioOutputDecodeLoop  AnalogAudioOutput = "DecodeLoop"
)

type NDIAudio string

const (
	NDIAudioEnabled  NDIAudio = "NDIAudioEn"
	NDIAudioDisabled NDIAudio = "NDIAudioDis"
)

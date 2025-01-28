package types

type VideoOutput string

const (
	VideoOutputSDI        VideoOutput = "sdi"
	VideoOutputHDMI       VideoOutput = "hdmi"
	VideoOutputLowLatency VideoOutput = "LowLatency"
	VideoOutputNormalMode VideoOutput = "NormalMode"
)

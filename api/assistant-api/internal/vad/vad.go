package internal_vad

// Activity represents a detected Audio segment.

type VADCallback func(*VadResult) error
type VadResult struct {
	StartSec float64
	EndSec   float64
}

func (a *VadResult) GetSpeechStartAt() float64 { return a.StartSec }
func (a *VadResult) GetSpeechEndAt() float64   { return a.EndSec }
func (a *VadResult) GetDuration() float64      { return a.EndSec - a.StartSec }

type Vad interface {
	Name() string
	Process(frame []byte) error
	Close() error
}

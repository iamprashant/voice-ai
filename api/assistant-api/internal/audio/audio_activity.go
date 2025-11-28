package internal_audio

type AudioActivitySegment struct {
	// The relative timestamp in seconds of when a speech segment begins.
	SpeechStartAt float64
	// The relative timestamp in seconds of when a speech segment ends.
	SpeechEndAt float64
	// The duration of the speech segment in seconds.
	Duration float64
	// The energy level of the speech segment, used to filter out noise.
	Energy float64
	// A confidence score indicating the likelihood of valid speech (0 to 1).
	Confidence float32
}

// GetSpeechStartAt returns the relative timestamp in seconds of when a speech segment begins.
func (d *AudioActivitySegment) GetSpeechStartAt() float64 {
	return d.SpeechStartAt
}

// GetSpeechEndAt returns the relative timestamp in seconds of when a speech segment ends.
func (d *AudioActivitySegment) GetSpeechEndAt() float64 {
	return d.SpeechEndAt
}

// GetDuration returns the duration of the speech segment in seconds.
func (d *AudioActivitySegment) GetDuration() float64 {
	return d.Duration
}

// GetEnergy returns the energy level of the speech segment, used to filter out noise.
func (d *AudioActivitySegment) GetEnergy() float64 {
	return d.Energy
}

// GetConfidence returns a confidence score indicating the likelihood of valid speech (0 to 1).
func (d *AudioActivitySegment) GetConfidence() float32 {
	return d.Confidence
}

type AudioActivity interface {
	Detect(pcm []float32) ([]AudioActivitySegment, error)
}

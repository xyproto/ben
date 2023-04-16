package midi

import (
	"sync"
	"time"
)

const ticksPerQuarterNote = 96

var (
	currentBPM = 120
	mut        sync.RWMutex
)

func SetTempo(bpm int) {
	mut.Lock()
	defer mut.Unlock()
	currentBPM = bpm
}

func Tempo() int {
	mut.RLock()
	defer mut.RUnlock()
	return currentBPM
}

// TicksToDuration converts a MIDI tick value to a time.Duration based on the given tempo in beats per minute
func TicksToDuration(tick int) time.Duration {
	return time.Duration(tick) * time.Minute / (time.Duration(currentBPM) * time.Duration(ticksPerQuarterNote))
}

// DurationToTicks converts a time.Duration value to a MIDI tick value based on the given tempo in beats per minute
func DurationToTicks(duration time.Duration) int {
	return int(duration.Minutes() * float64(currentBPM) * float64(ticksPerQuarterNote))
}

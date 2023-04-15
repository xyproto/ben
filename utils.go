package ben

import "time"

const ticksPerQuarterNote = 96

func uint16ToBytes(value uint16) []byte {
	return []byte{byte(value >> 8), byte(value & 0xFF)}
}

func uint32ToBytes(value uint32) []byte {
	return []byte{byte(value >> 24), byte((value >> 16) & 0xFF), byte((value >> 8) & 0xFF), byte(value & 0xFF)}
}

// TicksToDuration converts a MIDI tick value to a time.Duration based on the given tempo in beats per minute
func TicksToDuration(tick int) time.Duration {
	return time.Duration(tick) * time.Minute / (time.Duration(currentBPM) * time.Duration(ticksPerQuarterNote))
}

// DurationToTicks converts a time.Duration value to a MIDI tick value based on the given tempo in beats per minute
func DurationToTicks(duration time.Duration) int {
	return int(duration.Minutes() * float64(currentBPM) * float64(ticksPerQuarterNote))
}

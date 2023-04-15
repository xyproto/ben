package ben

import (
	"testing"
	"time"
)

func TestTicksToDuration(t *testing.T) {
	tick := 480 // 480 ticks represent a whole note
	currentBPM = 120
	expectedDuration := time.Duration(tick) * time.Minute / (time.Duration(currentBPM) * time.Duration(ticksPerQuarterNote))
	actualDuration := TicksToDuration(tick)
	if actualDuration != expectedDuration {
		t.Errorf("Expected duration: %v, but got: %v", expectedDuration, actualDuration)
	}
}

func TestDurationToTicks(t *testing.T) {
	duration := time.Duration(240) * time.Millisecond // a duration of 240 milliseconds
	currentBPM = 120
	expectedTicks := int(duration.Minutes() * float64(currentBPM) * float64(ticksPerQuarterNote))
	actualTicks := DurationToTicks(duration)
	if actualTicks != expectedTicks {
		t.Errorf("Expected ticks: %v, but got: %v", expectedTicks, actualTicks)
	}
}

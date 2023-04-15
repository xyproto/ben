package ben

import (
	"reflect"
	"testing"
)

func TestProcessBenTrack(t *testing.T) {
	benInput := "C^ D. F G! F. E^ D C"

	expectedMidiNotes := []MidiNote{
		{Duration: TicksToDuration(144)},
		{Duration: TicksToDuration(48)},
		{Duration: TicksToDuration(96)},
		{Duration: TicksToDuration(96)},
		{Duration: TicksToDuration(48)},
		{Duration: TicksToDuration(144)},
		{Duration: TicksToDuration(96)},
		{Duration: TicksToDuration(96)},
	}

	midiNotes := ProcessBenTrack(benInput)

	if len(midiNotes) != len(expectedMidiNotes) {
		t.Fatalf("Expected %d midi notes, got %d", len(expectedMidiNotes), len(midiNotes))
	}

	for i, expectedMidiNote := range expectedMidiNotes {
		if !reflect.DeepEqual(expectedMidiNote.Duration, midiNotes[i].Duration) {
			t.Errorf("MidiNote[%d] durations do not match. Expected %d, got %d", i, expectedMidiNote.Duration, midiNotes[i].Duration)
		}
	}
}

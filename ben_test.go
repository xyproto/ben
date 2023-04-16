package ben

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestConvertToMIDI(t *testing.T) {
	notes := []MidiNote{
		{Frequency: 440, Duration: time.Second, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
		{Frequency: 880, Duration: time.Second, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
		{Frequency: 220, Duration: time.Second, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
	}
	tracks := [][]MidiNote{notes}
	data, err := ConvertToMIDI(tracks)
	if err != nil {
		t.Errorf("ConvertToMIDI returned an error: %v", err)
	}
	// write to a file to verify that MIDI file is valid
	err = os.WriteFile("test.mid", data, 0644)
	if err != nil {
		t.Errorf("Failed to write MIDI data to file: %v", err)
	}
}

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

	midiNotes, err := ParseBenTrack(benInput)
	if err != nil {
		t.Errorf("ParseBenTrack: %v", err)
	}

	if len(midiNotes) != len(expectedMidiNotes) {
		t.Fatalf("Expected %d midi notes, got %d", len(expectedMidiNotes), len(midiNotes))
	}

	for i, expectedMidiNote := range expectedMidiNotes {
		if !reflect.DeepEqual(expectedMidiNote.Duration, midiNotes[i].Duration) {
			t.Errorf("MidiNote[%d] durations do not match. Expected %d, got %d", i, expectedMidiNote.Duration, midiNotes[i].Duration)
		}
	}
}

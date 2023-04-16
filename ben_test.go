package ben

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/xyproto/midi"
)

func TestConvertToMIDI(t *testing.T) {
	notes := []midi.Note{
		{Frequency: 440, Duration: time.Second, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
		{Frequency: 880, Duration: time.Second, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
		{Frequency: 220, Duration: time.Second, Velocity: 127, Channel: 1, Instrument: 0, Slur: false},
	}
	tracks := [][]midi.Note{notes}
	data, err := midi.ConvertToMIDI(tracks)
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

	expectedNotes := []midi.Note{
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

	if len(midiNotes) != len(expectedNotes) {
		t.Fatalf("Expected %d midi notes, got %d", len(expectedNotes), len(midiNotes))
	}

	for i, expectedNote := range expectedNotes {
		if !reflect.DeepEqual(expectedNote.Duration, midiNotes[i].Duration) {
			t.Errorf("midi.Note[%d] durations do not match. Expected %d, got %d", i, expectedNote.Duration, midiNotes[i].Duration)
		}
	}
}

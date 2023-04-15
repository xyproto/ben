package ben

import (
	"bytes"
	"testing"
	"time"
)

func TestWriteTrack(t *testing.T) {
	notes := []MidiNote{
		{Frequency: 261.63, Duration: 750 * time.Millisecond, Velocity: 127, Slur: false},
		{Frequency: 293.66, Duration: 250 * time.Millisecond, Velocity: 127, Slur: false},
		{Frequency: 349.23, Duration: 500 * time.Millisecond, Velocity: 127, Slur: false},
		{Frequency: 392.00, Duration: 500 * time.Millisecond, Velocity: 127, Slur: false},
		{Frequency: 349.23, Duration: 250 * time.Millisecond, Velocity: 127, Slur: false},
		{Frequency: 329.63, Duration: 750 * time.Millisecond, Velocity: 127, Slur: false},
		{Frequency: 293.66, Duration: 500 * time.Millisecond, Velocity: 127, Slur: false},
		{Frequency: 261.63, Duration: 500 * time.Millisecond, Velocity: 127, Slur: false},
	}

	var buf bytes.Buffer
	WriteTrack(&buf, notes)

	expected := []byte{
		77, 84, 114, 107, 0, 0, 0, 38,
		0, 0, 143, 60, 127, 0, 223, 2, 64,
		129, 16, 0, 127, 60, 0, 0, 0, 143, 62, 127, 0, 223, 126, 63,
		48, 0, 127, 62, 0, 0, 0, 143,
	}

	if !bytes.Equal(buf.Bytes(), expected) {
		t.Errorf("Incorrect MIDI data written. Expected:\n%v\nGot:\n%v", expected, buf.Bytes())
	}
}

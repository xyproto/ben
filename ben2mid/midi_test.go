package main

import (
	"bytes"
	"testing"
)

func TestWriteTrack(t *testing.T) {
	buf := new(bytes.Buffer)

	writeTrack(buf, []MidiNote{
		{Frequency: 261.63, Duration: 96, Velocity: 127},
		{Frequency: 293.66, Duration: 96, Velocity: 127},
		{Frequency: 329.63, Duration: 96, Velocity: 127},
	}, 1)

	expected := []byte{
		77, 84, 114, 107, // MTrk
		0, 0, 0, 40, // Track length
		0, 224, 2, 64, // Pitch bend for first note
		0, 144, 60, 127, // Note on for first note
		96, 128, 60, 0, // Note off for first note
		0, 224, 126, 63, // Pitch bend for second note
		0, 144, 62, 127, // Note on for second note
		96, 128, 62, 0, // Note off for second note
		0, 224, 1, 64, // Pitch bend for third note
		0, 144, 64, 127, // Note on for third note
		96, 128, 64, 0, // Note off for third note
		0, 255, 47, 0, // End of track
	}

	if !bytes.Equal(buf.Bytes(), expected) {
		t.Errorf("Incorrect MIDI data written. Expected:\n%v\nGot:\n%v", expected, buf.Bytes())
	}
}

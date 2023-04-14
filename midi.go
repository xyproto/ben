package main

import (
	"fmt"
	"io"
	"math"
)

type MidiNote struct {
	Frequency float64
	Duration  int
}

func frequencyToMidi(freq float64) (int, int) {
	midiNote := 69 + 12*math.Log2(freq/440.0)
	midiNoteRounded := int(math.Round(midiNote))

	pitchBend := int(math.Round(8192 * (midiNote - float64(midiNoteRounded))))
	pitchBend = pitchBend + 8192 // Center value for pitch bend

	return midiNoteRounded, pitchBend
}

func writeMidiFile(w io.Writer, tracks [][]MidiNote) {
	writeHeader(w, len(tracks))
	for i, track := range tracks {
		writeTrack(w, track, i+1)
	}
}

func writeHeader(w io.Writer, numTracks int) {
	w.Write([]byte("MThd"))                   // MIDI header chunk
	w.Write(uint32ToBytes(6))                 // Header length
	w.Write(uint16ToBytes(1))                 // Format type (1: multiple tracks)
	w.Write(uint16ToBytes(uint16(numTracks))) // Number of tracks
	w.Write(uint16ToBytes(96))                // Division (ticks per quarter note)
}

func writePitchBend(trackData []byte, pitchBend int) []byte {
	pitchBendLSB := pitchBend & 0x7F
	pitchBendMSB := (pitchBend >> 7) & 0x7F

	trackData = append(trackData, 0x00) // Delta time (0 ticks)
	trackData = append(trackData, 0xE0) // Pitch bend change, channel 1
	trackData = append(trackData, byte(pitchBendLSB))
	trackData = append(trackData, byte(pitchBendMSB))

	return trackData
}

func writeTrack(w io.Writer, notes []MidiNote, channel int) {
	// Prepare track data
	trackData := make([]byte, 0)
	for _, note := range notes {
		midiNote, pitchBend := frequencyToMidi(note.Frequency)

		trackData = writePitchBend(trackData, pitchBend)

		// Note on
		trackData = append(trackData, 0x00) // Delta time (0 ticks)
		trackData = append(trackData, 0x90) // Note on, channel 1
		// byte(0x90+channel-1)) // Note on, channel
		trackData = append(trackData, byte(midiNote))
		trackData = append(trackData, byte(0x60)) // Velocity (0x60)

		// Note off
		trackData = append(trackData, byte(note.Duration)) // Delta time (ticks)
		trackData = append(trackData, 0x80)                // Note off, channel 1
		//byte(0x80+channel-1)) // Note off, channel
		trackData = append(trackData, byte(midiNote))
		trackData = append(trackData, byte(0x00)) // Velocity (0x00)
	}

	// End of track event
	trackData = append(trackData, 0x00, 0xFF, 0x2F, 0x00)

	// Write track chunk
	w.Write([]byte("MTrk"))                        // MIDI track chunk
	w.Write(uint32ToBytes(uint32(len(trackData)))) // Track length
	w.Write(trackData)                             // Track data

	fmt.Println("WROTE A TRACK")
}

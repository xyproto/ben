package main

import (
	"fmt"
	"io"
	"math"
)

type MidiNote struct {
	Frequency float64
	Duration  int
	Velocity  byte
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

func writeNoteOn(trackData []byte, channel int, midiNote int, velocity byte) []byte {
	trackData = append(trackData, 0x00)                 // Delta time (0 ticks)
	trackData = append(trackData, byte(0x90+channel-1)) // Note on, channel
	trackData = append(trackData, byte(midiNote))
	trackData = append(trackData, velocity) // Velocity

	return trackData
}

func writeNoteOff(trackData []byte, channel int, midiNote int, ticks int) []byte {
	trackData = append(trackData, byte(ticks))          // Delta time (ticks)
	trackData = append(trackData, byte(0x80+channel-1)) // Note off, channel
	trackData = append(trackData, byte(midiNote))
	trackData = append(trackData, byte(0x00)) // Velocity

	return trackData
}

func writeTrack(w io.Writer, notes []MidiNote, channel int) {
	// Prepare track data
	trackData := make([]byte, 0)
	var lastNoteDuration int
	var lastMidiNote int
	for i, note := range notes {
		midiNote, pitchBend := frequencyToMidi(note.Frequency)

		trackData = writePitchBend(trackData, pitchBend)

		// Note on
		if i > 0 && notes[i-1].Duration == -1 {
			// Handle slurs
			writeNoteOn(trackData, channel, lastMidiNote, note.Velocity)
		} else {
			// Normal note on
			writeNoteOn(trackData, channel, midiNote, note.Velocity)
		}

		// Note off
		if note.Duration != -1 {
			writeNoteOff(trackData, channel, midiNote, note.Duration)
			lastNoteDuration = note.Duration
			lastMidiNote = midiNote
		} else {
			// Handle slurs
			if lastNoteDuration != -1 {
				trackData = append(trackData, byte(lastNoteDuration))
				trackData = append(trackData, 0x80+byte(channel-1))
				trackData = append(trackData, byte(lastMidiNote))
				trackData = append(trackData, byte(0x00))
			}
		}
	}

	// End of track event
	trackData = append(trackData, 0x00, 0xFF, 0x2F, 0x00)

	// Write track chunk
	w.Write([]byte("MTrk"))                        // MIDI track chunk
	w.Write(uint32ToBytes(uint32(len(trackData)))) // Track length
	w.Write(trackData)                             // Track data

	fmt.Println("WROTE A TRACK")
}

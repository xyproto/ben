package ben

import (
	"bytes"
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

func writePitchBend(w io.Writer, pitchBend int) {
	pitchBendLSB := pitchBend & 0x7F
	pitchBendMSB := (pitchBend >> 7) & 0x7F

	w.Write([]byte{0x00, 0xE0, byte(pitchBendLSB), byte(pitchBendMSB)})
}

func writeNoteOn(w io.Writer, channel int, midiNote int, velocity byte) {
	w.Write([]byte{0x00, byte(0x90 + channel - 1), byte(midiNote), velocity})
}

func writeNoteOff(w io.Writer, channel int, midiNote int, ticks int) {
	w.Write([]byte{byte(ticks), byte(0x80 + channel - 1), byte(midiNote), 0x00})
}

func writeTrack(w io.Writer, notes []MidiNote, channel int) {
	// Create a buffer to hold track data
	buf := new(bytes.Buffer)

	for i, note := range notes {
		midiNote, pitchBend := frequencyToMidi(note.Frequency)

		// Only call writePitchBend if pitchBend is not the center value
		if pitchBend != 8192 {
			writePitchBend(buf, pitchBend)
		}

		if i > 0 && notes[i-1].Duration == -1 {
			writeNoteOn(buf, channel, midiNote, note.Velocity)
		} else {
			writeNoteOn(buf, channel, midiNote, note.Velocity)
		}

		if note.Duration != -1 {
			writeNoteOff(buf, channel, midiNote, note.Duration)
		} else {
			if notes[i-1].Duration != -1 {
				buf.WriteByte(byte(notes[i-1].Duration))
				writeNoteOff(buf, channel, midiNote, 0)
			}
		}
	}

	// End of track event
	buf.Write([]byte{0x00, 0xFF, 0x2F, 0x00})

	// Write track chunk
	w.Write([]byte("MTrk"))                   // MIDI track chunk
	w.Write(uint32ToBytes(uint32(buf.Len()))) // Track length
	buf.WriteTo(w)                            // Track data

	fmt.Println("WROTE A TRACK")
}

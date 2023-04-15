package ben

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"time"
)

type MidiNote struct {
	Frequency  float64
	Duration   time.Duration
	Velocity   byte
	Channel    int
	Instrument int
	Slur       bool
}

func FrequencyToMidi(freq float64) (int, int) {
	midiNote := 69 + 12*math.Log2(freq/440.0)
	midiNoteRounded := int(math.Round(midiNote))

	pitchBend := int(math.Round(8192 * (midiNote - float64(midiNoteRounded))))
	pitchBend = pitchBend + 8192 // Center value for pitch bend

	return midiNoteRounded, pitchBend
}

func WriteMidiFile(w io.Writer, tracks [][]MidiNote) {
	WriteHeader(w, len(tracks))
	for _, track := range tracks {
		writeTrack(w, track)
	}
}

func WriteHeader(w io.Writer, numTracks int) {
	w.Write([]byte("MThd"))                   // MIDI header chunk
	w.Write(uint32ToBytes(6))                 // Header length
	w.Write(uint16ToBytes(1))                 // Format type (1: multiple tracks)
	w.Write(uint16ToBytes(uint16(numTracks))) // Number of tracks
	w.Write(uint16ToBytes(96))                // Division (ticks per quarter note)
}

func WritePitchBend(w io.Writer, channel int, pitchBend int) {
	pitchBendLSB := pitchBend & 0x7F
	pitchBendMSB := (pitchBend >> 7) & 0x7F

	w.Write([]byte{0x00, byte(0xE0 + channel - 1), byte(pitchBendLSB), byte(pitchBendMSB)})
}

func WriteNoteOn(w io.Writer, channel int, midiNote int, velocity byte) {
	w.Write([]byte{0x00, byte(0x90 + channel - 1), byte(midiNote), velocity})
}

func WriteNoteOff(w io.Writer, channel int, midiNote int, ticks int) {
	w.Write([]byte{byte(ticks), byte(0x80 + channel - 1), byte(midiNote), 0x00})
}

func WriteProgramChange(w io.Writer, channel int, program int) {
	w.Write([]byte{0x00, byte(0xC0 + channel - 1), byte(program)})
}

func writeTrack(w io.Writer, notes []MidiNote) {
	// Create a buffer to hold track data
	buf := new(bytes.Buffer)

	for _, note := range notes {
		midiNote, pitchBend := FrequencyToMidi(note.Frequency)

		// Set the instrument (program)
		if note.Instrument != -1 {
			WriteProgramChange(buf, note.Channel, note.Instrument)
		}

		// Set the pitch bend
		if pitchBend != 8192 {
			WritePitchBend(buf, note.Channel, pitchBend)
		}

		WriteNoteOn(buf, note.Channel, midiNote, note.Velocity)

		// Sleep for note duration
		buf.WriteByte(byte(note.Duration))

		WriteNoteOff(buf, note.Channel, midiNote, 0)

		// Reset the pitch bend
		if pitchBend != 8192 {
			WritePitchBend(buf, note.Channel, 8192)
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

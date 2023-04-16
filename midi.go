package ben

import (
	"bytes"
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

func ConvertToMIDI(tracks [][]MidiNote) ([]byte, error) {
	buf := new(bytes.Buffer)
	WriteHeader(buf, len(tracks))
	for _, track := range tracks {
		if err := WriteTrack(buf, track); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func ConvertToMIDITracks(parsedBEN []MidiNote) [][]MidiNote {
	return [][]MidiNote{parsedBEN}
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

func writeDeltaTime(w io.Writer, deltaTime int) error {
	var buf [4]byte
	endIndex := 0
	for i := 3; i >= 0; i-- {
		buf[i] = byte(deltaTime)&0x7F | 0x80
		deltaTime >>= 7
		if deltaTime == 0 {
			endIndex = i
			break
		}
	}
	buf[3] &= 0x7F
	_, err := w.Write(buf[endIndex:])
	return err
}
func deltaTimeLength(value int) int {
	if value < 0x80 {
		return 1
	}
	if value < 0x4000 {
		return 2
	}
	if value < 0x200000 {
		return 3
	}
	return 4
}

func writeBytes(w io.Writer, data []byte) error {
	_, err := w.Write(data)
	return err
}

func writeUint32(w io.Writer, value uint32) error {
	return writeBytes(w, uint32ToBytes(value))
}

func WriteTrack(w io.Writer, notes []MidiNote) error {
	midiNotes := make([]int, len(notes))

	// Calculate the MIDI note numbers and track length
	trackLength := uint32(0)
	for i, note := range notes {
		midiNote, _ := FrequencyToMidi(note.Frequency)
		midiNotes[i] = midiNote
		trackLength += 3 + uint32(deltaTimeLength(DurationToTicks(note.Duration)))
	}
	trackLength += 4 // For the end of track event

	// Write the track header
	if err := writeBytes(w, []byte{77, 84, 114, 107}); err != nil { // MTrk
		return err
	}
	if err := writeUint32(w, trackLength); err != nil {
		return err
	}

	lastNoteOff := -1
	for i, note := range notes {
		midiNote, pitchBend := FrequencyToMidi(note.Frequency)

		// Note Off
		if lastNoteOff != -1 {
			if err := writeDeltaTime(w, DurationToTicks(notes[lastNoteOff].Duration)); err != nil {
				return err
			}
			WriteNoteOff(w, 0, midiNotes[lastNoteOff], 0)
		}

		// Note On
		if err := writeDeltaTime(w, 0); err != nil {
			return err
		}
		WriteNoteOn(w, 0, midiNote, note.Velocity)
		WritePitchBend(w, 0, pitchBend)

		lastNoteOff = i
	}

	// Write the end of track event
	if err := writeDeltaTime(w, DurationToTicks(notes[lastNoteOff].Duration)); err != nil {
		return err
	}

	return writeBytes(w, []byte{0xFF, 0x2F, 0x00}) // End of track
}

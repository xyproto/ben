package midi

import (
	"bytes"
	"io"
	"math"
	"time"
)

type Note struct {
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

func NoteToFrequency(note string) float64 {
	noteNameToPitch := map[string]float64{
		"C":  0,
		"C#": 1,
		"Db": 1,
		"D":  2,
		"D#": 3,
		"Eb": 3,
		"E":  4,
		"F":  5,
		"F#": 6,
		"Gb": 6,
		"G":  7,
		"G#": 8,
		"Ab": 8,
		"A":  9,
		"A#": 10,
		"Bb": 10,
		"B":  11,
		"H":  11,
	}

	if len(note) < 2 || len(note) > 3 {
		return 0
	}

	noteName := note[:len(note)-1]
	pitch, valid := noteNameToPitch[noteName]
	if !valid {
		return 0
	}

	octave := float64(note[len(note)-1] - '0')
	midiNumber := 12*(octave+1) + pitch
	return 440 * math.Pow(2, (midiNumber-69)/12)
}

func ConvertToMIDI(tracks [][]Note) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write MIDI header
	WriteHeader(buf, len(tracks))

	// Write tracks
	for _, track := range tracks {
		trackData, err := ConvertToMIDITracks(track)
		if err != nil {
			return nil, err
		}
		buf.Write(trackData)
	}

	return buf.Bytes(), nil
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

// writeDeltaTime writes a delta time value in MIDI ticks as a variable-length quantity to an io.Writer
func writeDeltaTime(w io.Writer, deltaTime int) error {
	vlq := deltaTimeToVLQ(deltaTime)
	return writeBytes(w, vlq)
}

// deltaTimeToVLQ converts a delta time value in MIDI ticks to a variable-length quantity representation
func deltaTimeToVLQ(deltaTime int) []byte {
	var buf []byte
	current := deltaTime & 0x7F

	for deltaTime >>= 7; deltaTime > 0; deltaTime >>= 7 {
		buf = append([]byte{byte(current | 0x80)}, buf...)
		current = deltaTime & 0x7F
	}
	buf = append([]byte{byte(current)}, buf...)
	return buf
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

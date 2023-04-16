package midi

import (
	"bytes"
	"io"
)

var lastNoteOn Note

func ConvertToMIDITracks(tracks []Note) ([]byte, error) {
	buf := new(bytes.Buffer)

	// Write track header
	err := WriteTrack(buf, tracks)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func WriteTrack(w io.Writer, notes []Note) error {
	// Calculate track length in bytes
	trackLength := 0
	for i, note := range notes {
		deltaTime := 0
		if i > 0 {
			deltaTime = int(note.Duration.Seconds() * 96.0)
		}
		trackLength += deltaTimeLength(deltaTime)
		trackLength += 3 // Note on
		trackLength += 3 // Note off
	}

	// Write track header
	w.Write([]byte("MTrk"))
	w.Write(uint32ToBytes(uint32(trackLength)))

	// Write MIDI events
	for i, note := range notes {
		deltaTime := 0
		if i > 0 {
			deltaTime = int(note.Duration.Seconds() * 96.0)
		}

		// Write delta time
		err := writeDeltaTime(w, deltaTime)
		if err != nil {
			return err
		}

		// Write program change
		WriteProgramChange(w, note.Channel, note.Instrument)

		// Write note on
		midiNote, pitchBend := FrequencyToMidi(note.Frequency)
		WriteNoteOn(w, note.Channel, midiNote, note.Velocity)

		// Write pitch bend if needed
		if pitchBend != 8192 {
			WritePitchBend(w, note.Channel, pitchBend)
		}

		// Write note off
		if !note.Slur {
			noteOffDeltaTime := int(note.Duration.Seconds() * 96.0)
			WriteNoteOff(w, note.Channel, midiNote, noteOffDeltaTime)
		}

		lastNoteOn = note
	}

	return nil
}

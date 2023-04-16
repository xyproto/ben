package ben

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xyproto/midi"
)

func ListMIDIOutDevices() (string, error) {
	var sb strings.Builder
	devPath := "/dev"

	files, err := os.ReadDir(devPath)
	if err != nil {
		return "", nil
	}

	sb.WriteString("Available MIDI out devices:\n")
	deviceIndex := 0
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "midi") {
			sb.WriteString(fmt.Sprintf("[%d] %s\n", deviceIndex, file.Name()))
			deviceIndex++
		}
	}
	return sb.String(), nil
}

func PlayWithSynth(midiOutDevice int, tracks [][]midi.Note) error {
	devPath := "/dev"

	files, err := os.ReadDir(devPath)
	if err != nil {
		return err
	}

	midiDeviceFound := false
	deviceIndex := 0
	var deviceFile string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "midi") {
			if deviceIndex == midiOutDevice {
				midiDeviceFound = true
				deviceFile = file.Name()
				break
			}
			deviceIndex++
		}
	}

	if !midiDeviceFound {
		return fmt.Errorf("invalid MIDI out device index: %d", midiOutDevice)
	}

	devicePath := filepath.Join(devPath, deviceFile)
	device, err := os.OpenFile(devicePath, os.O_WRONLY, 0)
	if err != nil {
		return fmt.Errorf("error opening MIDI out device: %v", err)
	}
	defer device.Close()

	for trackIndex, track := range tracks {
		for _, note := range track {
			midiNote, pitchBend := midi.FrequencyToMidi(note.Frequency)

			if note.Channel != -1 {
				// Set the instrument (program)
				midi.WriteProgramChange(device, note.Channel, note.Instrument)

				// Set the pitch bend
				if pitchBend != 8192 {
					midi.WritePitchBend(device, note.Channel, pitchBend)
				}

				// Note On
				midi.WriteNoteOn(device, note.Channel, midiNote, note.Velocity)

				// Sleep for note duration
				time.Sleep(time.Duration(note.Duration) * time.Millisecond)

				// Note Off
				midi.WriteNoteOff(device, note.Channel, midiNote, DurationToTicks(note.Duration))

				// Reset the pitch bend
				if pitchBend != 8192 {
					midi.WritePitchBend(device, note.Channel, 8192)
				}
			}
		}

		// Sleep for 1 second between tracks
		if trackIndex < len(tracks)-1 {
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

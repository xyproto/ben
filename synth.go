package ben

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ListMIDIOutDevices() {
	devPath := "/dev"

	files, err := os.ReadDir(devPath)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	fmt.Println("Available MIDI out devices:")
	deviceIndex := 0
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "midi") {
			fmt.Printf("[%d] %s\n", deviceIndex, file.Name())
			deviceIndex++
		}
	}
}

func PlayWithSynth(midiOutDevice int, tracks [][]MidiNote) {
	devPath := "/dev"

	files, err := os.ReadDir(devPath)
	if err != nil {
		fmt.Println("Error: ", err)
		return
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
		fmt.Println("Invalid MIDI out device index")
		return
	}

	devicePath := filepath.Join(devPath, deviceFile)
	device, err := os.OpenFile(devicePath, os.O_WRONLY, 0)
	if err != nil {
		fmt.Println("Error opening MIDI out device: ", err)
		return
	}
	defer device.Close()

	for trackIndex, track := range tracks {
		for _, note := range track {
			midiNote, pitchBend := FrequencyToMidi(note.Frequency)

			if note.Channel != -1 {
				// Set the instrument (program)
				WriteProgramChange(device, note.Channel, note.Instrument)

				// Set the pitch bend
				if pitchBend != 8192 {
					WritePitchBend(device, note.Channel, pitchBend)
				}

				// Note On
				WriteNoteOn(device, note.Channel, midiNote, note.Velocity)

				// Sleep for note duration
				time.Sleep(time.Duration(note.Duration) * time.Millisecond)

				// Note Off
				WriteNoteOff(device, note.Channel, midiNote, DurationToTicks(note.Duration))

				// Reset the pitch bend
				if pitchBend != 8192 {
					WritePitchBend(device, note.Channel, 8192)
				}
			}
		}

		// Sleep for 1 second between tracks
		if trackIndex < len(tracks)-1 {
			time.Sleep(1 * time.Second)
		}
	}
}

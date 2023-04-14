package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gomidi/midi/mid"
	"github.com/gomidi/midi/midimessage/channel"
	"github.com/gomidi/midi/writer"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var input string
	for scanner.Scan() {
		input += scanner.Text()
	}
	if scanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", scanner.Err())
		os.Exit(1)
	}

	err := ben2mid(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting BEN to MIDI: %v\n", err)
		os.Exit(1)
	}
}

func ben2mid(ben string) error {
	midiFile := mid.NewSMF(0, mid.TicksPerQuarterNote(120))
	midiTrack := &mid.Track{}
	midiFile.AddTrack(midiTrack)

	writer := writer.New(midiTrack)

	// Tokenize the BEN input
	tokens := strings.Fields(ben)

	for _, token := range tokens {
		switch {
		case isNote(token):
			note, duration, err := parseNote(token)
			if err != nil {
				return fmt.Errorf("could not parse note: %w", err)
			}
			writer.NoteOn(note)
			writer.Rest(duration)
			writer.NoteOff(note)

		case isRest(token):
			duration, err := parseRest(token)
			if err != nil {
				return fmt.Errorf("could not parse rest: %w", err)
			}
			writer.Rest(duration)

		default:
			return fmt.Errorf("unknown token: %s", token)
		}
	}

	err := mid.WriteSMF(os.Stdout, midiFile)
	if err != nil {
		return fmt.Errorf("could not write MIDI file: %w", err)
	}

	return nil
}

func isNote(token string) bool {
    if len(token) == 0 {
        return false
    }
    return unicode.IsLetter(rune(token[0]))
}

func isRest(token string) bool {
    if len(token) == 0 {
        return false
    }
    return token[0] == ','
}

func parseNote(noteStr string) (note channel.NoteOn, duration uint32, err error) {
    var pitch string
    var octave int
    var articulation string
    var microtonality string
    var durStr string

    re := regexp.MustCompile(`([A-Ga-g][-+#]?)(\d)?([~>]*)([%$]?)(\d*\.\d+|\d+)?`)
    matches := re.FindStringSubmatch(noteStr)

    if len(matches) == 0 {
        return note, 0, fmt.Errorf("invalid note format: %s", noteStr)
    }

    pitch = matches[1]
    if matches[2] != "" {
        octave, err = strconv.Atoi(matches[2])
        if err != nil {
            return note, 0, fmt.Errorf("invalid octave value: %s", noteStr)
        }
    } else {
        octave = 4
    }
    articulation = matches[3]
    microtonality = matches[4]
    durStr = matches[5]

    midiNote, err := noteToMidi(pitch, octave)
    if err != nil {
        return note, 0, err
    }

    duration = 1
    if durStr != "" {
        duration, err = strconv.Atoi(durStr)
        if err != nil {
            return note, 0, fmt.Errorf("invalid duration value: %s", noteStr)
        }
    }

    note = channel.NoteOn(midiNote, 100)

    if articulation != "" {
        note = applyArticulation(note, articulation)
    }

    if microtonality != "" {
        note = applyMicrotonality(note, microtonality)
    }

    return note, uint32(duration), nil
}

func noteToMidi(pitch string, octave int) (uint8, error) {
    baseNotes := map[string]int{
        "C": 0,
        "D": 2,
        "E": 4,
        "F": 5,
        "G": 7,
        "A": 9,
        "B": 11,
        "H": 11,
    }

    pitchLength := len(pitch)
    if pitchLength < 1 || pitchLength > 3 {
        return 0, fmt.Errorf("invalid pitch format: %s", pitch)
    }

    baseNote := strings.ToUpper(pitch[:1])
    baseNoteValue, ok := baseNotes[baseNote]
    if !ok {
        return 0, fmt.Errorf("invalid note: %s", pitch)
    }

    noteValue := baseNoteValue
    if pitchLength > 1 {
        modifier := pitch[1:]
        switch modifier {
        case "#":
            noteValue++
        case "+":
            noteValue++
        case "-":
            noteValue--
        default:
            return 0, fmt.Errorf("invalid modifier: %s", pitch)
        }
    }

    midiValue := noteValue + (octave+1)*12
    if midiValue < 0 || midiValue > 127 {
        return 0, fmt.Errorf("MIDI value out of range: %d", midiValue)
    }

    return uint8(midiValue), nil
}

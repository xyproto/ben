package main

import (
	"fmt"
	"math"
	"strings"
)

var currentOctave = 4

func processBenTrack(benTrack string) []MidiNote {
	fmt.Println("PROCESSING: " + benTrack)
	var midiNotes []MidiNote
	benNotes := strings.Split(benTrack, " ")

	for _, benNote := range benNotes {
		if freq, duration, ok := benToFrequency(benNote); ok {
			midiNotes = append(midiNotes, MidiNote{Frequency: freq, Duration: duration})
		}
	}

	return midiNotes
}

func benToFrequency(benNote string) (float64, int, bool) {
	baseNotes := map[string]float64{
		"C": 261.63,
		"D": 293.66,
		"E": 329.63,
		"F": 349.23,
		"G": 392.00,
		"A": 440.00,
		"B": 493.88,
		"H": 493.88,
	}

	var frequency float64 = 0
	var duration int = 96 // quarter note duration
	var octaveMode bool = false

	// Parse note string
	for _, c := range benNote {
		if octaveMode {
			switch c {
			case '-':
				currentOctave--
			case '+':
				currentOctave++
			case '0':
				currentOctave = 0
			case '1':
				currentOctave = 1
			case '2':
				currentOctave = 2
			case '3':
				currentOctave = 3
			case '4':
				currentOctave = 4
			case '5':
				currentOctave = 5
			case '6':
				currentOctave = 6
			case '7':
				currentOctave = 7
			case '8':
				currentOctave = 8
			case '9':
				currentOctave = 9
			}
		}
		if frequency > 0 {
			switch c {
			case '#':
				frequency *= math.Pow(2, 1.0/12.0)
			case 'b':
				frequency /= math.Pow(2, 1.0/12.0)
				continue
			}
		}
		switch c {
		case 'C', 'D', 'E', 'F', 'G', 'A', 'B', 'H':
			frequency = baseNotes[string(c)]
			fmt.Println("TO IMPLEMENT: ALSO MODIFY THE DURATION")
		case 'c', 'd', 'e', 'f', 'g', 'a', 'b', 'h':
			frequency = baseNotes[strings.ToUpper(string(c))]
			fmt.Println("TO IMPLEMENT: ALSO MODIFY THE DURATION")
		case '{':
			frequency *= 0.98181818181 // 1.818% decrease
		case '}':
			frequency /= 0.98181818181 // 1.818% increase
		case '(':
			octaveMode = true
			continue
		case ')':
			octaveMode = false
			continue
		case '+':
			currentOctave++
		case '-':
			currentOctave--
		case ',':
			duration = 48 // eighth note rest
		case '~':
			fmt.Println("TO IMPLEMENT: SLURS AND TIES")
		case '.':
			duration /= 2 // staccato
		case '^':
			fmt.Println("TO IMPLEMENT: ACCENTS")
		case '/':
			fmt.Println("TO IMPLEMENT: WHOLE NOTES")
		default:
			return 0, 0, false
		}
	}

	// Calculate frequency based on current octave
	frequency *= math.Pow(2, float64(currentOctave-4))

	if frequency < 0 {
		frequency = 0
	}

	fmt.Printf("%s ==> %fHz, %d\n", benNote, frequency, duration)

	return frequency, duration, true
}

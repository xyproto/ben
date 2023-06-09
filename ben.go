package ben

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/xyproto/midi"
)

var (
	currentOctave = 4
	currentBPM    = 120
)

func ParseBenTrack(benTrack string) ([]midi.Note, error) {
	var midiNotes []midi.Note
	benNotes := strings.Split(strings.TrimSpace(benTrack), " ")
	for _, benNote := range benNotes {
		if freq, duration, velocity, channel, instrument, slur, ok := NoteToFrequency(strings.TrimSpace(benNote)); ok {
			midiNotes = append(midiNotes, midi.Note{Frequency: freq, Duration: duration, Velocity: velocity, Channel: channel, Instrument: instrument, Slur: slur})
		} else {
			return nil, fmt.Errorf("could not parse this note: %q", benNote)
		}
	}
	return midiNotes, nil
}

func NoteToFrequency(benNote string) (float64, time.Duration, byte, int, int, bool, bool) {
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
	var (
		channel    = -1
		instrument = -1
		frequency  float64
		duration        = TicksToDuration(96) // quarter note duration
		octaveMode      = false
		velocity   byte = 127 // the default is full velocity
		slur            = false
	)

	// Parse note string
	for i, c := range benNote {
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
			if i+1 < len(benNote) && benNote[i+1] == '/' {
				duration = TicksToDuration(192) // half note duration
			} else if i+1 < len(benNote) && c == rune(benNote[i+1]) {
				duration = TicksToDuration(384) // whole note duration
			} else {
				duration = TicksToDuration(96) // quarter note duration
			}
			if i+1 < len(benNote) && (benNote[i+1] == '#' || benNote[i+1] == 'b') && i+2 < len(benNote) && benNote[i+2] == '(' {
				octaveMode = true
			}
		case 'c', 'd', 'e', 'f', 'g', 'a', 'b', 'h':
			frequency = baseNotes[strings.ToUpper(string(c))]
			duration = TicksToDuration(48) // eighth note duration
			if i+1 < len(benNote) && (benNote[i+1] == '#' || benNote[i+1] == 'b') && i+2 < len(benNote) && benNote[i+2] == '(' {
				octaveMode = true
			}
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
		case ',':
			duration = TicksToDuration(48) // eighth note rest
		case '~':
			slur = true
		case '.':
			duration /= 2 // staccato
		case '^':
			duration += (duration / 2) // accent
		case 'v':
			velocity = byte(float64(velocity) * 0.9)
		case '-':
			if i+1 < len(benNote) && benNote[i+1] == '>' {
				// do nothing, handle in the next iteration
				continue
			} else {
				currentOctave--
			}
		case '>':
			if i > 0 && benNote[i-1] == '-' {
				slur = true // slur
			} else {
				duration += (duration / 2) // accent
			}
		case '!':
			if i+1 < len(benNote) {
				volume := int(benNote[i+1]) - '0'
				if volume >= 0 && volume <= 9 {
					velocity = byte((float64(volume) / 9) * 127)
					i++ // Skip the volume digit
				}
			}
		case '@':
			if i+1 < len(benNote) {
				channel = int(benNote[i+1]) - '0'
				if channel >= 0 && channel <= 9 {
					i++ // Skip the channel digit
				} else {
					channel = -1
				}
			}
		case '*':
			if i+1 < len(benNote) {
				instrument = int(benNote[i+1]) - '0'
				if instrument >= 0 && instrument <= 9 {
					i++ // Skip the instrument digit
				} else {
					instrument = -1
				}
			}
		case 'T':
			if i+1 < len(benNote) {
				tempo, err := strconv.Atoi(benNote[i+1:])
				if err == nil && tempo > 0 {
					currentBPM = tempo
				}
			}
		default:
			return 0, 0, 0, 0, 0, false, false
		}
	}

	// Calculate frequency based on current octave
	frequency *= math.Pow(2, float64(currentOctave-4))

	if frequency < 0 {
		frequency = 0
	}

	return frequency, duration, velocity, channel, instrument, slur, true
}

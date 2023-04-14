package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Check if there are input files provided as command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Usage: ben2mid <input1.ben> <input2.ben> ...")
		os.Exit(1)
	}

	// Read the .ben files from the command-line arguments
	benFilenames := os.Args[1:]
	var benTracks []string

	for _, filename := range benFilenames {
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filename, err)
			os.Exit(1)
		}
		benTracks = append(benTracks, strings.TrimSpace(string(content)))
	}

	// Process the .ben files and convert them to MIDI notes
	var midiTracks [][]MidiNote
	for _, benTrack := range benTracks {
		midiNotes := processBenTrack(benTrack)
		midiTracks = append(midiTracks, midiNotes)
	}

	// Write the MIDI tracks to a single output.mid file
	outputFile, err := os.Create("output.mid")
	if err != nil {
		fmt.Printf("Error creating output.mid: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	writeMidiFile(outputFile, midiTracks)
}

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	outputFilename = flag.String("o", "output.mid", "Specify the output MIDI filename")
	versionFlag    = flag.Bool("version", false, "Display the version information")
	helpFlag       = flag.Bool("help", false, "Display help information")
)

func main() {
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Println("ben2mid version 1.0.0")
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Usage: ben2mid [flags] <input1.ben> <input2.ben> ...")
		flag.PrintDefaults()
		os.Exit(1)
	}

	benFilenames := args
	var benTracks []string

	for _, filename := range benFilenames {
		content, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filename, err)
			os.Exit(1)
		}
		benTracks = append(benTracks, strings.TrimSpace(string(content)))
	}

	var midiTracks [][]MidiNote
	for _, benTrack := range benTracks {
		midiNotes := processBenTrack(benTrack)
		midiTracks = append(midiTracks, midiNotes)
	}

	outputFile, err := os.Create(*outputFilename)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", *outputFilename, err)
		os.Exit(1)
	}
	defer outputFile.Close()

	writeMidiFile(outputFile, midiTracks)
}

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/xyproto/ben"
)

func main() {
	input := flag.String("i", "", "Input file (.ben)")
	output := flag.String("o", "", "Output file (.mid)")
	listDevices := flag.Bool("l", false, "List MIDI devices")
	synthDevice := flag.String("s", "", "Play with external synth (device index, use -l to list)")
	flag.Parse()

	if *listDevices {
		ben.ListDevices()
		return
	}

	if *input == "" {
		fmt.Println("Please provide an input file.")
		flag.Usage()
		return
	}

	if *output == "" && *synthDevice == "" {
		fmt.Println("Please provide an output file or a synth device index.")
		flag.Usage()
		return
	}

	content, err := ioutil.ReadFile(*input)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		return
	}

	parsedBEN, err := ben.ParseBEN(string(content))
	if err != nil {
		fmt.Println("Error parsing BEN:", err)
		return
	}

	if *synthDevice != "" {
		deviceIndex, err := strconv.Atoi(*synthDevice)
		if err != nil {
			fmt.Println("Invalid synth device index:", err)
			return
		}
		err = ben.PlayWithSynth(parsedBEN, deviceIndex)
		if err != nil {
			fmt.Println("Error playing with synth:", err)
		}
		return
	}

	midiData, err := ben.ConvertToMIDI(parsedBEN)
	if err != nil {
		fmt.Println("Error converting BEN to MIDI:", err)
		return
	}

	err = ioutil.WriteFile(*output, midiData, 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}
}

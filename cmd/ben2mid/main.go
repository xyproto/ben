package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli/v2"
	"github.com/xyproto/ben"
	"github.com/xyproto/midi"
)

func main() {
	app := &cli.App{
		Name:  "ben2mid",
		Usage: "Convert BEN files to MIDI",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output file (.mid)",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "list-devices",
				Usage: "List MIDI devices",
			},
			&cli.StringFlag{
				Name:  "synth-device",
				Usage: "Play with external synth (device index, use --list-devices to list)",
			},
		},
		Action: func(c *cli.Context) error {
			output := c.String("output")
			listDevices := c.Bool("list-devices")
			synthDevice := c.String("synth-device")

			if listDevices {
				s, err := ben.ListMIDIOutDevices()
				if err != nil {
					return fmt.Errorf("error listing MIDI output devices: %v", err)
				}
				fmt.Println(strings.TrimSpace(s))
				return nil
			}

			if output == "" && synthDevice == "" {
				return fmt.Errorf("please provide an output file or a synth device index")
			}

			var midiTracks [][]midi.Note
			for _, inputFile := range c.Args().Slice() {
				content, err := os.ReadFile(inputFile)
				if err != nil {
					return fmt.Errorf("error reading input file: %v", err)
				}

				parsedBEN, err := ben.ParseBenTrack(string(content))
				if err != nil {
					return fmt.Errorf("error parsing %s: %v", inputFile, err)
				}
				midiTracks = append(midiTracks, parsedBEN)
			}

			if synthDevice != "" {
				deviceIndex, err := strconv.Atoi(synthDevice)
				if err != nil {
					return fmt.Errorf("invalid synth device index: %v", err)
				}
				return ben.PlayWithSynth(deviceIndex, midiTracks)
			}

			midiData, err := midi.ConvertToMIDI(midiTracks)
			if err != nil {
				return fmt.Errorf("error converting BEN to MIDI: %v", err)
			}

			err = ioutil.WriteFile(output, midiData, 0644)
			if err != nil {
				return fmt.Errorf("error writing output file: %v", err)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

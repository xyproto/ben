# Brief Expressive Notation (BEN) Format Specification

The Brief Expressive Notation (BEN) format is a concise and human-readable representation of musical notes, designed for easy creation and editing of simple melodies and tracks. This document describes the format and its various commands and notations.

## Basic Notation

A BEN file consists of one or more tracks, separated by newlines. Each track is a sequence of space-separated tokens, where each token represents a musical event such as a note, rest, or octave change.

### Notes

Notes are represented by uppercase or lowercase letters A-G (or H, as an alias for B), optionally followed by an accidental (# for sharp, b for flat). 

Examples:
- C: C note in the current octave
- G#: G sharp note in the current octave
- Eb: E flat note in the current octave

Uppercase letters represent quarter notes, while lowercase letters represent eighth notes. Notes can be lengthened or shortened using the following modifiers:
- /: half note
- //: whole note
- .: staccato (half the duration of the original note)
- ^: accent (1.5 times the duration of the original note)

### Rests

Rests are represented by a comma (,) for an eighth note rest.

### Octave Changes

To change the current octave, use the following commands:
- +: increase the current octave by 1
- -: decrease the current octave by 1
- (n): set the current octave to n, where n is an integer between 0 and 9

### Velocity Changes

To change the velocity of a note, use the following command:
- v: decrease the velocity by 10%
- !n: set the velocity to n * 10%, where n is a single hexadecimal digit (0-F)

### MIDI Channel and Program Changes

To change the MIDI channel and program, use the following commands:
- [n]: set the MIDI channel to n, where n is an integer between 1 and 16
- {n}: set the MIDI program to n, where n is an integer between 0 and 127

## Examples

Here's an example of a simple BEN melody:

C D E F G A B H

And a more complex example using various modifiers:

C +C -C c c/ c// c. c^ , !A C vC {10}C [2]C

## Implementation Notes

- The channel, program, and octave changes are "sticky" and persist across notes within a track, but do not carry over to other tracks.
- When using an external synth, the software will send MIDI events directly to the specified device, bypassing the creation of a MIDI file.

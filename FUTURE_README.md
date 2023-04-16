# Brief Expressive Notation

Brief Expressive Notation is a simple and efficient text-based music notation format designed for readability and expressiveness, which includes support for microtonality.

The ben2mid utility is included, which can convert .ben files to .mid (MIDI).

## Quick installation

    go install github.com/ben/cmd/ben2mid@latest

## Description of BEN

    Lowercase letters (e.g., "c", "d", "e") represent eighth notes (0.5 beats).
    Uppercase letters (e.g., "C", "D", "E") represent quarter notes (1 beat).
    Uppercase letters followed by a forward slash (e.g., "C/", "D/", "E/") represent half notes (2 beats).
    Repeating an uppercase letter twice (e.g., "CC", "DD", "EE") represents whole notes (4 beats).
    Sharp notes are indicated with "#".
    Flat notes are indicated with "b".
    Microtonal adjustments: Use "{" for a 1.818% decrease in pitch (which results in 440Hz being reduced to 432Hz), and "}" for the opposite (an increase of 1.818%). These percentages will be applied consistently across all notes for a more uniform effect.
    Stronger dynamics are indicated with an exclamation mark (!) followed by a single hexadecimal digit (0-F).
    Octave designations are given by appending a number in parentheses, e.g., "(3)" for the 3rd octave, "(5)" for the 5th octave, and omitting any designation for the normal (4th) octave. Use "-" for "the current octave - 1" and "+" for "the current octave + 1".
    Rests are indicated with a comma (","), representing a quarter-note rest.
    Slurs and ties are indicated by placing the connected notes in parentheses, e.g., (C~D) for a slur between C and D, or (C~C) for a tied C note.
    Articulations:
        Staccato: Add a dot (.) after the note, e.g., "C.". This indicates a note should be played shorter than its full duration, with a brief silence between notes.
        Accent: Add a caret (^) after the note, e.g., "C^". This indicates a note should be played with more emphasis or force.
        Tenuto: Add an arrow (->) after the note, e.g., "C->". This indicates a note should be played with full duration, connecting smoothly with the next note. Tenuto can also be used to emphasize the connection between notes in a slur.
        Velocity: Add a v after the note, e.g., "Cv". This indicates that the velocity should be reduced by 10% for this note.

## Examples

### Example 1: Simple Melody

    C D E F G F E D C

This simple melody consists of quarter notes starting from C and going up to G, then back down to C.

### Example 2: Simple Melody, with octave changes

    C D# g b C(6) C- C- C+ f g g# a# C

### Example 3: Melody with Articulations and Dynamics

    C^ D. E-> F G! F. E^ D C

In this example, the same melody from Example 1 has been enhanced with articulations and dynamics: the first note (C) is accented, the second note (D) is staccato, the third note (E) is tenuto, and the fifth note (G) is played with stronger dynamics (90%).

### Example 4: Melody with Octave Changes and Microtonal Adjustments

    C- D{(3) E}(5) F+ G{- F E(4) D C

This example demonstrates a melody with octave changes and microtonal adjustments. The melody starts in the low octave, moves to the 3rd and 5th octaves, and returns to the normal (4th) octave.

### Example 5: Indian Raga with Microtonal Adjustments

    !P108
    C E{^+ G{.+ B{ C/! D{^+ F{.+ A{ D/! C E{^+ G{.+ B{ C/! D{^+ F{.+ A{ D/!

This example demonstrates a simple representation of an Indian raga with microtonal adjustments. The melody moves in a pattern, ascending and descending the scale with the microtonal adjustments applied to E, G, B, D, F, and A notes. The pattern is played twice, with articulations and dynamics added for expressiveness.

Keep in mind that this is a simplified representation and may not fully capture the intricacies of the original raga.

### Example 6: ACID Synth Bass

    !P38
    C(3) G#(3) C(4) G#(4) C(4) G#(4) C(4) G#(4)

This example demonstrates an ACID synth bass pattern. The bass starts with a stronger dynamics on the C note and alternates between a low C and a higher G, using higher dynamics on the higher notes. The second half of the pattern moves through a sequence of notes in the current octave and the octave below, with stronger dynamics applied to the low C note.

## General info

* Version: 1.0.2
* License: BSD-3

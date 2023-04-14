# Brief Expressive Notation (BEN)

Brief Expressive Notation (BEN) is a simple and efficient text-based music notation format designed for readability and expressiveness.

GPT-4 has been used as an assistant to help create this format.

## Format

- Lowercase letters (e.g., "c", "d", "e") represent eighth notes (0.5 beats).
- Uppercase letters (e.g., "C", "D", "E") represent quarter notes (1 beat).
- Uppercase letters followed by a forward slash (e.g., "C/", "D/", "E/") represent half notes (2 beats).
- Repeating an uppercase letter twice (e.g., "CC", "DD", "EE") represents whole notes (4 beats).
- Sharp notes are indicated with "+".
- Flat notes are indicated with "-".
- Microtonal adjustments: Use "{" for a 1.818% decrease in pitch (which results in 440Hz being reduced to 432Hz), and "}" for the opposite (an increase of 1.818%). These percentages will be applied consistently across all notes for a more uniform effect.
- Stronger dynamics are indicated with an exclamation mark (!).
- Octave designations are given by appending a number in parentheses, e.g., "(3)" for the 3rd octave, "(5)" for the 5th octave, and omitting any designation for the normal (4th) octave. Use "-" for "the current octave - 1" and "+" for "the current octave + 1".
- Rests are indicated with a comma (","), representing a quarter-note rest.
- Slurs and ties are indicated by placing the connected notes in parentheses, e.g., "(C~D)" for a slur between C and D, or "(C~C)" for a tied C note.
- Articulations:
  - Staccato: Add a dot (`.`) after the note, e.g., "C.". This indicates a note should be played shorter than its full duration, with a brief silence between notes.
  - Accent: Add a caret (`^`) after the note, e.g., "C^". This indicates a note should be played with more emphasis or force.
  - Tenuto: Add an arrow (`->`) after the note, e.g., "C->". This indicates a note should be played with full duration, connecting smoothly with the next note. Tenuto can also be used to emphasize the connection between notes in a slur.
  - Velocity: Add a `v` after the note, e.g., "Cv". This indicates that the velocty should be reduced by 10% for this note.

## Examples

### Example 1: Simple Melody

    C D E F G F E D C

This simple melody consists of quarter notes starting from C and going up to G, then back down to C.

### Example 2: Melody with Articulations and Dynamics

    C^ D. E-> F G! F. E^ D C

In this example, the same melody from Example 1 has been enhanced with articulations and dynamics: the first note (C) is accented, the second note (D) is staccato, the third note (E) is tenuto, and the fifth note (G) is played with stronger dynamics.

### Example 3: Melody with Octave Changes and Microtonal Adjustments

    C- D{(3) E}(5) F+ G{- F E(4) D C

This example demonstrates a melody with octave changes and microtonal adjustments. The melody starts in the low octave, moves to the 3rd and 5th oct

### Example 4: Indian Raga with Microtonal Adjustments

    C (E{^}+ G{.}+ B{) C/! (D{^}+ F{.}+ A{) D/! C E{^}+ G{.}+ B{ C/! D{^}+ F{.}+ A{ D/!

This example demonstrates a simple representation of an Indian raga with microtonal adjustments. The melody moves in a pattern, ascending and descending the scale with the microtonal adjustments applied to E, G, B, D, F, and A notes. The pattern is played twice, with articulations and dynamics added for expressiveness.

Keep in mind that this is a simplified representation and may not fully capture the intricacies of the original raga.

### General info

* Version: 1.0.0
* License: BSD-3

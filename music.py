from librosa import midi_to_hz as note
from ahap import AHAP, freq
import mido
import sys


# Step 1: Define initial variables
time = 0.0
duration = 0.0

# Step 2: Read the MIDI file (obtain the filename from command line argument)
if len(sys.argv) < 2:
    print("Please provide the path to the MIDI file as a command line argument.")
    sys.exit(1)
filename=sys.argv[1]
midi_file = mido.MidiFile(filename)
ahap = AHAP(f"midi file {filename}", "midi to haptic generator")


# Step 3: Convert notes to haptics
note_state = {}  # Dictionary to track note states (on/off)
for msg in midi_file:
    time += msg.time
    if msg.is_meta and hasattr(msg, "note"): continue
    if msg.type == 'note_on' and msg.velocity>0:
        note_state[msg.note] = time
    elif msg.type == 'note_off' or (msg.type=='note_on' and msg.velocity==0):  # musescore doesn't do note_off, it does note on with velocity 0.
        if msg.note not in note_state:
            print(f"Warning: Found note_off message without a corresponding note_on for note {msg.note}")
        else:
            duration = time - note_state[msg.note]
            #print(duration)
            # Add a haptic event for the note
            ahap.add_haptic_continuous_event(note_state[msg.note], duration, 1.0, freq(note(msg.note)))


# Step 4: Export the haptics to an AHAP file
output_filename = sys.argv[1].split('.')[0] + '.ahap'
ahap.export(output_filename)

# Finished! You've converted the MIDI file to haptics and saved it as '[filename].ahap'

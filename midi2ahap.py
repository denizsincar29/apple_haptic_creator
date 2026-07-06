"""Converts a MIDI (.mid) file to an AHAP haptic pattern.

Same approach as the Go (cmd/midi2ahap) and Rust (ahap_rs) versions:
GM channel-10 (index 9) drum notes get instrument-appropriate haptic shapes
instead of one flat transient for every hit, and melodic notes on other
channels are mapped to sustained continuous events with pitch -> sharpness.

Requires the `mido` package (pip install mido).
"""

import argparse
import sys
from dataclasses import dataclass
from enum import Enum, auto
from typing import Dict, Optional

import mido

from ahap import AHAP, CurveParamID, HapticCurve, create_ease_in_out_curve, freq as freq_to_sharpness


class HapticKind(Enum):
    TRANSIENT = auto()  # instantaneous snap: snares, sticks, claves, closed hats
    THUMP = auto()      # short felt punch with a bit of body: kicks, toms, congas
    RINGING = auto()    # long decaying tail: cymbals, open hi-hat, tambourine, triangle


@dataclass
class DrumMapping:
    kind: HapticKind
    intensity: float
    sharpness: float
    duration: float = 0.0  # only used for THUMP/RINGING


# General MIDI drum mappings (channel 10). Bass drums/toms/congas are THUMP
# (short Continuous + decay envelope so they feel like a punch, not a click);
# cymbals/open hi-hat/tambourine/etc are RINGING (long Continuous + a fading
# intensity curve so they actually ring out); everything else stays a crisp
# instantaneous TRANSIENT.
DRUM_MAPPINGS: Dict[int, DrumMapping] = {
    35: DrumMapping(HapticKind.THUMP, 1.0, 0.15, 0.09),   # acoustic bass drum
    36: DrumMapping(HapticKind.THUMP, 1.0, 0.15, 0.09),   # bass drum 1
    38: DrumMapping(HapticKind.TRANSIENT, 0.95, 0.85),    # acoustic snare
    40: DrumMapping(HapticKind.TRANSIENT, 0.9, 0.9),      # electric snare
    41: DrumMapping(HapticKind.THUMP, 0.85, 0.30, 0.07),  # low floor tom
    43: DrumMapping(HapticKind.THUMP, 0.85, 0.35, 0.065), # high floor tom
    45: DrumMapping(HapticKind.THUMP, 0.85, 0.40, 0.06),  # low tom
    47: DrumMapping(HapticKind.THUMP, 0.85, 0.45, 0.055), # low-mid tom
    48: DrumMapping(HapticKind.THUMP, 0.85, 0.50, 0.05),  # hi-mid tom
    50: DrumMapping(HapticKind.THUMP, 0.85, 0.55, 0.045), # high tom
    42: DrumMapping(HapticKind.TRANSIENT, 0.5, 1.0),      # closed hi-hat
    44: DrumMapping(HapticKind.TRANSIENT, 0.5, 0.95),     # pedal hi-hat
    46: DrumMapping(HapticKind.RINGING, 0.6, 0.9, 0.25),  # open hi-hat
    49: DrumMapping(HapticKind.RINGING, 0.9, 0.85, 0.6),  # crash cymbal 1
    51: DrumMapping(HapticKind.RINGING, 0.6, 0.75, 0.35), # ride cymbal 1
    52: DrumMapping(HapticKind.RINGING, 0.85, 0.8, 0.55), # chinese cymbal
    53: DrumMapping(HapticKind.TRANSIENT, 0.65, 0.7),     # ride bell
    55: DrumMapping(HapticKind.RINGING, 0.75, 0.9, 0.3),  # splash cymbal
    57: DrumMapping(HapticKind.RINGING, 0.9, 0.85, 0.6),  # crash cymbal 2
    59: DrumMapping(HapticKind.RINGING, 0.6, 0.75, 0.35), # ride cymbal 2
    37: DrumMapping(HapticKind.TRANSIENT, 0.7, 0.95),     # side stick
    39: DrumMapping(HapticKind.TRANSIENT, 0.75, 0.8),     # hand clap
    54: DrumMapping(HapticKind.RINGING, 0.6, 0.85, 0.15), # tambourine
    56: DrumMapping(HapticKind.TRANSIENT, 0.7, 0.7),      # cowbell
    58: DrumMapping(HapticKind.RINGING, 0.65, 0.75, 0.3), # vibraslap
    60: DrumMapping(HapticKind.THUMP, 0.75, 0.6, 0.04),   # hi bongo
    61: DrumMapping(HapticKind.THUMP, 0.75, 0.5, 0.05),   # low bongo
    62: DrumMapping(HapticKind.TRANSIENT, 0.75, 0.65),    # mute hi conga
    63: DrumMapping(HapticKind.THUMP, 0.75, 0.6, 0.05),   # open hi conga
    64: DrumMapping(HapticKind.THUMP, 0.75, 0.55, 0.06),  # low conga
    65: DrumMapping(HapticKind.THUMP, 0.8, 0.7, 0.04),    # high timbale
    66: DrumMapping(HapticKind.THUMP, 0.8, 0.65, 0.05),   # low timbale
    67: DrumMapping(HapticKind.TRANSIENT, 0.7, 0.8),      # high agogo
    68: DrumMapping(HapticKind.TRANSIENT, 0.7, 0.75),     # low agogo
    69: DrumMapping(HapticKind.TRANSIENT, 0.55, 0.7),     # cabasa
    70: DrumMapping(HapticKind.TRANSIENT, 0.5, 0.85),     # maracas
    71: DrumMapping(HapticKind.TRANSIENT, 0.6, 0.9),      # short whistle
    72: DrumMapping(HapticKind.RINGING, 0.6, 0.85, 0.2),  # long whistle
    73: DrumMapping(HapticKind.TRANSIENT, 0.65, 0.75),    # short guiro
    74: DrumMapping(HapticKind.RINGING, 0.65, 0.7, 0.15), # long guiro
    75: DrumMapping(HapticKind.TRANSIENT, 0.7, 0.95),     # claves
    76: DrumMapping(HapticKind.TRANSIENT, 0.7, 0.8),      # hi wood block
    77: DrumMapping(HapticKind.TRANSIENT, 0.7, 0.75),     # low wood block
    78: DrumMapping(HapticKind.TRANSIENT, 0.65, 0.7),     # mute cuica
    79: DrumMapping(HapticKind.THUMP, 0.65, 0.65, 0.06),  # open cuica
    80: DrumMapping(HapticKind.TRANSIENT, 0.55, 0.9),     # mute triangle
    81: DrumMapping(HapticKind.RINGING, 0.55, 0.95, 0.3), # open triangle
}

DRUM_CHANNEL = 9  # GM channel 10, 0-indexed


def midi_note_to_freq(note: int) -> float:
    return 440.0 * 2 ** ((note - 69) / 12.0)


def notes_for_low_pitch(note: int, floor_hz: float = 80.0) -> list:
    """The Taptic Engine's continuous events only track frequency down to
    ~80 Hz; below that a single tone doesn't read as a pitch anymore. So
    for low notes, shift up by octaves until the root clears the floor,
    then add a fourth below it as a second simultaneous note - e.g. C2
    becomes C3+G2. Two notes a fourth apart perceptually still reads as
    "that low note" much better than one out-of-range tone.
    """
    root = note
    while midi_note_to_freq(root) < floor_hz:
        root += 12
    if root == note:
        return [note]
    return [root, root - 5]


def add_drum_hit(ahap: AHAP, t: float, mapping: DrumMapping, intensity: float) -> None:
    """Renders one drum hit according to its instrument kind - the core of
    "realistic" drums: kicks/toms get a short felt punch (Continuous + decay
    envelope), cymbals/open hi-hat/etc get a long Continuous event with a
    fading intensity curve, and only snares/sticks/etc stay a flat
    instantaneous Transient."""
    if mapping.kind is HapticKind.THUMP:
        decay = mapping.duration * 0.6
        release = mapping.duration * 0.4
        ahap.add_haptic_continuous_event(
            t, mapping.duration, intensity, mapping.sharpness,
            attack=0.0, decay=decay, release=release,
        )

    elif mapping.kind is HapticKind.RINGING:
        ahap.add_haptic_continuous_event(t, mapping.duration, intensity, mapping.sharpness)
        # HapticIntensityControl multiplies the event's base HapticIntensity
        # (output = intensity * curve), so this ramps 1.0 -> 0.0, not
        # intensity -> 0. Anchored at relative time 0 so the ring starts at
        # full strength immediately instead of holding at the first
        # interpolated point's value (create_curve only emits points
        # strictly after the start).
        anchor = HapticCurve(0.0, 1.0)
        ramp = create_ease_in_out_curve(0.0, mapping.duration, 1.0, 0.0, total=6)
        ahap.add_parameter_curve(CurveParamID.H_Intensity, t, [anchor] + ramp)

    else:  # HapticKind.TRANSIENT
        ahap.add_haptic_transient_event(t, intensity, mapping.sharpness)


def convert(
    input_path: str,
    output_path: str,
    no_drums: bool = False,
    drums_as_melody: bool = False,
    debug_channels: bool = False,
    indent: Optional[int] = None,
) -> dict:
    """Converts a MIDI file to AHAP. Returns a stats dict.

    Channel 10 (GM drums) handling, in order of precedence:
    - default: realistic per-instrument drum haptics (see add_drum_hit)
    - no_drums=True: channel 10 is fully ignored, no events at all from it
    - drums_as_melody=True: channel 10 notes are treated as regular melodic
      notes instead (mutually exclusive with no_drums)
    """
    if no_drums and drums_as_melody:
        raise ValueError("no_drums and drums_as_melody are mutually exclusive")

    mid = mido.MidiFile(input_path)
    ahap = AHAP(description=f"midi file {input_path}", created_by="midi to haptic generator (python)")

    drum_count = 0
    unknown_drum_count = 0
    melodic_count = 0
    channel_counts: Dict[int, int] = {}

    for track in mid.tracks:
        current_time = 0.0
        # Keyed by (channel, note), not just note: a MIDI file with several
        # channels active at once (very common - almost every real-world
        # song MIDI has one channel per instrument) can easily have two
        # channels holding the same pitch simultaneously. Keying by note
        # alone let a note-on on one channel clobber another channel's
        # still-open note of the same pitch, corrupting durations or
        # dropping notes outright.
        note_on_times: Dict[tuple, tuple] = {}  # (channel, note) -> (start_time, velocity)

        for msg in track:
            # mido gives delta times in ticks per message; convert using the
            # file's ticks_per_beat and mido's own tempo-aware tick2second,
            # which already integrates per tempo segment correctly.
            current_time += mido.tick2second(msg.time, mid.ticks_per_beat, _current_tempo(track))

            if msg.type == "set_tempo":
                _set_current_tempo(track, msg.tempo)

            if msg.type == "note_on" and msg.velocity > 0:
                is_drum_channel = msg.channel == DRUM_CHANNEL
                channel_counts[msg.channel] = channel_counts.get(msg.channel, 0) + 1

                if is_drum_channel and no_drums:
                    pass  # fully ignored: no event, no note_on_times entry
                elif is_drum_channel and not drums_as_melody:
                    velocity_scale = msg.velocity / 127.0
                    mapping = DRUM_MAPPINGS.get(msg.note)
                    if mapping is not None:
                        add_drum_hit(ahap, current_time, mapping, mapping.intensity * velocity_scale)
                        drum_count += 1
                    else:
                        ahap.add_haptic_transient_event(current_time, velocity_scale, 0.7)
                        drum_count += 1
                        unknown_drum_count += 1
                else:
                    # Either a melodic channel, or channel 10 with drums_as_melody.
                    note_on_times[(msg.channel, msg.note)] = (current_time, msg.velocity)

            # A note_on with velocity 0 is a note_off per the MIDI spec.
            elif (msg.type == "note_off") or (msg.type == "note_on" and msg.velocity == 0):
                is_drum_channel = msg.channel == DRUM_CHANNEL
                treated_as_melodic = not is_drum_channel or drums_as_melody
                if treated_as_melodic:
                    start = note_on_times.pop((msg.channel, msg.note), None)
                    if start is not None:
                        start_time, velocity = start
                        duration = current_time - start_time
                        if duration > 0:
                            for haptic_note in notes_for_low_pitch(msg.note):
                                try:
                                    sharpness = freq_to_sharpness(midi_note_to_freq(haptic_note))
                                except ValueError:
                                    sharpness = 0.5
                                ahap.add_haptic_continuous_event(start_time, duration, velocity / 127.0, sharpness)
                            melodic_count += 1

    if debug_channels:
        print("Note-on events per channel (1-indexed for readability):")
        for channel in sorted(channel_counts):
            count = channel_counts[channel]
            if count > 0:
                marker = " (GM drum channel)" if channel == DRUM_CHANNEL else ""
                print(f"  channel {channel + 1}: {count} events{marker}")

    ahap.export(output_path, path=".", indent=indent)

    return {
        "drum_events": drum_count,
        "unknown_drum_events": unknown_drum_count,
        "melodic_events": melodic_count,
        "total_events": drum_count + melodic_count,
    }


# mido's message.time is already in seconds if you iterate `mid` directly
# (mid.__iter__ merges tracks and does the tempo math for you); iterating a
# raw MidiTrack gives ticks instead, so we track tempo ourselves per track
# using a tiny mutable cache keyed by id(track) to avoid a global.
_tempo_cache: Dict[int, int] = {}


def _current_tempo(track) -> int:
    return _tempo_cache.get(id(track), 500000)  # default 120 BPM


def _set_current_tempo(track, tempo: int) -> None:
    _tempo_cache[id(track)] = tempo


def main() -> None:
    parser = argparse.ArgumentParser(description="Convert a MIDI file to an AHAP haptic pattern.")
    parser.add_argument("input", help="input .mid file")
    parser.add_argument("output", nargs="?", default=None, help="output .ahap file (default: <input>.ahap)")
    parser.add_argument("--no-drums", action="store_true", help="completely ignore channel 10 (GM drums) - no events at all from it")
    parser.add_argument("--drums-as-melody", action="store_true", help="treat channel 10 as regular melodic notes instead of GM drums (rather than dropping it - see --no-drums)")
    parser.add_argument("--debug-channels", action="store_true", help="print how many note-on events came from each channel")
    parser.add_argument("--indent", type=int, default=None, help="indent the JSON output for readability")
    args = parser.parse_args()

    output = args.output
    if output is None:
        stem = args.input.rsplit(".", 1)[0]
        output = f"{stem}.ahap"

    if args.no_drums and args.drums_as_melody:
        parser.error("--no-drums and --drums-as-melody are mutually exclusive")

    stats = convert(
        args.input, output,
        no_drums=args.no_drums,
        drums_as_melody=args.drums_as_melody,
        debug_channels=args.debug_channels,
        indent=args.indent,
    )

    print(f"Successfully created {output}")
    print("Conversion statistics:")
    print(f"  Drum events: {stats['drum_events']}")
    if stats["unknown_drum_events"]:
        print(f"    (including {stats['unknown_drum_events']} unmapped drum notes)")
    print(f"  Melodic events (continuous): {stats['melodic_events']}")
    print(f"  Total haptic events: {stats['total_events']}")


if __name__ == "__main__":
    main()

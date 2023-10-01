from librosa import note_to_hz as note
from ahap import AHAP, CurveParamID, HapticCurve, create_curve, freq


ahap=AHAP("do re mi", "Deniz Sincar")
time=0.0
dur=0.4
doremi=note(["e2", "f2", "g2", "a2", "b2", "c3", "d3", "e3", "f3", "g3", "a3"])

for i, anote  in enumerate(doremi):
    ahap.add_haptic_continuous_event(dur*i, dur, 1.0, freq(anote))

ahap.export("music.ahap")
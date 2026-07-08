import unittest
from ahap import freq, AHAP, create_curve, create_ease_in_out_curve

class TestFreq(unittest.TestCase):
    #def setUp(self) -> None:
    
    def test_freq(self):
        for i in range(80, 231):
            try:
                freq(i, False)
            except Exception as e:
                self.fail(f'the {i} hz freq not converted., exception: {e}')

    def test_raisefreq(self):
        with self.assertRaises(ValueError):
            freq(79, False)
            freq(231, False)


class TestEnvelope(unittest.TestCase):
    def test_transient_omits_none_envelope_fields(self):
        a = AHAP("t", "t")
        a.add_haptic_transient_event(0.0, 1.0, 0.5)
        params = a.data["Pattern"][0]["Event"]["EventParameters"]
        self.assertEqual(len(params), 2)

    def test_transient_includes_set_envelope_fields_only(self):
        a = AHAP("t", "t")
        a.add_haptic_transient_event(0.0, 1.0, 0.5, attack=0.01, decay=0.02)
        params = a.data["Pattern"][0]["Event"]["EventParameters"]
        ids = [p["ParameterID"] for p in params]
        self.assertIn("HapticAttackTime", ids)
        self.assertIn("HapticDecayTime", ids)
        self.assertNotIn("HapticReleaseTime", ids)

    def test_continuous_has_duration(self):
        a = AHAP("t", "t")
        a.add_haptic_continuous_event(1.0, 2.0, 0.7, 0.6)
        event = a.data["Pattern"][0]["Event"]
        self.assertEqual(event["EventDuration"], 2.0)


class TestCurves(unittest.TestCase):
    def test_linear_curve_endpoints(self):
        points = create_curve(0.0, 1.0, 0.0, 1.0, total=10)
        self.assertEqual(len(points), 10)
        self.assertAlmostEqual(points[0].time, 0.1)
        self.assertAlmostEqual(points[-1].time, 1.0)
        self.assertAlmostEqual(points[-1].parameter_value, 1.0)

    def test_ease_in_out_matches_linear_at_endpoints(self):
        points = create_ease_in_out_curve(0.0, 0.6, 1.0, 0.0, total=6)
        self.assertEqual(len(points), 6)
        self.assertAlmostEqual(points[-1].time, 0.6)
        self.assertAlmostEqual(points[-1].parameter_value, 0.0, places=9)
        # smoothstep midpoint should be exactly 0.5 for a symmetric ramp
        self.assertAlmostEqual(points[2].parameter_value, 0.5)


class TestGlobalControl(unittest.TestCase):
    def test_cc_values_map_to_fractions_and_offset(self):
        from midi2ahap import GlobalControl

        control = GlobalControl()
        self.assertTrue(control.apply_cc(73, 127))
        self.assertAlmostEqual(control.attack, 1.0)

        self.assertTrue(control.apply_cc(72, 0))
        self.assertAlmostEqual(control.release, 0.0)

        self.assertTrue(control.apply_cc(75, 64))
        self.assertTrue(0.49 < control.decay < 0.51)

        self.assertTrue(control.apply_cc(74, 64))  # ~center, near-zero offset
        self.assertTrue(0.45 < control.adjust_sharpness(0.5) < 0.55)

        self.assertFalse(control.apply_cc(1, 100))  # unrelated CC (mod wheel) ignored

    def test_no_cc_seen_means_no_override(self):
        from midi2ahap import GlobalControl

        control = GlobalControl()
        self.assertIsNone(control.attack)
        self.assertIsNone(control.decay)
        self.assertIsNone(control.release)
        self.assertEqual(control.adjust_sharpness(0.42), 0.42)

    def test_cc_release_never_exceeds_short_note_duration(self):
        # This is the exact bug scenario: CC72=100 arriving once at tick 0,
        # then short drum/melodic notes (0.15-0.36s) later in the file.
        from midi2ahap import GlobalControl

        control = GlobalControl()
        control.apply_cc(72, 100)

        short_note_duration = 0.15
        release = control.release_for(short_note_duration)
        self.assertLessEqual(release, short_note_duration)
        self.assertAlmostEqual(release, short_note_duration * (100.0 / 127.0))

        # A longer note gets a proportionally longer release, not the same
        # fixed absolute time as the short note.
        long_note_duration = 1.2
        long_release = control.release_for(long_note_duration)
        self.assertAlmostEqual(long_release, long_note_duration * (100.0 / 127.0))
        self.assertGreater(long_release, release)


class TestMidi2Ahap(unittest.TestCase):
    def test_low_note_splits_into_root_and_fourth(self):
        from midi2ahap import notes_for_low_pitch

        self.assertEqual(notes_for_low_pitch(36), [48, 43])  # C2 -> C3+G2
        self.assertEqual(notes_for_low_pitch(60), [60])       # C4 stays single

    def test_drum_kinds_produce_expected_shapes(self):
        from midi2ahap import add_drum_hit, DRUM_MAPPINGS, HapticKind, GlobalControl

        a = AHAP("t", "t")
        control = GlobalControl()
        add_drum_hit(a, 0.0, DRUM_MAPPINGS[36], 1.0, control)   # kick -> THUMP
        add_drum_hit(a, 0.25, DRUM_MAPPINGS[38], 0.9, control)  # snare -> TRANSIENT
        add_drum_hit(a, 0.5, DRUM_MAPPINGS[49], 0.9, control)   # crash -> RINGING

        kinds = [p["Event"]["EventType"] for p in a.data["Pattern"] if "Event" in p]
        self.assertEqual(kinds, ["HapticContinuous", "HapticTransient", "HapticContinuous"])

        curves = [p for p in a.data["Pattern"] if "ParameterCurve" in p]
        self.assertEqual(len(curves), 1, "only the ringing (crash) hit should add a decay curve")
        self.assertEqual(curves[0]["ParameterCurve"]["ParameterID"], "HapticIntensityControl")

if __name__=="__main__":
    unittest.main()
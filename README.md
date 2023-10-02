# Apple Haptic Creator

This Python class allows you to create Apple Haptic pattern files. With this class, you can define haptic and audio events, as well as parameter curves, to create custom haptic patterns for devices that support the Apple Haptic API.

## What are AHAP files?

AHAP files are JSON-formatted special Apple Haptic pattern files. They are commonly used in iOS games and applications to create immersive experiences. However, I recently discovered that you can play AHAP files directly from the Files app or any other apps that support Apple's Quick Look API. This means that you can freely share an AHAP file via platforms like WhatsApp or Telegram. It's worth noting that WhatsApp has limitations on loading large AHAP files. I wrote an article about Apple Haptics on [Applevis](https://applevis.com/forum/ios-ipados/now-possible-ios-17-can-play-haptic-signals-vibrations-special-ahap-apple-haptic). Feel free to check it out for more information.

## What's in the repo

- ahaps/: Examples folder.
- ahap.py: Module for creating AHAP (Apple Haptic) files.
- makeahap.py: A file that creates a motorcycle sound with vibrations.
- music.py: An attempt to create musical notes via haptics, but failed.

## Requirements

My script will run on Python 3.6+ and doesn't require any additional modules. However, if you want to run music.py, you can install the librosa module by running the following command:
```bash
pip install librosa
```

## How to Use
```python
# Create an instance of the AHAP class to start creating an AHAP file.
from ahap import AHAP, CurveParamID, HapticCurve, create_curve

ahap = AHAP()

# Add events to the pattern by using the available methods of the AHAP class. For example, to add a haptic continuous event:
ahap.add_haptic_continuous_event(time=0.5, event_duration=1.0, haptic_intensity=0.8, haptic_sharpness=0.5)

# Add parameter curves to dynamically change haptic or audio parameters over time. For example, to add a haptic sharpness curve:
curve = create_curve(start_time=0.0, end_time=1.0, start_value=0.4, end_value=0.8, total=10)
ahap.add_parameter_curve(CurveParamID.H_Sharpness, start_time=0.0, control_points=curve)

# Export the AHAP file by calling the export() method.
ahap.export(filename="example.ahap")
```

You can run the makeahap.py file to generate a sample AHAP file with a truly great motorcycle sound!

## Examples

The ahaps/ folder contains example AHAP files that you can use as a reference or starting point for creating your own haptic patterns.

## Known Limitations

- The music.py file does not currently generate musical notes via haptics as intended. Further development is required to achieve this functionality.

## Contributing

Contributions are welcome! If you have an idea for an improvement or found a bug, please open an issue on GitHub or submit a pull request.
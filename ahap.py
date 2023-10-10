from enum import Enum
import datetime
import math
import os
import json
from typing import Any, List, Tuple

class HapticCurve:
    """Represents the haptic curve"""
    def __init__(self, time: float, parameter_value: float):
        """
        Initialize a HapticCurve object.

        Args:
            time (float): The time value for the curve.
            parameter_value (float): The parameter value for the curve.
                Should be a float between 0 and 1.
        """
        self.time = time
        self.parameter_value = parameter_value

    def get_data(self):
        """
        Get the data of the HapticCurve.

        Returns:
            dict: A dictionary containing the time and parameter value of the curve.
        """
        data = {"Time": self.time, "ParameterValue": self.parameter_value}
        return data

    def __call__(self, *args: Any, **kwds: Any) -> Any:
        return self.get_data()
    def __repr__(self): return repr(self.get_data())

def curves(c: List[HapticCurve]) -> List[dict]:
    """
    Convert a list of HapticCurve objects to a list of dictionaries.

    Args:
        c (List[HapticCurve]): The list of HapticCurve objects.

    Returns:
        List[dict]: A list of dictionaries containing the time and parameter value of each curve.
    """
    return [i.get_data() for i in c]

class CurveParamID(Enum):
    H_Intensity = "HapticIntensityControl"
    H_Sharpness = "HapticSharpnessControl"
    H_AttackTime = "HapticAttackTimeControl"
    H_DecayTime = "HapticDecayTimeControl"
    H_ReleaseTime = "HapticReleaseTimeControl"
    A_Brightness = "AudioBrightnessControl"
    A_Pan = "AudioPanControl"
    A_Pitch = "AudioPitchControl"
    A_Volume = "AudioVolumeControl"
    A_AttackTime = "AudioAttackTimeControl"
    A_DecayTime = "AudioDecayTimeControl"
    A_ReleaseTime = "AudioReleaseTimeControl"

class ParamID(Enum):
    H_Intensity = "HapticIntensity"
    H_Sharpness = "HapticSharpness"
    H_AttackTime = "HapticAttackTime"
    H_DecayTime = "HapticDecayTime"
    H_ReleaseTime = "HapticReleaseTime"
    A_Brightness = "AudioBrightness"
    A_Pan = "AudioPan"
    A_Pitch = "AudioPitch"
    A_Volume = "AudioVolume"
    A_AttackTime = "AudioAttackTime"
    A_DecayTime = "AudioDecayTime"
    A_ReleaseTime = "AudioReleaseTime"

# soon we will do it a @classmethod, but it'll break compatibility so i'm lazy!
def create_curve(start_time: float, end_time: float, start_value: float, end_value: float, total=10):
    timediff=end_time-start_time
    valuediff=end_value-start_value
    timestep=timediff/total
    valuestep=valuediff/total
    curvelist=[]
    for i in range(total):
        curvelist.append(HapticCurve(start_time+timestep*(i+1), start_value+valuestep*(i+1)))
    #print("start time", start_time, "endtime", end_time)
    return curvelist


class AHAP:
    """_Class that allows to make Apple haptic signal files (.ahap)."""
    def __init__(self, description: str = "test AHAP file", created_by: str = "Deniz Sincar"):
        """
        Initialize an AHAP object.

        Args:
            description (str): The description of the AHAP file.
            created_by (str): The creator of the AHAP file.
        """
        self.data = {
            "Version": 1.0,
            "Metadata": {
                "Project": "Basis",
                "Created": str(datetime.datetime.now()),
                "Description": description,
                "Created By": created_by
            },
            "Pattern": []
        }

    def add_event(self, etype: str, time: float, parameters: List[dict], event_duration: float = None, event_waveform_path: str = None):
        """
        Adds an event to the pattern.

        Args:
            etype (str): The type of event.
                Possible values: "AudioContinuous", "AudioCustom", "HapticTransient", and "HapticContinuous".
            time (float): The time of the event in seconds.
            parameters (List[dict]): The event parameters as a list of dictionaries.
        """
        pattern = {
            "Event": {
                "Time": time,
                "EventType": etype,
                "EventParameters": parameters
            }
        }
        if event_duration is not None:
            pattern["Event"]["EventDuration"] = event_duration
        if event_waveform_path is not None:
            pattern["Event"]["EventWaveformPath"] = event_waveform_path
        self.data["Pattern"].append(pattern)

    def __rshift__(self, args: Tuple):
        self.add_event(*args)

    def add_haptic_transient_event(self, time: float, haptic_intensity: float = 0.5, haptic_sharpness: float = 0.5):
        """
        Adds a haptic transient event to the pattern.

        Args:
            time (float): The time of the event in seconds.
            haptic_intensity (float): The intensity of the haptic event.
                Should be a float between 0 and 1.
            haptic_sharpness (float): The sharpness of the haptic event.
                Should be a float between 0 and 1.
        """
        parameters = [
            {
                "ParameterID": ParamID.H_Intensity.value,
                "ParameterValue": haptic_intensity,
            },
            {
                "ParameterID": ParamID.H_Sharpness.value,
                "ParameterValue": haptic_sharpness,
            }
        ]

        self.add_event(etype="HapticTransient", time=time, parameters=parameters)

    def add_haptic_continuous_event(self, time: float, event_duration: float = 1, haptic_intensity: float = 0.5, haptic_sharpness: float = 0.5):
        """
        Adds a haptic continuous event to the pattern.

        Args:
            time (float): The time of the event in seconds.
            event_duration (float): The duration of the haptic event in seconds.
            haptic_intensity (float): The intensity of the haptic event.
                Should be a float between 0 and 1.
            haptic_sharpness (float): The sharpness of the haptic event.
                Should be a float between 0 and 1.
        """
        parameters = [
            {
                "ParameterID": ParamID.H_Intensity.value,
                "ParameterValue": haptic_intensity,
            },
            {
                "ParameterID": ParamID.H_Sharpness.value,
                "ParameterValue": haptic_sharpness,
            }
        ]

        self.add_event(etype="HapticContinuous", time=time, parameters=parameters, event_duration=event_duration)

    def add_audio_custom_event(self, time: float, wav_filepath: str, volume: float = 0.75):
        """
        Adds an audio custom event to the pattern.

        Args:
            time (float): The time of the event in seconds.
            wav_filepath (str): The path to the WAV file containing the sound.
            volume (float): The volume of the audio event.
                Should be a float between 0 and 1.
        """
        parameters = [
            {
                "ParameterID": ParamID.A_Volume.value,
                "ParameterValue": volume,
            }
        ]
        self.add_event(etype="AudioCustom", time=time, parameters=parameters, event_waveform_path=wav_filepath)

    def add_parameter_curve(self, parameter_id: CurveParamID, start_time: float, control_points: List[HapticCurve]):
        """
        Adds a parameter curve to the pattern.

        Args:
            parameter_id (CurveParamID): The parameter to dynamically change.
                Possible values are one of: H_Intensity, H_Sharpness, H_AttackTime, H_DecayTime, H_ReleaseTime,
                A_Brightness, A_Pan, A_Pitch, A_Volume, A_AttackTime, A_DecayTime, A_ReleaseTime.
            start_time (float): The time of the start of the curve in seconds.
            control_points (List[HapticCurve]): The list of control points for the curve.
                Should be a list of HapticCurve objects.
        """
        pattern = {
            "ParameterCurve": {
                "ParameterID": parameter_id.value,
                "Time": start_time,
                "ParameterCurveControlPoints": curves(control_points)
            }
        }

        self.data["Pattern"].append(pattern)

    def __repr__(self):
        """
        Print the data of the AHAP object.
        """
        repr(self.data)

    def export(self, filename: str, path: str = ".", **kwargs):
        """
        Export the AHAP object to a JSON file.

        Args:
            filename (str): The name of the output file.
            path (str): The path to the output directory.
            **kwargs: Extra arguments you want to pass on to json.dumps(). For example, indent=4 for a pretty formatted JSON. 
        """
        with open(os.path.join(path, filename), 'w') as f:
            f.write(json.dumps(self.data, **kwargs))

    def __call__(self, *args: Any, **kwds: Any) -> Any:
        self.export(*args, **kwds)

    def __add__(self, other: AHAP):
        """adds 2 ahap files. Attension, it smooshes them one on another, it doesn't work as expected now. Please don't use this method if you don't want to really smoosh them.

        Args:
            other (AHAP): another Ahap class.
        """
        data = {
            "Version": 1.0,
            "Metadata": {
                "Project": "Basis",
                "Created": str(datetime.datetime.now()),
                "Description": self.data["Metadata"]["Description"],
                "Created By": self.data["Metadata"]["Created By"]
            },
            "Pattern": self.data["Pattern"]+other.data["Pattern"]
        }

def freq(n: int, normalize: bool=True) -> float:
    """
    calculates the haptic sharpness value from frequency in hz.

    Args:
        n (int): The input frequency value.
        normalize (bool): if normalizing, all high frequencies will be 230 and all low will be 80 if value is too high or too low.
    Returns:
        float: The normalized frequency value between 0 and 1.

    Raises:
        ValueError: If the input frequency is less than 80 or greater than 230.
        ValueError: If the calculated normalized frequency is less than 0 or greater than 1.
    """
    if normalize and n>230: n=230
    if normalize and n<80: n=80
    if n < 80 or n > 230:
        raise ValueError(f"Incorrect frequency. Frequency must be between 80 and 230, but it is {n}")
    r = (math.log(n) - math.log(80)) / (math.log(230) - math.log(80))
    if r < 0 or r > 1:
        raise ValueError("The calculated normalized frequency is out of range. Result must be between 0 and 1.")
    return r
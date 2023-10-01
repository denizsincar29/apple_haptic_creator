from ahap import AHAP, CurveParamID, HapticCurve, create_curve


ahap=AHAP("bike sound", "Deniz Sincar")
time=0.0
dur=0.4
ahap.add_haptic_continuous_event(time, dur, 0.5, 0.4)
ahap.add_parameter_curve(CurveParamID.H_Sharpness, time, create_curve(0.0, 0.4, 0.4, 0.75, 10))
time=0.45
for i in range(7):
    ahap.add_haptic_transient_event(time, 1.0, 0.3)
    time+=0.05
ahap.add_haptic_continuous_event(time, 15.0, 0.75, 0.0)
for i in range(300):
ahap.add_haptic_transient_event(time+i*0.05 , 1.0, 1.0)
ahap.add_parameter_curve(CurveParamID.H_Sharpness, time, create_curve(0.0, 0.4, 0.0, 0.75, 10))
time+=0.4
ahap.add_parameter_curve(CurveParamID.H_Sharpness, time, create_curve(0.0, 0.8, 0.75, 0.2))
time+=0.8
ahap.add_parameter_curve(CurveParamID.H_Sharpness, time, create_curve(0, 3, 0.0, 0.5))
ahap.add_parameter_curve(CurveParamID.H_Sharpness, time+3, create_curve(0, 3, 0.2, 0.65))
ahap.add_parameter_curve(CurveParamID.H_Sharpness, time+6, create_curve(0, 4, 0.4, 1.0))
ahap.add_parameter_curve(CurveParamID.H_Sharpness, time+10, create_curve(0, 2, 1.0, 0.0))


ahap.export("bike.ahap")
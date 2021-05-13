# ADSK-SA-2020-0003-Tools# ADSK-SA-2020-0003-Tools

Scripts for detection and eradication of a third-party malicious script enbedded in Autodesk Maya scene files.

The malware is malicious and destructive, modifying the source scene file upon save by deleting script nodes within the scene, as well as modification of the user's environment, startup scripts, and multiple mel files in the Maya installation.

The issue was identified and a fix has been made available by Autodesk and labeled as Autodesk ID: ADSK-SA-2020-0003. The script can execute malicious code that can corrupt the Maya environment, causing data loss and instability, as well as spreading to other systems.

The code provided herein can be useful in both detection and eradication but should be used in conjunction with the fix provided by Autodesk as described in this Autodesk Security Advisory: https://www.autodesk.com/trust/security-advisories/adsk-sa-2020-0003

IMPORTANT limitations of the Autodesk code include the inability to detect malicous code within Maya binary (.mb) files, which is addressed in our solution.

You can access the Autodesk mitigation tools here: https://apps.autodesk.com/MAYA/en/Detail/Index?id=8637238041954239715

Our tools are provided as source and a single Windows, Linux and MacOS binaries which can be run on volumens containing infected files, but we advise you use this code to construct inline mitigation tools that automate this process within your environment.

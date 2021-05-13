# ADSK-SA-2020-0003-Tools

Scripts for detection and eradication of a third-party malicious script enbedded in Autodesk Maya scene files.

The malware is malicious and destructive, modifying the source scene file upon save by deleting script nodes within the scene, as well as modification of the user's environment, startup scripts, and multiple mel files in the Maya installation.

The issue was identified and a fix has been made available by Autodesk and labeled as Autodesk ID: ADSK-SA-2020-0003. The script can execute malicious code that can corrupt the Maya environment, cause data loss and instability, as well as spread to other systems.

The code provided herein can be useful in both detection and eradication but should be used in conjunction with the fix provided by Autodesk as described in this Autodesk Security Advisory: https://www.autodesk.com/trust/security-advisories/adsk-sa-2020-0003

You can access the Autodesk mitigation tools here: https://apps.autodesk.com/MAYA/en/Detail/Index?id=8637238041954239715

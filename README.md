# ADSK-SA-2020-0003-Tools

Source and Go binaries for detection and eradication of a third-party malicious script enbedded in Autodesk Maya scene files indicated as SEVERE.  A vulnerability, which if exploited, would directly impact the confidentiality, integrity or availability of user’s data or processing resources.

The malware is malicious and destructive, modifying the source scene file upon save by deleting script nodes within the scene, as well as modification of the user's environment, startup scripts, and multiple mel files in the Maya installation.

The issue was identified and a fix has been made available by Autodesk and labeled as Autodesk ID: ADSK-SA-2020-0003. The script can execute malicious code that can corrupt the Maya environment, causing data loss and instability, as well as spreading to other systems.

The code provided herein can be useful in both detection and eradication but should be used in conjunction with the fix provided by Autodesk as described in this Autodesk Security Advisory: https://www.autodesk.com/trust/security-advisories/adsk-sa-2020-0003

IMPORTANT limitations of the Autodesk code include the inability to detect malicous code within Maya binary (.mb) files, which is addressed in our solution.

You can access the Autodesk mitigation tools here: https://apps.autodesk.com/MAYA/en/Detail/Index?id=8637238041954239715

Our tools are provided as both source and as a single Windows, Linux or MacOS binary which can be run on volumens containing infected files.  Infected files in the current user's home directory will also be located.  

We advise you use this code to construct inline mitigation tools that automate this process within your environment.

# Usage

    dephage [-c | -v] [root-folder]

* Detects and optionally cleans the ADSK-SA-2020-0003 Autodesk Maya virus.

* Infected text .ma and .mb files will be cleaned and the original file
  renamed with a .INFECTED extension.

* Infected binary .mb files will NOT be cleaned and the file
  renamed with a .INFECTED extension.

## Flags
  -c	detect and clean (default is detect only)
  -v	version

## Examples

Detect from current folder.

    dephage

Detect from selected folder.

    dephage documents/maya

Detect and clean from current folder.

    dephage -c

Detect and clean selected folder.

    `dephage -c documents/maya`









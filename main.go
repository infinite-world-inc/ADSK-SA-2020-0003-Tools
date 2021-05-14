/*
Copyright 2021 DreamView Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/tools/godoc/util"
)

var (
	cleanFlag *bool
	wg        sync.WaitGroup
)

func main() {
	versionFlag := flag.Bool("v", false, "version")
	cleanFlag = flag.Bool("c", false, "detect and clean (default is detect only)")

	flag.Usage = func() {
		fmt.Println("Usage: dephage [-c | -v] <file-path>|<folder-path>")
		fmt.Println("Version: v1.2.0")
		fmt.Println()
		fmt.Println("Detects and optionally cleans the ADSK-SA-2020-0003 Autodesk Maya virus.")
		fmt.Println("\nInfected text .ma and .mb files will be cleaned and the original file")
		fmt.Println("  renamed with a .INFECTED extension.")
		fmt.Println("\nInfected binary .mb files will NOT be cleaned and the file")
		fmt.Println("  renamed with a .INFECTED extension.")
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  Detect file.")
		fmt.Println("    dephage documents/maya/file.ma")
		fmt.Println()
		fmt.Println("  Detect folder and all sub-folders.")
		fmt.Println("    dephage documents/maya")
		fmt.Println()
		fmt.Println("  Detect and clean file.")
		fmt.Println("    dephage -c documents/maya/file.ma")
		fmt.Println()
		fmt.Println("  Detect and clean folder and all subfolders.")
		fmt.Println("    dephage -c documents/maya")
	}

	flag.Parse()

	if *versionFlag {
		fmt.Println("dephage v1.2.0")
		return
	}

	pathArg := flag.Arg(0)
	if pathArg == "" {
		flag.Usage()
		return
	}

	fi, err := os.Stat(pathArg)
	if err != nil {
		fmt.Println("ERROR unable to read:", pathArg)
		return
	}

	absPath, _ := filepath.Abs(pathArg)
	fmt.Println("Processing:            ", absPath)

	homePath, _ := os.UserHomeDir()
	mayaPath := path.Join(homePath, "Documents", "maya", "scripts")
	fmt.Println("Processing maya folder:", mayaPath)
	fmt.Println()

	if detectHomeDir() {
		if *cleanFlag {
			fmt.Println("INFECTED and CLEANING: home folder")
			cleanHomeDir()
			if detectHomeDir() {
				fmt.Println("INFECTED unable to clean: home folder")
			}
		} else {
			fmt.Println("INFECTED: home folder")
		}
	}

	if fi.IsDir() {
		processDir(pathArg)
	} else {
		processFile(pathArg)
	}

	wg.Wait()
}

func processDir(pathArg string) {
	err := filepath.Walk(pathArg, func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() {
			return nil
		}

		processFile(path)
		return nil
	})

	if err != nil {
		fmt.Println("ERROR: unable to read folder:", err)
	}
}

func processFile(path string) {
	if !(strings.HasSuffix(path, ".ma") || strings.HasSuffix(path, ".mb")) {
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		isText, found := detectFile(path)
		if !found {
			return
		}
		if *cleanFlag {
			if isText {
				fmt.Println("INFECTED and CLEANING:", path)
			} else {
				fmt.Println("INFECTED and RENAMING:", path)
			}
			if err := cleanFile(path, isText); err != nil {
				fmt.Println("INFECTED unable to clean:", path)
			}
		} else {
			fmt.Println("INFECTED:", path)
		}
	}()
}

func detectHomeDir() bool {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Unable to check home folder")
		return false
	}

	f := path.Join(home, "Documents", "maya", "scripts", "vaccine.py")
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		return true
	}

	f = path.Join(home, "Documents", "maya", "scripts", "userSetup.py")
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		return true
	}

	f = path.Join(home, "Documents", "maya", "scripts", "userSetup.mel")
	if _, err := os.Stat(f); !os.IsNotExist(err) {
		return true
	}

	return false
}

func cleanHomeDir() {
	home, _ := os.UserHomeDir()
	os.Remove(path.Join(home, "Documents", "maya", "scripts", "vaccine.py"))
	os.Remove(path.Join(home, "Documents", "maya", "scripts", "userSetup.py"))
	os.Remove(path.Join(home, "Documents", "maya", "scripts", "userSetup.mel"))
}

func detectFile(file string) (bool, bool) {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("ERROR: unable to read", file)
		return false, false
	}

	found := bytes.Contains(content, []byte("phage"))
	isText := util.IsText(content)
	return isText, found
}

func cleanFile(file string, isText bool) error {
	infectedName := file + ".INFECTED"
	if err := os.Rename(file, infectedName); err != nil {
		return err
	}

	if !isText {
		return nil
	}

	fin, err := os.Open(infectedName)
	if err != nil {
		return err
	}
	defer fin.Close()

	fout, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fout.Close()

	scanner := bufio.NewScanner(fin)
	ignore := 0
	for scanner.Scan() {
		if ignore > 0 {
			ignore--
			continue
		}

		line := scanner.Text()
		if line == "createNode script -n \"vaccine_gene\";" {
			ignore = 7
			continue
		}

		if line == "createNode script -n \"breed_gene\";" {
			ignore = 4
			continue
		}

		fout.WriteString(line + "\n")
	}

	return nil
}

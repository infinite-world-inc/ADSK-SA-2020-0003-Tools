/*
Copyright 2021 Google LLC

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
	cleanFlag = flag.Bool("c", false, "detect and clean (default is detect only)")

	flag.Usage = func() {
		fmt.Println("Usage: dephage [-c] [root folder]")
		fmt.Println()
		fmt.Println("Detects and optionally cleans the ADSK-SA-2020-0003 Autodesk Maya virus.")
		fmt.Println()
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  Detect from current folder.")
		fmt.Println("    dephage")
		fmt.Println()
		fmt.Println("  Detect from selected folder.")
		fmt.Println("    dephage documents/maya")
		fmt.Println()
		fmt.Println("  Detect and clean from current folder.")
		fmt.Println("    dephage -c")
		fmt.Println()
		fmt.Println("  Detect and clean selected folder.")
		fmt.Println("    dephage -c documents/maya")
	}

	flag.Parse()

	pathArg := flag.Arg(0)
	if pathArg == "" {
		pathArg = "."
	}

	fi, err := os.Stat(pathArg)
	if err != nil {
		fmt.Println("ERROR unable to read folder:", pathArg)
		return
	}
	if !fi.IsDir() {
		fmt.Println(pathArg, "is not a folder")
		return
	}

	absPath, _ := filepath.Abs(pathArg)
	fmt.Println("Processing folder:", absPath)

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

	if err := filepath.Walk(pathArg, processFile); err != nil {
		fmt.Println("ERROR: unable to read folder:", err)
	}

	wg.Wait()
}

func processFile(path string, fi os.FileInfo, err error) error {
	if fi.IsDir() {
		return nil
	}

	if !(strings.HasSuffix(path, ".ma") || strings.HasSuffix(path, ".mb")) {
		return nil
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if detectFile(path) {
			if *cleanFlag {
				fmt.Println("INFECTED and CLEANING:", path)
				if err := cleanFile(path); err != nil {
					fmt.Println("INFECTED unable to clean:", path)
				}
			} else {
				fmt.Println("INFECTED:", path)
			}
		}
	}()

	return nil
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

func detectFile(file string) bool {
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("ERROR: unable to read", file)
	}
	found := bytes.Contains(content, []byte("phage"))
	if found && !util.IsText(content) {
		fmt.Println("INFECTED Unable to clean:", file)
		return false
	}
	return found
}

func cleanFile(file string) error {
	infectedName := file + ".INFECTED"
	if err := os.Rename(file, infectedName); err != nil {
		return err
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

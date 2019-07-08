// Copyright 2019 Max Godfrey
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package explorer

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

// List returns the contents of the directory which the explorer is currently in. Given a bool,
// if true it will include files and directories prefixed with a '.', otherwise it will not.
func (e *explorer) List(listAll bool) ([]string, error) {
	contents := []string{}
	if e.Path != "" {
		contents = append(contents, ".."+PathSep)
	}

	f, err := os.Open(e.GetPath())
	if err != nil {
		return contents, err
	}
	fileInfo, err := f.Readdir(0)
	f.Close()
	if err != nil {
		return contents, err
	}

	if listAll {
		for _, file := range fileInfo {
			if file.IsDir() {
				contents = append(contents, file.Name()+PathSep)
			} else {
				contents = append(contents, file.Name())
			}
		}
	} else {
		for _, file := range fileInfo {
			if file.Name()[0] != '.' {
				if file.IsDir() {
					contents = append(contents, file.Name()+PathSep)
				} else {
					contents = append(contents, file.Name())
				}
			}
		}
	}
	return contents, nil
}

// ListDirectories returns all directories within the current directory in which the explorer is
// located. Note that this function will return an array of directories exclusively. Files will
// be ignored. Given a bool, if true it will include directories with a leading '.', otherwise it
// will not.
func (e *explorer) ListDirectories(listAll bool) ([]string, error) {
	directories := []string{}
	if e.Path != "" {
		directories = append(directories, ".."+PathSep)
	}

	f, err := os.Open(e.GetPath())
	if err != nil {
		return directories, err
	}
	fileInfo, err := f.Readdir(0)
	f.Close()
	if err != nil {
		return directories, err
	}

	if listAll {
		for _, file := range fileInfo {
			if file.IsDir() {
				directories = append(directories, file.Name()+PathSep)
			}
		}
	} else {
		for _, file := range fileInfo {
			if file.IsDir() && file.Name()[0] != '.' {
				directories = append(directories, file.Name()+PathSep)
			}
		}
	}
	return directories, nil
}

// ListN returns the first N contents of the current directory. If the current directory contains
// fewer than N things, then the contents of the current directory will be returned.
func (e *explorer) ListN(curSelected string, n int, listAll bool) ([]string, error) {
	var contents []string
	if n <= 0 {
		return contents, nil
	}

	f, err := os.Open(e.GetPath() + curSelected)
	if err != nil {
		if strings.HasSuffix(err.Error(), "denied") {
			contents = append(contents, "PERMISSION DENIED")
			err = nil
		}
		return contents, err
	}
	defer f.Close()
	if listAll {
		fileInfo, err := f.Readdir(n)
		if err != nil {
			if err.Error() == "EOF" {
				contents = append(contents, "DIRECTORY IS EMPTY")
				err = nil
			}
			return contents, err
		}
		for _, file := range fileInfo {
			if file.IsDir() {
				contents = append(contents, file.Name()+PathSep)
			} else {
				contents = append(contents, file.Name())
			}
		}
		return contents, nil
	}

	fileInfo, err := f.Readdir(0)
	if err != nil {
		return contents, err
	}
	for _, file := range fileInfo {
		if len(contents) == n {
			break
		}
		if file.Name()[0] != '.' {
			if file.IsDir() {
				contents = append(contents, file.Name()+PathSep)
			} else {
				contents = append(contents, file.Name())
			}
		}
	}
	if len(contents) == 0 {
		contents = append(contents, "NO CONTENTS TO DISPLAY IN LIST MODE")
	}

	return contents, nil
}

// ListFiles returns all files within the current directory in which the explorer is located. Note
// that this function will return an array of files exclusively. No directory will be included.
// Given a bool, if true it will include files prefixed with a '.', otherwise it will not.
func (e *explorer) ListFiles(listAll bool) ([]string, error) {
	var files []string
	f, err := os.Open(e.GetPath())
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(0)
	f.Close()
	if err != nil {
		return files, err
	}

	if listAll {
		for _, file := range fileInfo {
			if !file.IsDir() {
				files = append(files, file.Name())
			}
		}
	} else {
		for _, file := range fileInfo {
			if !file.IsDir() && file.Name()[0] != '.' {
				files = append(files, file.Name())
			}
		}
	}
	return files, nil
}

// ReadN reads the first N lines of a
func (e *explorer) ReadN(fileName string, n int) ([]string, error) {
	var contents []string
	file, err := os.Open(e.GetPath() + fileName)
	if err != nil {
		return contents, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for i := 0; i < n && scanner.Scan(); i++ {
		contents = append(contents, scanner.Text())
	}

	file.Close()
	return contents, nil
}

// View will open an os-specific editor in which a file can be viewed (preferably for editing).
func (e *explorer) View(fileName string) error {
	cmd := exec.Command(TextEditor, e.GetPath()+fileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

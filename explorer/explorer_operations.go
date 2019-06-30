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
	"os"
)

// List returns the contents of the directory which the explorer is currently in.
func (e *explorer) List() ([]string, error) {
	var contents []string
	f, err := os.Open(e.Path)
	if err != nil {
		return contents, err
	}
	fileInfo, err := f.Readdir(0)
	f.Close()
	if err != nil {
		return contents, err
	}

	for _, file := range fileInfo {
		contents = append(contents, file.Name())
	}
	return contents, nil
}

// ListDirectories returns all directories within the current directory in which the explorer is
// located. Note that this function will return an array of directories exclusively. Files will
// be ignored.
func (e *explorer) ListDirectories() ([]string, error) {
	var directories []string
	f, err := os.Open(e.Path)
	if err != nil {
		return directories, err
	}
	fileInfo, err := f.Readdir(0)
	f.Close()
	if err != nil {
		return directories, err
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			directories = append(directories, file.Name())
		}
	}
	return directories, nil
}

// ListFiles returns all files within the current directory in which the explorer is located. Note
// that this function will return an array of files exclusively. No directory will be included.
func (e *explorer) ListFiles() ([]string, error) {
	var files []string
	f, err := os.Open(e.Path)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(0)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if !file.IsDir() {
			files = append(files, file.Name())
		}
	}
	return files, nil
}

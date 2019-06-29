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

package fileutils

import (
	"errors"
	"os"
	"strings"
)

// directoryExists determines whether or not a directory exists
func directoryExists(filePath string) error {
	fileType, err := os.Stat(filePath)
	if err != nil {
		return err
	}
	if !fileType.IsDir() {
		return errors.New("specified file is not a directory")
	}
	return nil
}

// Explorer represents a file explorer.
type Explorer struct {
	Path string
}

// Move will move the explorer to a given directory relative to the current working directory.
// The given directory must be adjacent to the directory that the explorer is currently in.
func (e *Explorer) MoveOne(nextDirectory string) error {
	// Remove trailing forward slashes from nextDirectory
	if nextDirectory[len(nextDirectory) - 1] == '/' {
		nextDirectory = nextDirectory[:len(nextDirectory) - 1]
	}
	
	if nextDirectory == ".." {
		lastForwardSlash := strings.LastIndexAny(e.Path, "/")
		if lastForwardSlash == -1 {
			return errors.New("already at top-level directory")
		}
		e.Path = e.Path[:lastForwardSlash]
		return nil
	}

	nextPath := e.Path + "/" + nextDirectory
	if err := directoryExists(nextPath); err != nil {
		return err
	}
	e.Path = nextPath
	return nil
}
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
	"errors"
	"os"
	"os/user"
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

// explorer represents a file explorer.
type explorer struct {
	Path        string
	CurrentUser *user.User
}

// MoveAbsolute will move the explorer to a specified absolute path.
// The path may begin with either a '~' or a '/'.
func (e *explorer) MoveAbsolute(path string) error {
	// Remove trailing forward slashes from the path
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if path == "~" {
		path = e.CurrentUser.HomeDir
	} else if path[0] == '~' {
		path = e.CurrentUser.HomeDir + path[1:]
	}

	if err := directoryExists(path); err != nil {
		return err
	}
	e.Path = path
	return nil
}

// MoveMultiple will move the explorer through a list of directories separated by '/' characters.
// Each directory must be adjacent to the directory that the explorer is currently in.
func (e *explorer) MoveMultiple(directories string) error {
	// Remove trailing forward slashes from directories
	if directories[len(directories)-1] == '/' {
		directories = directories[:len(directories)-1]
	}

	dirList := strings.Split(directories, "/")
	for _, dir := range dirList {
		err := e.MoveOne(dir)
		if err != nil {
			return err
		}
	}
	return nil
}

// Move will move the explorer to a given directory relative to the current working directory.
// The given directory must be adjacent to the directory that the explorer is currently in.
func (e *explorer) MoveOne(nextDirectory string) error {
	// Remove trailing forward slashes from nextDirectory
	if nextDirectory[len(nextDirectory)-1] == '/' {
		nextDirectory = nextDirectory[:len(nextDirectory)-1]
	}

	if nextDirectory == "." {
		return nil
	} else if nextDirectory == ".." {
		lastForwardSlash := strings.LastIndexAny(e.Path, "/")
		if lastForwardSlash == -1 {
			// We are at the top-level directory
			return nil
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

// New returns an explorer type with all member values initialised to their defaults.
func New() (e explorer) {
	e.Path = ""
	e.CurrentUser, _ = user.Current()
	return
}

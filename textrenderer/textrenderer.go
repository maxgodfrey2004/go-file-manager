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

package textrenderer

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type textrenderer struct {
	Text []string
}

// TerminalDimensions obtains the dimensions of the current terminal window and returns them.
// Dimensions are returned with height first and width second. An error is also returned.
func (t *textrenderer) TerminalDimensions() (int, int, error) {
	rowsOutput, err := exec.Command("tput", "lines").Output()
	if err != nil {
		return -1, -1, err
	}
	height, _ := strconv.Atoi(strings.TrimRight(string(rowsOutput), " \n"))

	columnsOutput, err := exec.Command("tput", "cols").Output()
	if err != nil {
		return -1, -1, err
	}
	width, _ := strconv.Atoi(strings.TrimRight(string(columnsOutput), " \n"))

	return height, width, nil
}

// ClearScreen clears the terminal screen.
func (t *textrenderer) ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// New returns a new instance of the textrenderer type.
func New() (t textrenderer) {
	return
}

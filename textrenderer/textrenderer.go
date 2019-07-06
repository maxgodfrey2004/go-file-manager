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

	"github.com/fatih/color"
)

// min returns the minimum of two integers. Strangely, math.Min takes two float64 variables as
// parameters.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type textrenderer struct {
	SelectedIndex int
	StartIndex    int
	Text          []string
}

// render displays the selected window of text on the terminal screen. The selected file will be
// displayed with a blue background, indicative of its selection. If told to show scrollback, it
// will in turn notify the ClearScreen method; and scrollback will be preserved.
func (t *textrenderer) render(showScrollback bool) error {
	t.ClearScreen(showScrollback)
	termHeight, _, err := t.TerminalDimensions()
	if err != nil {
		return err
	}

	selected := color.New(color.FgWhite, color.BgBlue)
	endIndex := min(t.StartIndex+termHeight, len(t.Text))
	for i := t.StartIndex; i < endIndex; i++ {
		if i == t.SelectedIndex {
			if _, err := selected.Println(t.Text[i]); err != nil {
				return err
			}
		} else {
			color.White(t.Text[i])
		}
	}

	return nil
}

// Display reassigns the lines which the textrenderer will be displaying, and then renders them
// on the terminal screen.
func (t *textrenderer) Display(text []string) error {
	t.Text = text
	t.SelectedIndex = 0
	t.StartIndex = 0

	if err := t.Render(); err != nil {
		return err
	}
	return nil
}

// Render displays the selected window of text on the terminal screen. No scrollback will be
// kept when the screen is cleared.
func (t *textrenderer) Render() error {
	if err := t.render(false); err != nil {
		return err
	}
	return nil
}

// RenderWithScrollback displays the selected window of text on the terminal screen. When the
// screen is cleared, scrollback will be preserved.
func (t *textrenderer) RenderWithScrollback() error {
	if err := t.render(true); err != nil {
		return err
	}
	return nil
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

// ClearScreen clears the terminal screen. If told to show scrollback, the method will append
// a -l flag to the system call "clear" being made. This will preserve scrollback.
func (t *textrenderer) ClearScreen(showScrollback bool) {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// New returns a new instance of the textrenderer type.
func New() (t textrenderer) {
	t.SelectedIndex = 0
	t.StartIndex = 0
	return
}

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

	"github.com/nsf/termbox-go"
)

// X positions to render various elements of the explorer.
const (
	CaretRenderX = 1
	FileRenderX  = 3
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

// CurrentSelected returns the element of the textrenderer's Text attribute which is currently
// selected.
func (t *textrenderer) CurrentSelected() string {
	return t.Text[t.SelectedIndex]
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

// render displays the selected window of text on the terminal screen. The selected file will be
// displayed with a blue background, indicative of its selection.
func (t *textrenderer) Render() error {
	_, termHeight := termbox.Size()
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}

	endIndex := min(t.StartIndex+termHeight, len(t.Text))
	for i := t.StartIndex; i < endIndex; i++ {
		bgColor := termbox.ColorDefault
		if i == t.SelectedIndex {
			termbox.SetCell(CaretRenderX, i-t.StartIndex, rune('>'), termbox.ColorDefault, termbox.ColorDefault)
		}
		fgColor := termbox.ColorDefault
		if t.Text[i][len(t.Text[i])-1] == '/' {
			fgColor = termbox.ColorBlue
		}
		for j := 0; j < len(t.Text[i]); j++ {
			termbox.SetCell(FileRenderX+j, i-t.StartIndex, rune(t.Text[i][j]), fgColor, bgColor)
		}
	}
	termbox.Flush()

	return nil
}

// ClearScreen clears the terminal screen.
func (t *textrenderer) ClearScreen() error {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// New returns a new instance of the textrenderer type.
func New() (t textrenderer) {
	t.SelectedIndex = 0
	t.StartIndex = 0
	return
}

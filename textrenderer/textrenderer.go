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
	"github.com/nsf/termbox-go"
)

// X positions to render various elements of the explorer.
const (
	CaretRenderX = 1
	FileRenderX  = 3
)

// Modifiers affecting the size of the view through which textrenderer.Text is displayed.
const (
	textHeightModifier = 1
	textWidthModifier  = 0
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
	Header        string
	SelectedIndex int
	StartIndex    int
	Text          []string
}

// CurrentSelected returns the element of the textrenderer's Text attribute which is currently
// selected.
func (t *textrenderer) CurrentSelected() string {
	return t.Text[t.SelectedIndex]
}

// Display reassigns the lines which the textrenderer will be displaying, and their respective
// header. It then renders them on the terminal screen.
func (t *textrenderer) Display(header string, text []string) {
	t.Header = header
	t.Text = text
	t.SelectedIndex = 0
	t.StartIndex = 0

	t.Render()
}

// Render displays the selected window of text and respective header on the terminal screen. The
// selected file will be displayed with a caret, indicative of its selection.
func (t *textrenderer) Render() {
	termWidth, termHeight := termbox.Size()
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}

	headerEnd := min(termWidth, len(t.Header))
	for i := 0; i < headerEnd; i++ {
		termbox.SetCell(i, 0, rune(t.Header[i]), termbox.ColorDefault, termbox.ColorDefault)
	}

	endIndex := min(t.StartIndex+termHeight-1, len(t.Text))
	for i := t.StartIndex; i < endIndex; i++ {
		bgColor := termbox.ColorDefault
		yCoord := i - t.StartIndex + 1
		if i == t.SelectedIndex {
			termbox.SetCell(CaretRenderX, yCoord, rune('>'), termbox.ColorDefault, termbox.ColorDefault)
		}
		fgColor := termbox.ColorDefault
		if t.Text[i][len(t.Text[i])-1] == '/' {
			fgColor = termbox.ColorBlue
		}
		for j := 0; j < len(t.Text[i]); j++ {
			termbox.SetCell(FileRenderX+j, yCoord, rune(t.Text[i][j]), fgColor, bgColor)
		}
	}
	termbox.Flush()
}

// TextViewSize returns the dimensions of the box in which textrenderer.Text is stored
func (t *textrenderer) TextViewSize() (int, int) {
	width, height := termbox.Size()
	return width - textWidthModifier, height - textHeightModifier
}

// New returns a new instance of the textrenderer type.
func New() (t textrenderer) {
	t.SelectedIndex = 0
	t.StartIndex = 0
	return
}

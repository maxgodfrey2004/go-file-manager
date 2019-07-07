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
	CaretRenderX       = 1
	FilePreviewRenderY = 2
	FileRenderX        = 3
)

// Modifiers affecting the size of the view through which textrenderer.Text is displayed.
const (
	textHeightModifier = 1
	textWidthModifier  = 0

	filePreviewHeightModifier = 1
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
	Header        string   // The string to render above Text.
	SelectedIndex int      // The selected index in Text.
	StartIndex    int      // Start rendering text from this index in Text.
	StopRight     int      // Stop rendering text past this point.
	Text          []string // The text which the renderer draws on the screen.
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

// RecalculateBounds recalculates the positions on the terminal at which textrenderer stops
// rendering text.
func (t *textrenderer) RecalculateBounds() {
	width, _ := termbox.Size()
	t.StopRight = width / 2
}

// Render displays the selected window of text and respective header on the terminal screen. The
// selected file will be displayed with a caret, indicative of its selection.
func (t *textrenderer) Render() {
	t.RecalculateBounds()
	_, termHeight := termbox.Size()
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}

	headerEnd := min(t.StopRight, len(t.Header))
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

func (t *textrenderer) RenderPreview(preview []string) {
	t.RecalculateBounds()
	previewX := t.StopRight + 2
	width, _ := termbox.Size()

	for i := 0; i < len(preview); i++ {
		y := i + FilePreviewRenderY
		for j := 0; j < min(width-previewX, len(preview[i])); j++ {
			termbox.SetCell(previewX+j, y, rune(preview[i][j]), termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	termbox.Flush()
}

// PreviewSize returns the dimensions of the box in which the file preview is rendered.
func (t *textrenderer) PreviewHeight() int {
	_, height := termbox.Size()
	return height - filePreviewHeightModifier - FilePreviewRenderY
}

// TextViewSize returns the dimensions of the box in which textrenderer.Text is stored.
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

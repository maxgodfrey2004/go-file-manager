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
	"github.com/maxgodfrey2004/go-file-manager/explorer"
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
	filePreviewWidthModifier  = 1
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
	KeyFunctions  []string // The function of each command, rendered at the bottom of the terminal.
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
func (t *textrenderer) Display(header string, text []string, preview []string) {
	t.Header = header
	t.Text = text
	t.SelectedIndex = 0
	t.StartIndex = 0

	t.Render(preview)
}

/// Init initialises the textrenderer with a header, and a body of text to display.
func (t *textrenderer) Init(header string, text []string) {
	t.Header = header
	t.Text = text
	t.SelectedIndex = 0
	t.StartIndex = 0
}

// RecalculateBounds recalculates the positions on the terminal at which textrenderer stops
// rendering text.
func (t *textrenderer) RecalculateBounds() {
	width, _ := termbox.Size()
	t.StopRight = width / 2
}

// Render displays the selected window of text and respective header on the terminal screen. The
// selected file will be displayed with a caret, indicative of its selection. A preview of the
// current selected item will also be displayed on the right hand side of the screen.
func (t *textrenderer) Render(preview []string) {
	t.RecalculateBounds()
	_, termHeight := termbox.Size()
	if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
		panic(err)
	}

	headerEnd := min(t.StopRight, len(t.Header))
	for i := 0; i < headerEnd; i++ {
		termbox.SetCell(i, 0, rune(t.Header[i]), termbox.ColorDefault, termbox.ColorDefault)
	}

	endIndex := min(t.StartIndex+termHeight-1-textHeightModifier, len(t.Text))
	for i := t.StartIndex; i < endIndex; i++ {
		bgColor := termbox.ColorDefault
		yCoord := i - t.StartIndex + 1
		if i == t.SelectedIndex {
			termbox.SetCell(CaretRenderX, yCoord, rune('>'), termbox.ColorDefault, termbox.ColorDefault)
		}
		fgColor := termbox.ColorDefault
		if t.Text[i][len(t.Text[i])-1] == explorer.PathSepChar {
			fgColor = termbox.ColorBlue
		}
		for j := 0; j < len(t.Text[i]); j++ {
			termbox.SetCell(FileRenderX+j, yCoord, rune(t.Text[i][j]), fgColor, bgColor)
		}
	}

	t.RenderKeyFunctions()
	t.RenderPreview(preview)
	termbox.HideCursor()
	termbox.Flush()
}

// RenderBox renders a box on the terminal whose upper left corner, width and height are specified.
func (t *textrenderer) RenderBox(topLeftX, topLeftY, width, height int) {
	if width <= 0 || height <= 0 {
		return
	}

	fgColor := termbox.ColorDefault
	bgColor := termbox.ColorDefault
	termbox.SetCell(topLeftX, topLeftY, rune('┌'), fgColor, bgColor)
	termbox.SetCell(topLeftX+width, topLeftY, rune('┐'), fgColor, bgColor)
	termbox.SetCell(topLeftX, topLeftY+height, rune('└'), fgColor, bgColor)
	termbox.SetCell(topLeftX+width, topLeftY+height, rune('┘'), fgColor, bgColor)

	for x := 1; x < width; x++ {
		termbox.SetCell(topLeftX+x, topLeftY, rune('─'), fgColor, bgColor)
		termbox.SetCell(topLeftX+x, topLeftY+height, rune('─'), fgColor, bgColor)
	}
	for y := 1; y < height; y++ {
		termbox.SetCell(topLeftX, topLeftY+y, rune('│'), fgColor, bgColor)
		termbox.SetCell(topLeftX+width, topLeftY+y, rune('│'), fgColor, bgColor)
	}
	termbox.HideCursor()
	termbox.Flush()
}

// RenderKeyFunctions renders the textrenderer's attribute KeyFunctions on the bottom line of the
// terminal screen.
func (t *textrenderer) RenderKeyFunctions() {
	width, height := termbox.Size()
	x := 1
	y := height - 1
	renderedFirst := 0

	fgColor := termbox.ColorCyan
	bgColor := termbox.ColorDefault
	for _, token := range t.KeyFunctions {
		if renderedFirst == 1 {
			termbox.SetCell(x-2, y, rune(','), fgColor, bgColor)
		}
		if x+len(token)+renderedFirst < width {
			for i := 0; i < len(token); i++ {
				termbox.SetCell(x+i, y, rune(token[i]), fgColor, bgColor)
			}
			x += len(token)
		} else {
			break
		}
		renderedFirst = 1
		x += 2
	}
	termbox.HideCursor()
	termbox.Flush()
}

// RenderPreview renders a preview of the current selected file (not a directory) on the right hand
// half of the terminal screen.
func (t *textrenderer) RenderPreview(preview []string) {
	t.RecalculateBounds()
	previewX := t.StopRight + 2
	width, _ := termbox.Size()
	boxWidth := width - previewX - filePreviewWidthModifier
	boxHeight := t.PreviewHeight() + 1
	t.RenderBox(previewX-1, FilePreviewRenderY-1, boxWidth, boxHeight)

	if preview == nil {
		return
	}
	for i := 0; i < len(preview); i++ {
		y := i + FilePreviewRenderY
		fgColor := termbox.ColorDefault
		bgColor := termbox.ColorDefault
		if preview[i] == "PERMISSION DENIED" {
			bgColor = termbox.ColorRed
		}
		for j := 0; j < min(boxWidth, len(preview[i])); j++ {
			termbox.SetCell(previewX+j, y, rune(preview[i][j]), fgColor, bgColor)
		}
	}

	termbox.HideCursor()
	termbox.Flush()
}

// PreviewSize returns the dimensions of the box in which the file preview is rendered.
func (t *textrenderer) PreviewHeight() int {
	_, height := termbox.Size()
	return height - filePreviewHeightModifier - FilePreviewRenderY - 1
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

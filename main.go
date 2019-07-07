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

package main

import (
	"github.com/maxgodfrey2004/go-file-manager/explorer"
	"github.com/maxgodfrey2004/go-file-manager/textrenderer"
	"github.com/nsf/termbox-go"
)

// KeyEvent enumerates the different actions a user may take when using the file manager.
type KeyEvent int

const (
	// Reselect represents the user pressing a navigation key, and requiring the selected file
	// to be rendered in a different place (in the same directory).
	Reselect KeyEvent = iota + 1

	// Move represents the user selecting a directory, and being taken to a new location in the
	// file system.
	Move

	// Quit represents the termination of the application.
	Quit
)

// Movement directions
const (
	Down = 1
	Up   = -1
)

// keypress represents a physical key being pressed on the keyboard.
type keypress struct {
	EventType KeyEvent
	Key       termbox.Key
}

var (
	// nav is used to traverse the user's file system.
	nav = explorer.New()
	// screen is used to render colored output on the terminal.
	screen = textrenderer.New()

	// keypressChan is used to register incoming keypresses.
	keypressChan = make(chan keypress)
)

// listenForKeypress indefinitely listens for keyboard input, and sends it to a specified channel.
// Note that ths method is intended to be called asynchronously (ie. as a goroutine).
func listenForKeypress(ch chan keypress) {
	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowDown:
				ch <- keypress{EventType: Reselect, Key: ev.Key}
			case termbox.KeyArrowUp:
				ch <- keypress{EventType: Reselect, Key: ev.Key}
			case termbox.KeyArrowLeft, termbox.KeyEnter:
				ch <- keypress{EventType: Move, Key: ev.Key}
			default:
				switch ev.Ch {
				case rune('Q'), rune('q'):
					ch <- keypress{EventType: Quit, Key: ev.Key}
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

// move moves nav, the explorer, to the current directory which the user has selected.
func move() {
	nextDir := screen.CurrentSelected()
	if err := explorer.DirectoryExists(nav.Path + "/" + nextDir); err == nil {
		if err := nav.MoveOne(nextDir); err != nil {
			panic(err)
		}
		dirContents, err := nav.List(false)
		if err != nil {
			panic(nav.Path)
		}
		if err := screen.Display(dirContents); err != nil {
			panic(err)
		}
	}
}

// keyToDirection converts a termbox Key code into a direction for the selected file or directory
// when the up or down arrow keys are pressed.
func keyToDirection(key termbox.Key) int {
	switch key {
	case termbox.KeyArrowDown:
		return Down
	case termbox.KeyArrowUp:
		return Up
	default:
		return 0
	}
}

// reselect moves the screen's display of files when the user presses either an up or down arrow
// key.
func reselect(ev keypress) {
	newIndex := screen.SelectedIndex + keyToDirection(ev.Key)
	if newIndex < 0 || newIndex >= len(screen.Text) {
		return
	}

	height, _, err := screen.TerminalDimensions()
	if err != nil {
		panic(err)
	}

	screen.SelectedIndex = newIndex
	if newIndex >= screen.StartIndex+height {
		screen.StartIndex++
	} else if newIndex < screen.StartIndex {
		screen.StartIndex--
	}
	if err := screen.Render(); err != nil {
		panic(err)
	}
}

// startExplorer runs the file manager until a Quit event is sent.
func startExplorer() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	if err := nav.MoveAbsolute("~"); err != nil {
		panic(err)
	}
	dirContents, err := nav.List(false)
	if err != nil {
		panic(err)
	}
	if err := screen.Display(dirContents); err != nil {
		panic(err)
	}

	go listenForKeypress(keypressChan)

mainloop:
	for {
		select {
		case ev := <-keypressChan:
			switch ev.EventType {
			case Reselect:
				reselect(ev)
			case Move:
				move()
			case Quit:
				break mainloop
			}
		}
	}
}

func main() {
	startExplorer()
}

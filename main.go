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
	"os"

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

	// Select represents the user pressing return or the right arrow, thus selecting the
	// current file or directory to be either viewed or moved to respectively.
	Select

	// ToggleListAll represents the user toggling the state of whether or not they want to see
	// directory contents whose names contain leading `.` characters.
	ToggleListAll

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
	keypressChan chan keypress

	// listAll is used to determine whether the user wishes to see directory contents with
	// a leading `.`. By default, we assume that they do not.
	listAll = false
)

// listenForEvents indefinitely listens for termbox events. Any events that take the form of
// keyboard input are sent to a specified channel, where they will be proecessed externally.
// Note that this method is intended to be called asynchronously (ie. as a goroutine).
func listenForEvents(ch chan keypress) {
	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowDown:
				ch <- keypress{EventType: Reselect, Key: ev.Key}
			case termbox.KeyArrowUp:
				ch <- keypress{EventType: Reselect, Key: ev.Key}
			case termbox.KeyArrowRight, termbox.KeyEnter:
				ch <- keypress{EventType: Select, Key: ev.Key}
			case termbox.KeyCtrlC:
				ch <- keypress{EventType: Quit, Key: ev.Key}
			default:
				switch ev.Ch {
				case rune('Q'), rune('q'):
					ch <- keypress{EventType: Quit, Key: ev.Key}
				case rune('A'), rune('a'):
					ch <- keypress{EventType: ToggleListAll, Key: ev.Key}
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		case termbox.EventInterrupt:
			return
		case termbox.EventResize:
			screen.Render(genPreview())
		}
	}
}

// moveDirectory moves nav, the explorer, to the current directory which the user has selected.
func moveDirectory() {
	nextDir := screen.CurrentSelected()
	if err := explorer.DirectoryExists(nav.GetPath() + nextDir); err == nil {
		if err := nav.MoveOne(nextDir); err != nil {
			panic(err)
		}
		dirContents, err := nav.List(listAll)
		if err != nil {
			panic(nav.Path)
		}
		screen.Init(nav.GetPath(), dirContents)
		screen.Display(nav.GetPath(), dirContents, genPreview())
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

// genPreview returns a preview of the current selected file or directory.
func genPreview() []string {
	curSelected := screen.CurrentSelected()
	var err error
	var preview []string
	if curSelected[len(curSelected)-1] != explorer.PathSepChar {
		preview, err = nav.ReadN(curSelected, screen.PreviewHeight())
		if err != nil {
			panic(err)
		}
	} else {
		preview, err = nav.ListN(curSelected, screen.PreviewHeight(), listAll)
		if err != nil {
			panic(err)
		}
	}

	return preview
}

// reselect moves the screen's display of files when the user presses either an up or down arrow
// key.
func reselect(ev keypress) {
	newIndex := screen.SelectedIndex + keyToDirection(ev.Key)
	if newIndex < 0 || newIndex >= len(screen.Text) {
		return
	}

	_, height := screen.TextViewSize()
	screen.SelectedIndex = newIndex
	if newIndex >= screen.StartIndex+height {
		screen.StartIndex++
	} else if newIndex < screen.StartIndex {
		screen.StartIndex--
	}
	screen.Render(genPreview())
}

// selectContents is called when the user selects either a file or a directory. It in turn will
// either open an editor with the selected file, or move to the selected directory.
func selectContents() {
	curSelected := screen.CurrentSelected()
	if curSelected[len(curSelected)-1] == explorer.PathSepChar {
		moveDirectory()
	} else {
		pathCopy := nav.GetPath()
		if err := nav.View(curSelected); err != nil {
			panic(err)
		}
		termbox.Interrupt()
		startExplorer(pathCopy)
		/*termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

		dirContents, err := nav.List(listAll)
		if err != nil {
			panic(err)
		}
		screen.Init(nav.GetPath(), dirContents)
		screen.Display(nav.GetPath(), dirContents, genPreview())*/
	}
}

// startExplorer runs the file manager until a Quit event is sent.
func startExplorer(startDirectory string) {
	if err := termbox.Init(); err != nil {
		panic(err)
	}

	keypressChan = make(chan keypress)
	if err := nav.MoveAbsolute(startDirectory); err != nil {
		panic(startDirectory)
	}
	dirContents, err := nav.List(listAll)
	if err != nil {
		panic(err)
	}
	screen.KeyFunctions = []string{
		"[Q|q: Quit]",
		"[A|a: List]",
	}
	screen.Init(nav.GetPath(), dirContents)
	screen.Display(nav.GetPath(), dirContents, genPreview())

	go listenForEvents(keypressChan)

	for {
		select {
		case ev := <-keypressChan:
			switch ev.EventType {
			case Reselect:
				reselect(ev)
			case Select:
				selectContents()
			case ToggleListAll:
				toggleListAll()
			case Quit:
				termbox.Close()
				os.Exit(0)
			}
		}
	}
}

func toggleListAll() {
	listAll = !listAll
	dirContents, err := nav.List(listAll)
	if err != nil {
		panic(err)
	}
	screen.Init(nav.GetPath(), dirContents)
	screen.Display(nav.GetPath(), dirContents, genPreview())
}

func main() {
	startExplorer("~")
}

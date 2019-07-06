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

type KeyEvent int

// Keyboard events
const (
	RESELECT KeyEvent = iota + 1
	MOVE
	QUIT
)

// Movement directions
const (
	DOWN = -1
	UP   = 1
)

type keypress struct {
	EventType KeyEvent
	Key       termbox.Key
}

var (
	nav    = explorer.New()
	screen = textrenderer.New()

	keypressChan = make(chan keypress)
)

func listenForKeypress(ch chan keypress) {
	termbox.SetInputMode(termbox.InputEsc)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowDown:
				ch <- keypress{EventType: RESELECT, Key: ev.Key}
			case termbox.KeyArrowUp:
				ch <- keypress{EventType: RESELECT, Key: ev.Key}
			case termbox.KeyArrowLeft, termbox.KeyEnter:
				ch <- keypress{EventType: MOVE, Key: ev.Key}
			default:
				switch ev.Ch {
				case rune('Q'), rune('q'):
					ch <- keypress{EventType: QUIT, Key: ev.Key}
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func keyToDirection(key termbox.Key) int {
	switch key {
	case termbox.KeyArrowDown:
		return 1
	case termbox.KeyArrowUp:
		return -1
	default:
		return 0
	}
}

func startExplorer() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	nav.MoveAbsolute("~")
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
			case RESELECT:
				newIndex := screen.SelectedIndex + keyToDirection(ev.Key)
				if newIndex < 0 || newIndex >= len(screen.Text) {
					break
				}
				_, height := termbox.Size()
				screen.SelectedIndex = newIndex
				if newIndex >= screen.StartIndex+height {
					screen.StartIndex++
				} else if newIndex < screen.StartIndex {
					screen.StartIndex--
				}
				screen.Render()
			case MOVE:
				nextDir := screen.CurrentSelected()
				if err := explorer.DirectoryExists(nav.Path + "/" + nextDir); err == nil {
					nav.MoveOne(nextDir)
					dirContents, err := nav.List(false)
					if err != nil {
						panic(nav.Path)
					}
					if err := screen.Display(dirContents); err != nil {
						break
					}
				}
			case QUIT:
				break mainloop
			}
		}
	}
}

func main() {
	startExplorer()
}

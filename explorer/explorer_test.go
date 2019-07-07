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
	"testing"
)

func TestDirectoryExists(t *testing.T) {
	t.Log(DirectoryExists("/"))              // exists
	t.Log(DirectoryExists("/doesnotexist/")) // does not exist
	t.Log(DirectoryExists("/bin"))           // exists
	t.Log(DirectoryExists("/bin/bash"))      // not a directory
}

func TestMoveOne(t *testing.T) {
	e := New()
	t.Log("e.Path:", e.Path)
	err1 := e.MoveOne("bin")
	t.Log("e.Path:", e.Path, "error:", err1)
	err2 := e.MoveOne("..")
	t.Log("e.Path:", e.Path, "error:", err2)

	err3 := e.MoveOne("bin/")
	t.Log("e.Path:", e.Path, "error:", err3)
	err4 := e.MoveOne(".")
	t.Log("e.Path:", e.Path, "error:", err4)
	err5 := e.MoveOne("doesnotexist")
	t.Log("e.Path:", e.Path, "error:", err5)
}

func TestMoveMultiple(t *testing.T) {
	e := New()
	err1 := e.MoveMultiple("../../")
	t.Log("e.Path:", e.Path, "error:", err1)

	err2 := e.MoveMultiple("bin/../bin/../")
	t.Log("e.Path:", e.Path, "error:", err2)
	err3 := e.MoveMultiple("doesnotexist/doesnotexist")
	t.Log("e.Path:", e.Path, "error:", err3)
}

func TestMoveAbsolute(t *testing.T) {
	e := New()
	err1 := e.MoveAbsolute("/mnt/c/")
	t.Log("e.Path:", e.Path, "error:", err1)
	err2 := e.MoveAbsolute("~/bin")
	t.Log("e.Path:", e.Path, "error:", err2)
	err3 := e.MoveAbsolute("~")
	t.Log("e.Path:", e.Path, "error:", err3)
}

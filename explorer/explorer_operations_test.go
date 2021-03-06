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

func TestList(t *testing.T) {
	e := New()
	if err := e.MoveAbsolute("~"); err != nil {
		t.Fatal(err)
	}
	all, _ := e.List(false)
	t.Log(all)
	all, _ = e.List(true)
	t.Log(all)

	e.Path = "/doesnotexist"
	all, err := e.List(false)
	t.Log(all, err)
}

func TestListDirectories(t *testing.T) {
	e := New()
	if err := e.MoveAbsolute("~"); err != nil {
		t.Fatal(err)
	}
	directories, _ := e.ListDirectories(false)
	t.Log(directories)
	directories, _ = e.ListDirectories(true)
	t.Log(directories)

	e.Path = "/doesnotexist"
	all, err := e.ListDirectories(false)
	t.Log(all, err)
}

func TestListFiles(t *testing.T) {
	e := New()
	if err := e.MoveAbsolute("~"); err != nil {
		t.Fatal(err)
	}
	files, _ := e.ListFiles(false)
	t.Log(files)
	files, _ = e.ListFiles(true)
	t.Log(files)

	e.Path = "/doesnotexist"
	all, err := e.ListFiles(false)
	t.Log(all, err)
}

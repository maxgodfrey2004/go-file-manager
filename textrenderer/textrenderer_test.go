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
	"strconv"
	"testing"
)

func TestRender(t *testing.T) {
	tr := New()
	var dispArray []string
	for i := 0; i <= 40; i++ {
		dispArray = append(dispArray, "Test "+strconv.Itoa(i))
	}
	tr.Text = dispArray
	tr.SelectedIndex = 2
	tr.StartIndex = 1

	tr.Render()
}

func TestNew(t *testing.T) {
	tr := New()
	t.Log(tr.Text)
	t.Log(tr.SelectedIndex)
	t.Log(tr.StartIndex)
}

/*
* Copyright (C) 2021 Henrik A. Christensen
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package page

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

type testGroup struct {
	GroupName string
	Members   []testPerson
}

type testPerson struct {
	Name string
	Age  int32
}

func TestWriteFilenameIsRequired(t *testing.T) {
	if err := Write("", JSON, struct{}{}); err == nil {
		t.Fatal("Expected error if filename is empty")
	}
}

func TestWriteFormat(t *testing.T) {
	dir, err := os.MkdirTemp("", "go-page")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tests := []struct {
		Filename string
		Format   string
		Expected string
	}{
		{
			Filename: filepath.Join(dir, "data"),
			Format:   JSON,
			Expected: filepath.Join(dir, "data.json"),
		},
		{
			Filename: filepath.Join(dir, "data.json"),
			Format:   JSON,
			Expected: filepath.Join(dir, "data.json"),
		},
		{
			Filename: filepath.Join(dir, "data.json"),
			Format:   "",
			Expected: filepath.Join(dir, "data.json"),
		},
		{
			Filename: filepath.Join(dir, "data.config"),
			Format:   JSON,
			Expected: filepath.Join(dir, "data.config.json"),
		},
		{
			Filename: filepath.Join(dir, "data"),
			Format:   "JSON",
			Expected: filepath.Join(dir, "data.json"),
		},
		{
			Filename: filepath.Join(dir, "data.json"),
			Format:   "JsON",
			Expected: filepath.Join(dir, "data.json"),
		},
	}

	for _, test := range tests {
		err = Write(test.Filename, test.Format, struct{}{})
		if err != nil {
			t.Fatal(err)
		}

		if _, err = os.Stat(test.Expected); os.IsNotExist(err) {
			t.Fatalf("Expected to find %s, but could not find the file\n", test.Expected)
		}

		os.Remove(test.Expected)
	}
}

func TestWriteCorrectData(t *testing.T) {
	dir, err := os.MkdirTemp("", "go-page")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	testData := testGroup{
		GroupName: "TheGophers",
		Members: []testPerson{
			{Name: "Alice", Age: 19},
			{Name: "Bob", Age: 52},
		},
	}

	tests := []struct {
		Input    interface{}
		Format   string
		Expected []byte
	}{
		{
			Input:    struct{}{},
			Format:   JSON,
			Expected: []byte(`{}`),
		},
		{
			Input:    testData,
			Format:   JSON,
			Expected: []byte(`{"GroupName":"TheGophers","Members":[{"Name":"Alice","Age":19},{"Name":"Bob","Age":52}]}`),
		},
	}

	filename := filepath.Join(dir, "data.json")
	for _, test := range tests {
		err := Write(filename, test.Format, &test.Input)
		if err != nil {
			t.Fatal(err)
		}

		data, err := os.ReadFile(filename)
		if err != nil {
			t.Fatal(err)
		}

		if bytes.Compare(data, test.Expected) != 0 {
			t.Fatalf("Expected: %s\nActual: %s\n", string(test.Expected), string(data))
		}
	}
}

func TestWriteInvalidFormat(t *testing.T) {
	dir, err := os.MkdirTemp("", "go-page")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tests := []struct {
		Filename     string
		Format       string
		ErrorMessage string
	}{
		{
			Filename:     "data",
			Format:       "",
			ErrorMessage: "Expected an error if filename has no extension and format is undefined",
		},
		{
			Filename:     "data.json",
			Format:       "invalidformat",
			ErrorMessage: "Expected an error if an invalid format is provided",
		},
	}

	for _, test := range tests {
		filename := filepath.Join(dir, test.Filename)
		err = Write(filename, test.Format, struct{}{})
		if err == nil {
			t.Fatal(test.ErrorMessage)
		}

		os.Remove(filename)
	}
}

// ##########
// Read tests
// ##########

func TestReadFilenameIsRequired(t *testing.T) {
	if err := Read("", nil); err == nil {
		t.Fatal("Expected an error if filename is missing")
	}
}

func TestReadDataCorrectly(t *testing.T) {
	dir, err := os.MkdirTemp("", "go-page")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tests := []struct {
		Filename string
		Format   string
		Input    testGroup
	}{
		{
			Filename: "data.json",
			Format:   "JSON",
			Input:    testGroup{},
		},
		{
			Filename: "data.json",
			Format:   "jsOn",
			Input: testGroup{
				GroupName: "TheGophers",
			},
		},
		{
			Filename: "data.json",
			Format:   JSON,
			Input: testGroup{
				GroupName: "TheGophers",
				Members: []testPerson{
					{Name: "Alice", Age: 72},
					{Name: "Bob", Age: 31},
				},
			},
		},
	}

	for _, test := range tests {
		filename := filepath.Join(dir, test.Filename)

		err = Write(filename, test.Format, test.Input)
		if err != nil {
			t.Fatal(err)
		}

		var actual testGroup

		err = Read(filename, &actual)
		if err != nil {
			t.Fatal(err)
		}

		if test.Input.GroupName != actual.GroupName || len(test.Input.Members) != len(actual.Members) {
			t.Fatalf("Expected: %s %v\nActual: %s %v\n",
				test.Input.GroupName,
				test.Input.Members,
				actual.GroupName,
				actual.Members,
			)
		}

		for i, member := range test.Input.Members {
			if member.Name != actual.Members[i].Name || member.Age != actual.Members[i].Age {
				t.Fatalf("Expected: %s %d\nActual: %s %d\n",
					member.Name,
					member.Age,
					actual.Members[i].Name,
					actual.Members[i].Age,
				)
			}
		}

		os.Remove(filename)
	}
}

func TestReadInvalidFormat(t *testing.T) {
	dir, err := os.MkdirTemp("", "go-page")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tests := []struct {
		Filename string
	}{
		{Filename: "data"},
		{Filename: "data."},
		{Filename: "data.invalid"},
	}

	for _, test := range tests {
		filename := filepath.Join(dir, test.Filename)

		err = os.WriteFile(filename, []byte(""), 0666)
		if err != nil {
			t.Fatal(err)
		}

		var tmp interface{}

		err = Read(filename, &tmp)
		if err == nil {
			t.Fatalf("Expected the extension %s to return an error\n", filepath.Ext(test.Filename))
		}
	}
}

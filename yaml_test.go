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

func TestWriteYAMLFilenameIsRequired(t *testing.T) {
	if err := WriteYAML("", struct{}{}); err == nil {
		t.Fatal("Expected error if filename is empty")
	}
}

func TestWriteYAML(t *testing.T) {
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
		Expected []byte
	}{
		{
			Input: testGroup{},
			Expected: []byte(`groupname: ""
members: []
`),
		},
		{
			Input: testData,
			Expected: []byte(`groupname: TheGophers
members:
    - name: Alice
      age: 19
    - name: Bob
      age: 52
`),
		},
	}

	filename := filepath.Join(dir, "data.yaml")
	for _, test := range tests {
		err := WriteYAML(filename, &test.Input)
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

// ##########
// Read tests
// ##########

func TestReadYAMLFilenameIsRequired(t *testing.T) {
	if err := ReadYAML("", nil); err == nil {
		t.Fatal("Expected an error if filename is missing")
	}
}

func TestReadYAML(t *testing.T) {
	dir, err := os.MkdirTemp("", "go-page")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tests := []struct {
		Input testGroup
	}{
		{
			Input: testGroup{},
		},
		{
			Input: testGroup{
				GroupName: "TheGophers",
			},
		},
		{
			Input: testGroup{
				GroupName: "TheGophers",
				Members: []testPerson{
					{Name: "Alice", Age: 72},
					{Name: "Bob", Age: 31},
				},
			},
		},
	}

	filename := filepath.Join(dir, "data.yaml")
	for _, test := range tests {
		err = WriteYAML(filename, test.Input)
		if err != nil {
			t.Fatal(err)
		}

		var actual testGroup

		err = ReadYAML(filename, &actual)
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

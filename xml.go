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
	"encoding/xml"
	"errors"
	"os"
)

// WriteXML writes the XML encoding of in to the named file.
// If the file does not exist, WriteXML creates it with permissions 0666.
// However, if the file already exist then WriteXML will truncate it before
// writing to it, without changing permissions.
func WriteXML(filename string, in interface{}) error {
	if filename == "" {
		return errors.New("filename is undefined")
	}

	data, err := xml.Marshal(&in)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

// ReadXML reads the specified file and parses the data and
// stores the result in the value pointed to by out.
func ReadXML(filename string, out interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(data, &out)
	if err != nil {
		return err
	}

	return nil
}

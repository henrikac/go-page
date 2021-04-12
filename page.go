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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	JSON = "JSON"
)

// Write Encodes in to the specified format and write it to file.
// If the file does not exist, Write creates it with permissions 0666.
// However, if the file already exist then Write will truncate it before
// writing to it, without changing permissions.
func Write(filename, format string, in interface{}) error {
	if filename == "" {
		return errors.New("filename is undefined")
	}

	ext := filepath.Ext(filename)

	if ext == "" && format != "" {
		filename = fmt.Sprintf("%s.%s", filename, strings.ToLower(format))
	} else if ext != "" {
		if format == "" {
			format = strings.ToUpper(strings.TrimPrefix(ext, "."))
		} else if format != strings.ToUpper(strings.TrimPrefix(ext, ".")) {
			filename = fmt.Sprintf("%s.%s", filename, strings.ToLower(format))
		}
	}

	var data []byte
	var err error

	switch format {
	case JSON:
		data, err = json.Marshal(in)
	default:
		err = fmt.Errorf("%s is not a supported format\n", strings.ToLower(format))
	}

	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0666)
	if err != nil {
		return err
	}

	return nil
}

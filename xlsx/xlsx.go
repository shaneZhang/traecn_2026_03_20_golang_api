// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xlsx

import (
	"bytes"
	"os"

	"github.com/tealeg/xlsx"
)

func ReadXlsxFile(path string) [][]string {
	file, err := xlsx.OpenFile(path)
	if err != nil {
		panic(err)
	}

	res := [][]string{}
	for _, sheet := range file.Sheets {
		for _, row := range sheet.Rows {
			line := []string{}
			for _, cell := range row.Cells {
				text := cell.String()
				line = append(line, text)
			}
			res = append(res, line)
		}
		break
	}

	return res
}

// ReadXlsxFileBytes reads xlsx file from bytes
func ReadXlsxFileBytes(data []byte) ([][]string, error) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "*.xlsx")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpFile.Name())

	// Write data
	_, err = tmpFile.Write(data)
	if err != nil {
		return nil, err
	}
	tmpFile.Close()

	// Read file
	file, err := xlsx.OpenFile(tmpFile.Name())
	if err != nil {
		return nil, err
	}

	res := [][]string{}
	for _, sheet := range file.Sheets {
		for _, row := range sheet.Rows {
			line := []string{}
			for _, cell := range row.Cells {
				text := cell.String()
				line = append(line, text)
			}
			res = append(res, line)
		}
		break
	}

	return res, nil
}

// WriteXlsxFileBytes writes data to xlsx format bytes
func WriteXlsxFileBytes(data [][]string) ([]byte, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Sheet1")
	if err != nil {
		return nil, err
	}

	for _, rowData := range data {
		row := sheet.AddRow()
		for _, cellData := range rowData {
			cell := row.AddCell()
			cell.Value = cellData
		}
	}

	var buf bytes.Buffer
	err = file.Write(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

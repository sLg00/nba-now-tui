package tui

import (
	"github.com/evertras/bubble-table/table"
	"reflect"
	"strconv"
	"strings"
)

// buildTables is a generic function that takes 'headers' as []string, rows as interface{}
// (to support multiline and single-line resultsets) and any type (to use as a sampleType for field mappings)
// and returns a table.Model object which uses the struct tags to define field visibility and more.
// It's extendable and usable across all TUI views that have tables
func buildTables[T any](headers []string, rows interface{}, sampleType T) table.Model {
	itemType := reflect.TypeOf(sampleType)

	isVisible := make(map[string]bool)
	isID := make(map[string]bool)
	displayMap := make(map[string]string)
	columnWidthMap := make(map[string]string)

	for i := 0; i < itemType.NumField(); i++ {
		field := itemType.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			continue
		}

		tagParts := strings.Split(jsonTag, ",")
		jsonName := tagParts[0]

		isVisibleTag := field.Tag.Get("isVisible")
		if isVisibleTag == "true" {
			isVisible[jsonName] = true
		}

		idTag := field.Tag.Get("isID")
		if idTag == "true" {
			isID[jsonName] = true
		}

		displayTag := field.Tag.Get("display")
		if displayTag != "" {
			displayMap[jsonName] = displayTag
		} else {
			displayMap[jsonName] = jsonName
		}

		columnWidthTag := field.Tag.Get("width")
		if columnWidthTag != "" {
			columnWidthMap[jsonName] = columnWidthTag
		}
	}

	var tableColumns []table.Column
	for _, header := range headers {
		if isVisible[header] && !isID[header] {
			displayName := header
			columnWidthInt := 13
			if name, ok := displayMap[header]; ok {
				displayName = name
			}
			if name, ok := columnWidthMap[header]; ok {
				columnWidthInt, _ = strconv.Atoi(name)
			}
			tableColumns = append(tableColumns, table.NewColumn(header, displayName, columnWidthInt))
		}
	}

	var tableRows []table.Row

	switch typedRows := rows.(type) {
	case [][]string:
		for _, row := range typedRows {
			rowMap := make(table.RowData)
			for i, value := range row {
				if i < len(headers) {
					headerName := headers[i]
					if isVisible[headers[i]] || isID[headers[i]] {
						rowMap[headerName] = value
					}
				}
			}
			tableRows = append(tableRows, table.NewRow(rowMap))
		}

	case []string:
		rowMap := make(table.RowData)
		for i, value := range typedRows {
			if i < len(headers) {
				headerName := headers[i]
				if isVisible[headers[i]] || isID[headers[i]] {
					rowMap[headerName] = value
				}
			}
		}
		tableRows = append(tableRows, table.NewRow(rowMap))
	}

	return table.New(tableColumns).WithRows(tableRows).
		SelectableRows(false).
		WithMaxTotalWidth(WindowSize.Width - 10).WithBaseStyle(TableStyle).
		Focused(false)
}

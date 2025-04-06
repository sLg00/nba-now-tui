package tui

import (
	"github.com/evertras/bubble-table/table"
	"log"
	"reflect"
	"strings"
)

func buildTables[T any](headers []string, rows [][]string, sampleType T) table.Model {
	itemType := reflect.TypeOf(sampleType)

	log.Println(headers)

	isVisible := make(map[string]bool)
	isID := make(map[string]bool)

	for i := 0; i < itemType.NumField(); i++ {
		field := itemType.Field(i)

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		tagParts := strings.Split(jsonTag, ",")
		jsonName := tagParts[0]

		isVisibleTag := field.Tag.Get("isVisible")
		idTag := field.Tag.Get("isID")

		if isVisibleTag == "true" {
			isVisible[jsonName] = true
		}

		if idTag == "true" {
			isID[jsonName] = true
		}
	}

	var tableColumns []table.Column
	for _, header := range headers {
		if isVisible[header] && !isID[header] {
			tableColumns = append(tableColumns, table.NewColumn(header, header, 15))
		}
	}

	var tableRows []table.Row
	for _, row := range rows {
		rowMap := make(table.RowData)

		for i, value := range row {
			if i < len(headers) {
				headerName := headers[i]

				if isVisible[headerName] || isID[headerName] {
					rowMap[headerName] = value
				}
			}
		}

		tableRows = append(tableRows, table.NewRow(rowMap))
	}

	return table.New(tableColumns).WithRows(tableRows).
		SelectableRows(false).
		WithMaxTotalWidth(100).WithBaseStyle(TableStyle).Focused(true)
}

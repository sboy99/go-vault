package ui

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/sboy99/go-vault/pkg/utils"
)

func RenderTable(headers []string, contents []interface{}) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)
	transformedContents, err := transformContents(contents)
	if err != nil {
		return err
	}
	table.AppendBulk(transformedContents) // Add Bulk Data
	table.Render()
	return nil
}

func transformContents(contents []interface{}) ([][]string, error) {
	var transformedContents [][]string
	for _, content := range contents {
		values, err := utils.GetStructValues(content)
		if err != nil {
			return nil, err
		}
		transformedContents = append(transformedContents, values)
	}
	return transformedContents, nil
}

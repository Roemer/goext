package goext

import "testing"

func TestTablePrinter(t *testing.T) {
	tablePrinter := NewTablePrinter(nil)
	tablePrinter.SetHeaders("Id", "Name", "Age", "City")
	tablePrinter.Columns[0].ValueAlignment = TABLE_PRINTER_ALIGNMENT_RIGHT
	tablePrinter.AddRows(
		[]any{1, "Alice", "30", "New York"},
		[]any{2, "Bob", "25", "Los Angeles"},
		[]any{3, "Charlie", "35", "Chicago"},
	)
	tablePrinter.PrintStdout()
}

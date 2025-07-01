package goext

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

type TablePrinterAlignment int

const (
	TABLE_PRINTER_ALIGNMENT_LEFT TablePrinterAlignment = iota
	TABLE_PRINTER_ALIGNMENT_RIGHT
)

var TablePrinterStyleDefault = &TablePrinterStyle{
	TopLeft:               "┌",
	TopRight:              "┐",
	TopIntersection:       "┬",
	TopSpacer:             "─",
	MiddleLeft:            "│",
	MiddleRight:           "│",
	MiddleIntersection:    "│",
	SeparatorLeft:         "├",
	SeparatorRight:        "┤",
	SeparatorIntersection: "┼",
	SeparatorSpacer:       "─",
	BottomLeft:            "└",
	BottomRight:           "┘",
	BottomIntersection:    "┴",
	BottomSpacer:          "─",
}

var TablePrinterStyleRounded = &TablePrinterStyle{
	TopLeft:               "╭",
	TopRight:              "╮",
	TopIntersection:       "┬",
	TopSpacer:             "─",
	MiddleLeft:            "│",
	MiddleRight:           "│",
	MiddleIntersection:    "│",
	SeparatorLeft:         "├",
	SeparatorRight:        "┤",
	SeparatorIntersection: "┼",
	SeparatorSpacer:       "─",
	BottomLeft:            "╰",
	BottomRight:           "╯",
	BottomIntersection:    "┴",
	BottomSpacer:          "─",
}

var TablePrinterStyleAscii = &TablePrinterStyle{
	TopLeft:               "+",
	TopRight:              "+",
	TopIntersection:       "+",
	TopSpacer:             "-",
	MiddleLeft:            "|",
	MiddleRight:           "|",
	MiddleIntersection:    "|",
	SeparatorLeft:         "+",
	SeparatorRight:        "+",
	SeparatorIntersection: "+",
	SeparatorSpacer:       "-",
	BottomLeft:            "+",
	BottomRight:           "+╯",
	BottomIntersection:    "+",
	BottomSpacer:          "-",
}

type TablePrinter struct {
	Options *TablePrinterOptions
	Columns []*TablePrinterColumn
	Rows    []*TablePrinterRow
}

type TablePrinterOptions struct {
	Padding int
	Style   *TablePrinterStyle
}

type TablePrinterStyle struct {
	TopLeft               string
	TopRight              string
	TopIntersection       string
	TopSpacer             string
	MiddleLeft            string
	MiddleRight           string
	MiddleIntersection    string
	SeparatorLeft         string
	SeparatorRight        string
	SeparatorIntersection string
	SeparatorSpacer       string
	BottomLeft            string
	BottomRight           string
	BottomIntersection    string
	BottomSpacer          string
}

type TablePrinterColumn struct {
	Header          string
	HeaderAlignment TablePrinterAlignment
	ValueAlignment  TablePrinterAlignment
	Hide            bool
	maxLength       int
}

type TablePrinterRow struct {
	Values []string
}

// Create a new TablePrinter with the given options or default options.
func NewTablePrinter(options *TablePrinterOptions) *TablePrinter {
	if options == nil {
		options = &TablePrinterOptions{
			Padding: 2,
		}
	}
	if options.Style == nil {
		options.Style = TablePrinterStyleDefault
	}
	return &TablePrinter{
		Options: options,
	}
}

// Set the headers for the TablePrinter.
func (tp *TablePrinter) SetHeaders(headers ...string) {
	columns := []*TablePrinterColumn{}
	for _, header := range headers {
		newColumn := &TablePrinterColumn{
			Header:          header,
			HeaderAlignment: TABLE_PRINTER_ALIGNMENT_LEFT,
			ValueAlignment:  TABLE_PRINTER_ALIGNMENT_LEFT,
		}
		columns = append(columns, newColumn)
	}
	tp.Columns = columns
}

// Add rows to the TablePrinter with the given values.
func (tp *TablePrinter) AddRows(values ...[]string) {
	for _, value := range values {
		newRow := &TablePrinterRow{
			Values: value,
		}
		tp.Rows = append(tp.Rows, newRow)
	}
}

// Print the table to stdout.
func (tp *TablePrinter) PrintStdout() {
	tp.Print(os.Stdout)
}

// Print the table to a file.
func (tp *TablePrinter) PrintToFile(filePath string) error {
	fileWriter, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer fileWriter.Close()
	tp.Print(fileWriter)
	return nil
}

// Print the table to the given writer.
func (tp *TablePrinter) Print(writer io.Writer) {
	// Remove hidden columns
	columns := slices.DeleteFunc(tp.Columns, func(column *TablePrinterColumn) bool {
		return column.Hide
	})

	// Collect the header texts
	headerTexts := []string{}
	for _, column := range columns {
		headerTexts = append(headerTexts, column.Header)
	}

	// Calculate the max length per column
	for i := range columns {
		// Add the length of the header
		maxLength := len(columns[i].Header)
		// Adjust the value if a row has a longer value
		for _, row := range tp.Rows {
			value := row.Values[i]
			if len(value) > maxLength {
				maxLength = len(value)
			}
		}
		columns[i].maxLength = maxLength
	}

	// Top
	tp.writeTableSpacer(writer, tp.Options.Style.TopLeft, tp.Options.Style.TopRight, tp.Options.Style.TopIntersection, tp.Options.Style.TopSpacer, columns)
	// Header
	tp.writeTableRow(writer, tp.Options.Style.MiddleLeft, tp.Options.Style.MiddleRight, tp.Options.Style.MiddleIntersection, columns, headerTexts, true)
	// Separator
	tp.writeTableSpacer(writer, tp.Options.Style.SeparatorLeft, tp.Options.Style.SeparatorRight, tp.Options.Style.SeparatorIntersection, tp.Options.Style.SeparatorSpacer, columns)
	// Rows
	for _, row := range tp.Rows {
		tp.writeTableRow(writer, tp.Options.Style.MiddleLeft, tp.Options.Style.MiddleRight, tp.Options.Style.MiddleIntersection, columns, row.Values, false)
	}
	// Bottom
	tp.writeTableSpacer(writer, tp.Options.Style.BottomLeft, tp.Options.Style.BottomRight, tp.Options.Style.BottomIntersection, tp.Options.Style.BottomSpacer, columns)
}

func (tp *TablePrinter) writeTableSpacer(writer io.Writer, left string, right string, intersection string, spacer string, columns []*TablePrinterColumn) {
	fmt.Fprint(writer, left)
	for i, column := range columns {
		fmt.Fprint(writer, strings.Repeat(spacer, column.maxLength+2*tp.Options.Padding))
		if i < len(columns)-1 {
			fmt.Fprint(writer, intersection)
		}
	}
	fmt.Fprint(writer, right)
	fmt.Fprintln(writer)
}

func (tp *TablePrinter) writeTableRow(writer io.Writer, left string, right string, intersection string, columns []*TablePrinterColumn, values []string, isHeader bool) {
	paddingString := strings.Repeat(" ", tp.Options.Padding)
	fmt.Fprint(writer, left)
	for i, column := range columns {
		alignment := column.ValueAlignment
		if isHeader {
			alignment = column.HeaderAlignment
		}
		valueString := ""
		switch alignment {
		case TABLE_PRINTER_ALIGNMENT_LEFT:
			valueString = fmt.Sprintf("%-*s", column.maxLength, values[i])
		case TABLE_PRINTER_ALIGNMENT_RIGHT:
			valueString = fmt.Sprintf("%*s", column.maxLength, values[i])
		}
		fmt.Fprint(writer, paddingString+valueString+paddingString)
		if i < len(columns)-1 {
			fmt.Fprint(writer, intersection)
		}
	}
	fmt.Fprint(writer, right)
	fmt.Fprintln(writer)
}

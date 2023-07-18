package iui

import "image"
import "image/color"

type Image = image.Image

type TableValue interface {
	TableValue() TableValue
}

// TableValueString is a TableValue that stores a string. TableString is
// used for displaying text in a Table.
type TableValueString string

func (s TableValueString) TableValue() TableValue {
	return s
}

// TableValueImage is a TableValue that represents an Image.
type TableValueImage struct {
	I Image
}

func (i TableValueImage) TableValue() TableValue {
	return i
}

// TableValueInt is a TableValue that stores integers. These are used for
// progressbars. Due to current limitations of libui, they also
// represent checkbox states, via TableFalse and TableTrue.
type TableValueInt int

// TableFalse and TableTrue are the Boolean constants for TableInt.
const (
	TableValueFalse TableValueInt = 0
	TableValueTrue  TableValueInt = 1
)

func (i TableValueInt) TableValue() TableValue {
	return i
}

// TableValueColor is a TableValue that represents a color.
type TableValueColor struct {
	color.RGBA
}

func (c TableValueColor) TableValue() TableValue {
	return c
}

type TableModel interface {
	Model() any
	Free()
	RowInserted(index int)
	RowChanged(index int)
	RowDeleted(index int)
}

// TableModelHandler defines the methods that TableModel
// calls when it needs data.
type TableModelHandler interface {
	// ColumnTypes returns a slice of value types of the data
	// stored in the model columns of the TableModel.
	// Each entry in the slice should ideally be a zero value for
	// the TableValue type of the column in question; the number
	// of elements in the slice determines the number of model
	// columns in the TableModel. The returned slice must remain
	// constant through the lifetime of the TableModel. This
	// method is not guaranteed to be called depending on the
	// system.
	ColumnTypes(m TableModel) []TableValue

	// NumRows returns the number or rows in the TableModel.
	// This value must be non-negative.
	NumRows(m TableModel) int

	// CellValue returns a TableValue corresponding to the model
	// cell at (row, column). The type of the returned TableValue
	// must match column's value type. Under some circumstances,
	// nil may be returned; refer to the various methods that add
	// columns to Table for details.
	CellValue(m TableModel, row, column int) TableValue

	// SetCellValue changes the model cell value at (row, column)
	// in the TableModel. Within this function, either do nothing
	// to keep the current cell value or save the new cell value as
	// appropriate. After SetCellValue is called, the Table will
	// itself reload the table cell. Under certain conditions, the
	// TableValue passed in can be nil; refer to the various
	// methods that add columns to Table for details.
	SetCellValue(m TableModel, row, column int, value TableValue)
}

// TableModelColumnNeverEditable and
// TableModelColumnAlwaysEditable are the value of an editable
const (
	TableModelColumnNeverEditable  = -1
	TableModelColumnAlwaysEditable = -2
)

type TableTextColumnOptionalParams struct {
}

type Table interface {
	AppendTextColumn(name string, textModelColumn int, textEditableModelColumn int, textParams *TableTextColumnOptionalParams)
	AppendImageTextColumn(name string, imageModelColumn int, textModelColumn int, textEditableModelColumn int, textParams *TableTextColumnOptionalParams)
	AppendCheckboxColumn(name string, checkboxModelColumn int, checkboxEditableModelColumn int)
	AppendCheckboxTextColumn(name string, checkboxModelColumn int, checkboxEditableModelColumn int, textModelColumn int, textEditableModelColumn int, textParams *TableTextColumnOptionalParams)
}

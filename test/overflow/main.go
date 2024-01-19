package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

type Model []Row

func NewModel(nrow, ncol int) Model {
	res := Model{}
	for j := 0; j < nrow; j++ {
		row := Row{}
		for i := 0; i < ncol; i++ {
			var val Value
			switch i % 3 {
			case 0:
				val = NewValue(fmt.Sprintf("日本語 %d:%d", i, j))
			case 1:
				val = NewValue(i%6 == 1)
			case 2:
				val = NewValue(i%6 == 2)
			}
			row = append(row, val)
		}
		res = append(res, row)
	}
	return res
}

func (m Model) NumRows() int {
	return len(m)
}

func (m Model) FetchRow(index int) Row {
	if index < 0 || index >= len(m) {
		return nil
	}
	return m[index]
}

func (m Model) UpdateRow(index int, to Row) {
	if index < 0 || index >= len(m) || to == nil {
		return
	}
	m[index] = to
}

func mainTableWithOverflow() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewBox()

	nrow := 150
	ncol := 10

	hbox := NewTray()
	hbox.Append(NewLabel("Above table"))

	box.Append(hbox)

	model := NewModel(nrow, ncol)
	table := NewTable(model)
	overflow := NewOverflow(table, 480, 320)

	box.Append(overflow)

	hbox2 := NewTray()
	hbox2.Append(NewLabel("Below table"))

	box.Append(hbox2)

	cols := []*Column{}

	for i := 0; i < ncol; i++ {
		var column *Column
		switch i % 3 {
		case 0:
			column = NewTextColumn(fmt.Sprintf("Col %d", i), i)
		case 1:
			column = NewCheckboxColumn(fmt.Sprintf("Col %d", i), i)
		case 2:
			column = NewButtonColumn(fmt.Sprintf("Col %d", i), i, "toolPencil")
		}
		table.AppendColumn(column)
		cols = append(cols, column)
	}

	fmt.Printf("numRows %t\n", table.NumRows() == nrow)

	for i := 0; i < ncol; i++ {
		column := cols[i]
		column.SetWidth(-1)
		width := column.Width()
		fmt.Printf("Column width: %d\n", width)
	}

	table.OnClicked(func(tab *Table, row, col int) {
		fmt.Printf("Table clicked at %d %d\n", row, col)
	})
	table.OnHeaderClicked(func(tab *Table, col int) {
		column := table.Column(col)
		marker := column.Marker()
		if marker == "" || marker == "up" {
			marker = "down"
		} else {
			marker = "up"
		}
		column.SetMarker(marker)
		fmt.Printf("Table header clicked at %d: %s\n", col, marker)
	})

	table.SetHeaderVisible(true)
	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v\n", wi)
		Exit(0)

	})

	w.SetChild(box)
	Main(w)
}

func main() {
	mainTableWithOverflow()
}

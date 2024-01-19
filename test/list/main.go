package main

import "fmt"
import . "github.com/bjorndm/golang-ui"
import "github.com/bjorndm/golang-ui/icon"

type Model []Row

func NewModel(nrow, ncol int) Model {
	res := Model{}
	for j := 0; j < nrow; j++ {
		row := Row{}
		for i := 0; i < ncol; i++ {
			var val Value
			switch i % 3 {
			case 0:
				val = NewValue(fmt.Sprintf("Card 日本語 %d:%d", i, j))
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

var window *Window

func inDialog() func(b *Button) {
	return func(b *Button) {
		tab := NewTab()
		box := NewBox()
		box.Append(NewLabel("Title"))
		tray := NewTray()
		model := NewModel(5, 4)
		list := NewList(model, cardTemplate)
		list.CreateCards()
		tray.Append(NewButton("Nothing"))
		tray.Append(NewButton("Happens"))
		box.Append(tray)
		box.Append(list)
		tab.Append("tab1", NewLabel("Not Here"))
		tab.Append("list", box)
		tab.Select(1)
		dialog := NewDialog("List in dialog in tab", tab)
		dialog.Display(window, nil)
	}
}

func cardTemplate(row Row) *Card {
	card := NewCard(row[0].(string))
	card.AppendButtonWithIcon(fmt.Sprintf("In dialog"), icon.Question).OnClicked(inDialog())

	for i := 1; i < len(row); i++ {
		switch i % 3 {
		case 0:
			card.AppendLabel(row[i].(string))
		case 1:
			check := card.AppendCheckbox("Card checkbox")
			check.SetChecked(row[1].(bool))
		case 2:
			card.AppendButtonWithIcon(fmt.Sprintf("Button %d", i), icon.ToolPencil)
		}
	}
	return card
}

func mainList() {
	Init()
	w := NewWindow("test window", 640, 480, false)
	window = w

	box := NewBox()

	nrow := 150
	ncol := 10

	hbox := NewTray()
	hbox.Append(NewLabel("Above table"))

	box.Append(hbox)

	model := NewModel(nrow, ncol)
	list := NewList(model, cardTemplate)
	list.CreateCards()

	box.Append(list)

	hbox2 := NewTray()
	hbox2.Append(NewLabel("Below table"))

	box.Append(hbox2)

	fmt.Printf("numRows %t\n", list.NumRows() == nrow)

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v\n", wi)
		Exit(0)

	})

	w.SetChild(box)
	Main(w)
}

func main() {
	mainList()
}

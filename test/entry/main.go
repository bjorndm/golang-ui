package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

func mainEntry() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewVerticalBox()
	w.SetChild(box)

	hbox := NewHorizontalBox()
	vbox := NewVerticalBox()

	nbut := 10
	combo := make([]*Entry, nbut)
	for i := 0; i < nbut; i++ {
		nr := i + 1
		switch i % 3 {
		case 0:
			combo[i] = NewEntry()
		case 1:
			combo[i] = NewPasswordEntry()
		case 2:
			combo[i] = NewSearchEntry()
		}
		combo[i].OnChanged(func(c *Entry) {
			fmt.Printf("Entry %d changed", nr)
			txt := c.Text()
			fmt.Printf("Text: %s\n", txt)
		})
		if i < 5 {
			combo[i].SetText(fmt.Sprintf("hello %d", nr))
			vbox.Append(combo[i])
		} else {
			combo[i].SetText(fmt.Sprintf("今日は %d", nr))

			// XXX: bug: the first combobox in the vbox takes up all space
			// making the others hard to see.
			hbox.Append(combo[i])
		}
	}

	box.Append(vbox)
	box.Append(hbox)

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v", wi)
		Exit(0)

	})
	Main(w)
}

func main() {
	mainEntry()
}

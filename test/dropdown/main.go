package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

func mainDropdown() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewVerticalBox()

	hbox := NewHorizontalBox()
	vbox := NewVerticalBox()

	box.Append(vbox)
	box.Append(hbox)

	selections := []string{"one", "two", "三", "four", "five"}

	nbut := 10
	drop := make([]*Dropdown, nbut)
	for i := 0; i < nbut; i++ {
		nr := i + 1
		drop[i] = NewDropdown()
		for _, t := range selections {
			drop[i].Append(t)
		}
		drop[i].OnSelected(func(c *Dropdown) {
			fmt.Printf("Dropdown %d clicked:\n", nr)
			si := c.Selected()
			txt := c.Text()
			fmt.Printf("Selected: %d: %\n", si, txt)
		})
		if i < 5 {
			drop[i].SetText("two")
			vbox.Append(drop[i])
		} else {
			// XXX: bug: the first dropbox in the vbox takes up all space
			// making the others hard to see.
			hbox.Append(drop[i])
		}
	}

	w.SetChild(box)

	w.OnClosing(func(wi *Window) {
		fmt.Println("Closing window: %v", wi)

	})
	Main(w)
}

func main() {
	mainDropdown()
}

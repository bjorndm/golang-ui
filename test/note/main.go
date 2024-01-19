package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

func mainNote() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewVerticalBox()

	hbox := NewHorizontalBox()
	vbox := NewVerticalBox()

	box.Append(vbox)
	box.Append(hbox)

	nbut := 10
	combo := make([]*Note, nbut)
	for i := 0; i < nbut; i++ {
		nr := i + 1
		combo[i] = NewNote()
		combo[i].OnChanged(func(c *Note) {
			fmt.Printf("Note %d changed", nr)
			txt := c.Text()
			fmt.Printf("Text: %s\n", txt)
		})
		if i < 5 {
			combo[i].SetText(fmt.Sprintf("hello %d\n今日は", nr))
			vbox.Append(combo[i])
		} else {
			combo[i].SetText(fmt.Sprintf("hello %d\n今日は", nr))
			hbox.Append(combo[i])
		}
	}

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v", wi)
		Exit(0)

	})
	w.SetChild(box)
	Main(w)
}

func main() {
	mainNote()
}

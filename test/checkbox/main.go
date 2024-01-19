package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

func mainCheckbox() {
	Init()
	w := NewWindow("test checkbox", 640, 480, false)

	box := NewVerticalBox()

	hbox := NewHorizontalBox()
	vbox := NewVerticalBox()

	box.Append(vbox)
	box.Append(hbox)

	nbut := 10
	checkbox := make([]*Checkbox, nbut)
	for i := 0; i < nbut; i++ {
		nr := i + 1
		checkbox[i] = NewCheckbox(fmt.Sprintf("button %d", nr))
		checkbox[i].SetChecked((nr % 2) == 0)
		checkbox[i].OnClicked(func(b *Checkbox) {
			fmt.Printf("Checkbox %d clicked: %v\n", nr, b)
		})
		if i < 5 {
			vbox.Append(checkbox[i])
		} else {
			hbox.Append(checkbox[i])
		}
	}

	w.SetChild(box)

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v\n", wi)

	})
	Main(w)
}

func main() {
	mainCheckbox()
}

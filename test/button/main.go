package main

import "fmt"
import . "github.com/bjorndm/golang-ui"
import "github.com/bjorndm/golang-ui/icon"

func mainButton() {
	Init()
	w := NewWindow("test button", 640, 480, false)

	box := NewVerticalBox()

	hbox := NewHorizontalBox()
	vbox := NewVerticalBox()

	box.Append(vbox)
	box.Append(hbox)

	nbut := 10
	button := make([]*Button, nbut)
	for i := 0; i < nbut; i++ {
		nr := i + 1
		switch i % 3 {
		case 0:
			button[i] = NewButton(fmt.Sprintf("button %d", nr))
		case 1:
			button[i] = NewButtonWithIcon(fmt.Sprintf("button %d", nr), icon.Medal1)
		case 2:
			button[i] = NewButtonWithIcon("", icon.Key)
		}
		button[i].OnClicked(func(b *Button) {
			fmt.Printf("Button %d clicked: %v\n", nr, b)
		})
		if i < 5 {
			vbox.Append(button[i])
		} else {
			hbox.Append(button[i])
		}
	}

	w.SetChild(box)

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v\n", wi)

	})
	Main(w)
}

func main() {
	mainButton()
}

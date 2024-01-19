package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

func mainTabs() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewTab()

	hbox := NewHorizontalBox()
	vbox := NewVerticalBox()
	vgroup := NewVerticalBox()
	group := NewGroup("Group")
	group.SetChild(vgroup)

	box.AppendWithIcon("user", "Tab 1", vbox)
	box.AppendWithIcon("users", "Tab 2", hbox)
	box.AppendWithIcon("organization", "Tab 3", group)

	nbut := 10
	button := make([]*Button, nbut)
	for i := 0; i < nbut; i++ {
		nr := i + 1
		button[i] = NewButton(fmt.Sprintf("button %d", nr))
		button[i].OnClicked(func(b *Button) {
			fmt.Printf("Button %d clicked: %v\n", nr, b)
		})
		if i < 4 {
			vbox.Append(button[i])
		} else if i < 8 {
			hbox.Append(button[i])
		} else {
			vgroup.Append(button[i])
		}
	}

	w.SetChild(box)

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v\n", wi)
		Exit(0)

	})
	Main(w)
}

func main() {
	mainTabs()
}

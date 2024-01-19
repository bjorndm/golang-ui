package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

func mainCheckbox() {
	Init()
	w := NewWindow("test toggle", 640, 480, false)

	box := NewVerticalBox()

	hbox := NewHorizontalBox()
	vbox := NewVerticalBox()
	hr := NewRadio()
	vr := NewRadio()

	hr.OnSelected(func(r *Radio) {
		fmt.Printf("Radio clicked: %s\n", r.Text())
	})

	vr.OnSelected(func(r *Radio) {
		fmt.Printf("Radio clicked: %s\n", r.Text())
	})

	box.Append(vbox)
	box.Append(hbox)

	nbut := 10
	toggle := make([]*Toggle, nbut)
	for i := 0; i < nbut; i++ {
		nr := i + 1
		radio := hr
		if i < 5 {
			radio = vr
		}

		toggle[i] = radio.Append(fmt.Sprintf("button %d", nr))
		toggle[i].OnClicked(func(t *Toggle) {
			fmt.Printf("Toggle %d clicked: %v\n", t)
		})
		if i < 5 {
			vbox.Append(toggle[i])
		} else {
			hbox.Append(toggle[i])
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

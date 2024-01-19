package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

func recursiveAlert(w *Window) func(res DialogResult) {
	var ral func(res DialogResult)

	ral = func(res DialogResult) {
		fmt.Printf("Recursive alert licked with alert result: %s\n", res)
		ShowConfirmAlert(w, "Recursive alert", "This is\na recursive alert.\nOK recurses.", func(res DialogResult) {
			fmt.Printf("Recursive alert clicked with alert result: %s\n", res)
			if res == DialogResultOK {
				ral(res)
			}
		})
	}
	return ral
}

func mainAlert() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewBox()

	hbox := NewTray()
	vbox := NewBox()

	nbut := 10
	button := make([]*Button, nbut)
	for i := 0; i < nbut; i++ {
		nr := i + 1

		button[i] = NewButton(fmt.Sprintf("button %d", nr))
		button[i].OnClicked(func(b *Button) {
			var alert *Alert
			switch (nr - 1) % 3 {
			case 0:
				alert = NewAlert(fmt.Sprintf("alert button %d", nr), "more details\nhere")
			case 1:
				alert = NewErrorAlert(fmt.Sprintf("error button %d", nr), "more error details\nhere\nthird line")
			case 2:
				alert = NewConfirmAlert(fmt.Sprintf("confirm button %d", nr), "more confirmation\ndetails here")
			}
			fmt.Printf("Button %d clicked\n", nr)
			alert.Display(w, func(res DialogResult) {
				if nr == 1 {
					alert := recursiveAlert(w)
					alert(res)
				}
				fmt.Printf("Button %d clicked with alert result: %s\n", nr, res)
			})
		})
		if i < 5 {
			vbox.Append(button[i])
		} else {
			hbox.Append(button[i])
		}
	}

	box.Append(vbox)
	box.Append(hbox)
	w.SetChild(box)

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v", wi)

	})

	Main(w)
}

func main() {
	mainAlert()
}

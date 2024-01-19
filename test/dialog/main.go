package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

const longLabel = `Cum multae res in philosophia nequaquam
satis adhuc explicatae sint, tum perdifficilis,
Brute, quod tu minime ignoras,
et perobscura quaestio est de natura deorum,
quae et ad cognitionem animi pulcherrima est
et ad moderandam religionem necessaria.`

func recursiveDialog(w *Window) func(res *Dialog) {
	var ral func(res *Dialog)

	ral = func(d1 *Dialog) {
		fmt.Printf("Recursive dialog clicked with dialog result: %s\n", d1.Result())
		label := NewLabel(longLabel)
		ShowDialogWithDialog(w, "Recursive dialog", label, func(d *Dialog) {
			fmt.Printf("Recursive dialog clicked with dialog result: %s\n", d.Result())
			if d.Result() == DialogResultOK {
				ral(d)
			}
		}).AddButton("Continue", DialogResultOK).AddButton("Stop", DialogResultCancel)

	}
	return ral
}

func mainDialog() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewBox()

	hbox := NewTray()
	vbox := NewBox()

	box.Append(vbox)
	box.Append(hbox)

	nbut := 10
	button := make([]*Button, nbut)
	dialog := make([]*Dialog, nbut)
	drops := []string{"one", "two", "three", "four", "five"}
	ndrops := 3

	for i := 0; i < nbut; i++ {
		nr := i + 1
		grid := NewGrid()
		label := NewLabel(longLabel)
		entry := NewEntry()
		grid.Append(label, 0, i, 1, AlignStart)
		grid.Append(entry, 1, i, 1, AlignStart)
		for j := 0; j < ndrops; j++ {
			d := NewDropdown()
			for _, drop := range drops {
				d.Append(drop)
			}
			grid.AppendWithLabel(fmt.Sprintf("label %d", j), d)
		}

		switch i % 3 {
		case 0:
			dialog[i] = NewDialog(fmt.Sprintf("dialog button %d", nr), grid).
				AddButton("Confirmed", DialogResultOK).
				AddButtonKeepOpen("Cancel No Close", DialogResultCancel)
		case 1:
			dialog[i] = NewDialog(fmt.Sprintf("error button %d", nr), grid)
			dialog[i].AddButton("Error", DialogResultOK)
		case 2:
			dialog[i] = NewDialog(fmt.Sprintf("confirm button %d", nr), grid)
			dialog[i].AddButton("Continue", DialogResultOK)
			dialog[i].AddButton("Stop", DialogResultCancel)
		}

		button[i] = NewButton(fmt.Sprintf("button %d", nr))
		button[i].OnClicked(func(b *Button) {
			dialog := dialog[nr-1]
			fmt.Printf("Button %d clicked\n", nr)
			dialog.DisplayDialog(w, func(d *Dialog) {
				if nr == 1 {
					dialog := recursiveDialog(w)
					dialog(d)
				}
				fmt.Printf("Button %d clicked with dialog result: %s\n", nr, d.Result())
			})
		})
		if i < 5 {
			vbox.Append(button[i])
		} else {
			hbox.Append(button[i])
		}
	}

	w.SetChild(box)

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v", wi)

	})

	Main(w)
}

func main() {
	mainDialog()
}

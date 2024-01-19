package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

const jaLabel = `本ソフトウェアの開発はExaWizards
のスポンサーによるものです。`

const longLabel = `Cum multae res in philosophia nequaquam satis adhuc explicatae sint,
tum perdifficilis, Brute, quod tu minime ignoras, et perobscura quaestio est de
natura deorum, quae et ad cognitionem animi pulcherrima est et ad moderandam
religionem necessaria.`

func mainMenu() {
	Init()
	w := NewWindow("test window with menu", 640, 480, true)

	box := NewVerticalBox()

	hbox := NewHorizontalBox()
	vbox := NewVerticalBox()

	box.Append(vbox)
	box.Append(hbox)
	vbox.Append(NewLabel(longLabel))
	vbox.Append(NewEntry())
	hbox.Append(NewLabel(jaLabel))
	hbox.Append(NewEntry())

	bar := NewMenuBar()
	menuHello := bar.AppendMenu("|Hello|")
	menuHello.AppendItem("World").OnClicked(func(it *MenuItem) {
		fmt.Printf("World clicked\n")
	})
	menu := bar.AppendMenu("|Another Menu|")
	var item7 *MenuItem

	nitem := 14

	for j := 0; j < nitem; j++ {
		nr := j + 1
		title := fmt.Sprintf("Item %d", nr)
		var item *MenuItem
		switch j % 3 {
		case 0:
			item = menu.AppendItem(title)
		case 1:
			item = menu.AppendCheckItem(title)
			if j%6 == 1 {
				item.SetChecked(true)
			}
		case 2:
			item = menu.AppendSeparator()
		}
		if nr == 7 {
			item.Disable()
			item7 = item
		}
		if nr == 8 {
			item.OnClicked(func(it *MenuItem) {
				item7.Enable()
				fmt.Printf("Menu item clicked, enabling item 7: %d, %v\n", nr, it)
			})
		} else {
			item.OnClicked(func(it *MenuItem) {
				fmt.Printf("Menu item clicked: %d, %v\n", nr, it)
			})

		}
	}

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v\n", wi)
		Exit(0)

	})

	w.SetChild(box)
	w.SetMenuBar(bar)

	Main(w)
}

func main() {
	mainMenu()
}

package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

const jaLabel = `本ソフトウェアの開発はExaWizards
のスポンサーによるものです。`

const longLabel = `Cum multae res in philosophia nequaquam satis adhuc explicatae sint,
tum perdifficilis, Brute, quod tu minime ignoras, et perobscura quaestio est de
natura deorum, quae et ad cognitionem animi pulcherrima est et ad moderandam
religionem necessaria.`

func mainJournal() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewVerticalBox()

	hjournal := NewJournal(false)
	vjournal := NewJournal(true)

	hbox := NewHorizontalBox()
	vbox := NewVerticalBox()
	hbox.Append(hjournal)
	vbox.Append(vjournal)
	box.Append(vbox)
	box.Append(hbox)

	nlog := 10

	for i := 0; i < nlog; i++ {
		nr := i + 1
		if i < 5 {
			hjournal.Append(fmt.Sprintf("hello %d\n今日は\n%s\n", nr, jaLabel))
		} else {
			vjournal.Append(fmt.Sprintf("hello %d\n今日は\n%s\n", nr, longLabel))
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
	mainJournal()
}

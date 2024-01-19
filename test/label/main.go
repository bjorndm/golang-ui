package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

const jaLabel = `本ソフトウェアの開発はExaWizards
のスポンサーによるものです。`

const longLabel = `Cum multae res in philosophia nequaquam satis adhuc explicatae sint,
tum perdifficilis, Brute, quod tu minime ignoras, et perobscura quaestio est de
natura deorum, quae et ad cognitionem animi pulcherrima est et ad moderandam
religionem necessaria.`

/* ` De qua [cum] tam variae sint doctissimorum hominum
tamque discrepantes sententiae, magno argumento esse debeat [ea] causa,
principium philosophiae ad h* scientiam, prudenterque Academici a rebus
incertis adsensionem cohibuisse. Quid est enim temeritate turpius aut quid
tam temerarium tamque indignum sapientis gravitate atque constantia quam aut
falsum sentire aut, quod non satis explorate perceptum sit et cognitum,
sine ulla dubitatione defendere?` */

func boxMain() {
	Init()
	w := NewWindow("test label", 800, 640, false)

	// label := NewLabel("hello")
	box := NewVerticalBox()

	vbox := NewVerticalBox()

	vbox2 := NewVerticalBox()

	nlab := 15
	label := make([]*Label, nlab)
	for i := 0; i < len(label); i++ {
		nr := i + 1
		switch i % 3 {
		case 0:
			label[i] = NewLabel(jaLabel)
		case 1:
			label[i] = NewLabel(longLabel)
		case 2:
			label[i] = NewLabel("hello")
		}
		if i < 5 {
			label[i].SetText(fmt.Sprintf("hello %d: %s", nr, label[i].Text()))
			box.Append(label[i])
		} else if i < 10 {
			label[i].SetText(fmt.Sprintf("hello %d: %s", nr, label[i].Text()))
			vbox.Append(label[i])
		} else {
			label[i].SetText(fmt.Sprintf("hello %d: %s", nr, label[i].Text()))
			vbox2.Append(label[i])
		}
	}

	box.Append(vbox)
	box.Append(vbox2)

	w.SetChild(box)

	w.OnClosing(func(wi *Window) {
		fmt.Sprintf("Closing window: %v", wi)
		Exit(1)

	})
	Main(w)
}

func boxAndTrayMain() {
	Init()
	w := NewWindow("test label", 800, 640, false)

	// label := NewLabel("hello")
	box := NewVerticalBox()

	hbox := NewHorizontalBox()
	vbox := NewVerticalBox()

	box.Append(vbox)
	box.Append(hbox)

	nbut := 10
	label := make([]*Label, nbut)
	for i := 0; i < nbut; i++ {
		nr := i + 1
		switch i % 3 {
		case 0:
			label[i] = NewLabel(jaLabel)
		case 1:
			label[i] = NewLabel(longLabel)
		case 2:
			label[i] = NewLabel("hello")
		}
		if i < 5 {
			label[i].SetText(fmt.Sprintf("hello %d: %s", nr, label[i].Text()))
			vbox.Append(label[i])
		} else {
			label[i].SetText(fmt.Sprintf("hello %d: %s", nr, label[i].Text()))
			hbox.Append(label[i])
		}
	}

	w.SetChild(box)

	w.OnClosing(func(wi *Window) {
		fmt.Sprintf("Closing window: %v", wi)
		Exit(1)

	})
	Main(w)
}

func main() {
	boxAndTrayMain()
}

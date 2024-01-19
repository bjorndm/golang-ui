package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

const jaLabel = `本ソフトウェアの開発はExaWizards
のスポンサーによるものです。`

const longLabel = `Cum multae res in philosophia nequaquam satis adhuc explicatae sint,
tum perdifficilis, Brute, quod tu minime ignoras, et perobscura quaestio est de
natura deorum, quae et ad cognitionem animi pulcherrima est et ad moderandam
religionem necessaria.`

/* `De qua [cum] tam variae sint doctissimorum hominum
tamque discrepantes sententiae, magno argumento esse debeat [ea] causa,
principium philosophiae ad h* scientiam, prudenterque Academici a rebus
incertis adsensionem cohibuisse. Quid est enim temeritate turpius aut quid
tam temerarium tamque indignum sapientis gravitate atque constantia quam aut
falsum sentire aut, quod non satis explorate perceptum sit et cognitum,
sine ulla dubitatione defendere?` */

func paneMain() {
	Init()
	w := NewWindow("test label", 800, 640, false)

	// label := NewLabel("hello")
	stack := NewStack()

	bar := NewMenuBar()
	menuHello := bar.AppendMenu("Greeting")
	menuHello.AppendItem("Hello World").OnClicked(func(it *MenuItem) {
		fmt.Printf("%s clicked\n", it.Text())
	})
	menuHello.AppendItem("今日は").OnClicked(func(it *MenuItem) {
		fmt.Printf("%s clicked\n", it.Text())
	})

	menu := bar.AppendMenu("Another Menu")
	menu.AppendCheckItem("Item1").OnClicked(func(it *MenuItem) {
		fmt.Printf("%s clicked\n", it.Text())
	})
	menu.AppendItem("Item2").OnClicked(func(it *MenuItem) {
		fmt.Printf("%s clicked\n", it.Text())
	})

	nlab := 5
	label := make([]*Label, nlab)
	entry := make([]*Entry, nlab)
	box := make([]*Box, nlab)
	drops := []string{"one", "two", "three", "four", "five"}
	ndrops := 3

	panes := make([]*Pane, nlab)
	for i := 0; i < len(label); i++ {
		nr := i + 1
		box[i] = NewBox()
		entry[i] = NewEntry()

		panes[i] = NewPane(fmt.Sprintf("pane %d", nr), 800-i*40, 640-i*40, false)
		panes[i].OnClosing(func(wi *Pane) {
			fmt.Printf("Closing pane: %d\n", nr)
		})
		switch i % 3 {
		case 0:
			label[i] = NewLabel(jaLabel)
		case 1:
			label[i] = NewLabel(longLabel)
		case 2:
			label[i] = NewLabel("hello")
		}
		label[i].SetText(fmt.Sprintf("hello %d: %s", nr, label[i].Text()))
		box[i].Append(label[i])
		box[i].Append(entry[i])
		grid := NewGrid()

		for j := 0; j < ndrops; j++ {
			d := NewDropdown()
			for _, drop := range drops {
				d.Append(drop)
			}
			grid.AppendWithLabel(fmt.Sprintf("label %d", j), d)
		}
		box[i].Append(grid)

		panes[i].SetChild(box[i])
		if i == 4 {
			panes[i].SetMenuBar(bar)
		}
		stack.Append(panes[i])

	}

	w.SetChild(stack)

	w.OnClosing(func(wi *Window) {
		fmt.Sprintf("Closing window: %v", wi)
		Exit(1)
	})
	Main(w)
}

func main() {
	paneMain()
}

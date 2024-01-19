package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

func mainMultiple() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewVerticalBox()

	hbox := NewHorizontalBox()
	hslab := NewSlab()
	vbox := NewVerticalBox()
	vslab := NewSlab()
	vbox.Append(vslab)
	hbox.Append(hslab)

	hbox2 := NewHorizontalBox()
	vbox2 := NewVerticalBox()

	box.Append(vbox)
	box.Append(hbox)
	box.Append(vbox2)
	box.Append(hbox2)

	nbut := 10
	entry := make([]*Entry, nbut)

	label := make([]*Label, nbut)
	for i := 0; i < nbut; i++ {
		nr := i + 1
		label[i] = NewLabel(fmt.Sprintf("Label 番号 %d", nr))

		switch i % 3 {
		case 0:
			entry[i] = NewEntry()
		case 1:
			entry[i] = NewPasswordEntry()
		case 2:
			entry[i] = NewSearchEntry()
		}

		entry[i].OnChanged(func(c *Entry) {
			fmt.Printf("Entry %d changed", nr)
			txt := c.Text()
			fmt.Printf("Text: %s\n", txt)
		})
		if i < 5 {
			entry[i].SetText(fmt.Sprintf("hello %d", nr))
			vslab.Append(label[i], 10, i*20)
			vslab.Append(entry[i], 100, i*20)
		} else {
			j := i - 5
			hslab.Append(label[i], j*100, 10)
			hslab.Append(entry[i], j*100, 40)
		}
	}
	fmt.Printf("hslab children: %d\n", hslab.NumChildren())

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v", wi)
		Exit(0)

	})

	w.SetChild(box)
	Main(w)
}

func mainSingle() {
	Init()
	w := NewWindow("test slab", 640, 480, false)

	slab := NewSlab()
	w.SetChild(slab)

	nentry := 10
	nlabel := 10
	entry := make([]*Entry, nentry)
	label := make([]*Label, nlabel)
	for i := 0; i < nlabel; i++ {
		nr := i + 1
		label[i] = NewLabel(fmt.Sprintf("Label 番号 %d", nr))
		switch i % 3 {
		case 0:
			entry[i] = NewEntry()
		case 1:
			entry[i] = NewPasswordEntry()
		case 2:
			entry[i] = NewSearchEntry()
		}

		entry[i].OnChanged(func(c *Entry) {
			fmt.Printf("Entry %d changed: ", nr)
			txt := c.Text()
			fmt.Printf("Text: %s\n", txt)
		})
		gy := i / 4
		gx := i % 4

		entry[i].SetText(fmt.Sprintf("hello %d", nr))
		slab.Append(label[i], gx*20, gy*20)
		slab.Append(entry[i], gx*20+100, gy*20)
	}

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v", wi)
		Exit(0)

	})
	Main(w)
}

func main() {
	// mainSingle()
	mainMultiple()
}

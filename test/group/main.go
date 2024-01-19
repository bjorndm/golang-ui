package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

func mainGroupSimple() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewBox()
	group := NewGroup("This is a group")
	group.SetChild(box)
	w.SetChild(group)

	vbox := NewBox()
	box.Append(vbox)

	nentry := 10
	entry := make([]*Entry, nentry)
	for i := 0; i < nentry; i++ {
		nr := i + 1
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
		entry[i].SetText(fmt.Sprintf("hello %d", nr))
		vbox.Append(entry[i])
	}

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v", wi)
		Exit(0)

	})
	Main(w)
}

func mainGroup() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewVerticalBox()

	hbox := NewTray()
	hgroup := NewGroup("Tray")
	vbox := NewBox()
	vgroup := NewGroup("Box")
	hgroup.SetChild(hbox)
	vgroup.SetChild(vbox)

	box.Append(hgroup)
	box.Append(vgroup)

	nentry := 10
	entry := make([]*Entry, nentry)
	for i := 0; i < nentry; i++ {
		nr := i + 1
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
			vbox.Append(entry[i])
		} else {
			// XXX: bug: the first entrybox in the vbox takes up all space
			// making the others hard to see.
			hbox.Append(entry[i])
		}
	}

	w.SetChild(box)

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v", wi)
		Exit(0)

	})
	Main(w)
}

func main() {
	// mainGroupSimple()
	mainGroup()
}

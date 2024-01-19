package main

import "fmt"
import . "github.com/bjorndm/golang-ui"
import "github.com/bjorndm/golang-ui/icon"

const longLabel = `Cum multae res in philosophia nequaquam
satis adhuc explicatae sint, tum perdifficilis,
Brute, quod tu minime ignoras,
et perobscura quaestio est de natura deorum,
quae et ad cognitionem animi pulcherrima est
et ad moderandam religionem necessaria.`

func mainMultiple() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewVerticalBox()

	hgrid := NewGrid()
	vgrid := NewGrid()

	form := NewGrid()
	form.AppendWithoutLabel(NewPictureWithIcon("", icon.Basket))

	selections := []string{"one", "two", "三", "four", "five"}
	nbut := 15
	entry := make([]*Entry, nbut)
	drop := make([]*Dropdown, nbut)

	label := make([]*Label, nbut)
	for i := 0; i < nbut; i++ {
		nr := i + 1
		label[i] = NewLabel(fmt.Sprintf("Label 番号 %d", nr))
		drop[i] = NewDropdown()
		for _, t := range selections {
			drop[i].Append(t)
		}
		switch i % 3 {
		case 0:
			entry[i] = NewEntry()
		case 1:
			entry[i] = NewPasswordEntry()
		case 2:
			entry[i] = NewSearchEntry()
		}
		align := StyleAlignJustify
		switch i % 4 {
		case 0:
			align = StyleAlignStart
		case 1:
			align = StyleAlignCenter
		case 2:
			align = StyleAlignEnd
		case 3:
			align = StyleAlignJustify
		}

		entry[i].OnChanged(func(c *Entry) {
			fmt.Printf("Entry %d changed", nr)
			txt := c.Text()
			fmt.Printf("Text: %s\n", txt)
		})
		if i < 5 {
			entry[i].SetText(fmt.Sprintf("hello %d", nr))
			vgrid.Append(label[i], 0, i, 1, align)
			vgrid.Append(entry[i], 1, i, 1, align)
			vgrid.Append(drop[i], 2, i, 1, align)
		} else if i < 10 {
			hgrid.Append(label[i], i-5, 0, 1, align)
			hgrid.Append(entry[i], i-5, 1, 1, align)
			hgrid.Append(drop[i], i-5, 2, 1, align)
		} else {
			switch i % 3 {
			case 0:
				form.AppendWithLabel(fmt.Sprintf("label %d", nr), entry[i])
			case 1:
				form.AppendWithLabel(fmt.Sprintf("label %d\n%s", nr, longLabel), entry[i])
			default:
				form.AppendWithoutLabel(entry[i])
			}
		}
	}
	fmt.Printf("hgrid children: %d\n", hgrid.NumChildren())

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v", wi)
		Exit(0)

	})

	box.Append(vgrid)
	box.Append(hgrid)
	box.Append(form)

	w.SetChild(box)
	Main(w)
}

func mainSingle() {
	Init()
	w := NewWindow("test grid", 640, 480, false)

	grid := NewGrid()
	w.SetChild(grid)

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
		align := StyleAlignJustify
		switch i % 4 {
		case 0:
			align = StyleAlignStart
		case 1:
			align = StyleAlignCenter
		case 2:
			align = StyleAlignEnd
		case 3:
			align = StyleAlignJustify
		}

		entry[i].OnChanged(func(c *Entry) {
			fmt.Printf("Entry %d changed: ", nr)
			txt := c.Text()
			fmt.Printf("Text: %s\n", txt)
		})
		gy := i / 4
		gx := i % 4

		entry[i].SetText(fmt.Sprintf("hello %d", nr))
		grid.Append(label[i], gx*2, gy, 1, align)
		grid.Append(entry[i], gx*2+1, gy, 1, align)
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

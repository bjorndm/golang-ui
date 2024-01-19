package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

func mainScroller() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewVerticalBox()

	hbox := NewHorizontalBox()
	vbox := NewVerticalBox()

	nbut := 10
	slide := make([]*Scroller, nbut)
	roll := make([]*Roller, nbut)
	for i := 0; i < nbut; i++ {
		nr := i + 1
		switch i % 3 {
		case 0:
			slide[i] = NewScroller(0, 10)
			roll[i] = NewRoller(0, 10)
		case 1:
			slide[i] = NewScroller(0, 100)
			roll[i] = NewRoller(0, 100)
		case 2:
			slide[i] = NewScroller(10, 50)
			roll[i] = NewRoller(10, 50)
		}

		slide[i].OnChanged(func(c *Scroller) {
			fmt.Printf("Scroller %d changed\n", nr)
			val := c.Value()
			fmt.Printf("Value: %d\n", val)
		})
		roll[i].OnChanged(func(c *Roller) {
			fmt.Printf("Roller %d changed\n", nr)
			val := c.Value()
			fmt.Printf("Value: %d\n", val)
		})
		if i < 5 {
			slide[i].SetValue(i * 10)
			roll[i].SetValue(i * 10)
		}

		hbox.Append(slide[i])
		vbox.Append(roll[i])
	}

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v", wi)
		Exit(0)

	})
	box.Append(vbox)
	box.Append(hbox)
	w.SetChild(box)
	Main(w)
}

func main() {
	mainScroller()
}

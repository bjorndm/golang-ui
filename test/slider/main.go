package main

import "fmt"
import . "github.com/bjorndm/golang-ui"

func mainSlider() {
	Init()
	w := NewWindow("test window", 640, 480, false)

	box := NewVerticalBox()

	hbox := NewHorizontalBox()
	vbox := NewVerticalBox()

	box.Append(vbox)
	box.Append(hbox)

	nbut := 10
	slide := make([]*Slider, nbut)
	for i := 0; i < nbut; i++ {
		nr := i + 1
		switch i % 3 {
		case 0:
			slide[i] = NewSlider(0, 10)
		case 1:
			slide[i] = NewSlider(0, 100)
		case 2:
			slide[i] = NewSlider(10, 50)
		}
		slide[i].OnChanged(func(c *Slider) {
			fmt.Printf("Slider %d changed\n", nr)
			val := c.Value()
			fmt.Printf("Value: %d\n", val)
		})
		if i < 5 {
			slide[i].SetValue(i * 10)
			vbox.Append(slide[i])
		} else {
			hbox.Append(slide[i])
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
	mainSlider()
}

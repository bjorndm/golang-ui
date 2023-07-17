package unix

import "github.com/gotk3/gotk3/gtk"

const gtkxMargin = 7

func setMargined(c gtk.Container, margined bool) {
	if margined {
		c.SetBorderWidth(gtkxMargin)
	} else {
		c.SetBorderWidth(0)
	}
}

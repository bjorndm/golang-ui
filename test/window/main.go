package main

import "log"

import "github.com/bjorndm/golang-ui"

func main() {
	ui.Init()
	w := ui.NewWindow("test window", 800, 640, false)

	w.OnClosing(func(wi *ui.Window) {
		log.Printf("Closing window: %v", wi)
	})
	ui.Main(w)
}

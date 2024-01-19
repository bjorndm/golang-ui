package main

import "image"
import "embed"
import "io"
import "fmt"
import . "github.com/bjorndm/golang-ui"

//go:embed resource
var resource embed.FS

// panics on failure
func loadResourceReader(name string) io.ReadCloser {
	rd, err := resource.Open(name)
	if err != nil {
		panic(err)
	}
	return rd
}

func loadResourceImage(name string) image.Image {
	rd := loadResourceReader(name)
	defer rd.Close()
	img, _, err := image.Decode(rd)
	if err != nil {
		panic(err)
	}
	return img
}

func mainPictureSimple() {
	Init()
	w := NewWindow("test window", 640, 480, false)
	img := loadResourceImage("resource/cheese.png")

	box := NewTray()

	npic := 5
	pic := make([]*Picture, npic)
	for i := 0; i < npic; i++ {
		nr := i + 1
		if i%2 == 0 {
			pic[i] = NewPicture(fmt.Sprintf("This is a picture of cheese %d", nr), img)
		} else {
			pic[i] = NewPicture("", img)
		}
		box.Append(pic[i])
	}

	w.OnClosing(func(wi *Window) {
		fmt.Printf("Closing window: %v", wi)
		Exit(0)

	})

	w.SetChild(box)
	Main(w)
}

func main() {
	mainPictureSimple()
}

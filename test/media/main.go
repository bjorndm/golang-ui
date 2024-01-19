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

func mainMediaPlayerSimple() {
	Init()
	w := NewWindow("test window", 640, 480, false)
	padded := false
	media, err := NewMediaFromFileSystem(resource, "resource/chimp.mpg")
	if err != nil {
		panic(err)
	}
	width, height := media.VideoSize()
	fmt.Printf("Media loaded: size: %d %d\n", width, height)
	media2, err := NewMediaFromFileSystem(resource, "resource/bunny.mpg")
	if err != nil {
		panic(err)
	}
	width2, height2 := media2.VideoSize()
	fmt.Printf("Media loaded: size: %d %d\n", width2, height2)

	box := NewTray()
	box.SetPadded(padded)

	npic := 2 // only works with 1 player for now.
	pic := make([]*MediaPlayer, npic)
	for i := 0; i < npic; i++ {
		nr := i + 1
		if i%2 == 0 {
			pic[i] = NewMediaPlayer(fmt.Sprintf("This is a video of a chimp %d", nr), media)
		} else {
			pic[i] = NewMediaPlayer(fmt.Sprintf("This is a video of a bunny %d", nr), media2)
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
	mainMediaPlayerSimple()
}

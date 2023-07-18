package unix

import "testing"
import "github.com/matryer/is"

func TestNewWindow(t *testing.T) {
	is := is.New(t)
	Init()
	w := NewWindow("test window", 640, 480, false)
	is.True(w != nil)
	box := NewVerticalBox()
	is.True(box != nil)

	b := NewButton("Hello!")
	is.True(b != nil)
	is.Equal(b.Text(), "Hello!")

	b.OnClicked(func(b *Button, d any) {
		t.Logf("Hello clicked")
	}, nil)

	box.Append(b, false)
	t.Logf("before setchild")
	w.SetChild(box)
	t.Logf("after setchild")

	Main()
}

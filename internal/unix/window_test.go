package unix

import "testing"
import "github.com/matryer/is"

func TestNewWindow(t *testing.T) {
	is := is.New(t)
	Init()
	w := NewWindow("test window", 640, 480, false)
	is.True(w != nil)
	Main()
}

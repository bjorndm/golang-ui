package unix

import "github.com/gotk3/gotk3/gtk"

type Button struct {
	BasicWidget
	button        *gtk.Button
	onClicked     func(*Button, any)
	onClickedData any
}

func (b *Button) whenClicked() func(button *gtk.Button) {
	return func(button *gtk.Button) {
		b2 := b
		if b2.onClicked != nil {
			b2.onClicked(b2, b2.onClickedData)
		}
	}
}

func (b *Button) Text() string {
	t, _ := b.button.GetLabel()
	return t
}

func (b *Button) SetText(text string) {
	b.button.SetLabel(text)
}

func (b *Button) OnClicked(f func(*Button, any), data any) {
	b.onClicked = f
	b.onClickedData = data
}

func NewButton(text string) *Button {
	b := &Button{}
	b.c = b
	b.button, _ = gtk.ButtonNewWithLabel(text)
	b.widget = &b.button.Widget
	setWidgetBackgroundColor(b.widget, "green")

	b.widget.Connect("clicked", b.whenClicked)
	b.button.Show()
	return b
}

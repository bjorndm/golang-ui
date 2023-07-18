package unix

import "github.com/gotk3/gotk3/gtk"

type boxChild struct {
	BasicWidget
	stretchy   bool
	oldhexpand bool
	oldhalign  gtk.Align
	oldvexpand bool
	oldvalign  gtk.Align
}

type BasicWidget struct {
	widget *gtk.Widget
	parent Control
	c      Control
}

func (w BasicWidget) SetControl(c Control) {
	w.c = c
}

func (w BasicWidget) Handle() any {
	return w.widget
}

func (w BasicWidget) Parent() Control {
	return w.parent
}

func (w *BasicWidget) SetParent(parent Control) {
	w.parent = parent
}

func (w *BasicWidget) Toplevel() bool {
	return false
}

func (w *BasicWidget) Show() {
	w.widget.Show()
}

func (w *BasicWidget) Hide() {
	w.widget.Hide()
}

func (w *BasicWidget) Enabled() bool {
	return w.widget.GetCanFocus()
}

func (w *BasicWidget) Enable() {
	w.widget.SetCanFocus(true)
}

func (w *BasicWidget) Disable() {
	w.widget.SetCanFocus(false)
}

func (w *BasicWidget) Destroy() {
	if w.widget != nil {
		w.widget.Unref()
	}
	w.widget = nil
}

type Box struct {
	BasicWidget
	controls      []Control
	container     *gtk.Container
	box           *gtk.Box
	vertical      bool
	padded        bool
	stretchygroup *gtk.SizeGroup // ensures all stretchy controls have the same size
}

func (b *Box) Destroy() {
	// kill the size group
	b.stretchygroup.Unref()
	// free all controls
	for i := 0; i < len(b.controls); i++ {
		bc := b.controls[i]
		bc.SetParent(nil)
		// and make sure the widget itself stays alive
		// XXX uiUnixControlSetContainer(uiUnixControl(bc.c), b.container, TRUE);
		bc.Destroy()
	}
	// and then ourselves
	b.widget.Unref()
}

func (b *Box) Append(c Control, stretchy bool) {
	bc := &boxChild{}
	var widget *gtk.Widget

	bc.c = c
	bc.stretchy = stretchy
	widget = bc.c.Handle().(*gtk.Widget)
	bc.oldhexpand = widget.GetHExpand()
	bc.oldhalign = widget.GetHAlign()
	bc.oldvexpand = widget.GetVExpand()
	bc.oldvalign = widget.GetHAlign()

	if bc.stretchy {
		if b.vertical {
			widget.SetVExpand(true)
			widget.SetVAlign(gtk.ALIGN_FILL)
		} else {
			widget.SetHExpand(true)
			widget.SetHAlign(gtk.ALIGN_FILL)
		}
		b.stretchygroup.AddWidget(widget)
	} else {
		if b.vertical {
			widget.SetVExpand(false)
		} else {
			widget.SetHExpand(false)
		}
	}
	// and make them fill the opposite direction
	if b.vertical {
		widget.SetHExpand(true)
		widget.SetHAlign(gtk.ALIGN_FILL)
	} else {
		widget.SetVExpand(true)
		widget.SetVAlign(gtk.ALIGN_FILL)
	}

	bc.c.SetParent(b)
	b.controls = append(b.controls, bc)
}

func (b *Box) Delete(index int) {
	var bc *boxChild
	var widget *gtk.Widget

	bc, _ = b.controls[index].(*boxChild)
	widget, _ = bc.c.Handle().(*gtk.Widget)

	bc.c.SetParent(nil)

	if bc.stretchy {
		b.stretchygroup.RemoveWidget(widget)
	}

	widget.SetHExpand(bc.oldhexpand)
	widget.SetHAlign(bc.oldhalign)
	widget.SetVExpand(bc.oldvexpand)
	widget.SetHAlign(bc.oldvalign)
	b.controls = append(b.controls[:index-1], b.controls[:index+1]...)
}

func (b Box) NumChildren() int {
	return len(b.controls)
}

func (b Box) Padded() bool {
	return b.padded
}

const gtkYPadding = 5
const gtkXPadding = 5

func (b *Box) SetPadded(padded bool) {
	b.padded = padded
	if b.padded {
		if b.vertical {
			b.box.SetSpacing(gtkYPadding)
		} else {
			b.box.SetSpacing(gtkXPadding)
		}
	} else {
		b.box.SetSpacing(0)
	}
}

func newBox(orientation gtk.Orientation) *Box {
	b := &Box{}
	b.c = b
	b.box, _ = gtk.BoxNew(orientation, 0)
	b.widget = &b.box.Widget

	b.vertical = orientation == gtk.ORIENTATION_VERTICAL

	if b.vertical {
		b.stretchygroup, _ = gtk.SizeGroupNew(gtk.SIZE_GROUP_VERTICAL)
	} else {
		b.stretchygroup, _ = gtk.SizeGroupNew(gtk.SIZE_GROUP_HORIZONTAL)
	}

	b.controls = []Control{}

	return b
}

func NewHorizontalBox() *Box {
	return newBox(gtk.ORIENTATION_HORIZONTAL)
}

func NewVerticalBox() *Box {
	return newBox(gtk.ORIENTATION_VERTICAL)
}

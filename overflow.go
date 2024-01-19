package ui

// Scrollable is an interface for the use of Overflow. If the target widget
// implements the Scrollable interface, then Overflow will use this for
// communicatiion with the control/widget.
//
// TODO: actually use this.
type Scrollable interface {
	ScrollWidget(y int)              // ScrollWidget will be called by Overflow when the widget is scrolled up or down to the given Y position.
	RollWidget(x int)                // RollWidget will be called by Overflow when the widget is rolled left or right to the given X position.
	ScrollSize() (width, height int) // ScrollSize sould return the desired scroll size of the widget. These dimensions are used for ScollWidget() and RollWidget() as the maximum X or Y.
}

// Overflow is a widget meant to be used as a wrapper around other
// widgets, or to be added as a member of other target widgets to help
// them manage overflow by displaying a Scroller or a Roller if needed,
// and by providing a sub-bitmap of the screen that will clip the output
// of the target widget.
//
// If the target implements the Scrollable interface, then
// Overflow will use this for communicatiion with the control/widget.
type Overflow struct {
	BasicWidget
	target    Control
	roller    *Roller
	scroller  *Scroller
	clip      *Graphic
	maxWidth  int
	maxHeight int
}

func NewOverflow(target Control, maxWidth, maxHeight int) *Overflow {
	o := &Overflow{}
	o.maxWidth = maxWidth
	o.maxHeight = maxHeight
	o.width = o.maxWidth
	o.height = o.maxHeight
	o.scroller = NewScroller(0, 100)
	o.roller = NewRoller(0, 100)
	o.target = target
	o.scroller.SetParent(o)
	o.roller.SetParent(o)
	o.target.SetParent(o)
	o.clip = NewGraphic(o.maxWidth, o.maxHeight)

	o.scroller.OnChanged(func(s *Scroller) {
		if o.target != nil {
			y := s.Value()
			if scrollable, ok := o.target.(Scrollable); ok {
				scrollable.ScrollWidget(y)
			} else {
				x, _ := o.target.WidgetAt()
				o.target.MoveWidget(x, -y)
			}
		}
	})

	o.roller.OnChanged(func(s *Roller) {
		if o.target != nil {
			x := s.Value()
			if scrollable, ok := o.target.(Scrollable); ok {
				scrollable.RollWidget(x)
			} else {
				_, y := o.target.WidgetAt()
				o.target.MoveWidget(-x, y)
			}
		}
	})

	return o
}

func (o Overflow) WidgetSize() (width, height int) {
	return o.maxWidth, o.maxHeight
}

func (o Overflow) LayoutWidget(parentW, parentH int) {
	o.width = o.maxWidth
	o.height = o.maxHeight

	if o.target != nil {
		o.target.LayoutWidget(parentW, parentH)
		targetW, targetH := o.target.WidgetSize()
		if scrollable, ok := o.target.(Scrollable); ok {
			targetW, targetH = scrollable.ScrollSize()
		}
		o.scroller.SetRange(0, targetH-o.height)
		o.roller.SetRange(0, targetW-o.width)
	}

	o.scroller.LayoutWidget(o.width, o.height)
	w, _ := o.scroller.WidgetSize()
	o.scroller.MoveWidget(o.width-w, 0)

	_, h := o.roller.WidgetSize()
	o.roller.LayoutWidget(o.width-w, o.height)
	o.roller.MoveWidget(0, o.height-h)
}

func (o Overflow) DrawWidget(screen *Graphic) {
	dx, dy := o.WidgetAbsolute()
	// tw, th := o.target.WidgetSize()

	if o.target != nil {
		wx, wy := o.target.WidgetAt()
		// Have to temporarily move the widget to origin 0,0
		// since DrawWidget adds the parent offset, but the clip is a blank
		// at offset 0, 0 and is positioned afterwards.
		o.target.MoveWidget(wx-dx, wy-dy)
		o.target.DrawWidget(o.clip)
		o.target.MoveWidget(wx, wy)
		DrawGraphicAt(screen, o.clip, dx, dy)
	}

	// if tw > o.maxWidth {
	o.roller.DrawWidget(screen)
	// }
	// if th > o.maxHeight {
	o.scroller.DrawWidget(screen)
	//}
	o.DrawDebug(screen, "OVE")
}

func (o *Overflow) HandleWidget(ev Event) {
	HandleWidgetFocus(ev, o, o.scroller, o.roller, o.target)

	if o.Focus() != nil {
		o.Focus().HandleWidget(ev)
	}
}

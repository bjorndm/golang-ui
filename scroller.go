package ui

type Scroller struct {
	BasicWidget
	RangeValue // Range and value of the Scroller
	active     bool
	onChanged  func(*Scroller)
}

func (s *Scroller) OnChanged(f func(*Scroller)) {
	s.onChanged = f
}

func newScroller(min, max int) *Scroller {
	b := &Scroller{}
	if min >= max {
		panic("NewScroller: min must be strictly smaller than max")
	}
	b.min = min
	b.value = min
	b.max = max

	b.SetStyle(theme.Scroller)

	return b
}

const scrollerWidth = 12
const scrollerHeight = 50

func (s *Scroller) LayoutWidget(width, height int) {
	margin := s.Style().Margin.Int()
	s.width, s.height = s.Style().Size.Width.Int(), s.Style().Size.Height.Int()
	if height-margin > s.height {
		s.height = height - margin
	}
	s.ClipTo(width, height)
}

func (s Scroller) DrawWidget(dst *Graphic) {
	dx, dy := s.WidgetAbsolute()

	FillFrameStyle(dst, dx, dy, s.width, s.height, s.Style())
	markY := (s.value - s.min) * s.height / (s.max - s.min)
	r := s.width / 2
	FillCircleStyle(dst, dx+s.width/2, dy+markY, r, s.Style())
	DrawLineStyle(dst, dx, dy+markY, s.width, 0, s.Style())

}

func NewScroller(min, max int) *Scroller {
	return newScroller(min, max)
}

func (s *Scroller) HandleKeyPress(kp *KeyPressEvent) {
	switch kp.Key {
	case KeyArrowLeft, KeyArrowUp:
		s.SetValue(s.value - 1)
		if s.onChanged != nil {
			s.onChanged(s)
		}

	case KeyArrowRight, KeyArrowDown:
		s.SetValue(s.value + 1)
		if s.onChanged != nil {
			s.onChanged(s)
		}

	case KeyPageUp:
		s.SetValue(s.value - (s.max-s.min)/10)
		if s.onChanged != nil {
			s.onChanged(s)
		}

	case KeyPageDown:
		s.SetValue(s.value + (s.max-s.min)/10)
		if s.onChanged != nil {
			s.onChanged(s)
		}

	case KeyHome:
		s.SetValue(s.min)
		if s.onChanged != nil {
			s.onChanged(s)
		}

	case KeyEnd:
		s.SetValue(s.max)
		if s.onChanged != nil {
			s.onChanged(s)
		}
	}
}

func (d *Scroller) HandleKeyRelease(ke *KeyReleaseEvent) {

}

func (s *Scroller) HandleMouseClick(mc *MouseClickEvent) {
	dprintln("Scroller.HandleMouseClick ", mc.X, mc.Y)

	if !s.active {
		s.active = true
	}

	if s.active {
		_, dy := s.WidgetAbsolute()
		delta := mc.Y - dy + 1
		value := (s.max - s.min) * delta / s.height
		value += s.min
		s.SetValue(value)
		if s.onChanged != nil {
			s.onChanged(s)
		}
	}
}

func (s *Scroller) HandleWidget(ev Event) {
	if _, ok := ev.(*AwayEvent); ok {
		// Deactivate ourselves because another widget should get
		// the focus in stead.
		s.active = false
		return
	}

	if wheel, ok := ev.(*WheelEvent); s.active && ok {
		if debugDisplay {
			dprintln("Scroller.HandleWidget: ", ev)
		}
		if wheel.WheelY > 0 {
			s.SetValue(s.value - 1)
			if s.onChanged != nil {
				s.onChanged(s)
			}
		} else if wheel.WheelY < 0 {
			s.SetValue(s.value + 1)
			if s.onChanged != nil {
				s.onChanged(s)
			}
		}
	}

	if mc, ok := ev.(*MouseClickEvent); ok {
		s.HandleMouseClick(mc)
	}
	if s.active {
		// allow keyboard selection
		if ke, ok := ev.(*KeyPressEvent); ok {
			s.HandleKeyPress(ke)
		}
		if ke, ok := ev.(*KeyReleaseEvent); ok {
			s.HandleKeyRelease(ke)
		}
	}
}

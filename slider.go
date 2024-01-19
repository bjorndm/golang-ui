package ui

// RangeValue is a value that is between a minimum and a maximum
type RangeValue struct {
	value int
	min   int
	max   int
}

func MakeRangeValue(min, max int) RangeValue {
	r := RangeValue{}
	r.SetRange(min, max)
	return r
}

func (r RangeValue) Range() (min, max int) {
	return r.min, r.max
}

func (r *RangeValue) SetRange(min, max int) {
	if min >= max {
		min, max = max, min
	}
	r.min = min
	r.max = max
	if r.value < r.min {
		r.value = r.min
	}
	if r.value > r.max {
		r.value = r.max
	}
}

func (r RangeValue) Value() int {
	return r.value
}

func (r *RangeValue) SetValue(value int) {
	if value < r.min {
		value = r.min
	}
	if value > r.max {
		value = r.max
	}

	r.value = value
}

type Slider struct {
	BasicWidget
	RangeValue // Range and value of the Slider
	text       TextWidget
	active     bool
	tickmarks  int
	onChanged  func(*Slider)
}

func (s *Slider) SetTitle(title string) {
	s.text.text = title
}

func (s *Slider) Title() string {
	return s.text.text
}

func (s *Slider) OnChanged(f func(*Slider)) {
	s.onChanged = f
}

func newSlider(min, max int) *Slider {
	b := &Slider{}
	if min >= max {
		panic("NewSlider: min must be strictly smaller than max")
	}
	b.min = min
	b.value = min
	b.max = max

	delta := max - min
	b.tickmarks = 0
	if delta < sliderDiscreteDelta {
		b.tickmarks = delta + 1
	}
	b.text.SetParent(b) // The label is a sub-widget.
	b.SetStyle(theme.Slider)

	return b
}

const sliderHeight = 12
const sliderWidth = 100

func (s *Slider) LayoutWidget(parentWidth, parentHeight int) {

	s.text.LayoutWidget(parentWidth, parentHeight)
	width, height := s.text.WidgetSize()
	height += sliderHeight

	if width < s.Style().Size.Width.Int() {
		width = s.Style().Size.Width.Int()
	}

	if height < s.Style().Size.Height.Int() {
		height = s.Style().Size.Height.Int()
	}

	s.width, s.height = width, height
	s.ClipTo(parentWidth, parentHeight)
}

func (s Slider) DrawWidget(dst *Graphic) {
	dx, dy := s.WidgetAbsolute()
	if s.text.text != "" {
		s.text.DrawWidget(dst)
	}
	DrawLineStyle(dst, dx, dy+s.height/2, s.width, 0, s.Style())
	if s.tickmarks > 0 {
		for i := 0; i < s.tickmarks; i++ {
			markX := i * s.width / (s.max - s.min)
			DrawLineStyle(dst, dx+markX, dy+s.height/4, 0, s.height/2, s.Style())
		}
	}
	markX := (s.value - s.min) * s.width / (s.max - s.min)
	r := s.height / 4
	FillCircleStyle(dst, dx+markX, dy+s.height/2, r, s.Style())
	DrawLineStyle(dst, dx+markX, dy, 0, s.height, s.Style())

}

const sliderDiscreteDelta = 21

func NewSlider(min, max int) *Slider {
	return newSlider(min, max)
}

func (s *Slider) HandleKeyPress(kp *KeyPressEvent) {
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

func (d *Slider) HandleKeyRelease(ke *KeyReleaseEvent) {

}

func (s *Slider) HandleMouseClick(mc *MouseClickEvent) {
	if !s.active {
		s.active = true
	}

	if s.active {
		dx, _ := s.WidgetAbsolute()
		delta := mc.X - dx + 1
		value := (s.max - s.min) * delta / s.width
		value += s.min
		s.SetValue(value)
		if s.onChanged != nil {
			s.onChanged(s)
		}
	}
}

func (s *Slider) HandleWidget(ev Event) {
	if _, ok := ev.(*AwayEvent); ok {
		// Deactivate ourselves because another widget should get
		// the focus in stead.
		s.active = false
		return
	}

	if wheel, ok := ev.(*WheelEvent); s.active && ok {
		if debugDisplay {
			dprintln("Slider.HandleWidget: ", ev)
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

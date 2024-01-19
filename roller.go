package ui

// Roller is a horizontal scroller or horizontal scroll bar.
type Roller struct {
	BasicWidget
	RangeValue // Range and value of the Roller
	active     bool
	onChanged  func(*Roller)
}

func (r *Roller) OnChanged(f func(*Roller)) {
	r.onChanged = f
}

func newRoller(min, max int) *Roller {
	b := &Roller{}
	if min >= max {
		panic("NewRoller: min must be strictly smaller than max")
	}
	b.min = min
	b.value = min
	b.max = max

	b.SetStyle(theme.Roller)

	return b
}

const rollerWidth = 50
const rollerHeight = 12

func (r *Roller) LayoutWidget(width, height int) {
	margin := r.Style().Margin.Int()
	r.width, r.height = r.Style().Size.Width.Int(), r.Style().Size.Height.Int()
	if width-margin > r.width {
		r.width = width - margin
	}
	r.ClipTo(width, height)
}

func (r Roller) DrawWidget(dst *Graphic) {
	dx, dy := r.WidgetAbsolute()

	FillFrameStyle(dst, dx, dy, r.width, r.height, r.Style())
	markX := (r.value - r.min) * r.width / (r.max - r.min)
	radius := r.height / 2
	FillCircleStyle(dst, dx+markX, dy+radius, radius, r.Style())
	DrawLineStyle(dst, dx+markX, dy, 0, r.height, r.Style())
}

func NewRoller(min, max int) *Roller {
	return newRoller(min, max)
}

func (r *Roller) HandleKeyPress(kp *KeyPressEvent) {
	switch kp.Key {
	case KeyArrowLeft, KeyArrowUp:
		r.SetValue(r.value - 1)
		if r.onChanged != nil {
			r.onChanged(r)
		}

	case KeyArrowRight, KeyArrowDown:
		r.SetValue(r.value + 1)
		if r.onChanged != nil {
			r.onChanged(r)
		}

	case KeyPageUp:
		r.SetValue(r.value - (r.max-r.min)/10)
		if r.onChanged != nil {
			r.onChanged(r)
		}

	case KeyPageDown:
		r.SetValue(r.value + (r.max-r.min)/10)
		if r.onChanged != nil {
			r.onChanged(r)
		}

	case KeyHome:
		r.SetValue(r.min)
		if r.onChanged != nil {
			r.onChanged(r)
		}

	case KeyEnd:
		r.SetValue(r.max)
		if r.onChanged != nil {
			r.onChanged(r)
		}
	}
}

func (d *Roller) HandleKeyRelease(ke *KeyReleaseEvent) {

}

func (r *Roller) HandleMouseClick(mc *MouseClickEvent) {
	dprintln("Roller.HandleMouseClick ", mc.X, mc.Y)
	if !r.active {
		r.active = true
	}

	if r.active {
		dx, _ := r.WidgetAbsolute()
		delta := mc.X - dx + 1
		value := (r.max - r.min) * delta / r.width
		value += r.min
		r.SetValue(value)
		if r.onChanged != nil {
			r.onChanged(r)
		}
	}
}

func (r *Roller) HandleWidget(ev Event) {
	if _, ok := ev.(*AwayEvent); ok {
		// Deactivate ourselves because another widget should get
		// the focus in stead.
		r.active = false
		return
	}

	if wheel, ok := ev.(*WheelEvent); r.active && ok {
		if debugDisplay {
			dprintln("Roller.HandleWidget: ", ev)
		}
		if wheel.WheelY > 0 || wheel.WheelX > 0 {
			r.SetValue(r.value - 1)
			if r.onChanged != nil {
				r.onChanged(r)
			}
		} else if wheel.WheelY < 0 || wheel.WheelX < 0 {
			r.SetValue(r.value + 1)
			if r.onChanged != nil {
				r.onChanged(r)
			}
		}
	}

	if mc, ok := ev.(*MouseClickEvent); ok {
		r.HandleMouseClick(mc)
	}
	if r.active {
		// allow keyboard selection
		if ke, ok := ev.(*KeyPressEvent); ok {
			r.HandleKeyPress(ke)
		}
		if ke, ok := ev.(*KeyReleaseEvent); ok {
			r.HandleKeyRelease(ke)
		}
	}
}

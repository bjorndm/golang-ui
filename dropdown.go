package ui

type Dropdown struct {
	TextWidget
	onSelected func(*Dropdown)
	committed  int
	active     bool
	overlay    dropdownOverlay
}

// widget for the overlay of the dropdown
type dropdownOverlay struct {
	BasicContainer // contains TextWidgets for each option of the drop down.
	selected       int
	dropdown       *Dropdown
}

func (d dropdownOverlay) HandleWidget(ev Event) {
	// XXX: This is a hack, we should handle the events properly
	// in  stead of kicking them back to the dropdown...
	d.dropdown.HandleWidget(ev)
}

func (d *dropdownOverlay) LayoutWidget(width, height int) {
	dy := 0
	dx := 0
	margin := d.Style().Margin.Int()
	d.height = d.Style().Font.Face.Metrics().Height.Round()
	d.width = width

	for _, control := range d.controls {
		control.LayoutWidget(d.width, d.height)
		ww, wh := control.WidgetSize()
		if ww > d.width {
			d.width = ww
		}
		wh += margin * 2
		dy += wh
		control.MoveWidget(dx, dy)
	}
	d.height += dy
	d.GrowToStyleSize()
	d.width += margin * 2
	d.height += margin * 2
}

func (d *dropdownOverlay) DrawWidget(dst *Graphic) {
	dx, dy := d.WidgetAbsolute()

	// draw list of text widgets if focused
	var (
		margin = d.Style().Margin.Int()
	)

	dy += margin
	ah := 0
	ay := 0
	for i, control := range d.controls {
		_, wh := control.WidgetSize()
		wh += margin
		ah += wh
		if i == 0 {
			ay = wh
		}
	}
	// Draw background frame.
	FillFrameStyle(dst, dx, dy+ay, d.width, ah, d.Style())
	dys := dy

	// Draw selections
	for i, control := range d.controls {
		_, wh := control.WidgetSize()
		wh += margin
		dys += wh
		if i == d.selected {
			FillFrameOptionalStyle(dst, dx, dys, d.width, wh, theme.Active)
		}
		control.DrawWidget(dst)
	}
}

func (b *Dropdown) OnSelected(f func(*Dropdown)) {
	b.onSelected = f
}

const dropdownDefaultText = "--------"
const dropdownWidth = 64
const dropdownHeight = 16

func NewDropdown() *Dropdown {
	d := &Dropdown{}
	d.customStyle = theme.Dropdown
	d.overlay.customStyle = theme.Dropdown
	d.overlay.SetParent(d)
	d.overlay.dropdown = d

	return d
}

func (c *Dropdown) Append(text string) {
	widget := NewTextWidget(text)
	c.overlay.AppendWithParent(widget, &c.overlay)
}

func (c *Dropdown) Clear() {
	c.overlay.Clear()
}

func (c *Dropdown) Delete(column int) {
	c.overlay.Delete(column)
}

func (c *Dropdown) InsertAt(text string, column int) {
	widget := NewTextWidget(text)
	style := widget.Style()
	style.Layout = StyleLayoutStretch
	c.overlay.InsertAt(widget, column)
}

func (c *Dropdown) NumItems() int {
	return int(c.overlay.NumItems())
}

func (c *Dropdown) Selected() int {
	return c.overlay.selected
}

func (c *Dropdown) SetSelected(index int) {
	c.overlay.selected = index
	c.committed = c.overlay.selected
}

func (c *Dropdown) textWidget(index int) *TextWidget {
	if widget, ok := c.overlay.Get(index).(*TextWidget); ok && widget != nil {
		return widget
	}
	return nil
}

func (c *Dropdown) Text() string {
	if c.overlay.selected < 0 {
		return ""
	}
	if widget := c.textWidget(c.overlay.selected); widget != nil {
		return widget.Text()
	}
	return ""
}

func (c *Dropdown) SetText(text string) {
	for i, child := range c.overlay.controls {
		if widget, ok := child.(*TextWidget); ok && widget != nil {
			if widget.Text() == text {
				c.overlay.selected = i
			}
		}
	}
}

func (d *Dropdown) LayoutWidget(width, height int) {
	txt := d.Text()
	if txt == "" {
		txt = dropdownDefaultText
	}
	d.width, d.height = oneLineTextSize(d.Style().Font.Face, txt)
	minh := d.Style().Font.Face.Metrics().Height.Round()
	if d.height < minh {
		d.height = minh
	}

	margin := d.Style().Margin.Int()
	d.GrowToStyleSize()
	d.width += margin * 2
	d.height += margin * 2

	// Also lay out the text widgets for the selection list overlay
	// below the widget so we don't need to do so when drawing them.
	d.overlay.LayoutWidget(d.width, height)
	d.ClipTo(width, height)
}

func (d *Dropdown) DrawWidget(dst *Graphic) {
	dx, dy := d.WidgetAbsolute()

	var (
		margin = d.Style().Margin.Int()
	)

	if d.active {
		FillFrameStyle(dst, dx, dy, d.width, d.height, d.Style())
	} else {
		FillFrameStyle(dst, dx, dy, d.width, d.height, d.Style())
		// draw icon
		minh := d.Style().Font.Face.Metrics().Height.Round()
		icon := d.Style().Icon.String()
		ix := dx + d.width - dropdownHeight - margin*2
		iy := dy + d.height - dropdownHeight - margin*2
		iconAtlas.DrawSprite(dst, ix, iy, minh, minh, icon)
	}

	sub := GraphicClipStyle(dst, dx, dy, d.width, d.height, d.Style())
	text := d.Text()
	if text == "" {
		text = dropdownDefaultText
	}

	TextDrawOffsetStyle(sub, text, dx, dy, d.Style())
	if d.active {
		d.overlay.DrawWidget(dst)
	}
	d.DrawDebug(dst, "DRO")
}

func (d *Dropdown) setActive(active bool) {
	d.active = active
	if d.active {
		StartOverlayWidget(d.Parent(), &d.overlay)
	} else {
		EndOverlayWidget(d.Parent(), &d.overlay)
	}
}

func (d *Dropdown) HandleKeyPress(kp *KeyPressEvent) {
	switch kp.Key {
	case KeyArrowLeft, KeyArrowUp:
		d.SetSelected(d.overlay.selected - 1)
	case KeyArrowRight, KeyArrowDown:
		d.SetSelected(d.overlay.selected + 1)
	case KeyHome:
		d.SetSelected(0)
	case KeyEnd:
		d.SetSelected(len(d.overlay.controls) - 1)
	case KeyEscape:
		d.overlay.selected = d.committed
		d.setActive(false)
	case KeyEnter:
		d.committed = d.overlay.selected
		if d.onSelected != nil {
			d.onSelected(d)
		}
		d.setActive(false)
	}
}

func (d *Dropdown) HandleKeyRelease(ke *KeyReleaseEvent) {

}

func (d *Dropdown) HandleMouseClick(mc *MouseClickEvent) {
	dprintln("Dropdown.HandleMouseClick ", mc.X, mc.Y)
	if !d.active {
		d.setActive(true)
		// If not inside it is not for us, ignore it.
		return // only activate, so can't select yet.
	}

	if d.active {
		// check if inside selection widgets
		for i, control := range d.overlay.controls {
			tw := control.(*TextWidget)
			if mc.Inside(tw) {
				d.SetSelected(i)
				dprintln("Dropdown.HandleMouseClick selected ", i)
				d.setActive(false)
				break
			}
		}
	}
}

func (e *Dropdown) HandleWidget(ev Event) {
	// If we get an away event, deactivate
	if _, ok := ev.(*AwayEvent); ok {
		dprintln("Dropdown.HandleWidget deactivate")
		e.setActive(false)
		return
	}
	if mc, ok := ev.(*MouseClickEvent); ok {
		e.HandleMouseClick(mc)
	}
	if e.active {
		// allow keyboard selection
		if ke, ok := ev.(*KeyPressEvent); ok {
			e.HandleKeyPress(ke)
		}
		if ke, ok := ev.(*KeyReleaseEvent); ok {
			e.HandleKeyRelease(ke)
		}
	}
}

package ui

// A Slab is a container with fixed layout, no border and no background.
// The minimum size fits all widgets, which are given their minimal size and
// never repositioned after adding them. Slab is useful for fixed layouts,
// but should normally only be used for implementing other widgets.
type Slab struct {
	BasicContainer
	padded bool
}

func (b *Slab) Destroy() {
	// free all controls
	for i := 0; i < len(b.controls); i++ {
		bc := b.controls[i]
		bc.SetParent(nil)
		bc.Destroy()
	}
}

func (b *Slab) Enable() {
	for _, child := range b.controls {
		bw := child.(*BasicWidget)
		bw.Enable()
	}
}

func (b *Slab) Disable() {
	for _, child := range b.controls {
		bw := child.(*BasicWidget)
		bw.Disable()
	}
}

func (b *Slab) Append(c Control, x, y int) {
	b.BasicContainer.AppendWithParent(c, b)
	b.LayoutWidget(b.width, b.height)
	margin := b.Style().Margin.Int()
	c.MoveWidget(x+margin, y+margin)
}

func (b Slab) Padded() bool {
	return b.padded
}

func (b *Slab) SetPadded(padded bool) {
	b.padded = padded
}

func newSlab() *Slab {
	b := &Slab{}
	b.controls = []Control{}
	b.SetStyle(theme.Slab)

	return b
}

func NewSlab() *Slab {
	return newSlab()
}

func (b *Slab) LayoutWidget(width, height int) {
	// The widget size of a slab is:
	// for the width, the leftmost extend of any widget.
	// for the height, the downmost extend of any widget.

	w := 0
	h := 0
	availableWidth := width
	availableHeight := height

	for _, child := range b.controls {
		if child.Hidden() {
			continue
		}
		// widget can layout freely, so use maximum dimensions
		child.LayoutWidget(availableWidth, availableHeight)
		leastWidth, leastHeight := child.WidgetSize()
		cx, cy := child.WidgetAt()
		leastHeight += cy
		leastWidth += cx
		if leastHeight > h {
			h = leastHeight
		}
		if leastWidth > w {
			w = leastWidth
		}
	}
	margin := b.Style().Margin.Int()

	h += margin * 2
	w += margin * 2
	b.width, b.height = w, h
	b.ClipTo(width, height)
}

func (b Slab) DrawWidget(g *Graphic) {
	// dx, dy := b.WidgetAbsolute()
	// FillFrameStyle(g, dx, dy, b.width, b.height, b.Style())
	b.BasicContainer.DrawWidget(g)
	b.DrawDebug(g, "SLA")
}

var _ Control = &Slab{}

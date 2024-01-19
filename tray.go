package ui

// A Tray is a container with a horizontal widget layout.
type Tray struct {
	BasicContainer
	padded bool
}

func (b *Tray) Destroy() {
	// free all controls
	for i := 0; i < len(b.controls); i++ {
		bc := b.controls[i]
		bc.SetParent(nil)
		bc.Destroy()
	}
}

func (b *Tray) Enable() {
	for _, child := range b.controls {
		bw := child.(*BasicWidget)
		bw.Enable()
	}
}

func (b *Tray) Disable() {
	for _, child := range b.controls {
		bw := child.(*BasicWidget)
		bw.Disable()
	}
}

func (b *Tray) Append(c Control) {
	b.BasicContainer.AppendWithParent(c, b)
	if b.width > 0 && b.height > 0 {
		b.LayoutWidget(b.width, b.height)
	}
}

func (b Tray) Padded() bool {
	return b.padded
}

func (b *Tray) SetPadded(padded bool) {
	b.padded = padded
}

func newTray() *Tray {
	b := &Tray{}
	b.controls = []Control{}

	return b
}

func NewTray() *Tray {
	return newTray()
}

// LayoutWidget for a Tray places all widgets the one next to the other.
// Vertically, the child widgets are limited to the available height
// of the parent, but horizontally they are not constrained.
func (b *Tray) LayoutWidget(width, height int) {
	dprintln("Box.LayoutWidget", len(b.controls), width, height)

	margin := b.Style().Margin.Int()
	x := margin
	y := margin

	availableHeight := height - margin*2
	availableWidth := width - margin*2

	highest := 0
	b.width = margin * 2
	for _, child := range b.controls {
		if child.Hidden() {
			continue
		}
		// No limits on the width in a tray.
		child.LayoutWidget(availableWidth, availableHeight)
		child.MoveWidget(x, y)
		childWidth, childHeight := child.WidgetSize()
		x += childWidth
		b.width += childWidth
		if childHeight > highest {
			highest = childHeight
		}
	}

	// Take on the height of the heighest child/ plus margins.
	b.height = highest + margin*2

	// Finally clip to desired size.
	b.ClipTo(width, height)
	b.BasicContainer.UpdateOrdered()

	dprintln("Tray.LayoutWidget done", len(b.controls), b.width, b.height)
}

func (b Tray) DrawWidget(g *Graphic) {
	dx, dy := b.WidgetAbsolute()

	FillFrameStyle(g, dx, dy, b.width, b.height, b.Style())
	b.BasicContainer.DrawWidget(g)
	b.DrawDebug(g, "TRA")
}

var _ Control = &Tray{}

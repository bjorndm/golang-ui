package ui

// A Box is a container a vertical widget layout.
type Box struct {
	BasicContainer
}

func (b *Box) Destroy() {
	// free all controls
	for i := 0; i < len(b.controls); i++ {
		bc := b.controls[i]
		bc.SetParent(nil)
		bc.Destroy()
	}
}

func (b *Box) Enable() {
	for _, child := range b.controls {
		bw := child.(*BasicWidget)
		bw.Enable()
	}
}

func (b *Box) Disable() {
	for _, child := range b.controls {
		bw := child.(*BasicWidget)
		bw.Disable()
	}
}

func (b *Box) Append(c Control) {
	b.BasicContainer.AppendWithParent(c, b)
}

func (b *Box) Delete(index int) {
	dprintln("Box.delete")

	if index < 0 || index >= len(b.controls) {
		return
	}
	bc := b.controls[index]

	bc.SetParent(nil)
	b.controls = append(b.controls[:index], b.controls[index+1:]...)
	if b.width > 0 && b.height > 0 {
		b.LayoutWidget(b.width, b.height)
	}
}

func (b Box) NumChildren() int {
	return len(b.controls)
}

func newBox() *Box {
	b := &Box{}
	b.controls = []Control{}
	return b
}

func NewHorizontalBox() *Tray {
	return newTray()
}

func NewVerticalBox() *Box {
	return newBox()
}

func NewBox() *Box {
	return newBox()
}

// LayoutWidget for a Box places all widgets the one below the other.
// Horizontally, the child widgets are limited to the available width
// of the parent, but vertically they are not constrained.
func (b *Box) LayoutWidget(width, height int) {
	dprintln("Box.LayoutWidget", len(b.controls), width, height)

	margin := b.Style().Margin.Int()
	x := margin
	y := margin

	availableWidth := width - margin*2
	availableHeight := height - margin*2

	widest := 0
	b.height = margin * 2
	for _, child := range b.controls {
		if child.Hidden() {
			continue
		}
		// No limits on the height in a Box.
		child.LayoutWidget(availableWidth, availableHeight)
		child.MoveWidget(x, y)
		childWidth, childHeight := child.WidgetSize()
		y += childHeight
		b.height += childHeight
		if childWidth > widest {
			widest = childWidth
		}
	}

	// Take on the widest width of the child plus margins as our width.
	b.width = widest + margin*2

	// Finally clip to desired size.
	b.ClipTo(width, height)
	b.BasicContainer.UpdateOrdered()

	dprintln("Box.LayoutWidget done", len(b.controls), b.width, b.height)
}

func (b Box) DrawWidget(g *Graphic) {
	dx, dy := b.WidgetAbsolute()
	FillFrameStyle(g, dx, dy, b.width, b.height, b.Style())
	b.BasicContainer.DrawWidget(g)
	b.DrawDebug(g, "BOX")
}

var _ Control = &Box{}

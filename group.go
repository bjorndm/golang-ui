package ui

// Group is a group of widgets wrapped in a box with a label on top.
// Useful on largerUI screens to group related widgets together

type Group struct {
	BasicWidget
	child      Control
	title      TextWidget
	borderless bool
}

func NewGroup(title string) *Group {
	w := &Group{}
	w.customStyle = theme.Group
	w.title.SetText(title)
	w.title.SetParent(w)
	w.title.customStyle = theme.Group
	w.SetStyle(theme.Group)
	return w
}

func (w *Group) SetTitle(title string) {
	w.title.SetText(title)
}

func (w *Group) Title() string {
	return w.title.Text()
}

func (g *Group) LayoutWidget(width, height int) {
	var (
		margin = g.Style().Margin.Int()
		x, y   = margin, margin
	)

	availableWidth := width - margin*2
	availableHeight := height - margin*2

	g.width = width
	g.height = margin

	if g.title.Text() != "" {
		// Layout title widget.
		g.title.LayoutWidget(width-margin*2, height-margin*2)
		g.title.MoveWidget(x, y)
		tw, th := g.title.WidgetSize()
		y += th + margin
		g.height += tw + margin
		g.width = tw + margin*2
	}

	if g.child != nil {
		// Lay out child.
		g.child.LayoutWidget(availableWidth, availableHeight)
		g.child.MoveWidget(x, y)
		cw, ch := g.child.WidgetSize()
		g.height += ch + margin
		if g.width < cw {
			g.width = cw + margin*2
		}
	}

	g.ClipTo(width, height)
}

func (w *Group) SetChild(child Control) {
	if w.child != nil {
		w.child.SetParent(nil)
	}
	w.child = child
	w.child.SetParent(w)
	w.LayoutWidget(w.width, w.height)
}

func (w *Group) Destroy() {
	// first hide ourselves
	w.Hide()
	// now destroy the child
	if w.child != nil {
		w.child.SetParent(nil)
		w.child.Destroy()
	}
	// and finally free ourselves
}

func (w Group) Borderless() bool {
	return w.borderless
}

func (w *Group) SetBorderless(borderless bool) {
	w.borderless = borderless
}

func (w *Group) DrawWidget(screen *Graphic) {
	dx, dy := w.WidgetAbsolute()

	if !w.borderless {
		FillFrameStyle(screen, dx, dy, w.width, w.height, w.Style())
	}

	if w.title.Text() != "" {
		w.title.DrawWidget(screen)
	}

	if w.child != nil {
		w.child.DrawWidget(screen)
	}
	w.DrawDebug(screen, "GRO %d %d", dx, dy)
}

func (p *Group) HandleWidget(ev Event) {
	if p.child != nil {
		p.child.HandleWidget(ev)
	}
}

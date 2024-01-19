package ui

// Pane is a virtual sub-window that can be controlled like a normal window
// inside a physical Window. Useful on platforms that don't allow multiple
// real windows (like ebui or mobile platforms), and for use in a Stack.
type Pane struct {
	BasicWidget

	menuBar   *MenuBar
	child     Control
	margined  bool
	onClosing func(*Pane)

	suppressPositionChanged bool
	focused                 bool
	needLayout              bool
	closed                  bool
	minimized               bool
	dragging                bool
	resizing                bool
	title                   string
	Ability                 // Ability lets Pane inherit abilities.

	BasicOverlayer
}

func NewPane(title string, width, height int, hasMenubar bool) *Pane {
	w := &Pane{}
	w.width = width
	w.height = height
	if hasMenubar {
		w.menuBar = NewMenuBar()
		w.menuBar.SetParent(w)
	} else {
		w.menuBar = nil
	}
	w.title = title
	w.needLayout = true
	w.SetStyle(theme.Pane)
	return w
}

func (w *Pane) SetTitle(title string) {
	w.title = title
}

var paneHeaderHeight = 24
var paneLine = 1

func (p *Pane) layoutContents(width, height int) (int, int) {
	margin := p.Style().Margin.Int()
	childWidth := width - (2 * margin)
	childHeight := height - (2 * margin)

	if width < 1 || height < 1 {
		childWidth = LayoutUnlimited
		childHeight = LayoutUnlimited
	}

	cwidth, cheight := 0, paneHeaderHeight
	if p.menuBar != nil {
		p.menuBar.LayoutWidget(childWidth, childHeight)
		bw, bh := p.menuBar.WidgetSize()
		p.menuBar.MoveWidget(0, paneHeaderHeight)

		cwidth = bw
		cheight += bh
	}

	childHeight -= cheight

	if p.child != nil {
		p.child.LayoutWidget(childWidth, childHeight)
		p.child.MoveWidget(margin, cheight+margin)
		cw, ch := p.child.WidgetSize()
		cheight += ch
		if cw > cwidth {
			cwidth = cw
		}
	}
	return cwidth, cheight
}

func (p *Pane) LayoutWidget(width, height int) {
	margin := p.Style().Margin.Int()
	cw, ch := p.layoutContents(width, height)
	p.width, p.height = cw, ch

	p.GrowToStyleSize()
	p.width += 2 * margin
	p.height += 2 * margin

	p.ClipTo(width, height)
}

func (w *Pane) SetChild(child Control) {
	if w.child != nil {
		w.child.SetParent(nil)
	}
	w.child = child
	w.child.SetParent(w)
	w.LayoutWidget(w.width, w.height)
}

func (w *Pane) Destroy() {
	// first hide ourselves
	w.Hide()
	// now destroy the child
	if w.child != nil {
		w.child.SetParent(nil)
		w.child.Destroy()
	}
	// and finally free ourselves
}

func (w *Pane) Title() string {
	return w.title
}

func (w *Pane) OnClosing(f func(*Pane)) {
	w.onClosing = f
}

func (w *Pane) DrawWidget(screen *Graphic) {
	dx, dy := w.WidgetAbsolute()
	var (
		icons     = theme.Icons
		textFace  = w.Style().Font.Face
		textColor = w.Style().Color.RGBA()
	)
	style := w.Style()
	if w.minimized {
		style = *theme.Disable
	}

	if !w.minimized && !w.Plain() {
		FillFrameStyle(screen, dx, dy, w.width, w.height, style)
	} else {
		FillFrameStyle(screen, dx, dy, w.width, paneHeaderHeight, style)
	}
	if w.title != "" {
		tw, _ := oneLineTextSize(textFace, w.title)
		TextDrawOffset(screen, w.title, textFace, dx+w.width/2-tw/2-paneHeaderHeight*3/2, dy, textColor)
	}

	iconAtlas.DrawSprite(screen, dx+w.width-paneHeaderHeight, dy, paneHeaderHeight, paneHeaderHeight, icons.Close.String())
	iconAtlas.DrawSprite(screen, dx+w.width-paneHeaderHeight*2, dy, paneHeaderHeight, paneHeaderHeight, icons.Minimize.String())
	iconAtlas.DrawSprite(screen, dx+w.width-paneHeaderHeight*3, dy, paneHeaderHeight, paneHeaderHeight, icons.Maximize.String())

	if w.child != nil && !w.minimized {
		w.child.DrawWidget(screen)
	}

	// draw overlays
	w.DrawOverlays(screen)

	// draw menu bar over everything else.
	if w.menuBar != nil {
		w.menuBar.DrawWidget(screen)
	}

	DrawDebug(screen, dx, dy, w.width, w.height, "PAN")
}

func BringToTop(p Control) {
	if parent, ok := p.Parent().(Container); ok {
		ordered := parent.Ordered()
		if len(ordered) > 1 {
			last := ordered[len(ordered)-1]
			if last != p {
				last.RaiseWidget(-containerLayerOffset)
				p.RaiseWidget(containerLayerOffset)
			}
		}
	}
}

func (p *Pane) closePaneWithCallback() {
	if p.onClosing != nil {
		p.onClosing(p)
		// Permanent Panes cannot be closed.
		if p.Permanent() {
			return
		}
	}
	p.closePane()
}

func (p *Pane) closePane() {
	p.Hide()
	UnfocusParentIfNeeded(p)
	// Preserved Panes will only be hidden and set to closed, and not destroyed.
	if !p.Preserved() {
		p.Destroy()
	}
	p.closed = true
}

func (p *Pane) HandleWidget(ev Event) {

	// Handle menu bar with priority.
	if p.menuBar != nil {
		if used := HandleContainerIfNeeded(ev, p.menuBar); used {
			return
		}
	}

	// After that, overlays get the events. Stop if the event was used.
	if used := p.HandleEventForOverlays(ev); used {
		return
	}

	// Ok, maybe it is a pane maniplation then.
	if mc, ok := ev.(*MouseClickEvent); ok {
		if p.MouseInsidePart(p.width-paneHeaderHeight, 0, paneHeaderHeight, paneHeaderHeight, mc.MouseEvent) {
			// close button
			dprintln("Pane.HandleWidget: close")
			p.closePaneWithCallback()
			SetCursorShape(CursorShapeDefault)
			return
		} else if p.MouseInsidePart(p.width-paneHeaderHeight*2, 0, paneHeaderHeight, paneHeaderHeight, mc.MouseEvent) {
			// minimize button
			dprintln("Pane.HandleWidget: minimize")
			p.minimized = true
		} else if p.MouseInsidePart(p.width-paneHeaderHeight*3, 0, paneHeaderHeight, paneHeaderHeight, mc.MouseEvent) {
			if p.minimized {
				p.minimized = false
			} else {
				if p.Parent() != nil {
					pw, ph := p.Parent().WidgetSize()
					p.MoveWidget(0, 0)
					p.LayoutWidget(pw, ph)
				}
			}
			// maximize button
			dprintln("Pane.HandleWidget: maximize")
		} else if p.MouseInsidePart(0, 0, p.width-paneHeaderHeight*3, paneHeaderHeight, mc.MouseEvent) {
			dprintln("Pane.HandleWidget: drag")
			p.dragging = true
			SetCursorShape(CursorShapeMove)
			BringToTop(p)
		} else if p.MouseInsidePart(p.width-paneHeaderHeight, p.height-paneHeaderHeight, paneHeaderHeight, paneHeaderHeight, mc.MouseEvent) {
			dprintln("Pane.HandleWidget: resize")
			p.resizing = true
			SetCursorShape(CursorShapeNWSEResize)
			BringToTop(p)
		} else {
			BringToTop(p)
			SetCursorShape(CursorShapeDefault)
		}
	}

	if p.dragging {
		if _, ok := ev.(*MouseReleaseEvent); ok {
			dprintln("Pane.HandleWidget: drop")
			p.dragging = false
			SetCursorShape(CursorShapeDefault)
		}
		if mm, ok := ev.(*MouseMoveEvent); ok {
			dprintln("Pane.HandleWidget: move")
			p.x += mm.MoveX
			p.y += mm.MoveY
		}
	} else if p.resizing {
		if _, ok := ev.(*MouseReleaseEvent); ok {
			dprintln("Pane.HandleWidget: done resizing")
			p.resizing = false
			p.layoutContents(p.width, p.height)
			SetCursorShape(CursorShapeDefault)
		}
		if mm, ok := ev.(*MouseMoveEvent); ok {
			dprintln("Pane.HandleWidget: resize")
			p.width += mm.MoveX
			p.height += mm.MoveY
		}
	} else {
		if mm, ok := ev.(*MouseMoveEvent); ok {
			if p.MouseInsidePart(0, 0, p.width-paneHeaderHeight*3, paneHeaderHeight, mm.MouseEvent) {
				SetCursorShape(CursorShapeMove)
			} else if p.MouseInsidePart(p.width-paneHeaderHeight, p.height-paneHeaderHeight, paneHeaderHeight, paneHeaderHeight, mm.MouseEvent) {
				SetCursorShape(CursorShapeNWSEResize)
			}
		}
	}

	// Pass on to the child widget.
	if p.child != nil && !p.minimized {
		p.child.HandleWidget(ev)
	}
}

var _ Overlayer = &Pane{}

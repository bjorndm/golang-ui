package ui

import "github.com/hajimehoshi/ebiten/v2"
import "golang.org/x/exp/slices"

var DebugColor = true

type HasChildren interface {
	Children() []Control
	Ordered() []Control
}

type Container interface {
	Control
	HasChildren
}

func CompareControls(a, b Control) int {
	al := a.WidgetLayer()
	bl := b.WidgetLayer()
	if al != bl {
		return al - bl
	}
	ax, ay := a.WidgetAt()
	bx, by := a.WidgetAt()
	if ay != by {
		return ay - by
	}
	return ax - bx
}

type BasicContainer struct {
	BasicWidget
	controls []Control
	// Ordered are the controls ordered by ascending layer, to facilitate
	// drawing and event handling.
	ordered         []Control
	tab             int
	mouseFocusIndex int
	mouseFocus      Control
}

const containerLayerOffset = 100000

func (c BasicContainer) WidgetLayer() int {
	highest := 0
	if len(c.ordered) > 0 {
		highest = c.ordered[len(c.ordered)-1].WidgetLayer()
	}
	return c.z*containerLayerOffset + highest
}

func (c BasicContainer) DrawWidget(g *Graphic) {
	dx, dy := c.WidgetAbsolute() // NOTE static inheritance !

	if c.tab > 0 && c.tab < len(c.controls) {
		tabChild := c.controls[c.tab-1]
		wx, wy := tabChild.WidgetAt()
		ww, wh := tabChild.WidgetSize()
		DrawFrameOptionalStyle(g, dx+wx, dy+wy, ww, wh, theme.Focus)
	}

	// draw widgets in drawing order
	for _, child := range c.ordered {
		if child.Hidden() {
			continue
		}
		child.DrawWidget(g)
	}
}

func (c *BasicContainer) HandleAway(ae *AwayEvent) {
	for i := len(c.ordered) - 1; i >= 0; i-- {
		child := c.ordered[i]
		if child.Hidden() {
			continue
		}
		child.HandleWidget(ae)
	}
	// also remove our own focus, since the focus went away
	c.tab = 0
	c.SetFocus(nil) // XXX static inheritance !
	// sort the children in drawing order again
	c.UpdateOrdered() // XXX static inheritance !
}

func (c *BasicContainer) SetFocus(w Control) {
	c.BasicWidget.SetFocus(w)
	if w != nil {
		idx := slices.IndexFunc(c.controls, func(wi Control) bool { return wi == w })
		c.tab = idx + 1
	} else {
		c.tab = 0
	}
}

func (c *BasicContainer) HandleWidget(ev Event) {
	// If we get an away event, pass it on to all children.
	if ae, ok := ev.(*AwayEvent); ok {
		c.HandleAway(ae) // NOTE static inheritance !
		return
	}

	if kr, ok := ev.(*KeyReleaseEvent); ok {
		if debugDisplay {
			dprintln("Box.HandleWidget: key release", kr.Name(), kr.Key)
		}
		if kr.Key == ebiten.KeyF2 {
			debugDisplay = !debugDisplay
		}
		if kr.Key == ebiten.KeyTab {
			if len(c.controls) > 0 && c.tab > 0 && c.Focus() != nil { // XXX static inheritance !
				c.tab++
				if c.tab >= len(c.controls) {
					c.tab = 1
				}
				control := c.controls[c.tab]
				SetNewFocus(kr.Event(), c, control, c.controls...) // XXX static inheritance !
				return
			}
		}
	}

	HandleContainerIfNeeded(ev, c)

	// sort the children in drawing order again
	c.UpdateOrdered() // Note static inheritance
}

func (c *BasicContainer) Children() []Control {
	return c.controls
}

func (c *BasicContainer) Ordered() []Control {
	return c.ordered
}

func (b *BasicContainer) Delete(index int) {
	dprintln("BasicContainer.Delete")

	if index < 0 || index >= len(b.controls) {
		return
	}
	bc := b.controls[index]
	UnfocusParentIfNeeded(bc)
	b.controls = slices.Delete(b.controls, index, index+1)
	b.UpdateOrdered() // NOTE static inheritance

	NeedLayout(b) // NOTE static inheritance
}

func (b BasicContainer) NumChildren() int {
	return len(b.controls)
}

func (b *BasicContainer) UpdateOrdered() {
	b.ordered = slices.Clone(b.controls)
	slices.SortFunc(b.ordered, CompareControls)
}

// AppendAppendControl appends the control to the basic container.
// It does not set the parent as this is not possible to do this correctly
// with static inheritance and this function signature.
func (b *BasicContainer) AppendControl(control Control) {
	if b == control {
		panic("append cycle detected")
	}
	b.controls = append(b.controls, control)
	// Although these two are statically inherited due to the way they work
	// they are OK. UpdateOrdered only touches b, and NeedLayout bubbles up to
	// the toplevel window.
	b.UpdateOrdered()
	NeedLayout(b)
}

// AppendWithParent appends the control to the basic container. It also sets the parent.
func (b *BasicContainer) AppendWithParent(control, parent Control) {
	if parent == control {
		panic("append cycle detected")
	}
	control.SetParent(parent)
	b.AppendControl(control)
}

func (c *BasicContainer) Clear() {
	c.tab = -1
	for _, child := range c.controls {
		child.Destroy()
	}

	c.controls = nil
}

func (c *BasicContainer) InsertAt(control Control, index int) {
	control.SetParent(c)
	slices.Insert(c.controls, index, control)
}

func (c *BasicContainer) NumItems() int {
	return len(c.controls)
}

func (b *BasicContainer) Get(index int) Control {
	if index < 0 || index >= len(b.controls) {
		return nil
	}
	return b.controls[index]
}

func BasicContainerGet[T Control](b *BasicContainer, index int) (T, bool) {
	elem := b.Get(index) // NOTE static inheritance
	if elem == nil {
		var zero T
		return zero, false
	}
	t, ok := elem.(T)
	return t, ok
}

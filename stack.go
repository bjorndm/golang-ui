package ui

// import "github.com/hajimehoshi/ebiten/v2"

const stackLayerOffset = 1000

// A Stack is a container with widgets layed out stacked over each other.
// Unlike most widgest stack does not move the widgets, so exceptionally
// any child widgets can position themselves if needed.
// A Stack also stretches to the size given in LayoutWidget.
// Normally Window has a Stack for dialogs of the same size of Window.
type Stack struct {
	BasicContainer
	padded bool
}

func (b *Stack) Destroy() {
	// free all controls
	for i := 0; i < len(b.controls); i++ {
		bc := b.controls[i]
		bc.SetParent(nil)
		bc.Destroy()
	}
}

func (b *Stack) Enable() {
	for _, child := range b.controls {
		bw := child.(*BasicWidget)
		bw.Enable()
	}
}

func (b *Stack) Disable() {
	for _, child := range b.controls {
		bw := child.(*BasicWidget)
		bw.Disable()
	}
}

func (b *Stack) Append(c Control) {
	b.BasicContainer.AppendWithParent(c, b)
	c.RaiseWidget(stackLayerOffset * b.NumChildren())
	if b.width > 0 && b.height > 0 {
		b.LayoutWidget(b.width, b.height)
	}
}

func (b *Stack) deleteNoLayout(index int) {
	dprintln("Stack.delete")

	if index < 0 || index >= len(b.controls) {
		return
	}
	bc := b.controls[index]

	bc.SetParent(nil)
	b.controls = append(b.controls[:index-1], b.controls[:index+1]...)
}

func (b Stack) Padded() bool {
	return b.padded
}

func (b *Stack) SetPadded(padded bool) {
	b.padded = padded
}

func newStack() *Stack {
	b := &Stack{}
	b.controls = []Control{}
	return b
}

// NewStack creates a new stack. Stacks always stretch to the size
// passed to them in LayoutWidget.
func NewStack() *Stack {
	return newStack()
}

func (s *Stack) LayoutWidget(width, height int) {
	if width < 1 || height < 1 {
		panic("Stack: cannot lay out without space.")
	}
	// The size of a stack is simply the available space.
	// The largest of the least width and least height of the stacked items.
	w := 0
	h := 0
	for i, child := range s.controls {
		if child.Hidden() {
			continue
		}
		child.LayoutWidget(width, height)
		realWidth, realHeight := child.WidgetSize()
		if realWidth > w {
			w = realWidth
		}
		if realHeight > h {
			h = realHeight
		}
		if pane, ok := child.(*Pane); ok {
			if pane.needLayout {
				// The stacked items should stay where they are, but
				// still position new panes.
				// Center the stacked items, but if i > 0 shift it down a bit
				// to prevent overlaps.
				// Also wrap around to prevent invisible widgets.
				x := (width - realWidth) / 2
				remaining := (height - realHeight)
				if remaining < realHeight {
					remaining = realHeight + 1
				}
				y := (i * paneHeaderHeight) % (remaining)
				y += paneHeaderHeight * 2
				x += (i * paneHeaderHeight) / (remaining)
				child.MoveWidget(x, y)
				pane.needLayout = false
			}
		}
	}
	s.width, s.height = width, height
	s.ClipTo(width, height)
}

func (b Stack) DrawWidget(g *Graphic) {
	b.BasicContainer.DrawWidget(g)
	b.DrawDebug(g, "STA")
}

func (c *Stack) CleanupClosePanes() {
	closed := []int{}
	// Handle any closed panes.
	for i := len(c.controls) - 1; i >= 0; i-- {
		child := c.controls[i]
		if pane, ok := child.(*Pane); ok && pane.closed {
			closed = append(closed, i)
			break
		}
	}

	for _, toClose := range closed {
		c.Delete(toClose)
	}
	// sort the children in drawing order again
	c.UpdateOrdered()
}

func (c *Stack) HandleWidget(ev Event) {

	c.BasicContainer.HandleWidget(ev)
	c.CleanupClosePanes()
}

var _ Control = &Stack{}

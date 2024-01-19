package ui

import "golang.org/x/exp/slices"

// Mesh is an element of a Grid.
type Mesh struct {
	Control
	align StyleAlign
	left  int
	top   int
	span  int
}

func NewMesh(child Control, left, top, span int, align StyleAlign) *Mesh {
	if left < 0 || top < 0 || span < 1 {
		panic("NewMesh: out of range")
	}
	mesh := &Mesh{Control: child,
		left: left, top: top,
		span: span, align: align}
	return mesh
}

func (a StyleAlign) Position(value, size, space int) int {
	switch a {
	case StyleAlignRight:
		return value + space - size
	case StyleAlignMiddle:
		return value + (space-size)/2
	default:
		return value
	}
}

func MoveWidgetStyleAligned(w Control, x, y, width, height int, align StyleAlign) {
	widgetWidth, _ := w.WidgetSize()
	x = align.Position(x, widgetWidth, width)
	w.MoveWidget(x, y)
}

func (m *Mesh) MoveStyleAligned(x, y, width, height int) {
	if m.Control == nil {
		return
	}
	MoveWidgetStyleAligned(m.Control, x, y, width, height, m.align)
}

// A Grid is a grid of equally spaced columns with rows of varying height.
// The widgets in the grid are layed out in rows of columns of fixed space.
// It is possible to make wigets span multiple columns.
// It is also possible to align the Mesh in the columns to the start, center
// or end of the column they are in.
// Grid will automatically expand vertically but not horizontally.
type Grid struct {
	BasicContainer
	meshes  []*Mesh
	rows    int
	columns int
}

func (g *Grid) Destroy() {
	// free all controls
	for i := 0; i < len(g.controls); i++ {
		bc := g.controls[i]
		bc.SetParent(nil)
		bc.Destroy()
	}
	g.controls = []Control{}
	g.meshes = []*Mesh{}
}

func (g Grid) NumberOfColumns() int {
	return g.columns
}

func (g Grid) NumberOfRows() int {
	return g.rows
}

func (g *Grid) indexMesh(left, top int) int {
	if (left < 0) || (top < 0) || (left >= g.NumberOfColumns()) || (top >= g.NumberOfRows()) {
		panic("indexMesh in grid out of range")
	}

	index := slices.IndexFunc(g.meshes, func(mesh *Mesh) bool {
		return mesh.left == left && mesh.top == top
	})
	return index
}

func (g *Grid) getMesh(left, top int) *Mesh {
	index := g.indexMesh(left, top)
	if index < 0 {
		return nil
	}

	return g.meshes[index]
}

func (g *Grid) Merge(left, top, span int) {
	if (left < 0) || (top < 0) || (left >= g.NumberOfColumns()) || (top >= g.NumberOfRows()) {
		panic("Merge in grid: index out of range")
	}
	if span < 1 {
		panic("Merge in grid: span out of range")
	}
	right := left + span
	bottom := top
	if (right < 0) || (bottom < 0) || (right >= g.NumberOfColumns()) || (bottom >= g.NumberOfRows()) {
		panic("Merge in grid: span out of range")
	}

	mesh := g.getMesh(left, top)
	mesh.span = span
}

func (g *Grid) Put(child Control, left, top int) {
	if (left < 0) || (top < 0) || (left >= g.NumberOfColumns()) || (top >= g.NumberOfRows()) {
		panic("Put in grid out of range")
	}
	mesh := g.getMesh(left, top)
	if mesh == nil {
		// New mesh
		mesh = &Mesh{left: left, top: top, span: 1}
		g.meshes = append(g.meshes, mesh)
	} else {
		// Already exists
		if mesh.Control != nil {
			mesh.Control.Destroy()
			mesh.Control = nil
		}
	}
	mesh.Control = child

	g.BasicContainer.AppendWithParent(child, g)

	if g.width > 0 && g.height > 0 {
		g.LayoutWidget(g.width, g.height)
	}
}

func (g *Grid) putMesh(mesh *Mesh) {
	if (mesh.left < 0) || (mesh.top < 0) || (mesh.left >= g.NumberOfColumns()) || (mesh.top >= g.NumberOfRows()) {
		panic("Put in grid out of range")
	}
	old := g.indexMesh(mesh.left, mesh.top)
	if old < 0 {
		g.meshes = append(g.meshes, mesh)
		g.BasicContainer.AppendWithParent(mesh.Control, g)
	} else {
		if g.meshes[old].Control != nil {
			g.meshes[old].Control.Destroy()
		}
		g.meshes[old] = mesh
		g.BasicContainer.controls[old] = mesh.Control
		mesh.Control.SetParent(g)
	}

	if g.width > 0 && g.height > 0 {
		g.LayoutWidget(g.width, g.height)
	}
}

func (g *Grid) SetLayout(left, top int, align StyleAlign) {
	mesh := g.getMesh(left, top)
	if mesh == nil {
		panic("SetLayout in grid out of range")
	}

	mesh.align = align
}

func (g *Grid) Append(child Control, left, top, span int, align StyleAlign) {
	mesh := NewMesh(child, left, top, span, align)
	g.AppendMesh(mesh)
}

func (g *Grid) AppendMesh(mesh *Mesh) {
	// Increase column and row count if needed.
	if mesh.left+mesh.span > g.columns {
		g.columns = mesh.left + mesh.span
	}

	if mesh.top+1 > g.rows {
		g.rows = mesh.top + 1
	}

	g.putMesh(mesh)
	if g.width > 0 && g.height > 0 {
		g.LayoutWidget(g.width, g.height)
	}
}

func (g *Grid) AppendWithLabel(label string, widget Control) {
	row := g.NumberOfRows()
	g.Append(NewLabel(label), 0, row, 1, StyleAlignLeft)
	g.Append(widget, 1, row, 1, StyleAlignLeft)
}

func (g *Grid) AppendWithoutLabel(widget Control) {
	row := g.NumberOfRows()
	g.Append(widget, 1, row, 1, StyleAlignLeft)
}

// AppendMany will append many labeled controls to the grid.
// While the pairs arguments is of type any, it must be passed as
// string / ui.Control pairs. The string will be then converted to the label.
// However if the string is empty no label will be used and the control will be
// stretched.
// The function will panic if the arguments are not correctly typed.
func (g *Grid) AppendMany(pairs ...any) {
	for i := 1; i < len(pairs); i += 2 {
		label := pairs[i-1].(string)
		widget := pairs[i].(Control)
		if label == "" {
			g.AppendWithoutLabel(widget)
		} else {
			g.AppendWithLabel(label, widget)
		}
	}
}

func newGrid() *Grid {
	b := &Grid{}
	b.tab = 1
	b.SetStyle(theme.Grid)
	return b
}

func NewGrid() *Grid {
	return newGrid()
}

func (g *Grid) LayoutWidget(width, height int) {
	var margin = g.Style().Margin.Int()

	if g.columns == 0 || g.rows == 0 {
		g.GrowToStyleSize()
		g.width += 2 * margin
		g.height += 2 * margin
		g.ClipTo(width, height)
		return
	}

	var (
		y           = 0
		columnWidth = (width - margin*2) / g.columns
		rowHeight   = (height - margin*2) / g.rows
	)

	g.width = (width - margin*2)
	g.height = margin * 2

	for row := 0; row < g.rows; row++ {
		highest := 0
		for col := 0; col < g.columns; col++ {
			mesh := g.getMesh(col, row)
			if mesh == nil || mesh.Control == nil || mesh.Control.Hidden() {
				continue
			}
			// X position is fixed by the columns.
			x := col * columnWidth
			maxWidth := columnWidth * mesh.span
			maxHeight := rowHeight
			mesh.Control.LayoutWidget(maxWidth, maxHeight)
			// Stretch row if widget overflows it.
			_, oh := mesh.Control.WidgetOverflow()
			if oh > 0 {
				maxHeight := height
				mesh.Control.LayoutWidget(maxWidth, maxHeight)
			}

			mesh.MoveStyleAligned(x, y, maxWidth, maxHeight)

			ww, wh := mesh.Control.WidgetSize()
			if wh > highest {
				highest = wh
			}
			x += ww
		}
		// Y position is shifted by the highest in the row, which
		// must be smaller than rowHeight unless spanned
		y += highest
		g.height += highest
	}
	g.BasicContainer.UpdateOrdered()
	g.ClipTo(width, height)
}

func (b Grid) DrawWidget(g *Graphic) {
	dx, dy := b.WidgetAbsolute()

	FillFrameStyle(g, dx, dy, b.width, b.height, b.Style())
	b.BasicContainer.DrawWidget(g)
	b.DrawDebug(g, "GRI")
}

var _ Control = &Grid{}

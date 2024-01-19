package ui

type (
	ValueConstraint interface {
		~string | ~int64 | ~bool | RGBA | ~float64 | *Graphic
	}
	Value any
)

func NewValue[T ValueConstraint](v T) Value {
	return v
}

// Row is an array of values.
type Row []Value

func NewRow(values ...Value) Row {
	return values
}

func (r Row) Value(index int) Value {
	if index < 0 || index >= len(r) {
		return nil
	}
	return r[index]
}

func (r *Row) SetValue(index int, v Value) {
	if index < 0 || index > len(*r) {
		return
	}
	(*r)[index] = v
}

func (r Row) Values() []Value {
	return r
}

func (r Row) NumValues() int {
	return len(r)
}

// The table model is used to fetch the rows needed to display,
// and to react to edits.
type TableModel interface {
	NumRows() int
	FetchRow(index int) Row
	UpdateRow(index int, updated Row)
}

type columnKind int

const (
	columnKindText columnKind = iota
	columnKindButton
	columnKindCheckbox
	columnKindEntry
	columnKindPicture
	columnKindColor
)

// Column is a column in a table.
type Column struct {
	BasicWidget
	table     *Table
	caption   *TextWidget
	name      string
	index     int
	kind      columnKind
	marker    string // marker icon on header, for example to indicate sorting.
	button    string // button icon
	onClicked func(col *Column, row int)
}

// OnClicked sets a callback function f that is called, if not nil,
// if the column is clicked. This overrides the default behavior.
// If f is nil the default behavior is restored.
func (c *Column) OnClicked(f func(col *Column, row int)) {
	c.onClicked = f
}

const minColWidth = 16
const minRowHeight = 16

func newColumn(ct columnKind, name, caption string, index int) *Column {

	switch ct {
	case columnKindText:
	case columnKindButton:
	case columnKindCheckbox:
	case columnKindPicture:
	default:
		panic("unknown or not implemented column type")
	}

	col := &Column{
		name:  name,
		index: index,
		kind:  ct,
	}
	col.width = theme.Column.Size.Width.Int()
	col.SetStyle(theme.Column)
	col.caption = NewTextWidget(caption)
	col.caption.SetParent(col)
	col.caption.width = col.width
	return col
}

func NewButtonColumn(name string, index int, icon string) *Column {
	col := newColumn(columnKindButton, name, name, index)
	col.button = icon
	return col
}

func NewCheckboxColumn(name string, index int) *Column {
	return newColumn(columnKindCheckbox, name, name, index)
}

func NewPictureColumn(name string, index int) *Column {
	return newColumn(columnKindPicture, name, name, index)
}

// NewTextColumn creates a new table column. It is a text display column.
// index is the index in the row to use for the display of this column.
func NewTextColumn(name string, index int) *Column {
	return newColumn(columnKindText, name, name, index)
}

func (c Column) drawCellText(dst *Graphic, dx, dy, row int, value Value) {
	face := c.Style().Font.Face
	col := c.Style().Color.RGBA()
	if !c.Enabled() {
		col = theme.Disable.Color.RGBA()
	}

	if text, ok := value.(string); ok {
		TextDrawOffset(dst, text, face, dx, dy, col)
	}
}

func (c Column) drawCellButton(dst *Graphic, dx, dy, row int, value Value) {
	face := c.Style().Font.Face
	textColor := c.Style().Color.RGBA()
	marginButton := c.Style().Margin.Int()

	style := theme.Button

	if bvalue, ok := value.(bool); ok && bvalue {
		dx += marginButton / 2
		dy += marginButton / 2
		style = theme.Active
	}

	if !c.Enabled() {
		style = theme.Disable
	}

	h := style.Font.Face.Metrics().Height.Round()
	h -= (marginButton) * 2

	TextDrawOffset(dst, c.name, face, dx, dy, textColor)
	_, nw := oneLineTextSize(face, c.name)
	if c.button != "" {
		dx += marginButton + nw + h
		iconAtlas.DrawSprite(dst, dx, dy, h, h, c.button)
	}
}

func (c Column) drawCellCheckbox(dst *Graphic, dx, dy, row int, value Value) {
	style := c.Style()
	marginCheckbox := c.Style().Margin.Int()
	checkboxHeight := c.Style().Size.Height.Int()
	checkboxWidth := checkboxHeight

	dx += marginCheckbox
	dy += marginCheckbox
	FillFrameStyle(dst, dx, dy, checkboxWidth, checkboxHeight, style)
	if bvalue, ok := value.(bool); ok && bvalue {
		iconAtlas.DrawSprite(dst, dx, dy, checkboxWidth, checkboxHeight, theme.Icons.Check.String())
	}
}

func (c Column) drawCellEntry(dst *Graphic, dx, dy, row int, value Value) {
	// TODO: editing
	face := c.Style().Font.Face
	col := c.Style().Color.RGBA()
	if !c.Enabled() {
		col = theme.Disable.Color.RGBA()
	}

	if text, ok := value.(string); ok {
		TextDrawOffset(dst, text, face, dx, dy, col)
	}
}

func (c Column) drawCellPicture(dst *Graphic, dx, dy, row int, value Value) {
	var (
		fillColor = theme.Error.Fill.Color.RGBA()
		margin    = c.Style().Margin.Int()
	)
	h := c.Style().Font.Face.Metrics().Height.Round()
	ww, wh := h+margin*2, h+margin*2

	if graphic, ok := value.(*Graphic); ok && graphic != nil {
		iw, ih := graphic.Size()
		ix, iy := dx+margin, dy+margin+h
		sx, sy := float64(ww)/float64(iw), float64(wh)/float64(ih)
		DrawGraphicAtScale(dst, graphic, ix, iy, sx, sy)
	} else {
		FillRect(dst, dx, dy, ww, wh, fillColor)
	}
}

func (c Column) drawCellColor(dst *Graphic, dx, dy, row int, value Value) {
	var (
		fillColor = theme.Error.Fill.Color.RGBA()
		margin    = c.Style().Margin.Int()
	)
	h := c.Style().Font.Face.Metrics().Height.Round()
	ww, wh := h+margin*2, h+margin*2

	if color, ok := value.(RGBA); ok {
		fillColor = color
	}
	FillRect(dst, dx+margin, dy+margin, ww, wh, fillColor)
}

func (c Column) drawCellContents(dst *Graphic, dx, dy, row int, value Value) {
	switch c.kind {
	case columnKindText:
		c.drawCellText(dst, dx, dy, row, value)
	case columnKindButton:
		c.drawCellButton(dst, dx, dy, row, value)
	case columnKindCheckbox:
		c.drawCellCheckbox(dst, dx, dy, row, value)
	case columnKindEntry:
		c.drawCellEntry(dst, dx, dy, row, value)
	case columnKindPicture:
		c.drawCellPicture(dst, dx, dy, row, value)
	case columnKindColor:
		c.drawCellColor(dst, dx, dy, row, value)
	}
}

var missingValue = NewValue("###???###")

func (c Column) DrawWidget(screen *Graphic) {
	dx, dy := c.WidgetAbsolute()
	rows := c.table.NumRows()
	rowh := c.table.RowHeight()
	if !c.caption.Hidden() {
		c.caption.DrawWidget(screen)
		capw, caph := c.caption.WidgetSize()
		_ = capw
		if c.marker != "" {
			iconAtlas.DrawSprite(screen, dx+capw-rowh, dy, rowh, rowh, c.marker)
		}
		dy += caph
	}

	// Show the rows that should be visible due to scrolling only.
	start := c.table.from
	stop := rows + start
	if c.table.shown > 0 && c.table.shown < rows {
		stop = c.table.shown + start
		if stop > rows {
			stop = rows
		}
	}

	for i := start; i < stop; i++ {
		FillFrameStyle(screen, dx, dy, c.width, rowh, c.Style())
		row := c.table.FetchRow(i)
		var value Value = nil
		if row != nil {
			value = row.Value(c.index)
		}
		if value != nil {
			c.drawCellContents(screen, dx, dy, i, value)
		} else {
			c.drawCellContents(screen, dx, dy, i, missingValue)
		}
		dy += rowh
	}
}

type Table struct {
	Tray       // use a tray to lay out the columns.
	TableModel // table model for fetching the data.
	columns    []*Column

	from            int
	shown           int
	onHeaderClicked func(*Table, int)
	onClicked       func(*Table, int, int)
	rowHeight       int
}

func (t Table) RowHeight() int {
	h := t.Style().Font.Face.Metrics().Height.Round()
	p := t.Style().Margin.Int()
	m := t.Style().Margin.Int()
	rh := h + 2*p + 2*m

	if t.rowHeight > rh {
		return t.rowHeight
	}
	return rh
}

func (t Table) SetRowHeight(rh int) {
	if rh < 0 {
		h := t.Style().Font.Face.Metrics().Height.Round()
		p := t.Style().Margin.Int()
		m := t.Style().Margin.Int()
		t.rowHeight = h + 2*p + 2*m
	} else {
		t.rowHeight = rh
	}
}

func (t *Table) OnHeaderClicked(f func(*Table, int)) {
	t.onHeaderClicked = f
}

func (t *Table) OnClicked(f func(*Table, int, int)) {
	t.onClicked = f
}

func NewTable(model TableModel) *Table {
	g := &Table{TableModel: model}
	g.SetStyle(theme.Table)
	g.SetRowHeight(-1)
	return g
}

const tableMargin = 3

const tableMinRows = 10
const tableMinCols = 10

func (t *Table) LayoutWidget(width, height int) {
	w, h := 0, 0
	hw, hh := 0, 0

	for _, column := range t.columns {
		cw, ch := column.WidgetSize()
		w += cw
		// Highest column.
		if ch > h {
			h = ch
		}
	}
	h += hh
	if w < hw {
		w = hw
	}
	lh := t.Style().Size.Height.Int()
	lw := t.Style().Size.Width.Int()
	if w < lw {
		w = lw
	}
	if h < lh {
		h = lh
	}

	tw, th := w, h

	t.Tray.LayoutWidget(width, height-hh)
	t.Tray.MoveWidget(0, hh)
	t.width = tw
	t.height = th

	shownHeight := t.height
	if shownHeight > height {
		shownHeight = height
	}

	t.shown = shownHeight / t.RowHeight()
	if t.from >= t.shown {
		t.from = t.shown - 1
	}
	if t.from < 0 {
		t.from = 0
	}

	println("Table.LayoutWidget", t.width, t.height, t.Tray.width, t.Tray.height, shownHeight, t.RowHeight(), t.shown, height)
}

func (t Table) DrawWidget(dst *Graphic) {
	t.Tray.DrawWidget(dst)
}

// AppendColumn appends a column to the table.
func (g *Table) AppendColumn(col *Column) {
	idx := len(g.columns)
	col.index = idx
	col.table = g
	g.columns = append(g.columns, col)
	g.Tray.Append(col)
}

// ModelRowUpdated should be called whenever a row in the data model was updated.
func (g *Table) ModelRowUpdated(index int) {
}

// ModelRowCreated should be called whenever a row in the data model was created.
// For an append index may be equal to model.NumRows()
func (g *Table) ModelRowCreated(index int) {
}

// ModelRowDeleted should be called whenever a row in the data model was deleted.
func (g *Table) ModelRowDeleted(index int) {
}

func (t *Table) SetHeaderVisible(visible bool) {
	if visible {
		for _, col := range t.columns {
			col.caption.Show()
		}
	} else {
		for _, col := range t.columns {
			col.caption.Hide()
		}
	}
}

func (t *Table) Column(index int) *Column {
	if index < 0 || index > len(t.columns) {
		return nil
	}
	return t.columns[index]
}

// SetWidth sets the width of the column.
// If width is negative, this will attempt to auto-size the column.
// This is not guaranteed to work.
func (c *Column) SetWidth(width int) {
	if width < 0 {
		row := c.table.FetchRow(0)
		if row != nil {
			value := row.Value(c.index)
			if svalue, ok := value.(string); ok {
				c.width, _ = oneLineTextSize(c.Style().Font.Face, svalue)
				c.caption.width = c.width
				return
			}
		}
		c.width = c.Style().Size.Width.Int()
		c.caption.width = c.width
	} else {
		c.width = width
		c.caption.width = c.width
	}
}

// Returns the width of the column.
func (c Column) Width() int {
	return c.width
}

// NumColumns returns the amount of columns the table has.
func (t Table) NumColumns() int {
	return len(t.columns)
}

func (c *Column) LayoutWidget(parentWidth, parentHeight int) {
	var width, height int
	rows := c.table.NumRows()
	rowHeight := c.table.RowHeight()
	width = c.caption.width
	height = 0
	if rowHeight < 1 {
		p := c.Style().Margin.Int()
		m := c.Style().Margin.Int()
		h := c.Style().Font.Face.Metrics().Height.Round()
		rowHeight = 2*p + 2*m + h
	}
	height += rows * rowHeight

	lheight := c.Style().Size.Height.Int()
	lwidth := c.Style().Size.Width.Int()
	if !c.caption.Hidden() {
		_, caph := c.caption.WidgetSize()
		height += caph
	}

	if width < lwidth {
		width = lwidth
	}
	if height < lheight {
		height = lheight
	}

	minw, minh := width, height
	c.caption.LayoutWidget(parentWidth, parentHeight)
	capw, caph := c.caption.WidgetSize()
	if capw < minw {
		capw = minw
	}
	c.caption.LayoutWidget(capw, caph)
	c.caption.width = capw
	c.caption.height = caph
	// c.caption.MoveWidget((colw-capw)/2, (colh-caph)/2)
	c.height = minh
	c.width = minw
	c.ClipTo(parentWidth, parentHeight)
}

func (c *Column) headerClicked() {
	if c.table.onHeaderClicked != nil {
		c.table.onHeaderClicked(c.table, c.index)
	}
}

func (c Column) cellClickedText(index int, row Row) {
	// Do nothing for now.
}

func (c Column) cellClickedButton(index int, row Row) {
	value := row.Value(c.index)

	// Clicking the button toggles a boolean.
	if bvalue, ok := value.(bool); ok {
		// notify model
		newValue := NewValue(!bvalue)
		row.SetValue(c.index, newValue)
		c.table.UpdateRow(index, row)
	}
}

func (c Column) cellClickedCheckbox(index int, row Row) {
	value := row.Value(c.index)
	// toggle boolean.
	if bvalue, ok := value.(bool); ok {
		// notify model
		newValue := NewValue(!bvalue)
		row.SetValue(c.index, newValue)
		c.table.UpdateRow(index, row)
	}
}

func (c Column) cellClickedEntry(index int, row Row) {
	// TODO: editing, now just sends the update without editing.
	value := row.Value(c.index)

	if svalue, ok := value.(string); ok {
		// notify model
		newValue := NewValue(svalue)
		row.SetValue(c.index, newValue)
		c.table.UpdateRow(index, row)
	}
}

func (c Column) cellClickedPicture(index int, row Row) {
	// do nothing, may implement image update later.
}

func (c Column) cellClickedColor(index int, row Row) {
	// do nothing, may implement color update later.
}

func (c *Column) cellClicked(index int) {
	var row Row

	// Override click if onClicked is set.
	if c.onClicked != nil {
		c.onClicked(c, index)
		return
	}

	if c.table != nil || c.table.TableModel != nil {
		row = c.table.FetchRow(index)
	}
	switch c.kind {
	case columnKindText:

		c.cellClickedText(index, row)
	case columnKindButton:
		c.cellClickedButton(index, row)
	case columnKindCheckbox:
		c.cellClickedCheckbox(index, row)
	case columnKindEntry:
		c.cellClickedEntry(index, row)
	case columnKindPicture:
		c.cellClickedPicture(index, row)
	case columnKindColor:
		c.cellClickedColor(index, row)
	}
	if c.table.onClicked != nil {
		c.table.onClicked(c.table, c.index, index)
	}
}

func (c *Column) HandleWidget(ev Event) {
	if mc, ok := ev.(*MouseClickEvent); ok {
		_, ay := c.WidgetAbsolute()
		dy := mc.Y - ay
		caph := 0
		if !c.caption.Hidden() {
			if mc.Inside(c.caption) {
				c.headerClicked()
				return
			}
			_, caph = c.caption.WidgetSize()
		}
		dy -= caph
		rh := c.table.RowHeight()
		i := (dy / rh) - c.table.from
		c.cellClicked(i)
	}

}

// SetMarker sets the marker icon on the column header.
func (c *Column) SetMarker(marker string) {
	c.marker = marker
}

// Marker returns the marker icon on the column header.
func (c Column) Marker() string {
	return c.marker
}

func (t *Table) ScrollWidget(y int) {
	t.from = y / t.RowHeight()
	if t.from < 0 {
		t.from = 0
	}
	t.y = 0 // -y % t.TableModel.NumRows()
	println("ScrollWidget", y, t.y, t.from, t.shown, t.TableModel.NumRows())
}

func (t *Table) RollWidget(x int) {
	t.x = -x
}

func (t *Table) ScrollSize() (width, height int) {
	w, _ := t.WidgetSize()
	h := t.TableModel.NumRows() * t.RowHeight()
	return w, h
}

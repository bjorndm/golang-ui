package ui

import "strings"

import "golang.org/x/exp/slices"

type TextPoint struct {
	X int
	Y int
}

type TextCursor struct {
	TextPoint
}

type TextSelection struct {
	From TextPoint
	To   TextPoint
}

// A note widget allows editing multiple lines of plain text.
// The new line character is \n, nor \r or \r\n on all platforms.
type Note struct {
	TextWidget
	onChanged       func(*Note)
	readonly        bool
	placeholder     string
	lines           [][]rune // stored here as a silce or une slice lines
	cursor          TextCursor
	selection       TextSelection
	textInputState  chan TextInputState
	closeInputState func()
	lastInputState  TextInputState
	active          bool
}

type PasswordNote struct {
	Note
}

type SearchNote struct {
	Note
}

func (e *Note) OnChanged(cb func(e *Note)) {
	e.onChanged = cb
}

const noteMaxWidth = 200.0
const noteMinWidth = 100.0

func newNote() *Note {
	e := &Note{}
	e.lines = [][]rune{[]rune{}}
	e.cursor.X = 0
	e.cursor.Y = 0
	// Set placeholder
	e.placeholder = "placeholder"
	e.SetStyle(theme.Note)
	return e
}

const noteHeight = 14 * 3
const noteWidth = 64

func (e *Note) SetText(text string) {
	e.TextWidget.SetText(text)
	lines := strings.Split(text, "\n")
	e.lines = make([][]rune, len(lines)+1)
	for i, line := range lines {
		e.lines[i] = []rune(line)
	}

	last := e.lines[len(e.lines)-1]
	e.cursor.X = len([]rune(last))
	e.cursor.Y = len(e.lines) - 1
}

func (e *Note) Text() string {
	res := ""
	sep := ""
	for _, line := range e.lines {
		res += sep + string(line)
		sep = "\n"
	}
	return res
}

func (e *Note) LayoutWidget(width, height int) {
	margin := e.Style().Margin.Int()
	textFace := e.Style().Font.Face

	w, h := 0, 0
	for _, line := range e.lines {
		strline := string(line)
		curW, curH := oneLineTextSize(textFace, strline)
		curW += margin
		if curW > w {
			w = curW
		}
		h += curH + margin
	}

	if w < e.Style().Size.Width.Int() {
		w = e.Style().Size.Width.Int()
	}

	if h < e.Style().Size.Height.Int() {
		h = e.Style().Size.Height.Int()
	}

	w += margin * 2
	h += margin * 2

	e.width, e.height = w, h
	e.ClipTo(width, height)
}

func (e Note) DrawWidget(dst *Graphic) {
	dx, dy := e.WidgetAbsolute()
	margin := e.Style().Margin.Int()
	textFace := e.Style().Font.Face
	lineColorCursor := theme.Cursor.Color.RGBA()
	cursorThick := theme.Cursor.Size.Int()

	FillFrameStyle(dst, dx, dy, e.width, e.height, e.Style())
	sub := GraphicClipStyle(dst, dx, dy, e.width, e.height, e.Style())

	curX := dx + margin
	curY := dy

	for i, line := range e.lines {
		strline := string(line)
		TextDrawOffsetStyle(sub, strline, curX, curY, e.Style())
		_, curH := oneLineTextSize(textFace, strline)
		if e.active && i == e.cursor.Y {
			cut := string(line[0:e.cursor.X])
			if e.cursor.X >= len(strline) {
				cut = strline
			}
			curW, curH := oneLineTextSize(textFace, cut)
			curW += margin
			StrokeLine(sub, curX+curW+cursorThick, curY, 0, curH, cursorThick, lineColorCursor)
			// Draw input method candidate if available
			if e.lastInputState.Text != "" {
				// Draw candidate text
				TextDrawOffsetStyle(dst, e.lastInputState.Text, curX, curY, e.Style())
				// Underline it
				cw, ch := oneLineTextSize(e.Style().Font.Face, e.lastInputState.Text)
				StrokeLine(dst, curX+curW, curY+ch*3/4, cw, 0, cursorThick, lineColorCursor)
			}
		}
		curY += curH + margin
	}

	e.DrawDebug(dst, "NOT")
}

// BUG: Note works but in a vertical box the width overflows.
func NewNote() *Note {
	return newNote()
}

func (c *Note) ReadOnly() bool {
	return c.readonly // cached value
}

func (c *Note) SetReadOnly(readonly bool) {
	c.readonly = readonly
}

func (c *Note) Placeholder() string {
	return c.placeholder
}

func (c *Note) SetPlaceholder(text string) {
	c.placeholder = text
}

func (n *Note) HandleTextInputState(state TextInputState) {
	if state.Committed {
		runes := []rune(state.Text)
		if len(runes) > 0 && runes[0] != 127 { // ebiten inserts DEL characters sfor some reason.
			n.insertRunes(runes)
		}
	} else {
		n.lastInputState = state
	}
}

func (n *Note) startTextInput() {
	// start ime if not started/active yet
	if n.closeInputState == nil {
		absx, abxy := n.WidgetAbsolute()
		n.textInputState, n.closeInputState = StartTextInput(absx, abxy+n.height)
	}
}

func (n *Note) doTextInput() {
	if n.closeInputState == nil {
		absx, abxy := n.WidgetAbsolute()
		n.textInputState, n.closeInputState = StartTextInput(absx, abxy+n.height)
		if n.textInputState == nil {
			return
		}
	}
	for {
		select {
		case state, ok := <-n.textInputState:
			if ok {
				// channel was not closed and input available.
				n.HandleTextInputState(state)
			} else {
				// If the channel was closed we don't have to call the closer anymore.
				n.textInputState = nil
				n.closeInputState = nil
			}
			return
		default:
			// no input available.
			return
		}
	}
}

func (n *Note) closeTextInput() {
	if n.closeInputState != nil {
		n.closeInputState()
		n.closeInputState = nil
	}
	if n.textInputState != nil {
		n.textInputState = nil
	}
	n.lastInputState.Text = ""
}

func (e *Note) insertRunes(runes []rune) {
	if len(e.lines) <= e.cursor.Y {
		e.lines = append(e.lines, []rune{})
	}

	if e.cursor.X >= len(e.lines[e.cursor.Y]) {
		e.lines[e.cursor.Y] = append(e.lines[e.cursor.Y], runes...)
	} else {
		e.lines[e.cursor.Y] = slices.Insert(e.lines[e.cursor.Y], e.cursor.X, runes...)
	}
	e.cursor.X += len(runes)
}

func (e *Note) insertRunesMultipleLines(runes []rune) {
	nlpos := slices.Index(runes, '\n')
	for nlpos >= 0 && len(runes) > 0 {
		sub := runes[:nlpos]
		e.insertRunes(sub)
		e.cursor.X = 0
		e.cursor.Y++
		runes = runes[nlpos:]
		nlpos = slices.Index(runes, '\n')
	}
	if len(runes) > 0 {
		e.insertRunes(runes)
	}
}

func (e *Note) HandleWidget(ev Event) {
	if _, ok := ev.(*AwayEvent); ok {
		dprintln("Note.HandleWidget deactivate")
		e.active = false
		e.closeTextInput()
		return
	}

	if _, ok := ev.(*MouseClickEvent); ok {
		dprintln("Note.HandleWidget activate")
		e.active = true
		e.startTextInput()
	}
	if e.active {
		e.doTextInput()
		// we have to start the text input every time we have the focus,
		// as it might have been closed automatically.
		if ce, ok := ev.(*CharEvent); ok {
			e.insertRunes(ce.Runes)
		}
		if ke, ok := ev.(*KeyPressEvent); ok {
			e.HandleKeyPress(ke)
		}
		if ke, ok := ev.(*KeyReleaseEvent); ok {
			e.HandleKeyRelease(ke)
		}
	}
}

func (e *Note) setCursor(x, y int) {
	if y >= len(e.lines) {
		y = len(e.lines) - 1
	}

	if y < 0 {
		y = 0
	}

	if x > len(e.lines[y]) {
		x = len(e.lines[y])
	}

	if x < 0 {
		x = 0
	}

	e.cursor.X = x
	e.cursor.Y = y
}

func (e *Note) HandleKeyPress(kp *KeyPressEvent) {
	if e.lastInputState.Text != "" {
		switch kp.Key {
		case KeyEscape:
			e.active = false
			e.closeTextInput()
		}
		return
	}

	switch kp.Key {
	case KeyArrowLeft:
		e.setCursor(e.cursor.X-1, e.cursor.Y)
	case KeyArrowRight:
		e.setCursor(e.cursor.X+1, e.cursor.Y)
	case KeyArrowUp:
		e.setCursor(e.cursor.X, e.cursor.Y-1)
	case KeyArrowDown:
		e.setCursor(e.cursor.X, e.cursor.Y+1)
	case KeyHome:
		e.setCursor(0, e.cursor.Y)
	case KeyEnd:
		e.setCursor(len(e.lines[e.cursor.Y]), e.cursor.Y)
	case KeyDelete:
		if e.cursor.X < len(e.lines[e.cursor.Y]) {
			e.lines[e.cursor.Y] = slices.Delete(e.lines[e.cursor.Y], e.cursor.X, e.cursor.X+1)
		} else if e.cursor.Y+1 < len(e.lines) {
			line := e.lines[e.cursor.Y+1]
			e.lines[e.cursor.Y] = append(e.lines[e.cursor.Y], line...)
			e.lines = slices.Delete(e.lines, e.cursor.Y+1, e.cursor.Y+2)
			if len(e.lines) < e.cursor.Y {
				e.lines = append(e.lines, []rune{})
			}
		}
	case KeyBackspace:
		if e.cursor.X > 0 {
			e.lines[e.cursor.Y] = slices.Delete(e.lines[e.cursor.Y], e.cursor.X-1, e.cursor.X)
			if len(e.lines) < e.cursor.Y {
				e.lines = append(e.lines, []rune{})
			}
			e.setCursor(e.cursor.X-1, e.cursor.Y)
		} else if e.cursor.Y > 0 {
			line := e.lines[e.cursor.Y]
			pos := len(e.lines[e.cursor.Y-1])
			e.lines[e.cursor.Y-1] = append(e.lines[e.cursor.Y-1], line...)
			e.lines = slices.Delete(e.lines, e.cursor.Y, e.cursor.Y+1)
			if len(e.lines) < e.cursor.Y {
				e.lines = append(e.lines, []rune{})
			}
			e.setCursor(pos, e.cursor.Y-1)
		}
	case KeyF7:
		e.startTextInput()
	case KeyF8:
		e.closeTextInput()
	case KeyEnter:
		// Insert empty line slices
		if len(e.lines) <= e.cursor.Y {
			e.lines = append(e.lines, []rune{})
		}
		if e.cursor.X >= len(e.lines[e.cursor.Y]) {
			e.lines = slices.Insert(e.lines, e.cursor.Y+1, []rune{})
			e.cursor.X = 0
			e.cursor.Y++
		} else {
			before := slices.Clone(e.lines[e.cursor.Y][:e.cursor.X])
			after := slices.Clone(e.lines[e.cursor.Y][e.cursor.X:])
			e.lines[e.cursor.Y] = before
			e.lines = slices.Insert(e.lines, e.cursor.Y+1, after)
			e.cursor.X = 0
			e.cursor.Y++
		}
		if e.onChanged != nil {
			e.onChanged(e)
		}
	case KeyEscape:
		e.active = false
		e.closeTextInput()
	case KeyC:
		if kp.Modifiers().Control {
			CopyToClipboard(ClipboardFormatText, []byte(string(e.Text())))
		}
	case KeyV:
		if kp.Modifiers().Control {
			buf := CopyFromClipboard(ClipboardFormatText)
			e.insertRunesMultipleLines([]rune(string(buf)))
		}
	}
}

func (e *Note) HandleKeyRelease(kp *KeyReleaseEvent) {

}

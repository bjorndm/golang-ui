package ui

import "strings"

import "golang.org/x/exp/slices"

// Journal allows displaying multiple lines of plain text.
// It is intended for displaying log traces.
// The log is normally displayed in reversed order (latest on top).
// The new line character is \n, nor \r or \r\n on all platforms.
type Journal struct {
	BasicWidget
	active      bool
	onChanged   func(*Journal)
	readonly    bool
	placeholder string
	lines       [][]rune // stored here as a silce or une slice lines
	cursor      TextCursor
	reversed    bool
}

func (e *Journal) OnChanged(cb func(e *Journal)) {
	e.onChanged = cb
}

func newJournal(reversed bool) *Journal {
	e := &Journal{}
	e.lines = [][]rune{[]rune{}}
	e.SetStyle(theme.Journal)
	e.reversed = reversed
	return e
}

func (e *Journal) SetText(text string) {
	lines := strings.Split(text, "\n")
	e.lines = make([][]rune, len(lines)+1)
	for i, line := range lines {
		e.lines[i] = []rune(line)
	}
}

func (e *Journal) Append(text string) {
	if e.lines == nil {
		e.SetText(text)
		return
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		e.lines = append(e.lines, []rune(line))
	}
}

func (e *Journal) Text() string {
	res := ""
	sep := ""
	for _, line := range e.lines {
		res += sep + string(line)
		sep = "\n"
	}
	return res
}

func (e *Journal) LayoutWidget(width, height int) {
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

	e.width, e.height = w, h
	e.ClipTo(width, height)
}

func (e Journal) DrawWidget(dst *Graphic) {
	dx, dy := e.WidgetAbsolute()
	margin := e.Style().Margin.Int()
	textFace := e.Style().Font.Face

	FillFrameStyle(dst, dx, dy, e.width, e.height, e.Style())
	sub := GraphicClipStyle(dst, dx, dy, e.width, e.height, e.Style())

	curX := dx + margin
	curY := dy
	if e.reversed {
		for i := len(e.lines) - 1; i >= 0; i-- {
			line := e.lines[i]
			strline := string(line)
			TextDrawOffsetStyle(sub, strline, curX, curY, e.Style())
			_, curH := oneLineTextSize(textFace, strline)
			curY += curH
		}
	} else {
		for i := 0; i < len(e.lines); i++ {
			line := e.lines[i]
			strline := string(line)
			TextDrawOffsetStyle(sub, strline, curX, curY, e.Style())
			_, curH := oneLineTextSize(textFace, strline)
			curY += curH
		}
	}

	e.DrawDebug(dst, "JOU")
}

func NewJournal(reversed bool) *Journal {
	return newJournal(reversed)
}

func (e *Journal) HandleWidget(ev Event) {
	if _, ok := ev.(*AwayEvent); ok {
		dprintln("Journal.HandleWidget  deactivate")
		e.active = false
		e.RaiseWidget(-10)
		return
	}
	if _, ok := ev.(*MouseClickEvent); ok {
		dprintln("Journal.HandleWidget activate")
		e.active = true
	}
	if e.active {
		if ke, ok := ev.(*KeyPressEvent); ok {
			e.HandleKeyPress(ke)
		}
		if ke, ok := ev.(*KeyReleaseEvent); ok {
			e.HandleKeyRelease(ke)
		}
	}
}

func (e *Journal) setCursor(x, y int) {
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
}

func (e *Journal) HandleKeyPress(kp *KeyPressEvent) {
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
	case KeyEscape:
		e.active = false
	case KeyC:
		if kp.Modifiers().Control {
			CopyToClipboard(ClipboardFormatText, []byte(string(e.Text())))
		}
	}
}

func (e *Journal) HandleKeyRelease(kp *KeyReleaseEvent) {

}

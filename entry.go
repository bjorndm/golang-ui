package ui

import "strings"

import "golang.org/x/exp/slices"

type entryKind int

const (
	entryKindNormal entryKind = iota + 1
	entryKindPassword
	entryKindSearch
)

type Entry struct {
	TextWidget
	onChanged   func(*Entry)
	readonly    bool
	placeholder string
	entryKind
	input           []rune // stored here as runes
	cursor          int
	textInputState  chan TextInputState
	closeInputState func()
	lastInputState  TextInputState
	active          bool
}

type PasswordEntry struct {
	Entry
}

type SearchEntry struct {
	Entry
}

func (e *Entry) OnChanged(cb func(e *Entry)) {
	e.onChanged = cb
}

const entryMaxWidth = 200.0
const entryMinWidth = 100.0

func newEntry(kind entryKind) *Entry {
	e := &Entry{}
	e.entryKind = kind

	// Set placeholder
	if kind == entryKindNormal {
		e.placeholder = "placeholder"
	}
	e.SetStyle(theme.Entry)
	return e
}

const entryHeight = 14
const entryWidth = 64

func (e *Entry) SetText(text string) {
	e.TextWidget.SetText(text)
	e.input = []rune(text)
	e.cursor = len(e.input)
	if e.entryKind == entryKindPassword {
		e.text = strings.Repeat("*", len(e.input))
	}
}

// setRunes does not adjust the cursor.
func (e *Entry) setRunes(runes []rune) {
	e.input = runes
	e.text = string(runes)
	if e.entryKind == entryKindPassword {
		e.text = strings.Repeat("*", len(e.input))
	}
	if e.entryKind == entryKindSearch && e.onChanged != nil {
		e.onChanged(e)
	}
}

func (e *Entry) LayoutWidget(width, height int) {
	e.width, e.height = oneLineTextSize(e.Style().Font.Face, e.text)
	minh := e.Style().Font.Face.Metrics().Height.Round()
	if e.height < minh {
		e.height = minh
	}
	margin := e.Style().Margin.Int()
	e.GrowToStyleSize()
	e.width += margin * 2
	e.height += margin * 2

	e.ClipTo(width, height)
}

func (e Entry) DrawWidget(dst *Graphic) {
	dx, dy := e.WidgetAbsolute()
	margin := e.Style().Margin.Int()
	textFace := e.Style().Font.Face
	cursorThick := theme.Cursor.Size.Int()
	lineColorCursor := theme.Cursor.Color.RGBA()

	FillFrameStyle(dst, dx, dy, e.width, e.height, e.Style())
	sub := GraphicClipStyle(dst, dx, dy, e.width, e.height, e.Style())

	TextDrawOffsetStyle(sub, e.text, dx, dy, e.Style())
	// draw cursor if focused
	if e.active {
		cut := string(e.input[:e.cursor])
		if e.entryKind == entryKindPassword {
			cut = strings.Repeat("*", e.cursor) // XXX this allocates too much I think.
		}
		curW, curH := oneLineTextSize(textFace, cut)
		curX := dx + curW + margin
		curY := dy + margin

		StrokeLine(sub, curX+cursorThick, curY, 0, curH, cursorThick, lineColorCursor)
		// Draw input method candidate if available
		if e.lastInputState.Text != "" {
			// Draw candidate text
			TextDrawOffsetStyle(dst, e.lastInputState.Text, curX, curY, e.Style())
			// Underline it
			cw, ch := oneLineTextSize(e.Style().Font.Face, e.lastInputState.Text)
			StrokeLine(dst, curX, curY+ch*3/4, cw, 0, cursorThick, lineColorCursor)
		}
	}
	e.DrawDebug(dst, "ENT")
}

// BUG: Entry works but in a vertical box the width overflows.
func NewEntry() *Entry {
	return newEntry(entryKindNormal)
}

func NewPasswordEntry() *Entry {
	return newEntry(entryKindPassword)
}

func (e Entry) Text() string {
	if e.entryKind == entryKindPassword {
		return string(e.input)
	}
	return e.text
}

func NewSearchEntry() *Entry {
	return newEntry(entryKindSearch)
}

func (c *Entry) ReadOnly() bool {
	return c.readonly // cached value
}

func (c *Entry) SetReadOnly(readonly bool) {
	c.readonly = readonly
}

func (c *Entry) Placeholder() string {
	return c.placeholder
}

func (c *Entry) SetPlaceholder(text string) {
	c.placeholder = text
}

func (e *Entry) HandleTextInputState(state TextInputState) {
	if state.Committed {
		runes := []rune(state.Text)
		if len(runes) > 0 && runes[0] != 127 { // ebiten inserts DEL characters sfor some reason.
			e.insertRunes(runes)
		}
	} else {
		e.lastInputState = state
	}
}

func (e *Entry) startTextInput() {
	// start ime if not started/active yet
	if e.closeInputState == nil {
		absx, abxy := e.WidgetAbsolute()
		e.textInputState, e.closeInputState = StartTextInput(absx, abxy+e.height)
	}
}

func (e *Entry) doTextInput() {
	if e.closeInputState == nil {
		absx, abxy := e.WidgetAbsolute()
		e.textInputState, e.closeInputState = StartTextInput(absx, abxy+e.height)
		if e.textInputState == nil {
			return
		}
	}
	for {
		select {
		case state, ok := <-e.textInputState:
			if ok {
				// channel was not closed and input available.
				e.HandleTextInputState(state)
			} else {
				// If the channel was closed we don't have to call the closer anymore.
				e.textInputState = nil
				e.closeInputState = nil
			}
			return
		default:
			// no input available.
			return
		}
	}
}

func (e *Entry) closeTextInput() {
	if e.closeInputState != nil {
		e.closeInputState()
		e.closeInputState = nil
	}
	if e.textInputState != nil {
		e.textInputState = nil
	}
	e.lastInputState.Text = ""
}

func (e *Entry) insertRunes(runes []rune) {
	e.input = slices.Insert(e.input, e.cursor, runes...)
	e.text = string(e.input)
	e.setCursor(e.cursor + len(runes))
	if e.entryKind == entryKindPassword {
		e.text = strings.Repeat("*", len(e.input))
	}
	if e.entryKind == entryKindSearch && e.onChanged != nil {
		e.onChanged(e)
	}
}

func (e *Entry) HandleWidget(ev Event) {
	// If we get an away event, deactivate
	if _, ok := ev.(*AwayEvent); ok {
		dprintln("Entry.HandleWidget deactivate")
		e.active = false
		e.closeTextInput()
		return
	}

	if _, ok := ev.(*MouseClickEvent); ok {
		dprintln("Entry.HandleWidget activate")
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

func (e *Entry) setCursor(c int) {
	if c < 0 {
		c = 0
	}
	if c > len(e.input) {
		c = len(e.input)
	}
	e.cursor = c
}

func (e *Entry) HandleKeyPress(kp *KeyPressEvent) {
	switch kp.Key {
	case KeyArrowLeft:
		e.setCursor(e.cursor - 1)
	case KeyArrowRight:
		e.setCursor(e.cursor + 1)
	case KeyHome:
		e.setCursor(0)
	case KeyEnd:
		e.setCursor(len(e.input))
	case KeyDelete:
		if len(e.input) >= e.cursor+1 {
			e.setRunes(append(e.input[0:e.cursor], e.input[e.cursor+1:]...))
		}
	case KeyBackspace:
		if e.cursor > 0 {
			e.setRunes(append(e.input[0:e.cursor-1], e.input[e.cursor:]...))
			e.setCursor(e.cursor - 1)
		}
	case KeyF7:
		e.startTextInput()
	case KeyF8:
		e.closeTextInput()
	case KeyEnter:
		if e.onChanged != nil {
			e.onChanged(e)
		}
		e.active = false
	case KeyC:
		if kp.Modifiers().Control {
			CopyToClipboard(ClipboardFormatText, []byte(string(e.input)))
		}
	case KeyV:
		if kp.Modifiers().Control {
			buf := CopyFromClipboard(ClipboardFormatText)
			e.insertRunes([]rune(string(buf)))
		}
	}
}

func (e *Entry) HandleKeyRelease(kp *KeyReleaseEvent) {

}

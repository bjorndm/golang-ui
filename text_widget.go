package ui

// TextWidget is a common parent widget for widgets that have a text.
type TextWidget struct {
	BasicWidget
	text string
}

func (t TextWidget) Text() string {
	return t.text
}

func (t *TextWidget) SetText(text string) {
	t.text = text
	minw := t.Style().Size.Width.Int()
	minh := t.Style().Size.Height.Int()
	t.LayoutWidget(minh, minw)
}

func NewTextWidget(text string) *TextWidget {
	b := &TextWidget{}
	b.SetText(text)
	return b
}

func (t *TextWidget) LayoutWidget(parentWidth, parentHeight int) {
	t.width, t.height = multiLineTextSize(t.Style().Font.Face, t.text)
	minh := t.Style().Font.Face.Metrics().Height.Round()
	if t.height < minh {
		t.height = minh
	}

	if t.Style().Layout == StyleLayoutStretch {
		t.width = parentWidth
	}

	dprintln("TextWidget.LayoutWidget: ", t.width, t.height)
}

func (t TextWidget) DrawWidget(dst *Graphic) {
	dx, dy := t.WidgetAbsolute()

	// dx, dy, _, _ = t.Style().ApplyMarginPadding(&t, dx, dy)
	face := t.Style().Font.Face
	col := t.Style().Color.RGBA()
	if !t.Enabled() {
		col = theme.Disable.Color.RGBA()
	}

	if t.Style().Align == StyleAlignMiddle {
		dx = dx + t.width/2
	} else if t.Style().Align == StyleAlignRight {
		dx = dx + t.width
	}

	TextDrawOffset(dst, t.text, face, dx, dy, col)
	t.DrawDebug(dst, "TXT")
}

// IconTextWidget is a common parent widget for widgets that have a text
// and an icon in front of the text.
type IconTextWidget struct {
	TextWidget
	icon string
}

func NewIconTextWidget(icon, text string) *IconTextWidget {
	i := &IconTextWidget{}
	i.SetText(text)
	i.SetIcon(icon)
	return i
}

func (t IconTextWidget) Icon() string {
	return t.icon
}

func (t *IconTextWidget) SetIcon(icon string) {
	t.icon = icon
	minw := t.Style().Size.Width.Int()
	minh := t.Style().Size.Height.Int()
	t.LayoutWidget(minh, minw)
}

func (t *IconTextWidget) LayoutWidget(parentWidth, parentHeight int) {

	t.TextWidget.LayoutWidget(parentWidth, parentHeight)
	width, height := t.TextWidget.WidgetSize()
	margin := t.Style().Margin.Int()

	// If set, add room for the icon, which will have the same
	// size as the height of the widget, and for padding.
	if t.icon != "" {
		width += t.height + margin
	}
	width += margin * 4

	t.width, t.height = width, height

	dprintln("IconTextWidget.LayoutWidget: ", t.width, t.height)
}

func (t IconTextWidget) DrawWidget(dst *Graphic) {
	dx, dy := t.WidgetAbsolute()

	face := t.Style().Font.Face
	col := t.Style().Color.RGBA()
	margin := t.Style().Margin.Int()

	if !t.Enabled() {
		col = theme.Disable.Color.RGBA()
	}

	dxt := dx
	if t.icon != "" {
		iconAtlas.DrawSprite(dst, dx, dy, t.height, t.height, t.icon)
		dxt += t.height + margin
	}

	TextDrawOffset(dst, t.text, face, dxt, dy, col)
	t.DrawDebug(dst, "TXT")
}

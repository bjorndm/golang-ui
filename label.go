package ui

import "github.com/hajimehoshi/ebiten/v2/text"
import "golang.org/x/image/font"
import "strings"

type Label struct {
	BasicWidget
	text string
}

type labelKind int

const labelMaxWidth = 200.0
const labelMinWidth = 50.0

func NewLabel(text string) *Label {
	l := &Label{}
	l.text = text
	l.SetStyle(theme.Label)
	return l
}

func (l *Label) Text() string {
	return l.text
}

func (l *Label) SetText(text string) {
	l.text = text
}

func oneLineTextSize(face Face, text string) (width, height int) {
	height = face.Metrics().Height.Round()
	bounds := font.MeasureString(face, text)
	lineWidth := bounds.Round()
	return lineWidth, height
}

func multiLineTextSize(face Face, text string) (width, height int) {
	lines := strings.Split(text, "\n")
	height = face.Metrics().Height.Round() * len(lines)
	width = 0
	for _, line := range lines {
		bounds, _ := font.BoundString(face, line)
		lineWidth := FixedWidth(bounds)
		if lineWidth > width {
			width = lineWidth
		}
	}
	dprintln("multiLineTextSize: ", width, height)
	return width, height
}

func (l *Label) LayoutWidget(width, height int) {
	textFace := l.Style().Font.Face
	margin := l.Style().Margin.Int()

	l.width, l.height = multiLineTextSize(textFace, l.text)
	fh := textFace.Metrics().Height.Round()
	if l.width < labelMinWidth {
		l.width = labelMinWidth
	}
	if l.height < fh {
		l.height = fh
	}

	l.width += int(2 * margin)
	l.height += int(2 * margin)
	l.ClipTo(width, height)
}

func (l Label) DrawWidget(dst *Graphic) {
	dx, dy := l.WidgetAbsolute()

	textFace := l.Style().Font.Face
	textColor := l.Style().Color.RGBA()
	widgetMargin := l.Style().Margin.Int()

	dx += widgetMargin
	dy += widgetMargin + textFace.Metrics().Ascent.Round()

	text.Draw(dst, l.text, textFace, dx, dy, textColor)
	l.DrawDebug(dst, "LAB")
}

var _ Control = &Label{}

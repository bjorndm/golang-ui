package ui

import "github.com/hajimehoshi/ebiten/v2/text"

type Button struct {
	TextWidget
	onClicked func(*Button)
	pressed   bool
	icon      string
}

func (b *Button) OnClicked(f func(*Button)) {
	b.onClicked = f
}

func NewButton(text string) *Button {
	b := &Button{}
	b.SetText(text)
	b.customStyle = theme.Button
	return b
}

func NewButtonWithIcon(text, icon string) *Button {
	b := NewButton(text)
	b.icon = icon
	return b
}

func (b *Button) LayoutWidget(width, height int) {
	textFace := b.Style().Font.Face
	margin := b.Style().Margin.Int()

	b.width, b.height = multiLineTextSize(textFace, b.text)

	if b.icon != "" {
		b.width += b.height + margin
	}

	b.GrowToStyleSize()
	b.width += 2 * margin
	b.height += 2 * margin
	b.ClipTo(width, height)
}

func (b Button) DrawWidget(dst *Graphic) {
	dx, dy := b.WidgetAbsolute()

	textFace := b.Style().Font.Face
	textColor := b.Style().Color.RGBA()
	margin := b.Style().Margin.Int()
	style := b.Style()

	if b.pressed {
		dx += margin / 2
		dy += margin / 2
		active := *theme.Active
		// only use the color, not the sprite
		style.Fill.Color = active.Fill.Color
	}

	FillFrameStyle(dst, dx, dy, b.width, b.height, style)

	ix, iy := dx+b.width-margin-b.height, dy
	if b.text != "" {
		dx += margin
		dy += margin
		dy += textFace.Metrics().Ascent.Round()
		text.Draw(dst, b.text, textFace, dx, dy, textColor)
	}

	if b.icon != "" {
		iconAtlas.DrawSprite(dst, ix, iy, b.height, b.height, b.icon)
	}

	b.DrawDebug(dst, "BUT")
}

func (b *Button) HandleWidget(ev Event) {
	if _, ok := ev.(*MouseClickEvent); ok {
		dprintln("Button.HandleWidget: ")
		if b.onClicked != nil {
			b.onClicked(b)
		}
		b.pressed = true
	}
	if _, ok := ev.(*MouseReleaseEvent); ok {
		dprintln("Button.HandleWidget: ")
		b.pressed = false
	}
	if kr, ok := ev.(*KeyPressEvent); ok {
		if kr.Key != KeySpace {
			return
		}
		dprintln("Box.HandleWidget: key release on focused button", kr.Name(), kr.Key)
		if b.onClicked != nil {
			b.onClicked(b)
		}
		b.pressed = true
	}
	if kr, ok := ev.(*KeyReleaseEvent); ok {
		if kr.Key != KeySpace {
			return
		}
		dprintln("Box.HandleWidget: key release on focused button", kr.Name(), kr.Key)
		b.pressed = false
	}
}

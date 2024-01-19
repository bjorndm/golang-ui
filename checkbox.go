package ui

type Checkbox struct {
	TextWidget
	onClicked func(*Checkbox)
	checked   bool
}

func (b *Checkbox) Checked() bool {
	return b.checked
}

func (b *Checkbox) SetChecked(checked bool) {
	b.checked = checked
}

func (b *Checkbox) OnClicked(f func(*Checkbox)) {
	b.onClicked = f
}

func NewCheckbox(text string) *Checkbox {
	b := &Checkbox{}
	b.SetText(text)
	b.SetChecked(false)
	b.customStyle = theme.Checkbox
	return b
}

func (b *Checkbox) LayoutWidget(width, height int) {
	textFace := b.Style().Font.Face
	margin := b.Style().Margin.Int()
	checkboxWidth := b.Style().Size.Width.Int()
	checkboxHeight := b.Style().Size.Height.Int()

	b.width, b.height = multiLineTextSize(textFace, b.text)
	b.GrowToStyleSize()
	if b.height < checkboxHeight {
		b.height = checkboxHeight
	}

	b.width += int(4*margin) + checkboxWidth
	b.height += int(2 * margin)
	b.ClipTo(width, height)
}

func (b Checkbox) DrawWidget(dst *Graphic) {
	dx, dy := b.WidgetAbsolute()

	style := b.Style()

	margin := style.Margin.Int()
	checkboxWidth := b.Style().Size.Width.Int()
	checkboxHeight := b.Style().Size.Height.Int()

	dx += margin
	dy += margin
	FillFrameStyle(dst, dx, dy, checkboxWidth, checkboxHeight, style)
	if b.checked {
		iconAtlas.DrawSprite(dst, dx, dy, checkboxWidth, checkboxHeight, theme.Icons.Check.String())
	}

	dx += checkboxWidth + margin
	dy += checkboxHeight
	TextDrawStyle(dst, b.text, dx, dy, style)
	b.DrawDebug(dst, "CHE")
}

func (b *Checkbox) HandleWidget(ev Event) {
	if _, ok := ev.(*MouseClickEvent); ok {
		dprintln("Checkbox.HandleWidget: ")
		b.checked = !b.checked
		if b.onClicked != nil {
			b.onClicked(b)
		}
	}
	if kr, ok := ev.(*KeyPressEvent); ok {
		dprintln("Box.HandleWidget: key release on focused button", kr.Name(), kr.Key)
		if kr.Key != KeySpace {
			return
		}
		b.checked = !b.checked
		if b.onClicked != nil {
			b.onClicked(b)
		}
	}
}

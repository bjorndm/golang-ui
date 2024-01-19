package ui

import "golang.org/x/exp/slices"

// Radio a group of linked radio buttons, named Toggle, of which only one can be turned on at the same time.
// Radio is not the parent of the Toggles.
// The Toggles of a Radio need not be placed in the same place in the UI hierarchy so they can be laid out freely.
type Radio struct {
	toggles    []*Toggle
	selected   int
	onSelected func(*Radio)
}

func NewRadio() *Radio {
	r := &Radio{selected: -1}
	return r
}

func (r Radio) Selected() int {
	return r.selected
}

func (r *Radio) OnSelected(cb func(*Radio)) {
	r.onSelected = cb
}

func (r *Radio) SetSelected(index int) {
	if index < 0 || index >= len(r.toggles) {
		return
	}
	r.selected = index
	for i, toggle := range r.toggles {
		toggle.SetChecked(i == r.selected)
	}
	if r.onSelected != nil {
		r.onSelected(r)
	}
}

func (r *Radio) SetSelectedToggle(toggle *Toggle) {
	idx := slices.Index(r.toggles, toggle)
	r.SetSelected(idx)
}

func (r *Radio) AppendToggle(toggle *Toggle) {
	r.toggles = append(r.toggles, toggle)
}

// Append adds a new toggle for the Radio.
// Don't forget to also add it to a container widget to display it in the correct place!
func (r *Radio) Append(text string) *Toggle {
	toggle := NewToggle(text, r)
	r.AppendToggle(toggle)
	return toggle
}

func (r *Radio) Clear() {
	r.toggles = []*Toggle{}
}

func (r *Radio) Delete(column int) {
	r.toggles = slices.Delete(r.toggles, column, column+1)
}

func (r *Radio) InsertAt(text string, index int) *Toggle {
	toggle := NewToggle(text, r)
	r.toggles = slices.Insert(r.toggles, index, toggle)
	return toggle
}

func (r *Radio) NumItems() int {
	return len(r.toggles)
}

func (r *Radio) Text() string {
	if r.selected < 0 {
		return ""
	}
	if widget := r.toggles[r.selected]; widget != nil {
		return widget.Text()
	}
	return ""
}

func (r *Radio) SetText(text string) {
	for i, toggle := range r.toggles {
		if toggle != nil && toggle.Text() == text {
			r.SetSelected(i)
			return
		}
	}
}

// Toggle is a button that gets automatically unchecked if
// another linked to the same Radio is checked.
type Toggle struct {
	TextWidget
	onClicked func(*Toggle)
	checked   bool
	radio     *Radio
}

func (b *Toggle) Checked() bool {
	return b.checked
}

func (b *Toggle) SetChecked(checked bool) {
	b.checked = checked
}

func (b *Toggle) OnClicked(f func(*Toggle)) {
	b.onClicked = f
}

func NewToggle(text string, radio *Radio) *Toggle {
	b := &Toggle{radio: radio}
	b.SetText(text)
	b.SetChecked(false)
	b.customStyle = theme.Radio
	return b
}

func (b *Toggle) LayoutWidget(width, height int) {
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

func (b Toggle) DrawWidget(dst *Graphic) {
	dx, dy := b.WidgetAbsolute()

	style := b.Style()

	margin := style.Margin.Int()
	checkboxWidth := b.Style().Size.Width.Int()
	checkboxHeight := b.Style().Size.Height.Int()

	dx += margin
	dy += margin
	FillFrameStyle(dst, dx, dy, checkboxWidth, checkboxHeight, style)
	if b.checked {
		iconAtlas.DrawSprite(dst, dx, dy, checkboxWidth, checkboxHeight, theme.Icons.Radio.String())
	}

	dx += checkboxWidth + margin
	dy += checkboxHeight
	TextDrawStyle(dst, b.text, dx, dy, style)
	b.DrawDebug(dst, "CHE")
}

func (b *Toggle) toggle() {
	b.checked = true
	if b.onClicked != nil {
		b.onClicked(b)
	}
	b.radio.SetSelectedToggle(b)
}

func (b *Toggle) HandleWidget(ev Event) {
	if _, ok := ev.(*MouseClickEvent); ok {
		dprintln("Toggle.HandleWidget: ")
		b.toggle()
	}
	if kr, ok := ev.(*KeyPressEvent); ok {
		dprintln("Box.HandleWidget: key release on focused button", kr.Name(), kr.Key)
		if kr.Key != KeySpace {
			return
		}
		b.toggle()
	}
}

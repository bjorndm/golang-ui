package ui

import "image"
import "image/png"
import _ "image/jpeg"
import _ "image/gif"
import "bytes"

// Picture is a widget that displays an image. If the picture is editable,
// then it can be copy/pasted in or loaded from a file.
type Picture struct {
	BasicWidget
	Image
	graphic    *Graphic
	title      TextWidget
	icon       string
	borderless bool
	editable   bool
	active     bool
}

func NewPicture(title string, img Image) *Picture {
	w := &Picture{}
	w.title = *NewTextWidget(title)
	w.title.SetParent(w)
	w.SetImage(img)
	w.SetStyle(theme.Picture)
	return w
}

func NewPictureWithIcon(title string, icon string) *Picture {
	w := &Picture{}
	w.title = *NewTextWidget(title)
	w.title.SetParent(w)
	w.SetIcon(icon)
	w.SetStyle(theme.Picture)
	return w
}

func (w *Picture) SetTitle(title string) {
	w.title.SetText(title)
}

func (w *Picture) SetIcon(icon string) {
	w.icon = icon
	if w.Image != nil && w.graphic != nil {
		w.graphic.Dispose()
		w.graphic = nil
	}
	if w.icon != "" {
		w.graphic = iconAtlas.FindGraphic(w.icon)
	}
}

func (w *Picture) SetImage(img Image) {
	w.Image = img
	if w.graphic != nil && w.icon == "" {
		w.graphic.Dispose()
		w.graphic = nil
	}
	if w.Image != nil {
		w.graphic = NewGraphicFromImage(w.Image)
	}
}

func (w *Picture) Title() string {
	return w.title.Text()
}

func (g *Picture) LayoutWidget(width, height int) {
	g.width, g.height = 0, 0
	if g.graphic != nil {
		g.width, g.height = g.graphic.Size()
	}
	margin := g.Style().Margin.Int()
	g.title.LayoutWidget(width-margin*2, height-margin*2)
	g.title.MoveWidget(margin, margin)

	tw, th := 0, 0
	if g.title.Text() != "" {
		tw, th = g.title.WidgetSize()
	}

	if g.width < tw {
		g.width = tw
	}
	g.height += th

	g.GrowToStyleSize()
	g.width += margin * 2
	g.height += margin * 2

	if g.width < 1 || g.height < 1 {
		panic("Picture too small")
	}

	g.ClipTo(width, height)
}

func (w *Picture) Destroy() {
	if w.icon == "" {
		// destroy the child graphic if it is not an icon
		if w.graphic != nil {
			w.graphic.Dispose()
		}
	}
}

func (w Picture) Borderless() bool {
	return w.borderless
}

func (w *Picture) SetBorderless(borderless bool) {
	w.borderless = borderless
}

func (w *Picture) DrawWidget(screen *Graphic) {
	dx, dy := w.WidgetAbsolute()

	margin := w.Style().Margin.Int()

	/*if !w.borderless*/
	{
		FillFrameStyle(screen, dx, dy, w.width, w.height, w.Style())
	}

	th := 0
	if w.title.Text() != "" {
		w.title.DrawWidget(screen)
		_, th = w.title.WidgetSize()
	}
	dx += margin
	dy += th + margin

	if w.graphic != nil {
		ww, wh := w.WidgetSize()
		wh -= th
		wh -= 2 * margin
		ww -= 2 * margin
		iw, ih := w.graphic.Size()
		if ih == 0 {
			ih = 1
		}
		if iw == 0 {
			iw = 1
		}
		sx, sy := float64(ww)/float64(iw), float64(wh)/float64(ih)
		DrawGraphicAtScale(screen, w.graphic, dx, dy, sx, sy)
	}

	w.DrawDebug(screen, "PIC %d %d", dx, dy)
}

func (p *Picture) HandleWidget(ev Event) {
	if mc, ok := ev.(*MouseClickEvent); ok {
		if mc.MouseEvent.Inside(p) {
			p.active = true
		} else {
			p.active = false
		}
	}
	if p.active {
		if ke, ok := ev.(*KeyPressEvent); ok {
			p.HandleKeyPress(ke)
		}
	}
}

func (p *Picture) HandleKeyPress(kp *KeyPressEvent) {
	// Allow copy pasting of images.
	switch kp.Key {
	case KeyC:
		if kp.Modifiers().Control {
			wr := &bytes.Buffer{}
			err := png.Encode(wr, p.Image)
			if err == nil {
				CopyToClipboard(ClipboardFormatImage, wr.Bytes())
			} else {
				p.SetTitle(p.Title() + err.Error())
			}
		}
	case KeyV:
		if kp.Modifiers().Control {
			buf := CopyFromClipboard(ClipboardFormatImage)
			rd := bytes.NewBuffer(buf)
			img, _, err := image.Decode(rd)
			if err == nil {
				p.SetImage(img)
				p.LayoutWidget(p.WidgetSize())
			} else {
				p.SetTitle(p.Title() + err.Error())
			}
		}
	}
}

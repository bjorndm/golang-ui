package ui

import "image"
import "path"
import _ "image/png"
import "github.com/hajimehoshi/ebiten/v2"

type NineSlice struct {
	Slice  [9]*Graphic `json:"-"`
	Border int         `json:"-"`
}

type AtlasSprite struct {
	Name          string `json:"name"`
	X             int    `json:"x"`
	Y             int    `json:"y"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	Border        int    `json:"border,omitempty"`
	*ebiten.Image `json:"-"`
	NineSlice
}

// Atlas is a collection of AtlasSprites, or sub-images loaded from a single PNG file.
// In order to be able to find the location atlas sprites in that PNG file, a JSON file can be loaded.
type Atlas struct {
	Filename string         `json:"filename"`
	Border   int            `json:"border,omitempty"`
	Sprites  []*AtlasSprite `json:"sprites"`
	Image    *Graphic
	ByName   map[string]*AtlasSprite
}

func (a Atlas) FindSprite(name string) *AtlasSprite {
	if sprite, ok := a.ByName[name]; ok {
		return sprite
	}
	return nil
}

func (a Atlas) FindGraphic(name string) *Graphic {
	sprite := a.FindSprite(name)
	if sprite == nil {
		return nil
	}
	return sprite.Image
}

func (a Atlas) DrawSprite(dst *Graphic, x, y, w, h int, name string) {
	if sub, ok := a.ByName[name]; ok {
		opts := ebiten.DrawImageOptions{}
		sx, sy := sub.Image.Size()
		opts.GeoM.Scale(float64(w)/float64(sx), float64(h)/float64(sy))
		opts.GeoM.Translate(float64(x), float64(y))
		dst.DrawImage(sub.Image, &opts)
	} else {
		FillFrame(dst, x, y, w, h, 1, theme.Error.Fill.Color.RGBA(), theme.Error.Fill.Color.RGBA())
		dprintln("Warning: no such sprite", name)
	}
	DrawDebug(dst, x, y, w, h, "ASP")
}

func (a Atlas) DrawColoredSprite(dst *Graphic, x, y, w, h int, name string, col Color) {
	if sprite, ok := a.ByName[name]; ok {
		if !sprite.NineSlice.OK() {
			DrawSpriteAtScaleColor(dst, sprite.Image, x, y, w, h, col)
		} else {
			sprite.NineSlice.Draw(dst, x, y, w, h, col)
		}
	} else {
		TextDrawOffsetStyle(dst, name, x, y, *theme.Error)
		dprintln("Warning: no such sprite", name)
	}
	DrawDebug(dst, x, y, w, h, "ASP")
}

func (a Atlas) DrawColoredSprite2(dst *Graphic, x, y, w, h int, name string, col Color) {
	if sprite, ok := a.ByName[name]; ok {
		DrawSpriteAtScaleColor(dst, sprite.Image, x, y, w, h, col)
	} else {
		TextDrawOffsetStyle(dst, name, x, y, *theme.Error)
		dprintln("Warning: no such sprite", name)
	}
	DrawDebug(dst, x, y, w, h, "ASP")
}

func loadAtlas(name string) *Atlas {
	atlas := loadResourceJSON[Atlas](name)
	dir := path.Dir(name)
	imgName := path.Join(dir, atlas.Filename)
	atlas.Image = loadResourceImage(imgName)
	atlas.ByName = make(map[string]*AtlasSprite)
	for _, sprite := range atlas.Sprites {
		key := sprite.Name
		rect := image.Rect(sprite.X, sprite.Y, sprite.X+sprite.Width, sprite.Y+sprite.Height)
		value := atlas.Image.SubImage(rect).(*ebiten.Image)
		sprite.Image = value
		atlas.ByName[key] = sprite
	}
	if atlas.Border > 0 {
		atlas.NineSlice(atlas.Border)
	}

	return atlas
}

func (atlas *Atlas) NineSliceSprite(border int, sprite *AtlasSprite) {
	if sprite.Border < 0 { // don't nineslice.
		return
	}

	if sprite.Border > 0 {
		border = sprite.Border
	}
	xdiff := []int{0, border, sprite.Width - border}
	ydiff := []int{0, border, sprite.Height - border}
	wtab := []int{border, sprite.Width - 2*border, border}
	htab := []int{border, sprite.Height - 2*border, border}

	for i := 0; i < len(sprite.Slice); i++ {
		var rect image.Rectangle
		ix := i % 3
		iy := i / 3
		dx := sprite.X + xdiff[ix]
		dy := sprite.Y + ydiff[iy]
		dw := wtab[ix]
		dh := htab[iy]
		rect = image.Rect(dx, dy, dx+dw, dy+dh)
		slice := atlas.Image.SubImage(rect).(*ebiten.Image)
		sprite.Slice[i] = slice
	}
	sprite.NineSlice.Border = border
	sprite.Border = border
}

func (atlas *Atlas) NineSlice(border int) {
	for _, sprite := range atlas.Sprites {
		atlas.NineSliceSprite(border, sprite)
	}
}

func DrawSpriteAtScaleColor(dst, src *Graphic, x, y, w, h int, col Color) {
	opts := ebiten.DrawImageOptions{}
	sx, sy := src.Size()
	opts.GeoM.Scale(float64(w)/float64(sx), float64(h)/float64(sy))
	opts.GeoM.Translate(float64(x), float64(y))
	r, g, b, a := col.RGBA()
	rf := float32(r) / float32(0xffff)
	gf := float32(g) / float32(0xffff)
	bf := float32(b) / float32(0xffff)
	af := float32(a) / float32(0xffff)
	opts.ColorScale.Scale(rf, gf, bf, af)
	dst.DrawImage(src, &opts)
}

func (n NineSlice) OK() bool {
	return n.Slice[0] != nil
}

func (n NineSlice) Draw(dst *Graphic, x, y, w, h int, col Color) {
	xdiff := []int{0, n.Border, w - n.Border}
	ydiff := []int{0, n.Border, h - n.Border}
	wtab := []int{n.Border, w - 2*n.Border, n.Border}
	htab := []int{n.Border, h - 2*n.Border, n.Border}
	for i := 0; i < len(n.Slice); i++ {
		ix := i % 3
		iy := i / 3
		dx := x + xdiff[ix]
		dy := y + ydiff[iy]
		dw := wtab[ix]
		dh := htab[iy]
		DrawSpriteAtScaleColor(dst, n.Slice[i], dx, dy, dw, dh, col)
	}
}

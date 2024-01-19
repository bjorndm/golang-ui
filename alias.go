package ui

import "context"
import "fmt"
import "image"
import "image/color"

import "golang.org/x/image/font"
import "golang.org/x/image/font/opentype"
import "golang.org/x/image/math/fixed"

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/hajimehoshi/ebiten/v2/exp/textinput"
import "github.com/hajimehoshi/ebiten/v2/text"

// The vector package is convenient but unfortunately quite low in performance.
import "github.com/hajimehoshi/ebiten/v2/vector"

import "golang.design/x/clipboard"

type (
	Color   = color.Color
	Font    = opentype.Font
	Face    = font.Face
	Graphic = ebiten.Image
	RGBA    = color.RGBA
	Key     = ebiten.Key
	Image   = image.Image // Image is an image.Image
)

const (
	KeyA              = ebiten.KeyA
	KeyB              = ebiten.KeyB
	KeyC              = ebiten.KeyC
	KeyD              = ebiten.KeyD
	KeyE              = ebiten.KeyE
	KeyF              = ebiten.KeyF
	KeyG              = ebiten.KeyG
	KeyH              = ebiten.KeyH
	KeyI              = ebiten.KeyI
	KeyJ              = ebiten.KeyJ
	KeyK              = ebiten.KeyK
	KeyL              = ebiten.KeyL
	KeyM              = ebiten.KeyM
	KeyN              = ebiten.KeyN
	KeyO              = ebiten.KeyO
	KeyP              = ebiten.KeyP
	KeyQ              = ebiten.KeyQ
	KeyR              = ebiten.KeyR
	KeyS              = ebiten.KeyS
	KeyT              = ebiten.KeyT
	KeyU              = ebiten.KeyU
	KeyV              = ebiten.KeyV
	KeyW              = ebiten.KeyW
	KeyX              = ebiten.KeyX
	KeyY              = ebiten.KeyY
	KeyZ              = ebiten.KeyZ
	KeyAltLeft        = ebiten.KeyAltLeft
	KeyAltRight       = ebiten.KeyAltRight
	KeyArrowDown      = ebiten.KeyArrowDown
	KeyArrowLeft      = ebiten.KeyArrowLeft
	KeyArrowRight     = ebiten.KeyArrowRight
	KeyArrowUp        = ebiten.KeyArrowUp
	KeyBackquote      = ebiten.KeyBackquote
	KeyBackslash      = ebiten.KeyBackslash
	KeyBackspace      = ebiten.KeyBackspace
	KeyBracketLeft    = ebiten.KeyBracketLeft
	KeyBracketRight   = ebiten.KeyBracketRight
	KeyCapsLock       = ebiten.KeyCapsLock
	KeyComma          = ebiten.KeyComma
	KeyContextMenu    = ebiten.KeyContextMenu
	KeyControlLeft    = ebiten.KeyControlLeft
	KeyControlRight   = ebiten.KeyControlRight
	KeyDelete         = ebiten.KeyDelete
	KeyDigit0         = ebiten.KeyDigit0
	KeyDigit1         = ebiten.KeyDigit1
	KeyDigit2         = ebiten.KeyDigit2
	KeyDigit3         = ebiten.KeyDigit3
	KeyDigit4         = ebiten.KeyDigit4
	KeyDigit5         = ebiten.KeyDigit5
	KeyDigit6         = ebiten.KeyDigit6
	KeyDigit7         = ebiten.KeyDigit7
	KeyDigit8         = ebiten.KeyDigit8
	KeyDigit9         = ebiten.KeyDigit9
	KeyEnd            = ebiten.KeyEnd
	KeyEnter          = ebiten.KeyEnter
	KeyEqual          = ebiten.KeyEqual
	KeyEscape         = ebiten.KeyEscape
	KeyF1             = ebiten.KeyF1
	KeyF2             = ebiten.KeyF2
	KeyF3             = ebiten.KeyF3
	KeyF4             = ebiten.KeyF4
	KeyF5             = ebiten.KeyF5
	KeyF6             = ebiten.KeyF6
	KeyF7             = ebiten.KeyF7
	KeyF8             = ebiten.KeyF8
	KeyF9             = ebiten.KeyF9
	KeyF10            = ebiten.KeyF10
	KeyF11            = ebiten.KeyF11
	KeyF12            = ebiten.KeyF12
	KeyHome           = ebiten.KeyHome
	KeyInsert         = ebiten.KeyInsert
	KeyMetaLeft       = ebiten.KeyMetaLeft
	KeyMetaRight      = ebiten.KeyMetaRight
	KeyMinus          = ebiten.KeyMinus
	KeyNumLock        = ebiten.KeyNumLock
	KeyNumpad0        = ebiten.KeyNumpad0
	KeyNumpad1        = ebiten.KeyNumpad1
	KeyNumpad2        = ebiten.KeyNumpad2
	KeyNumpad3        = ebiten.KeyNumpad3
	KeyNumpad4        = ebiten.KeyNumpad4
	KeyNumpad5        = ebiten.KeyNumpad5
	KeyNumpad6        = ebiten.KeyNumpad6
	KeyNumpad7        = ebiten.KeyNumpad7
	KeyNumpad8        = ebiten.KeyNumpad8
	KeyNumpad9        = ebiten.KeyNumpad9
	KeyNumpadAdd      = ebiten.KeyNumpadAdd
	KeyNumpadDecimal  = ebiten.KeyNumpadDecimal
	KeyNumpadDivide   = ebiten.KeyNumpadDivide
	KeyNumpadEnter    = ebiten.KeyNumpadEnter
	KeyNumpadEqual    = ebiten.KeyNumpadEqual
	KeyNumpadMultiply = ebiten.KeyNumpadMultiply
	KeyNumpadSubtract = ebiten.KeyNumpadSubtract
	KeyPageDown       = ebiten.KeyPageDown
	KeyPageUp         = ebiten.KeyPageUp
	KeyPause          = ebiten.KeyPause
	KeyPeriod         = ebiten.KeyPeriod
	KeyPrintScreen    = ebiten.KeyPrintScreen
	KeyQuote          = ebiten.KeyQuote
	KeyScrollLock     = ebiten.KeyScrollLock
	KeySemicolon      = ebiten.KeySemicolon
	KeyShiftLeft      = ebiten.KeyShiftLeft
	KeyShiftRight     = ebiten.KeyShiftRight
	KeySlash          = ebiten.KeySlash
	KeySpace          = ebiten.KeySpace
	KeyTab            = ebiten.KeyTab
	KeyAlt            = ebiten.KeyAlt
	KeyControl        = ebiten.KeyControl
	KeyShift          = ebiten.KeyShift
	KeyMeta           = ebiten.KeyMeta
	KeyMax            = ebiten.KeyMax
)

func FixedWidth(r fixed.Rectangle26_6) int {
	return r.Max.Sub(r.Min).X.Round()
}

func FixedHeight(r fixed.Rectangle26_6) int {
	return r.Max.Sub(r.Min).Y.Round()
}

func FillRect(g *Graphic, x, y, w, h int, color Color) {
	vector.DrawFilledRect(g, float32(x), float32(y),
		float32(w), float32(h), color, true)
}

func StrokeRect(g *Graphic, x, y, w, h, t int, color Color) {
	vector.StrokeRect(g, float32(x), float32(y),
		float32(w), float32(h), float32(t), color, true)
}

func StrokeLine(dst *Graphic, x, y, w, h, t int, color Color) {
	vector.StrokeLine(dst, float32(x), float32(y),
		float32(x+w), float32(y+h), float32(t), color, true)
}

func StrokeCircle(dst *Graphic, cx, cy, r, t int, color Color) {
	vector.StrokeCircle(dst, float32(cx), float32(cy),
		float32(r), float32(t), color, true)
}

func FillCircle(dst *Graphic, cx, cy, r int, color Color) {
	vector.DrawFilledCircle(dst, float32(cx), float32(cy),
		float32(r), color, true)
}

func FillFrame(g *Graphic, x, y, w, h, thick int, fill, border Color) {
	if thick > 0 {
		// we don't use StrokeRect because it is significantly slower than
		// FillRect. So we draw the border as a solid rectangle and then
		// draw the filling as a smaller rectangle inside of it.
		FillRect(g, x, y, w, h, border)
	}

	FillRect(g, x+thick, y+thick, w-2*thick, h-2*thick, fill)
}

func DrawDebug(dst *Graphic, x, y, w, h int, form string, args ...any) {
	if !debugDisplay {
		return
	}

	r := 0
	g := 0
	b := 0
	a := 255

	if len(form) > 0 {
		r = ('Z' - int(form[0])) * 255 / ('Z' - 'A')
	}
	if len(form) > 1 {
		g = ('Z' - int(form[1])) * 255 / ('Z' - 'A')
	}
	if len(form) > 2 {
		b = ('Z' - int(form[2])) * 255 / ('Z' - 'A')
	}
	fillColor := RGBA{byte(r) / 2, byte(b) / 2, byte(g) / 2, byte(a) / 2}
	lineColor := RGBA{byte(r) * 3 / 4, byte(b) * 3 / 4, byte(g) * 3 / 4, byte(a) * 3 / 4}

	FillFrame(dst, x, y, w, h, 1, fillColor, lineColor)

	y += textFaceDebug.Metrics().Ascent.Round()
	str := fmt.Sprintf(form, args...)
	text.Draw(dst, str, textFaceDebug, x, y, textColorDebug)
}

func TextDraw(dst *Graphic, str string, face Face, x, y int, col Color) {
	text.Draw(dst, str, face, x, y, col)
}

func TextDrawStyle(dst *Graphic, str string, x, y int, style Style) {
	var (
		face = style.Font.Face
		col  = style.Color.RGBA()
	)
	TextDraw(dst, str, face, x, y, col)
}

func TextDrawOffset(dst *Graphic, str string, face Face, x, y int, col Color) {
	y += face.Metrics().Ascent.Round()
	TextDraw(dst, str, face, x, y, col)
}

func TextDrawOffsetStyle(dst *Graphic, str string, x, y int, style Style) {
	var (
		face   = style.Font.Face
		col    = style.Color.RGBA()
		margin = style.Margin.Int()
	)
	x += margin
	y += margin
	TextDrawOffset(dst, str, face, x, y, col)
}

func TextDrawOffsetStyleWithoutMargin(dst *Graphic, str string, x, y int, style Style) {
	var (
		face = style.Font.Face
		col  = style.Color.RGBA()
	)
	TextDrawOffset(dst, str, face, x, y, col)
}

func GraphicClipStyle(dst *Graphic, x, y, w, h int, style Style) *Graphic {
	var (
		margin = style.Margin.Int()
	)
	x += margin
	y += margin
	w -= margin * 2
	h -= margin * 2
	rect := image.Rect(x, y, x+w, y+h)
	return dst.SubImage(rect).(*ebiten.Image)
}

func TextDrawHeight(dst *Graphic, str string, face Face, x, y int, col Color) {
	h := face.Metrics().Height.Round()
	if h < 1 {
		dprintln("face height", h)
		panic("face metrics broken")
	}
	y += h

	TextDraw(dst, str, face, x, y, col)
}

func TextDrawHeightStyle(dst *Graphic, str string, x, y int, style Style) {
	var (
		face   = style.Font.Face
		col    = style.Color.RGBA()
		margin = style.Margin.Int()
	)
	x += margin
	y += margin
	TextDrawHeight(dst, str, face, x, y, col)
}

func FillRectClip(g *Graphic, x, y, w, h int, col Color) {
	rect := image.Rect(x, y, x+w, y+h)
	sub := g.SubImage(rect).(*ebiten.Image)
	sub.Fill(col)
	sub.Dispose()
}

func FillFrameStyle(g *Graphic, x, y, w, h int, style Style) {
	var fill = style.Fill.Color.RGBA()
	var sprite = style.Fill.Sprite.String()

	uiAtlas.DrawColoredSprite(g, x, y, w, h, sprite, fill)
}

func FillFrameOptionalStyle(g *Graphic, x, y, w, h int, style *Style) {
	if style == nil {
		style = &theme.Style
	}
	FillFrameStyle(g, x, y, w, h, *style)
}

func DrawFrameOptionalStyle(g *Graphic, x, y, w, h int, style *Style) {
	FillFrameOptionalStyle(g, x, y, w, h, style)
}

type TextInputState = textinput.State

func StartTextInput(x, y int) (states chan TextInputState, close func()) {
	return textinput.Start(x, y)
}

var clipboardAvailable = false

func initClipBoard() {
	clipboardAvailable = clipboard.Init() == nil
}

type ClipboardFormat = clipboard.Format

const ClipboardFormatText = clipboard.FmtText
const ClipboardFormatImage = clipboard.FmtImage

func CopyFromClipboard(format ClipboardFormat) []byte {
	if !clipboardAvailable {
		return nil
	}
	return clipboard.Read(format)
}

func CopyToClipboard(format ClipboardFormat, buf []byte) <-chan struct{} {
	if !clipboardAvailable {
		return nil
	}
	return clipboard.Write(format, buf)
}

func WatchClipboard(ctx context.Context, format ClipboardFormat) <-chan []byte {
	if !clipboardAvailable {
		return nil
	}
	return clipboard.Watch(ctx, format)
}

func NewGraphicFromImage(img Image) *Graphic {
	return ebiten.NewImageFromImage(img)
}

func NewGraphic(width, height int) *Graphic {
	return ebiten.NewImage(width, height)
}

func DrawGraphicAt(target *Graphic, source *Graphic, x, y int) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(x), float64(y))
	target.DrawImage(source, &opts)
}

func DrawGraphicAtScale(target *Graphic, source *Graphic, x, y int, sx, sy float64) {
	opts := ebiten.DrawImageOptions{}
	opts.GeoM.Scale(sx, sy)
	opts.GeoM.Translate(float64(x), float64(y))
	target.DrawImage(source, &opts)
}

func DrawLineStyle(g *Graphic, x, y, w, h int, style Style) {
	var (
		thickness = 1 // XXX use line sprite in stead
		border    = style.Color.RGBA()
	)
	StrokeLine(g, x, y, w, h, thickness, border)
}

func DrawCircleStyle(g *Graphic, cx, cy, r int, style LineStyle) {
	var (
		thickness = style.Size.Int()
		border    = style.Color.RGBA()
	)
	if thickness < 1 {
		thickness = 1
	}
	StrokeCircle(g, cx, cy, r, thickness, border)
}

func FillCircleStyle(g *Graphic, cx, cy, r int, style Style) {
	var (
		thickness = 1
		border    = style.Color.RGBA()
		fill      = style.Fill.Color.RGBA()
	)
	// XXX: draw a circle sprite in stead.
	FillCircle(g, cx, cy, r, fill)
	if thickness > 0 {
		StrokeCircle(g, cx, cy, r, thickness, border)
	}
}

type CursorShapeType = ebiten.CursorShapeType

const (
	CursorShapeDefault    = (ebiten.CursorShapeDefault)
	CursorShapeText       = (ebiten.CursorShapeText)
	CursorShapeCrosshair  = (ebiten.CursorShapeCrosshair)
	CursorShapePointer    = (ebiten.CursorShapePointer)
	CursorShapeEWResize   = (ebiten.CursorShapeEWResize)
	CursorShapeNSResize   = (ebiten.CursorShapeNSResize)
	CursorShapeNESWResize = (ebiten.CursorShapeNESWResize)
	CursorShapeNWSEResize = (ebiten.CursorShapeNWSEResize)
	CursorShapeMove       = (ebiten.CursorShapeMove)
	CursorShapeNotAllowed = (ebiten.CursorShapeNotAllowed)
)

func CursorShape() CursorShapeType {
	return ebiten.CursorShape()
}

func SetCursorShape(shape ebiten.CursorShapeType) {
	ebiten.SetCursorShape(shape)
}

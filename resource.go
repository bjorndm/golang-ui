package ui

import "io"
import "io/fs"
import "image"
import _ "image/png"
import "encoding/json"
import "embed"
import "golang.org/x/image/font/opentype"
import "github.com/hajimehoshi/ebiten/v2"

//go:embed resource
var resource embed.FS

type OverlayFS []fs.FS

func (o OverlayFS) Open(name string) (fs.File, error) {
	for i := len(o) - 1; i >= 0; i-- {
		sub := o[i]
		f, err := sub.Open(name)
		if f != nil && err == nil {
			return f, nil
		}
	}
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}

var resources OverlayFS

const defaultFontBoldName = "resource/font/GoNotoCurrent-Bold.ttf"
const defaultFontName = "resource/font/GoNotoCurrent-Regular.ttf"
const iconAtlasName = "resource/icon/icon_atlas.json"
const uiAtlasName = "resource/theme/ui_atlas.json"
const themeName = "resource/theme/default_theme.json"

var defaultFontBold *Font
var defaultFont *Font
var iconAtlas *Atlas
var uiAtlas *Atlas

// panics on failure
func loadResourceReader(name string) io.ReadCloser {
	rd, err := resources.Open(name)
	if err != nil {
		panic(err)
	}
	return rd
}

// panics on failure
func loadResourceBuffer(name string) []byte {
	rd := loadResourceReader(name)
	defer rd.Close()
	buf, err := io.ReadAll(rd)
	if err != nil {
		panic(err)
	}
	return buf
}

// returns nil on failure
func loadResourceReaderOptional(name string) io.ReadCloser {
	rd, err := resources.Open(name)
	if err != nil {
		return nil
	}
	return rd
}

// returns nil on failure
func loadResourceBufferOptional(name string) []byte {
	rd := loadResourceReaderOptional(name)
	if rd == nil {
		return nil
	}
	defer rd.Close()
	buf, err := io.ReadAll(rd)
	if err != nil {
		return nil
	}
	return buf
}

func loadResourceImage(name string) *ebiten.Image {
	rd := loadResourceReader(name)
	defer rd.Close()
	img, _, err := image.Decode(rd)
	if err != nil {
		panic(err)
	}
	return ebiten.NewImageFromImage(img)
}

func loadResourceJSON[T any](name string) *T {
	var obj T
	buf := loadResourceBuffer(name)
	err := json.Unmarshal(buf, &obj)
	if err != nil {
		panic(err)
	}
	return &obj
}

// panics on failure
func loadResourceFont(name string) *Font {
	buf := loadResourceBuffer(name)
	fnt, err := opentype.Parse(buf)
	if err != nil {
		panic(err)
	}
	return fnt
}

// returns nil on failure
func loadResourceFontOptional(name string) *Font {
	buf := loadResourceBuffer(name)
	if buf == nil {
		return nil
	}
	fnt, err := opentype.Parse(buf)
	if err != nil {
		return nil
	}
	return fnt
}

const defaultDPI = 90

func fontFace(font *Font, size int) Face {
	options := opentype.FaceOptions{
		Size: float64(size),
		DPI:  defaultDPI,
		// Hinting
	}
	if theme != nil && theme.DPI > 0 {
		options.DPI = theme.DPI.Float()
	}
	face, err := opentype.NewFace(font, &options)
	if err != nil {
		panic(err)
	}
	return face
}

func initResource() {
	resources = OverlayFS{resource}
	defaultFontBold = loadResourceFont(defaultFontBoldName)
	defaultFont = loadResourceFont(defaultFontName)
	textFaceDebug = fontFace(defaultFont, textSizeDebug)
	iconAtlas = loadAtlas(iconAtlasName)
	uiAtlas = loadAtlas(uiAtlasName)
	initTheme()
}

func MountResources(sys fs.FS) {
	resources = append(resources, sys)
}

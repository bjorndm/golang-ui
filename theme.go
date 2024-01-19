package ui

import "encoding/json"
import "fmt"
import "strings"
import "golang.org/x/image/colornames"

var textColorDebug = RGBA{0, 0, 0, 255}
var textSizeDebug = 6
var textFaceDebug Face

type ColorMod func(StyleColor) StyleColor

var ColorMethods StyleColor

func (StyleColor) Transparent(level float64) func(StyleColor) StyleColor {
	return func(s StyleColor) StyleColor {
		s.R = byte(float64(s.R) * level)
		s.G = byte(float64(s.G) * level)
		s.B = byte(float64(s.B) * level)
		s.A = byte(float64(s.A) * level)
		return s
	}
}

func (StyleColor) Light(level float64) func(StyleColor) StyleColor {
	return func(s StyleColor) StyleColor {
		s.R = byte(float64(s.R) * level)
		s.G = byte(float64(s.G) * level)
		s.B = byte(float64(s.B) * level)
		return s
	}
}

var ColorModMap = map[string]ColorMod{
	"crystal":    ColorMethods.Transparent(0.4),
	"diaphane":   ColorMethods.Transparent(0.5),
	"limpid":     ColorMethods.Transparent(0.6),
	"pale":       ColorMethods.Transparent(0.7),
	"sheer":      ColorMethods.Transparent(0.8),
	"translucid": ColorMethods.Transparent(0.9),
	"clear":      ColorMethods.Light(1.3),
	"light":      ColorMethods.Light(1.5),
	"bright":     ColorMethods.Light(1.7),
}

type StyleSprite string

func (s StyleSprite) WithDefault(def StyleSprite) StyleSprite {
	if s == "" {
		return def
	}
	return s
}

func (s StyleSprite) String() string {
	return string(s)
}

type StyleSize int

func (s StyleSize) WithDefault(def StyleSize) StyleSize {
	if s < 0 {
		return def
	}
	return s
}

func (s StyleSize) Int() int {
	return int(s)
}

func (s StyleSize) Float() float64 {
	return float64(s)
}

type StyleColor RGBA

func NewStyleColor(r, g, b, a byte) StyleColor {
	return StyleColor{r, g, b, a}
}

func (s StyleColor) RGBA() RGBA {
	return RGBA(s)
}

func (s StyleColor) String() string {
	// Return default for the zero color.
	if s.IsZero() {
		return "clear"
	}

	// look for named color
	for name, col := range colornames.Map {
		if col.R == s.R && col.G == s.G && col.B == s.B {
			if col.A == s.A {
				return name
			} else {
				return fmt.Sprintf("%s %02x", name, s.A)
			}
		}
	}
	return `#` + fmt.Sprintf("%02x%02x%02x%02x", s.R, s.G, s.B, s.A)
}

func (s StyleColor) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *StyleColor) UnmarshalText(buf []byte) error {
	scol := string(buf)
	if scol == "" || scol == "clear" { // allow empty/clear as unset color.
		s.R = 0
		s.G = 0
		s.B = 0
		s.A = 0
		return nil
	}

	parts := strings.Split(scol, " ")
	name := parts[0]

	// try color names.
	if ncol, ok := colornames.Map[strings.ToLower(name)]; ok {
		s.R = ncol.R
		s.G = ncol.G
		s.B = ncol.B
		s.A = ncol.A
		if len(parts) > 1 {
			_, err := fmt.Sscanf(parts[1], `%2x`, &s.A)
			if err != nil {
				dprintln("error: ", string(buf))
			}
			f := float64(s.A) / 255
			s.R = byte(float64(s.R) * f)
			s.G = byte(float64(s.G) * f)
			s.B = byte(float64(s.B) * f)
			return err
		}
		return nil
	}

	_, err := fmt.Sscanf(scol, `#%2x%2x%2x%2x`, &s.R, &s.G, &s.B, &s.A)
	if err != nil {
		dprintln("error: ", string(buf))
	}
	return err
}

func (s StyleColor) IsZero() bool {
	return s.R == 0 && s.G == 0 && s.B == 0 && s.A == 0
}

func (s StyleColor) WithDefault(def StyleColor) StyleColor {
	if s.IsZero() {
		if def.IsZero() {
			return NewStyleColor(0, 0, 0, 255)
		}
		return def
	}
	return s
}

type fontStyle struct {
	Family string     `json:"family,omitempty"`
	Size   StyleSize  `json:"size,omitempty"`
	Color  StyleColor `json:"color,omitempty"`
}

type FontStyle struct {
	Family string    `json:"family,omitempty"`
	Size   StyleSize `json:"size,omitempty"`
	Font   *Font     `json:"-"`
	Face   Face      `json:"-"`
}

func (s *FontStyle) UnmarshalJSON(buf []byte) error {
	fs := fontStyle{}
	err := json.Unmarshal(buf, &fs)
	if err != nil {
		return err
	}
	s.Family = fs.Family
	s.Size = fs.Size
	if s.Family == "" || s.Family == "default" {
		s.Font = defaultFont
	} else {
		font := loadResourceFontOptional("resource/font/" + s.Family + ".ttf")
		if font == nil {
			s.Font = defaultFont
		} else {
			s.Font = font
		}
		if s.Size < 1 {
			s.Size = 12
		}
		s.Face = fontFace(s.Font, s.Size.Int())
	}
	return nil
}

func (f FontStyle) WithDefault(def FontStyle) FontStyle {
	f.Size = f.Size.WithDefault(def.Size)

	if f.Family == "" {
		f.Family = def.Family
	}
	if f.Font == defaultFont || f.Font == nil {
		f.Font = def.Font
	}
	if f.Face == nil {
		f.Face = def.Face
	}
	return f
}

type LineStyle struct {
	Size  StyleSize  `json:"size,omitempty"`
	Color StyleColor `json:"color,omitempty"`
}

func (l LineStyle) WithDefault(def LineStyle) LineStyle {
	l.Size = l.Size.WithDefault(def.Size)
	l.Color = l.Color.WithDefault(def.Color)
	return l
}

func (l *LineStyle) WithDefaultPointer(def LineStyle) *LineStyle {
	if l == nil {
		return &def
	} else {
		res := l.WithDefault(def)
		return &res
	}
}

type FillStyle struct {
	Color  StyleColor  `json:"color"`
	Sprite StyleSprite `json:"sprite"` // background sprite to use from the ui_atlas sheet.
}

func (f FillStyle) WithDefault(def FillStyle) FillStyle {
	f.Color = f.Color.WithDefault(def.Color)
	f.Sprite = f.Sprite.WithDefault(def.Sprite)
	return f
}

type StyleRect struct {
	Width  StyleSize `json:"width,omitempty"`
	Height StyleSize `json:"height,omitempty"`
}

func (r StyleRect) WithDefault(def StyleRect) StyleRect {
	r.Width = r.Width.WithDefault(def.Width)
	r.Height = r.Height.WithDefault(def.Height)
	return r
}

// StyleSprites is a set of sprites per usage
type StyleSprites struct {
	Close    StyleSprite `json:"close,omitempty"`
	Check    StyleSprite `json:"check,omitempty"`
	Radio    StyleSprite `json:"radio,omitempty"`
	Minimize StyleSprite `json:"minimize,omitempty"`
	Maximize StyleSprite `json:"maximize,omitempty"`
}

type StyleLayout int

const (
	StyleLayoutDefault StyleLayout = iota
	StyleLayoutCompact             // Compact layout, the default.
	StyleLayoutStretch             // Stretch to the maximum size possible.
)

func (s StyleLayout) String() string {
	switch s {
	case StyleLayoutCompact:
		return "compact"
	case StyleLayoutStretch:
		return "stretch"
	default:
		return ""
	}
}

func (s StyleLayout) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *StyleLayout) UnmarshalText(buf []byte) error {
	sal := string(buf)
	switch sal {
	case "compact":
		*s = StyleLayoutCompact
	case "stretch":
		*s = StyleLayoutStretch
	case "":
		*s = StyleLayoutDefault
	default:
		return fmt.Errorf("Unknown alignment: %s", sal)
	}
	return nil
}

func (s StyleLayout) WithDefault(def StyleLayout) StyleLayout {
	if s == StyleLayoutDefault {
		if def == StyleLayoutDefault {
			return StyleLayoutCompact
		}
		return def
	}
	return s
}

type StyleAlign int

const (
	StyleAlignDefault StyleAlign = iota
	StyleAlignLeft
	StyleAlignMiddle
	StyleAlignRight
	StyleAlignJustify
	StyleAlignStart  = StyleAlignLeft
	StyleAlignCenter = StyleAlignMiddle
	StyleAlignEnd    = StyleAlignRight
	AlignStart       = StyleAlignLeft
	AlignCenter      = StyleAlignMiddle
	AlignEnd         = StyleAlignRight
)

func (s StyleAlign) String() string {
	switch s {
	case StyleAlignLeft:
		return "left"
	case StyleAlignMiddle:
		return "middle"
	case StyleAlignRight:
		return "right"
	case StyleAlignJustify:
		return "justify"
	default:
		return ""
	}
}

func (s StyleAlign) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *StyleAlign) UnmarshalText(buf []byte) error {
	sal := string(buf)
	switch sal {
	case "left":
		*s = StyleAlignLeft
	case "middle":
		*s = StyleAlignMiddle
	case "right":
		*s = StyleAlignRight
	case "justify":
		*s = StyleAlignJustify
	case "":
		*s = StyleAlignDefault
	default:
		return fmt.Errorf("Unknown alignment: %s", sal)
	}
	return nil
}

func (s StyleAlign) WithDefault(def StyleAlign) StyleAlign {
	if s == StyleAlignDefault {
		if def == StyleAlignDefault {
			return StyleAlignLeft
		}
		return def
	}
	return s
}

// Style is a set of colors, fonts, sizes, icons, and sprites that apply either
// for certain widgets or for certain states.
type Style struct {
	Color  StyleColor  `json:"color,omitempty"`
	Font   FontStyle   `json:"font,omitempty"`
	Fill   FillStyle   `json:"fill,omitempty"`
	Margin StyleSize   `json:"margin,omitempty"`
	Icon   StyleSprite `json:"icon,omitempty"`
	Size   StyleRect   `json:"size,omitempty"`
	Align  StyleAlign  `json:"align,omitempty"`
	Layout StyleLayout `json:"layout,omitempty"`
}

func (l Style) WithDefault(def Style) Style {
	l.Color = l.Color.WithDefault(def.Color)
	l.Font = l.Font.WithDefault(def.Font)
	l.Fill = l.Fill.WithDefault(def.Fill)
	l.Margin = l.Margin.WithDefault(def.Margin)
	l.Icon = l.Icon.WithDefault(def.Icon)
	l.Size = l.Size.WithDefault(def.Size)
	return l
}

func (l *Style) WithDefaultPointer(def Style) *Style {
	if l == nil {
		return &def
	} else {
		res := l.WithDefault(def)
		return &res
	}
}

// Theme is a set of styles.
type Theme struct {
	Style                     // default style.
	Active   *Style           `json:"active,omitempty"`
	Alert    *Style           `json:"alert,omitempty"`
	Disable  *Style           `json:"disable,omitempty"`
	Focus    *Style           `json:"focus,omitempty"`
	Hover    *Style           `json:"hover,omitempty"`
	Box      *Style           `json:"box,omitempty"`
	Button   *Style           `json:"button,omitempty"`
	Checkbox *Style           `json:"checkbox,omitempty"`
	Column   *Style           `json:"column,omitempty"`
	Dropdown *Style           `json:"dropdown,omitempty"`
	Entry    *Style           `json:"entry,omitempty"`
	Grid     *Style           `json:"grid,omitempty"`
	Group    *Style           `json:"group,omitempty"`
	Journal  *Style           `json:"journal,omitempty"`
	Label    *Style           `json:"label,omitempty"`
	Media    *Style           `json:"media,omitempty"`
	Note     *Style           `json:"note,omitempty"`
	Pane     *Style           `json:"pane,omitempty"`
	Picture  *Style           `json:"picture,omitempty"`
	Tab      *Style           `json:"tab,omitempty"`
	Table    *Style           `json:"table,omitempty"`
	Tray     *Style           `json:"tray,omitempty"`
	Error    *Style           `json:"error,omitempty"`
	Menu     *Style           `json:"menu,omitempty"`
	Slab     *Style           `json:"slab,omitempty"`
	Slider   *Style           `json:"slider,omitempty"`
	Scroller *Style           `json:"scroller,omitempty"`
	Radio    *Style           `json:"radio,omitempty"`
	Roller   *Style           `json:"roller,omitempty"`
	Card     *Style           `json:"card,omitempty"`
	List     *Style           `json:"list,omitempty"`
	Cursor   *LineStyle       `json:"cursor,omitempty"`
	Icons    StyleSprites     `json:"icons,omitempty"`
	DPI      StyleSize        `json:"dpi,omitempty"`
	Custom   map[string]Style `json:"custom,omitempty"`
}

func (s Style) ApplyMargin(c Control, x, y int) (dx, dy, dw, dh int) {
	margin := s.Margin.Int()
	dx, dy = c.WidgetAt()
	dw, dh = c.WidgetSize()
	dx += x + margin
	dy += y + margin
	dw -= 2 * margin
	dh -= 2 * margin
	return dx, dy, dw, dh
}

func (t Theme) WithDefault() Theme {
	t.Style = t.Style.WithDefault(t.Style)
	t.Active = t.Active.WithDefaultPointer(t.Style)
	t.Disable = t.Disable.WithDefaultPointer(t.Style)
	t.Focus = t.Focus.WithDefaultPointer(t.Style)
	t.Hover = t.Hover.WithDefaultPointer(t.Style)
	t.Box = t.Box.WithDefaultPointer(t.Style)
	t.Button = t.Button.WithDefaultPointer(t.Style)
	t.Checkbox = t.Checkbox.WithDefaultPointer(t.Style)
	t.Column = t.Column.WithDefaultPointer(t.Style)
	t.Dropdown = t.Dropdown.WithDefaultPointer(t.Style)
	t.Entry = t.Entry.WithDefaultPointer(t.Style)
	t.Journal = t.Journal.WithDefaultPointer(t.Style)
	t.Label = t.Label.WithDefaultPointer(t.Style)
	t.Pane = t.Pane.WithDefaultPointer(t.Style)
	t.Tray = t.Tray.WithDefaultPointer(t.Style)
	t.Tab = t.Tab.WithDefaultPointer(t.Style)
	t.Table = t.Table.WithDefaultPointer(t.Style)
	t.Error = t.Error.WithDefaultPointer(t.Style)
	t.Media = t.Media.WithDefaultPointer(t.Style)
	t.Menu = t.Menu.WithDefaultPointer(t.Style)
	t.Picture = t.Picture.WithDefaultPointer(t.Style)
	t.Note = t.Note.WithDefaultPointer(t.Style)
	t.Menu = t.Menu.WithDefaultPointer(t.Style)
	t.Slab = t.Picture.WithDefaultPointer(t.Style)
	t.Slider = t.Slider.WithDefaultPointer(t.Style)
	t.Scroller = t.Scroller.WithDefaultPointer(t.Style)
	t.Radio = t.Radio.WithDefaultPointer(t.Style)
	t.Roller = t.Roller.WithDefaultPointer(t.Style)
	t.Card = t.Card.WithDefaultPointer(t.Style)
	t.List = t.List.WithDefaultPointer(t.Style)
	t.Alert = t.Alert.WithDefaultPointer(t.Style)
	defaultCursor := LineStyle{Color: t.Style.Color, Size: 1}
	t.Cursor = t.Cursor.WithDefaultPointer(defaultCursor)
	return t
}

var theme *Theme

func ShowTheme() {
	buf, _ := json.MarshalIndent(theme, "", "    ")
	dprintf("theme:\n%s", string(buf))
}

func initTheme() {
	theme = loadResourceJSON[Theme](themeName)
	ShowTheme()
	themeDefaults := theme.WithDefault()
	theme = &themeDefaults
	ShowTheme()
}

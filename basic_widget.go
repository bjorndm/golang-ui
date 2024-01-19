package ui

// LayoutUnlimited can be used to impose no constrains on a layout.
const LayoutUnlimited = 1 << 31

// Control is a specific ebui control or widget.
type Control interface {
	// HandleWidget handles the event.
	// This function should only be called with events that the widget is meant to process.
	// The parent widget, or the caller of this method, is responsible for filtering out events that this widget does not need to receive.
	// This is also why this function has no return value, each event *will* be handled.
	HandleWidget(e Event)
	// DrawWidget draws the widget and it's children at the current location.
	DrawWidget(g *Graphic)

	// LayoutWidget lays out the widget, in particular the child widgets.
	// Here, width is the available width for layout, and
	// height is the available height for layout.
	//
	// To indicate there is no constraint, pass in the number obtained from
	// the parent, then the widget should be able to layout itself with its
	// preferred size. In case the size is not known yet,
	// pass in LayoutUnlimited. It is an error to pass 0 or negative
	// either for the width or the height, since this will cause the widget to
	// layout incorrectly since t has no space available.
	//
	// The widget should not change its own position, but rather lay out the
	// position of it's child widgets. The widget may change its size,
	// but it may not become larger than the given width or height.
	// If it was to become larger, it must clip, shrink, or provide a way
	// to scroll so it will fit. However, floating elements do not need to
	// clip themselves.
	//
	// After calling LayoutWidget, WidgetSize must return the layout width
	// and layout height of the widget.
	LayoutWidget(width, height int)

	// MoveWidget moves the widget to the given x and y.
	// The position is relative to the parent.
	MoveWidget(x, y int)
	// WidgetAt returns the current position of the widget relative to the parent.
	WidgetAt() (x, y int)
	// WidgetSize returns the current layout size of the widget.
	WidgetSize() (width, height int)
	// Style returns the current style of the widget
	Style() Style
	// SetStyle sets a custom style on the widget. Set nil to use the default
	// style. However, it is recommended to change the theme in stead if the
	// style should apply to multiple widgets.
	SetStyle(*Style)

	// WidgetLayer returns the later the widget is in.
	WidgetLayer() int
	// Changes the layer of the widget by delta
	RaiseWidget(delta int)
	// Hide hides the widget.
	Hide()
	// Show show the widget.
	Show()

	WidgetOverflow() (width, height int)

	Parent() Control
	SetParent(Control)
	Destroy()

	Hidden() bool
	Enabled() bool

	Focus() Control
	SetFocus(Control)
}

type BasicWidget struct {
	screen       *Graphic
	parent       Control
	focused      Control
	width        int
	height       int
	x            int
	y            int
	z            int
	wfull        int
	hfull        int
	wantHidden   bool
	wantDisabled bool
	tooltip      string
	customStyle  *Style
	sub          *Image // sub image for clipping
	floating     Control
}

type DialogStarter interface {
	StartDialog(dialog Control, title string, modal bool)
}

func (w BasicWidget) Focus() Control {
	return w.focused
}

func (w *BasicWidget) SetFocus(focused Control) {
	w.focused = focused
}

func (w BasicWidget) Parent() Control {
	return w.parent
}

func (w *BasicWidget) SetParent(parent Control) {
	w.parent = parent
}

func (w *BasicWidget) Toplevel() bool {
	return false
}

func (w *BasicWidget) Show() {
	w.wantHidden = false
}

func (w *BasicWidget) Hide() {
	w.wantHidden = true
	UnfocusParentIfNeeded(w)
}

func (w *BasicWidget) Hidden() bool {
	return w.wantHidden
}

func (w *BasicWidget) Enabled() bool {
	return !w.wantDisabled
}

func (w *BasicWidget) Enable() {
	w.wantDisabled = false
}

func (w *BasicWidget) Disable() {
	w.wantDisabled = true
}

// UnfocusParentIfNeeded, will, if the widget was focused by a parent,
// set the focus of that parent to nothing.
func UnfocusParentIfNeeded(c Control) {
	parent := c.Parent()
	if parent != nil {
		focus := parent.Focus()
		if focus == c {
			parent.SetFocus(nil)
		}
	}
}

func (w *BasicWidget) Destroy() {
	UnfocusParentIfNeeded(w) // NOTE static inheritance !
}

func (w BasicWidget) SetToolTip(value string) {
	w.tooltip = value
}

func (w *BasicWidget) ToolTip() string {
	return w.tooltip
}

func (s Style) SizeSize() (width, height int) {
	return s.Size.Width.Int(), s.Size.Height.Int()
}

// GrowToStyleSize makes the basic widget at least as big as the defined style size.
func (w *BasicWidget) GrowToStyleSize() {
	sw, sh := w.Style().SizeSize() // NOTE static inheritance on Style() !

	if w.width < sw {
		w.width = sw
	}

	if w.height < sh {
		w.height = sh
	}
}

// ClipTo clips the widgets to the given sizes if they are > 0.
// It als sets wfull and hfull to the original unclipped size.
func (w *BasicWidget) ClipTo(width, height int) {
	w.wfull = w.width
	w.hfull = w.height

	if width > 0 && w.width > width {
		w.width = width
	}

	if height > 0 && w.height > height {
		w.height = height
	}
}

func (w *BasicWidget) LayoutWidget(width, height int) {
	// By default a widget takes its style size, shinking if needed.
	w.width, w.height = w.Style().SizeSize() // NOTE static inheritance on Style() !
	w.ClipTo(width, height)
	dprintln("BasicWidget.LayoutWidget: ", w.width, w.height)
}

func (w *BasicWidget) ApplyMargin(x, y int) (dx, dy, dw, dh int) {
	style := w.Style()
	margin := style.Margin.Int()
	dx = x + w.x + margin
	dy = y + w.y + margin
	dw = w.width - 2*margin
	dh = w.width - 2*margin
	return dx, dy, dw, dh
}

func (b *BasicWidget) DrawWidget(g *Graphic) {
	x, y := b.WidgetAbsolute()
	w, h := b.WidgetSize()
	FillFrameStyle(g, x, y, w, h, b.Style())
}

func (w *BasicWidget) OverlayWidget(g *Graphic, layer, x, y int) {
}

func (w *BasicWidget) DrawDebug(g *Graphic, form string, args ...any) {
	x, y := w.WidgetAbsolute() // NOTE static inheritance !
	args = append(args, x, y, w.width, w.height)
	DrawDebug(g, x, y, w.width, w.height, form+"+ %d %d %d %d", args...)
}

func (w *BasicWidget) HandleWidget(ev Event) {
	if debugDisplay {
		dprintln("BasicWidget.HandleWidget: ", ev)
	}
}

func (w *BasicWidget) MoveWidget(x, y int) {
	w.x = x
	w.y = y
}

func (w BasicWidget) WidgetAt() (x, y int) {
	return w.x, w.y
}

func (w BasicWidget) WidgetSize() (width, height int) {
	return w.width, w.height
}

var _ Control = &BasicWidget{}

func ControlAbsolute(c Control) (x, y int) {
	x, y = c.WidgetAt()
	for parent := c.Parent(); parent != nil; parent = parent.Parent() {
		px, py := parent.WidgetAt()
		x += px
		y += py
	}
	return x, y
}

func (w BasicWidget) WidgetAbsolute() (x, y int) {
	x, y = w.WidgetAt() // NOTE static inheritance !
	for parent := w.Parent(); parent != nil; parent = parent.Parent() {
		px, py := parent.WidgetAt()
		x += px
		y += py
	}
	return x, y
}

// InsideBounds returns whether or not x and y are inside the rectangle
// defined by bx, by, bx+bw, by+bh
func InsideBounds(bx, by, bw, bh, x, y int) bool {
	return x >= bx && y >= by &&
		x <= bx+bw && y <= by+bh
}

func (w BasicWidget) Inside(x, y int) bool {
	absX, absY := w.WidgetAbsolute() // NOTE static inheritance !
	return x >= absX && y >= absY &&
		x <= absX+w.width && y <= absY+w.height
}

func (w BasicWidget) InsidePart(px, py, pw, ph, x, y int) bool {
	absX, absY := w.WidgetAbsolute() // NOTE static inheritance !
	absX += px
	absY += py
	return x >= absX && y >= absY &&
		x <= absX+pw && y <= absY+ph
}

func (w BasicWidget) MouseInside(me MouseEvent) bool {
	return w.Inside(me.X, me.Y) // NOTE static inheritance !
}

func (w BasicWidget) MouseInsidePart(px, py, pw, ph int, me MouseEvent) bool {
	return w.InsidePart(px, py, pw, ph, me.X, me.Y) // NOTE static inheritance !
}

func (w BasicWidget) Style() Style {
	if w.customStyle != nil {
		return *w.customStyle
	}
	return theme.Style
}

func (w *BasicWidget) SetStyle(style *Style) {
	w.customStyle = style
}

func (w BasicWidget) WidgetLayer() int {
	return w.z
}

func (w *BasicWidget) RaiseWidget(delta int) {
	w.z += delta
}

func (w BasicWidget) WidgetOverflow() (width, height int) {
	height = w.hfull - w.height
	width = w.wfull - w.width
	return height, width
}

func FindDialogStarter(control Control) DialogStarter {
	for parent := control; parent != nil; parent = parent.Parent() {
		if dialogStarter, ok := parent.(DialogStarter); ok {
			return dialogStarter
		}
	}
	return nil
}

func ControlStartDialog(above, dialog Control, title string, modal bool) {
	starter := FindDialogStarter(above)
	if starter != nil {
		starter.StartDialog(dialog, title, modal)
		return
	}
	panic("Could not start dialog: " + title)
}

func MoveWidgetBy(w Control, dx, dy int) {
	ox, oy := w.WidgetAt()
	w.MoveWidget(dx+ox, dy+oy)
}

func NeedLayout(control Control) {
	for parent := control; parent != nil; parent = parent.Parent() {
		if window, ok := parent.(*Window); ok {
			window.Relayout()
		}
	}
}

type Overlayer interface {
	// StartOverlay requests that the widget c will become an overlay in the Overlayer.
	StartOverlay(c Control)
	// End Overlay requests that the widget c will not be an overlay anomore.
	EndOverlay(c Control)
}

func FindOverlayer(control Control) Overlayer {
	for parent := control; parent != nil; parent = parent.Parent() {
		if overlayer, ok := parent.(Overlayer); ok {
			return overlayer
		}
	}
	return nil
}

func StartOverlayWidget(above, overlay Control) Overlayer {
	overlayer := FindOverlayer(above)
	if overlayer != nil {
		overlayer.StartOverlay(overlay)
		return overlayer
	}
	return nil
}

func EndOverlayWidget(above, overlay Control) Overlayer {
	overlayer := FindOverlayer(above)
	if overlayer != nil {
		overlayer.EndOverlay(overlay)
		return overlayer
	}
	return nil
}

package unix

import "github.com/gotk3/gotk3/gtk"
import "fmt"

type Control interface {
	Parent() Control
	SetParent(Control)
	SetControl(Control)
	Destroy()
	Handle() any
}

type Window struct {
	c           Control
	window      *gtk.Window
	vbox        *gtk.Box
	childHolder *gtk.Box
	menubar     *gtk.MenuBar

	child Control

	margined   bool
	resizeable bool
	focused    bool

	onClosing                func(*Window, any)
	onClosingData            any
	onContentSizeChanged     func(*Window, any)
	onContentSizeChangedData any
	onFocusChanged           func(*Window, any)
	onFocusChangedData       any
	fullscreen               bool
	onPositionChanged        func(*Window, any)
	onPositionChangedData    any
	changingPosition         bool
	changingSize             bool

	cachedPosX   int
	cachedPosY   int
	cachedWidth  int
	cachedHeight int
}

func (w *Window) Handle() any {
	return w.window
}

func (w *Window) SetControl(c Control) {
	w.c = c
}

func (w *Window) whenClosing() func(win *gtk.Window) bool {

	return func(win *gtk.Window) bool {
		w2 := w

		// manually destroy the window ourselves; don't let the destroy-event handler do it
		if w2.onClosing != nil {
			w2.onClosing(w2, w2.onClosingData)
		}
		w2.Destroy()
		// don't continue to the default destroy-event handler; we destroyed the window by now
		return true
	}
}

func (w *Window) Destroy() {

	// first hide ourselves
	w.window.Hide()
	// now destroy the child
	if w.child != nil {
		w.child.SetParent(nil)
		// uiUnixControlSetContainer(uiUnixControl(w.child), w.childHolderContainer, TRUE);
		w.child.Destroy()
	}
	// now destroy the menus, if any
	if w.menubar != nil {
		w.menubar.Destroy()
	}

	w.childHolder.Destroy()
	w.vbox.Destroy()
	// and finally free ourselves
	w.window.Destroy()
}

func (w *Window) whenSizeAllocate() func(*gtk.Box) {
	return func(allocation *gtk.Box) {
		w2 := w

		if (!w2.changingSize) && w2.onContentSizeChanged != nil {
			w2.onContentSizeChanged(w, w2.onContentSizeChangedData)
		}

		if w2.changingSize {
			w2.changingSize = false
		}
	}
}

func (w *Window) whenGetFocus() func(win *gtk.Window) bool {
	return func(win *gtk.Window) bool {
		w2 := w
		w2.focused = true
		if w.onFocusChanged != nil {
			w.onFocusChanged(w2, w2.onFocusChangedData)
		}
		return false
	}
}

func (w *Window) whenLoseFocus() func(win *gtk.Window) bool {
	return func(win *gtk.Window) bool {
		w2 := w
		w2.focused = false
		if w.onFocusChanged != nil {
			w.onFocusChanged(w2, w2.onFocusChangedData)
		}
		return false
	}
}

func (w *Window) whenConfigure() func(win *gtk.Window) bool {
	return func(win *gtk.Window) bool {
		w2 := w
		x, y := w2.window.GetPosition()
		if x != w2.cachedPosX || y != w2.cachedPosY {
			w2.cachedPosX = x
			w2.cachedPosY = y
			if (!w2.changingPosition) && w2.onPositionChanged != nil {
				w2.onPositionChanged(w2, w2.onPositionChangedData)
			}
		}
		if w2.changingPosition {
			w2.changingPosition = false
		}
		return false
	}
}

func (w *Window) Parent() Control {
	return nil
}

func (w *Window) SetParent(parent Control) {
	panic("cannot set parent on top level Window")
}

func (w *Window) Toplevel() bool {
	return true
}

func (w *Window) Show() {
	// don't use gtk_widget_show_all() as that will show all children, regardless of user settings
	// don't use gtk_widget_show(); that doesn't bring to front or give keyboard focus
	// (gtk_window_present() does call gtk_widget_show() though)
	w.window.Present()
}

func (w *Window) Hide() {
	w.window.Hide()
}

func (w *Window) Enabled() bool {
	return w.window.GetAcceptFocus()
}

func (w *Window) Enable() {
	w.window.SetAcceptFocus(true)
}

func (w *Window) Disable() {
	w.window.SetAcceptFocus(false)
}

func (w *Window) Title() string {
	t, _ := w.window.GetTitle()
	return t
}

func (w *Window) SetTitle(title string) {
	w.window.SetTitle(title)
}

func (w Window) Position() (x, y int) {
	return w.window.GetPosition()
}

func (w *Window) SetPosition(x, y int) {
	w.changingPosition = true
	w.window.Move(x, y)

	w.changingPosition = true
	w.window.Move(x, y)
	// gtk_window_move() is asynchronous. Wait for the configure-event
	for w.changingPosition {
		// if (!uiMainStep(1)) break; TODO
	}
}

func (w *Window) OnPositionChanged(f func(*Window, any), data any) {
	w.onPositionChanged = f
	w.onPositionChangedData = data
}

func (w *Window) ContentSize() (width, height int) {
	allocation := w.childHolder.GetAllocation()
	width = allocation.GetWidth()
	height = allocation.GetHeight()
	return width, height
}

func (w *Window) SetContentSize(width, height int) {
	// we need to resize the child holder widget to the given size
	// we can't resize that without running the event loop
	// but we can do gtk_window_set_size()
	// so how do we deal with the differences in sizes?
	// simple arithmetic, of course!

	// from what I can tell, the return from gtk_widget_get_allocation(w.window) and gtk_window_get_size(w.window) will be the same
	// this is not affected by Wayland and not affected by GTK+ builtin CSD
	// so we can safely juse use them to get the real window size!
	// since we're using gtk_window_resize(), use the latter
	winWidth, winHeight := w.window.GetSize()

	// now get the child holder widget's current allocation
	childAlloc := w.childHolder.GetAllocation()
	// and punch that out of the window size
	winWidth -= childAlloc.GetWidth()
	winHeight -= childAlloc.GetHeight()

	// now we just need to add the new size back in
	winWidth += width
	winHeight += height

	w.changingSize = true
	w.window.Resize(winWidth, winHeight)
	// gtk_window_resize may be asynchronous. Wait for the size-allocate event.
	for w.changingSize {
		// if (!uiMainStep(1)) break; TODO
	}
}

func (w Window) Fullscreen() bool {
	return w.fullscreen
}

// TODO use window-state-event to track
// TODO does this send an extra size changed?
// TODO what behavior do we want?
func (w *Window) SetFullscreen(fullscreen bool) {
	w.fullscreen = fullscreen
	if w.fullscreen {
		w.window.Fullscreen()
	} else {
		w.window.Unfullscreen()
	}
}

func (w *Window) OnContentSizeChanged(f func(*Window, any), data any) {
	w.onContentSizeChanged = f
	w.onContentSizeChangedData = data
}

func (w *Window) OnClosing(f func(*Window, any), data any) {
	w.onClosing = f
	w.onClosingData = data
}

func (w Window) Focused() bool {
	return w.focused
}

func (w Window) OnFocusChanged(f func(*Window, any), data any) {
	w.onFocusChanged = f
	w.onFocusChangedData = data
}

func (w Window) Borderless() bool {
	return w.window.GetDecorated()
}

func (w *Window) SetBorderless(borderless bool) {
	w.window.SetDecorated(borderless)
}

// TODO save and restore expands and aligns
func (w *Window) SetChild(child Control) {
	fmt.Println("set child")
	if w.child != nil {
		w.child.SetParent(nil)
		w.childHolder.Remove(w.child.Handle().(*gtk.Widget))
	}
	w.child = child
	if w.child != nil {
		w.child.SetParent(w)
		w.childHolder.Add(child.Handle().(*gtk.Widget))
	}
}

func (w Window) Margined() bool {
	return w.margined
}

func (w *Window) SetMargined(margined bool) {
	w.margined = margined
	setMargined(w.childHolder.Container, w.margined)
}

func (w *Window) Resizeable() bool {
	return w.resizeable
}

func setWidgetBackgroundColor(widget *gtk.Widget, name string) {
	provider, _ := gtk.CssProviderNew()
	provider.LoadFromData(fmt.Sprintf(".bgcol { background-image: none; background-color: %s; }", name))
	sc, _ := widget.GetStyleContext()
	sc.AddProvider(provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
	sc.AddClass("bgcol")
}

func (w *Window) SetResizeable(resizeable bool) {
	// workaround for https://gitlab.gnome.org/GNOME/gtk/-/issues/4945
	// calling gtk_window_set_resizable(w.window, 0) will cause the window to resize to default size
	// (default is smallest size here because we're using gtk_window_resize() when creating the window)
	// to prevent this we call gtk_window_set_default_size() on the current window size so that it doesn't resize
	if !resizeable {
		width, height := w.window.GetSize()
		w.window.SetDefaultSize(width, height)
	}

	w.resizeable = resizeable
	w.window.SetResizable(resizeable)
}

func NewWindow(title string, width int, height int, hasMenubar bool) *Window {
	w := &Window{}
	// w.Control = NewControl(w);

	w.resizeable = true
	w.window, _ = gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

	w.window.SetTitle(title)
	w.window.Resize(width, height)

	w.vbox, _ = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)

	// set the vbox as the GtkWindow child
	w.window.Add(w.vbox)

	if hasMenubar {
		w.menubar, _ = gtk.MenuBarNew()
		w.vbox.Add(w.menubar)
	}

	w.childHolder, _ = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)

	w.childHolder.SetHExpand(true)
	w.childHolder.SetHAlign(gtk.ALIGN_FILL)

	w.childHolder.SetVExpand(true)
	w.childHolder.SetVAlign(gtk.ALIGN_FILL)

	w.vbox.Add(w.childHolder)

	// and connect our events
	w.window.Connect("destroy", w.whenClosing())
	w.childHolder.Connect("size-allocate", w.whenSizeAllocate())
	w.window.Connect("focus-in", w.whenGetFocus())
	w.window.Connect("focus-out", w.whenLoseFocus())
	w.window.Connect("configure", w.whenConfigure())
	setWidgetBackgroundColor(&w.childHolder.Widget, "blue")

	/*
		g_signal_connect(w.widget, "delete-event", G_CALLBACK(onClosing), w);
		g_signal_connect(w.childHolderWidget, "size-allocate", G_CALLBACK(onSizeAllocate), w);
		g_signal_connect(w.widget, "focus-in-event", G_CALLBACK(onGetFocus), w);
		g_signal_connect(w.widget, "focus-out-event", G_CALLBACK(onLoseFocus), w);
		g_signal_connect(w.widget, "configure-event", G_CALLBACK(onConfigure), w);
	*/

	/*
		uiWindowOnClosing(w, defaultOnClosing, nil);
		uiWindowOnContentSizeChanged(w, defaultOnPositionContentSizeChanged, nil);
		uiWindowOnFocusChanged(w, defaultOnFocusChanged, nil);
		uiWindowOnPositionChanged(w, defaultOnPositionContentSizeChanged, nil);
	*/

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	w.window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	w.window.ShowAll()

	return w
}

func Init() {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)
}

func TestInit() {
	// Initialize GTK without parsing any command line arguments.
	gtk.TestInit(nil)
}

func Main() {
	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

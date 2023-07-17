package unix

import "github.com/gotk3/gotk3/gtk"

type Control interface{}

type Window struct {
	Control

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

/*
func ( ) onClosing(GtkWidget *win, GdkEvent *e, gpointer data) {
	uiWindow *w = uiWindow(data);

	// manually destroy the window ourselves; don't let the delete-event handler do it
	if ((*(w->onClosing))(w, w->onClosingData))
		uiControlDestroy(uiControl(w));
	// don't continue to the default delete-event handler; we destroyed the window by now
	return TRUE;
}

static anyonSizeAllocate(GtkWidget *widget, GdkRectangle *allocation, gpointer data)
{
	int width, height;
	uiWindow *w = uiWindow(data);

	// Ignore spurious size-allocate events
	uiWindowContentSize(w, &width, &height);
	if (width != w->cachedWidth || height != w->cachedHeight) {
		w->cachedWidth = width;
		w->cachedHeight = height;
		if (!w->changingSize)
			(*(w->onContentSizeChanged))(w, w->onContentSizeChangedData);
	}

	if (w->changingSize)
		w->changingSize = FALSE;
}

static gboolean onGetFocus(GtkWidget *win, GdkEvent *e, gpointer data)
{
	uiWindow *w = uiWindow(data);
	w->focused = 1;
	w->onFocusChanged(w, w->onFocusChangedData);
	return FALSE;
}

static gboolean onLoseFocus(GtkWidget *win, GdkEvent *e, gpointer data)
{
	uiWindow *w = uiWindow(data);
	w->focused = 0;
	w->onFocusChanged(w, w->onFocusChangedData);
	return FALSE;
}

static gboolean onConfigure(GtkWidget *win, GdkEvent *e, gpointer data)
{
	uiWindow *w = uiWindow(data);

	int x, y;

	// Ignore resize events
	uiWindowPosition(w, &x, &y);
	if (x != w->cachedPosX || y != w->cachedPosY) {
		w->cachedPosX = x;
		w->cachedPosY = y;
		if (!w->changingPosition)
			(*(w->onPositionChanged))(w, w->onPositionChangedData);
	}

	if (w->changingPosition)
		w->changingPosition = FALSE;

	return FALSE;
}

static int defaultOnClosing(uiWindow *w, any*data)
{
	return 0;
}

static anydefaultOnPositionContentSizeChanged(uiWindow *w, any*data)
{
	// do nothing
}

static anydefaultOnFocusChanged(uiWindow *w, any*data)
{
	// do nothing
}

static anyuiWindowDestroy(uiControl *c)
{
	uiWindow *w = uiWindow(c);

	// first hide ourselves
	gtk_widget_hide(w->widget);
	// now destroy the child
	if (w->child != NULL) {
		uiControlSetParent(w->child, NULL);
		uiUnixControlSetContainer(uiUnixControl(w->child), w->childHolderContainer, TRUE);
		uiControlDestroy(w->child);
	}
	// now destroy the menus, if any
	if (w->menubar != NULL)
		uiprivFreeMenubar(w->menubar);
	gtk_widget_destroy(w->childHolderWidget);
	gtk_widget_destroy(w->vboxWidget);
	// and finally free ourselves
	// use gtk_widget_destroy() instead of g_object_unref() because GTK+ has internal references (see #165)
	gtk_widget_destroy(w->widget);
	uiFreeControl(uiControl(w));
}

uiUnixControlDefaultHandle(uiWindow)

uiControl *uiWindowParent(uiControl *c)
{
	return NULL;
}

anyuiWindowSetParent(uiControl *c, uiControl *parent)
{
	uiUserBugCannotSetParentOnToplevel("uiWindow");
}

static int uiWindowToplevel(uiControl *c)
{
	return 1;
}

uiUnixControlDefaultVisible(uiWindow)

static anyuiWindowShow(uiControl *c)
{
	uiWindow *w = uiWindow(c);

	// don't use gtk_widget_show_all() as that will show all children, regardless of user settings
	// don't use gtk_widget_show(); that doesn't bring to front or give keyboard focus
	// (gtk_window_present() does call gtk_widget_show() though)
	gtk_window_present(w->window);
}

uiUnixControlDefaultHide(uiWindow)
uiUnixControlDefaultEnabled(uiWindow)
uiUnixControlDefaultEnable(uiWindow)
uiUnixControlDefaultDisable(uiWindow)
// TODO?
uiUnixControlDefaultSetContainer(uiWindow)

char *uiWindowTitle(uiWindow *w)
{
	return uiUnixStrdupText(gtk_window_get_title(w->window));
}

anyuiWindowSetTitle(uiWindow *w, const char *title)
{
	gtk_window_set_title(w->window, title);
}

anyuiWindowPosition(uiWindow *w, int *x, int *y)
{
	gtk_window_get_position(w->window, x, y);
}

anyuiWindowSetPosition(uiWindow *w, int x, int y)
{
	w->changingPosition = TRUE;
	gtk_window_move(w->window, x, y);
	// gtk_window_move() is asynchronous. Wait for the configure-event
	while (w->changingPosition)
		if (!uiMainStep(1))
			break;
}

anyuiWindowOnPositionChanged(uiWindow *w, any(*f)(uiWindow *, any*), any*data)
{
	w->onPositionChanged = f;
	w->onPositionChangedData = data;
}

anyuiWindowContentSize(uiWindow *w, int *width, int *height)
{
	GtkAllocation allocation;

	gtk_widget_get_allocation(w->childHolderWidget, &allocation);
	*width = allocation.width;
	*height = allocation.height;
}

anyuiWindowSetContentSize(uiWindow *w, int width, int height)
{
	GtkAllocation childAlloc;
	gint winWidth, winHeight;

	// we need to resize the child holder widget to the given size
	// we can't resize that without running the event loop
	// but we can do gtk_window_set_size()
	// so how do we deal with the differences in sizes?
	// simple arithmetic, of course!

	// from what I can tell, the return from gtk_widget_get_allocation(w->window) and gtk_window_get_size(w->window) will be the same
	// this is not affected by Wayland and not affected by GTK+ builtin CSD
	// so we can safely juse use them to get the real window size!
	// since we're using gtk_window_resize(), use the latter
	gtk_window_get_size(w->window, &winWidth, &winHeight);

	// now get the child holder widget's current allocation
	gtk_widget_get_allocation(w->childHolderWidget, &childAlloc);
	// and punch that out of the window size
	winWidth -= childAlloc.width;
	winHeight -= childAlloc.height;

	// now we just need to add the new size back in
	winWidth += width;
	winHeight += height;

	w->changingSize = TRUE;
	gtk_window_resize(w->window, winWidth, winHeight);
	// gtk_window_resize may be asynchronous. Wait for the size-allocate event.
	while (w->changingSize)
		if (!uiMainStep(1))
			break;
}

int uiWindowFullscreen(uiWindow *w)
{
	return w->fullscreen;
}

// TODO use window-state-event to track
// TODO does this send an extra size changed?
// TODO what behavior do we want?
anyuiWindowSetFullscreen(uiWindow *w, int fullscreen)
{
	w->fullscreen = fullscreen;
	if (w->fullscreen)
		gtk_window_fullscreen(w->window);
	else
		gtk_window_unfullscreen(w->window);
}

anyuiWindowOnContentSizeChanged(uiWindow *w, any(*f)(uiWindow *, any*), any*data)
{
	w->onContentSizeChanged = f;
	w->onContentSizeChangedData = data;
}

anyuiWindowOnClosing(uiWindow *w, int (*f)(uiWindow *, any*), any*data)
{
	w->onClosing = f;
	w->onClosingData = data;
}

int uiWindowFocused(uiWindow *w)
{
	return w->focused;
}

anyuiWindowOnFocusChanged(uiWindow *w, any(*f)(uiWindow *, any*), any*data)
{
	w->onFocusChanged = f;
	w->onFocusChangedData = data;
}

int uiWindowBorderless(uiWindow *w)
{
	return gtk_window_get_decorated(w->window) == FALSE;
}

anyuiWindowSetBorderless(uiWindow *w, int borderless)
{
	gtk_window_set_decorated(w->window, borderless == 0);
}

// TODO save and restore expands and aligns
anyuiWindowSetChild(uiWindow *w, uiControl *child)
{
	if (w->child != NULL) {
		uiControlSetParent(w->child, NULL);
		uiUnixControlSetContainer(uiUnixControl(w->child), w->childHolderContainer, TRUE);
	}
	w->child = child;
	if (w->child != NULL) {
		uiControlSetParent(w->child, uiControl(w));
		uiUnixControlSetContainer(uiUnixControl(w->child), w->childHolderContainer, FALSE);
	}
}

int uiWindowMargined(uiWindow *w)
{
	return w->margined;
}

anyuiWindowSetMargined(uiWindow *w, int margined)
{
	w->margined = margined;
	uiprivSetMargined(w->childHolderContainer, w->margined);
}

int uiWindowResizeable(uiWindow *w)
{
	return w->resizeable;
}
*/

func (w *Window) SetResizeable(resizeable bool) {
	// workaround for https://gitlab.gnome.org/GNOME/gtk/-/issues/4945
	// calling gtk_window_set_resizable(w->window, 0) will cause the window to resize to default size
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
	/*
		g_signal_connect(w->widget, "delete-event", G_CALLBACK(onClosing), w);
		g_signal_connect(w->childHolderWidget, "size-allocate", G_CALLBACK(onSizeAllocate), w);
		g_signal_connect(w->widget, "focus-in-event", G_CALLBACK(onGetFocus), w);
		g_signal_connect(w->widget, "focus-out-event", G_CALLBACK(onLoseFocus), w);
		g_signal_connect(w->widget, "configure-event", G_CALLBACK(onConfigure), w);
	*/

	/*
		uiWindowOnClosing(w, defaultOnClosing, NULL);
		uiWindowOnContentSizeChanged(w, defaultOnPositionContentSizeChanged, NULL);
		uiWindowOnFocusChanged(w, defaultOnFocusChanged, NULL);
		uiWindowOnPositionChanged(w, defaultOnPositionContentSizeChanged, NULL);
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

func Main() {
	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

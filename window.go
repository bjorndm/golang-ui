package ui

import (
	"fmt"
	"log"
	"os"
	"runtime/pprof"
)

import "github.com/hajimehoshi/ebiten/v2"
import "github.com/hajimehoshi/ebiten/v2/ebitenutil"

type Window struct {
	BasicWidget
	// menus   []*Menu

	shownOnce        bool
	changingSize     bool
	changingPosition bool

	child     Control
	margined  bool
	onClosing func(*Window)

	onContentSizeChanged func(*Window)

	suppressSizeChanged bool

	onFocusChanged func(*Window)

	onPositionChanged       func(*Window)
	onPositionChangedData   any
	suppressPositionChanged bool
	focused                 bool
	needLayout              bool
	drawDebug               bool
	title                   string
	focusControl            Control
	menuBar                 *MenuBar
	dialogs                 *Stack
	inputState
	BasicOverlayer
	Ability // Ability lets Window inherit abilities.
}

func NewWindow(title string, width, height int, hasMenubar bool) *Window {
	w := &Window{}
	w.width = width
	w.height = height
	if hasMenubar {
		w.menuBar = NewMenuBar()
	} else {
		w.menuBar = nil
	}
	w.title = title
	w.needLayout = true
	if v, ok := os.LookupEnv("EBUI_DEBUG"); ok && v == "1" {
		w.drawDebug = true
	}
	w.SetStyle(&theme.Style)
	w.dialogs = NewStack()
	w.dialogs.SetParent(w)

	return w
}

func (w *Window) SetTitle(title string) {
	ebiten.SetWindowTitle(title)
	w.title = title
}

func (w *Window) onAbilityChanged(a *Ability) {
	if a.Rigid() {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	} else {
		ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	}
	ebiten.SetWindowDecorated(!a.Plain())
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

const windowMargins = 4

func (w *Window) LayoutWidget(parentWidth, parentHeight int) {
	childY := 0
	if w.dialogs != nil {
		w.dialogs.LayoutWidget(parentWidth, parentHeight)
		w.dialogs.MoveWidget(0, 0)
	}
	if w.menuBar != nil {
		w.menuBar.LayoutWidget(parentWidth, parentHeight)
		w.menuBar.MoveWidget(0, 0)
		_, childY = w.menuBar.WidgetSize()
	}
	if w.child == nil {
		return
	}
	// size did not change and no relayout needed
	if w.width == parentWidth && w.height == parentHeight && !w.needLayout {
		return
	}
	if w.onContentSizeChanged != nil {
		w.onContentSizeChanged(w)
	}
	w.width = parentWidth
	w.height = parentHeight

	childWidth := w.width - (2 * windowMargins)
	childHeight := w.height - (2 * windowMargins)
	dprintln("LayoutWidget ", w.width, w.height)
	w.child.LayoutWidget(childWidth, childHeight)
	w.child.MoveWidget(windowMargins, windowMargins+childY)
	w.needLayout = false
}

func (w *Window) SetChild(child Control) {
	if w.child != nil {
		w.child.SetParent(nil)
	}
	w.child = child
	w.child.SetParent(w)
	width, height := ebiten.WindowSize()
	w.needLayout = true
	w.LayoutWidget(width, height)
}

func (w *Window) Destroy() {
	// first hide ourselves
	w.Hide()
	// If not preserverd, destroy the child
	if !w.Preserved() && w.child != nil {
		w.child.SetParent(nil)
		w.child.Destroy()
	}
}

func (w *Window) Relayout() {
	w.needLayout = true
}

func (w *Window) Hide() {
	ebiten.MinimizeWindow()
}

func (w *Window) Show() {
	ebiten.RestoreWindow()
}

func (w *Window) Title() string {
	return w.title
}

func (w Window) Position() (x, y int) {
	x, y = ebiten.WindowPosition()
	return x, y
}

func (w *Window) SetPosition(x, y int) {
	ebiten.SetWindowPosition(x, y)
	if w.onPositionChanged != nil {
		w.onPositionChanged(w)
	}
}

func (w *Window) OnPositionChanged(f func(*Window)) {
	w.onPositionChanged = f
}

func (w *Window) OnContentSizeChanged(f func(*Window)) {
	w.onContentSizeChanged = f
}

func (w *Window) OnClosing(f func(*Window)) {
	w.onClosing = f
	ebiten.SetWindowClosingHandled(f != nil)
}

func (w Window) Focused() bool {
	return w.focused
}

func (w Window) OnFocusChanged(f func(*Window)) {
	w.onFocusChanged = f
}

func Init() {
	initResource()
	initClipBoard()
}

func TestInit() {
}

func Exit(code int) {
}

// Window also implements ebiten.Game interface.
// Update proceeds the state.
// Update is called every tick (1/60 [s] by default).
func (w *Window) Update() error {
	if ebiten.IsWindowBeingClosed() {
		log.Print("update: closing", w)
		if w.onClosing != nil {
			w.onClosing(w)
		}
		if w.Permanent() {
			return nil
		} else {
			return ebiten.Termination
		}
	}
	if !w.Enabled() {
		return nil
	}
	w.inputState.convertInputToEvents(w, func(e Event) {
		w.HandleWidget(e)
	})

	return nil
}

// Draw draws the screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (w *Window) Draw(screen *ebiten.Image) {
	w.DrawWidget(screen)
	if w.drawDebug {
		msg := fmt.Sprintf(`TPS: %0.2f
	FPS: %0.2f`, ebiten.ActualTPS(), ebiten.ActualFPS())
		ebitenutil.DebugPrint(screen, msg)
	}
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
func (w *Window) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if w.width != outsideWidth || w.height != outsideHeight || w.needLayout {
		w.LayoutWidget(outsideWidth, outsideHeight)
	}
	return w.width, w.height
}

// Constrains the window to be at least the size of the main widget.
func (w *Window) ConstrainToWidgetSize() {
	width, height := 0, 0
	if w.menuBar != nil {
		bw, bh := w.menuBar.WidgetSize()
		width += bw
		height += bh
	}
	if w.child != nil {
		cw, ch := w.child.WidgetSize()
		width += cw
		height += ch
	}

	ebiten.SetWindowSizeLimits(width, height, -1, -1)
}

func (w *Window) HandleWidget(e Event) {
	if debugDisplay {
		log.Printf("event: %#v\n", e)
	}

	// dialogs have highest priority
	if w.dialogs != nil {
		if used := HandleContainerIfNeeded(e, w.dialogs); used {
			// Might need to clean up some closed panes that were focued here.
			w.dialogs.CleanupClosePanes()
			return
		}
	}

	// handle menu bar with priority.
	if w.menuBar != nil {
		if used := HandleContainerIfNeeded(e, w.menuBar); used {
			return
		}
	}

	// After that, overlays get the events. Stop if the event was used.
	if used := w.HandleEventForOverlays(e); used {
		return
	}

	// Finally the child widget.
	if w.child != nil {
		w.child.HandleWidget(e)
	}
}

func (w *Window) DrawWidget(screen *Graphic) {
	FillFrameStyle(screen, 0, 0, w.width, w.height, w.Style())

	if w.child != nil {
		w.child.DrawWidget(screen)
	}

	// Draw overlays over the other widgets, but under the menu bar.
	w.DrawOverlays(screen)

	// draw menu bar over other widgets.
	if w.menuBar != nil {
		w.menuBar.DrawWidget(screen)
	}

	// draw dialogs over the main window
	// dialogs have highest priority
	if w.dialogs != nil && w.dialogs.NumChildren() > 0 {
		w.dialogs.DrawWidget(screen)
	}
}

func (w Window) StartDialog(dialog Control, title string, modal bool) {
	if pane, ok := dialog.(*Pane); ok {
		// if it is a pane use it as is.
		pane.SetModal(modal)
		pane.SetTitle(title)
		pane.LayoutWidget(w.width, w.height)
		pw, ph := pane.WidgetSize()
		pane.MoveWidget(w.width/2-pw/2, w.height/2-ph/2)
		w.dialogs.Append(pane)
	} else {
		// Otherwide create a pane for the dialog.
		pad := w.Style().Margin.Int()
		pane = NewPane(title, w.width-pad*2, w.height-pad*2, false)
		pane.modal = modal
		pane.SetChild(dialog)
		pane.LayoutWidget(w.width, w.height)
		w.dialogs.Append(pane)
	}
}

func Main(w *Window) {
	// enable profiling if the env variable is set.
	if pf, ok := os.LookupEnv("EBUI_PPROF"); ok {
		f, err := os.Create(pf)
		if err != nil {
			fmt.Printf("error: %s", err.Error())
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	ebiten.SetRunnableOnUnfocused(true)
	w.Ability.setOnChanged(w.onAbilityChanged)

	ebiten.SetWindowSize(w.width, w.height)
	ebiten.SetWindowTitle(w.title)
	w.SetFixed(false)
	dprintln("Window: ", w.width, w.height, w.title)
	w.LayoutWidget(w.width, w.height)
	err := ebiten.RunGame(w)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Window: done")
}

/*
func MainWithMenus(w *Window, menus ...*Menu) {
	w.setupMenus(menus...)
	ebiten.RunGame(w)
}
*/

var _ Overlayer = &Window{}

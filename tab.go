package ui

import "golang.org/x/exp/slices"

// Tab contains several tabs of which only one is active,
// with a tab bar on top.
type Tab struct {
	BasicContainer              // Tab is a container
	bar            Tray         // the tab bar on top
	tabs           []*tabHeader // The tab headers in the bar.
	active         int          // index of active tab.
}

func NewTab() *Tab {
	t := &Tab{}
	t.bar.SetParent(t)
	// Use plain style for the tab bar. The tabs themselves will be styled.
	t.SetStyle(&theme.Style)
	t.bar.SetStyle(&theme.Style)
	return t
}

// tabHeader is the top clickable part of a tab.
// it is a small wrapper widget to make IconTextWidget clickable
type tabHeader struct {
	*IconTextWidget
	tab     *Tab
	index   int
	pressed bool
}

func newTabHeader(iwt *IconTextWidget, tab *Tab, index int) *tabHeader {
	th := &tabHeader{IconTextWidget: iwt, tab: tab, index: index, pressed: false}
	th.IconTextWidget.SetStyle(theme.Tab)
	return th
}

func (h *tabHeader) HandleWidget(ev Event) {
	// If clicked switch to the header.
	if mc, ok := ev.(*MouseClickEvent); ok {
		if mc.Inside(h) {
			dprintln("tabHeader selected", h.index)
			h.tab.Select(h.index)
		}
	}
}

func (h *tabHeader) DrawWidget(screen *Graphic) {
	dx, dy := h.WidgetAbsolute()
	style := h.Style()

	if h.pressed {
		// if active use the active color.
		active := *theme.Active
		style.Fill.Color = active.Fill.Color
	}

	FillFrameStyle(screen, dx, dy, h.width, h.height, style)
	h.IconTextWidget.DrawWidget(screen)
}

const tabPadding = 4
const tabMinHeight = 37

func (t *Tab) AppendWithIcon(icon, name string, child Control) {
	t.BasicContainer.AppendWithParent(child, t)
	iwt := NewIconTextWidget(icon, name)
	index := t.BasicContainer.NumChildren() - 1
	th := newTabHeader(iwt, t, index)
	t.tabs = append(t.tabs, th)
	t.bar.Append(th)
}

func (t *Tab) Append(name string, child Control) {
	t.AppendWithIcon("", name, child)
}

func (t *Tab) Delete(index int) {
	if index < 0 {
		return
	}
	if index >= t.NumChildren() {
		return
	}
	t.BasicContainer.Delete(index)
	t.bar.Delete(index)
	t.tabs = slices.Delete(t.tabs, index, index+1)
}

func (t Tab) NumPages() int {
	return t.NumChildren()
}

func (t *Tab) Select(index int) {
	if index < 0 || index >= len(t.controls) {
		return
	}
	t.active = index
	for i := 0; i < len(t.tabs); i++ {
		t.tabs[i].pressed = (i == t.active)
	}
	for i, control := range t.Children() {
		if i == t.active {
			control.Show()
		} else {
			control.Hide()
		}
	}
}

func (t *Tab) LayoutWidget(parentWidth, parentHeight int) {
	// We use height as the height for the tab bar, not for our height,
	// which will become the parent height minus the tab bar height
	tabHeight := t.Style().Size.Height.Int() + t.Style().Margin.Int()
	// lay out the bar as a tray, then...
	t.bar.LayoutWidget(parentWidth, tabHeight)
	t.bar.MoveWidget(0, 0)

	// lay out the container widgets
	for i, control := range t.BasicContainer.Children() {
		control.LayoutWidget(parentWidth, parentHeight-tabHeight)
		control.MoveWidget(0, tabHeight)
		if i != t.active {
			control.Hide()
		}
	}
	// and adjust the size.
	// XXX the tab might be higher and need a scrollbar
	t.width = parentWidth
	t.height = parentHeight - tabHeight
	t.ClipTo(parentWidth, parentHeight)
}

func (t *Tab) DrawWidget(screen *Graphic) {
	t.bar.DrawWidget(screen)
	t.BasicContainer.DrawWidget(screen)
}

func (t *Tab) HandleWidget(ev Event) {
	// If clicked switch to the header.
	if mc, ok := ev.(*MouseClickEvent); ok {
		if mc.Inside(&t.bar) {
			t.bar.HandleWidget(ev)
			return
		}
	}
	t.BasicContainer.HandleWidget(ev)
}

// Returns the index of the first tab with the given name.
// Returns -1 if not found.
func (t Tab) IndexOfName(name string) int {
	return slices.IndexFunc(t.tabs, func(t *tabHeader) bool { return t.Text() == name })
}

// Selects the first tab with the given name.
// Does nothing if not found.
func (t *Tab) SelectName(name string) {
	i := t.IndexOfName(name)
	if i < 0 {
		return
	}
	t.Select(i)
}

// Deletes the first tab with the given name.
// Does nothing if not found.
func (t *Tab) DeleteName(name string) {
	i := t.IndexOfName(name)
	if i < 0 {
		return
	}
	t.Delete(i)
}

// Returns the index of the first tab with the given icon.
// Returns -1 if not found.
func (t Tab) IndexOfIcon(icon string) int {
	return slices.IndexFunc(t.tabs, func(t *tabHeader) bool { return t.icon == icon })
}

// Selects the first tab with the given icon.
// Does nothing if not found.
func (t *Tab) SelectIcon(icon string) {
	i := t.IndexOfIcon(icon)
	if i < 0 {
		return
	}
	t.Select(i)
}

// Deletes the first tab with the given icon.
// Does nothing if not found.
func (t *Tab) DeleteIcon(icon string) {
	i := t.IndexOfIcon(icon)
	if i < 0 {
		return
	}
	t.Delete(i)
}

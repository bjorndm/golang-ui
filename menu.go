package ui

// MenuBar is a bar on top of a Window or Pane with menus in them.
// Since Ebitengine provides no access to the platform's menus, and
// since some platforms like mobile, web, or console don't even have menus
// EBUI simulates menu bars and menus.
// A MenuBar can be used as a normal widget and added everywhere, however,
// this is not reccomended.
type MenuBar struct {
	Tray  // Embed a tray
	menus []*Menu
}

func NewMenuBar() *MenuBar {
	b := &MenuBar{}
	b.SetStyle(theme.Menu)
	b.Tray.SetStyle(theme.Menu)
	return b
}

func (b *MenuBar) AppendMenu(title string) *Menu {
	menu := NewMenu(title)
	b.menus = append(b.menus, menu)
	b.Tray.Append(menu)
	return menu
}

func (b *MenuBar) LayoutWidget(parentWidth, parentHeight int) {
	margin := b.Style().Margin.Int()

	availableWidth := parentWidth
	availableHeight := b.Style().Size.Height.Int() + 2*margin

	if b.Parent() != nil {
		b.width, _ = b.Parent().WidgetSize()
	}

	// lay out as tray, then...
	b.Tray.LayoutWidget(availableWidth, availableHeight)
	// .. And then stretch the menu bar to the avaialble size.
	b.height = availableHeight
	b.width = availableWidth

}

func (b *MenuBar) HandleWidget(e Event) {
	b.Tray.HandleWidget(e)
	b.Tray.SetFocus(nil)
}

// Menu represents a drop down menu on a menu bar.
type Menu struct {
	TextWidget     // Embed text widget for the title
	box        Box // And a box asa static child widget for the menu items.
	title      string
	items      []*MenuItem
	disabled   bool
}

func NewMenu(name string) *Menu {
	m := &Menu{}
	m.SetText(name)
	m.SetStyle(theme.Menu)
	m.box.SetStyle(theme.Menu)
	m.box.SetParent(m) // the box is a child widget
	m.box.Hide()       // hide the menu item box by default
	return m
}

func (m *Menu) Title() string {
	return m.Text()
}

func (m *Menu) LayoutWidget(width, height int) {
	margin := m.Style().Margin.Int()
	h := m.Style().Size.Height.Int()
	m.TextWidget.LayoutWidget(width, height)
	tw, th := m.TextWidget.WidgetSize()

	w := tw
	if th > h {
		h = th
	}
	h += margin * 2

	m.box.LayoutWidget(width-2*margin, height-2*margin-h)
	bw, _ := m.box.WidgetSize()
	if bw > w {
		w = bw
	}
	w += margin * 2

	m.width, m.height = w, h

	// It is not allowed to move self, that is m.Textwidget during layout.
	// But we have to move the box, as it is set as a child widget.
	m.box.MoveWidget(margin, th+margin)

	m.ClipTo(width, height)
}

func (m Menu) DrawWidget(dst *Graphic) {
	m.TextWidget.DrawWidget(dst)
	if !m.box.Hidden() {
		m.box.DrawWidget(dst)
	}
}

func (m *Menu) closeMenu() {
	m.box.Hide()
	m.RaiseWidget(-menuLayer)
	m.SetFocus(nil)
	parent := m.Parent()
	if parent != nil {
		m.parent.SetFocus(nil)
	}
}

func (i *MenuItem) OnClicked(f func(*MenuItem)) {
	i.onClicked = f
}

type menuItemKind int

const (
	menuItemNormal menuItemKind = iota
	menuItemChecked
	menuItemSeparator
)

var menuItemId int

// MenuItem represents an item in a Menu.
type MenuItem struct {
	TextWidget
	// Has unexported fields.
	menu      *Menu
	onClicked func(*MenuItem)
	kind      menuItemKind
	title     string
	id        int
	checked   bool
	disabled  bool
}

func (i *MenuItem) LayoutWidget(width, height int) {
	// A menu bar has a fixed height
	margin := i.Style().Margin.Int()
	checkboxWidth := theme.Checkbox.Size.Width.Int()
	checkboxHeight := theme.Checkbox.Size.Height.Int()

	i.TextWidget.LayoutWidget(width-2*margin, height-2*margin)
	w, h := i.TextWidget.WidgetSize()
	if i.kind == menuItemChecked {
		w += checkboxWidth + margin
		if h < checkboxHeight {
			h = checkboxHeight
		}
	} else if i.kind == menuItemSeparator {
		w = i.Style().Size.Width.Int()
	}

	i.width, i.height = w+2*margin, h+2*margin
	i.ClipTo(width, height)
}

func (m *Menu) AppendItem(title string) *MenuItem {
	item := newMenuItem(menuItemNormal, title, m)
	m.items = append(m.items, item)
	return item
}

func (m *Menu) AppendCheckItem(title string) *MenuItem {
	item := newMenuItem(menuItemChecked, title, m)
	m.items = append(m.items, item)
	return item
}

func (m *Menu) AppendSeparator() *MenuItem {
	item := newMenuItem(menuItemSeparator, "", m)
	m.items = append(m.items, item)
	return item
}

const menuLayer = 10

func (m *Menu) Floating() Control {
	if !m.box.Hidden() {
		return &m.box
	}
	return nil
}

func (m *Menu) HandleWidget(ev Event) {
	if _, ok := ev.(*AwayEvent); ok {
		m.box.RaiseWidget(-menuLayer)
		m.box.Hide() // not for us
		return
	}

	if !m.box.Hidden() {
		m.box.HandleWidget(ev)
	}

	if mc, ok := ev.(*MouseClickEvent); ok {
		if mc.Inside(&m.TextWidget) && m.box.Hidden() {
			dprintln("Menu.HandleWidget: ")
			m.box.Show()
			m.box.RaiseWidget(menuLayer)
		}
	}
}

func newMenuItem(kind menuItemKind, title string, menu *Menu) *MenuItem {
	item := &MenuItem{}
	item.kind = kind
	item.SetText(title)
	item.id = menuItemId
	item.SetStyle(theme.Menu)
	item.menu = menu

	menuItemId++

	menu.box.Append(item)

	return item
}

func (i *MenuItem) Checked() bool {
	return i.checked
}

func (i *MenuItem) SetChecked(checked bool) {
	i.checked = checked
}

func (i *MenuItem) Disable() {
	i.wantDisabled = true
}

func (i *MenuItem) Enable() {
	i.wantDisabled = false
}

func (i MenuItem) DrawWidget(dst *Graphic) {
	dx, dy := i.WidgetAbsolute()
	margin := i.Style().Margin.Int()

	i.TextWidget.DrawWidget(dst)
	if i.kind == menuItemSeparator {
		uiAtlas.DrawSprite(dst, dx, dy+i.height/2, i.width, 3, "hsep")
	} else if i.kind == menuItemChecked {

		cbStyle := theme.Checkbox
		checkboxWidth := cbStyle.Size.Width.Int()
		checkboxHeight := cbStyle.Size.Height.Int()
		checkSprite := theme.Icons.Check.String()

		dx += (i.width - checkboxWidth - margin)
		dy += margin
		FillFrameStyle(dst, dx, dy, checkboxWidth, checkboxHeight, *cbStyle)
		if i.checked {
			iconAtlas.DrawSprite(dst, dx, dy, checkboxWidth, checkboxHeight, checkSprite)
		}
	}
}

func (i *MenuItem) HandleWidget(ev Event) {
	if mc, ok := ev.(*MouseClickEvent); ok {
		if mc.Inside(i) {
			if !i.Enabled() {
				// do nothing.
			} else if i.kind == menuItemChecked {
				i.checked = !i.checked
				if i.onClicked != nil {
					i.onClicked(i)
				}
			} else if i.kind == menuItemNormal {
				if i.onClicked != nil {
					i.onClicked(i)
				}
				if i.menu != nil {
					i.menu.closeMenu()
				}
			}
		}
	}
}

func (w *Window) SetMenuBar(bar *MenuBar) {
	w.menuBar = bar
	if w.menuBar != nil {
		w.menuBar.SetParent(w)
	}
}

func (w *Window) MenuBar() *MenuBar {
	return w.menuBar
}

func (p *Pane) SetMenuBar(bar *MenuBar) {
	p.menuBar = bar
	if p.menuBar != nil {
		p.menuBar.SetParent(p)
	}
}

func (p *Pane) MenuBar() *MenuBar {
	return p.menuBar
}

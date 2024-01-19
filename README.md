# golang-ui

## Introduction

The golang-ui library is portable Go library for graphical user interfaces.
The API is inspired by Andlab's golang-ui. It is based on ebitengine.

It was partially sponsored and open sourced by ExaWizards (https://exawizards.com/en/).

The golang-ui library strives to have an easy to use, portable API,
with automatic layout. It has support for internationalization and input methods.
Currently this is focused on support for European languages, Japanese and other
languages that don't require complex text output.

Unlike other GUI libraries for ebitengine, like ebitengine-ui, golang-ui supports
automatic layout and has an api similar to golang-ui, so you do not have to
style and lay out the widgets manually. However it is possible to style the
widgets if you really need to, and a theme engine is available.

The golang-ui API is designed to be as easy to use as possible with no need for
positioning. The API allow to do high level layouts to order the screen,
but apart from that, everything should work automatically. It is possible to
modify the style of the widgets using the golang-ui API. By default a "professional"
gray colored style will be used. The focus of go-ui is in making it useful for
business applications, so the default theme is not very exciting.

The API is designed to be partially backwards compatible with golang-ui,
however, in several instances the api of golang-ui was simplified, and in other
cases enhanced. The API is currently version 0 so not stable yet.

## Credits

* Icon pack icon_atlas.png by Kenney Vleugels (www.kenney.nl) under CC0 license.
* GoNoto*.ttf fonts by The Noto Project Authors, used under the
  SIL OPEN FONT LICENSE Version 1.1.
* RelaxedTypingMonoJP-*.ttf by Mshiota, used under the SIL OPEN FONT LICENSE
  Version 1.1.

## API and widgets

### Widgets / Controls

Each Widget/Control has a parent widget. Only for the top level Window
widget, the Parent is nil. Each Widget / Control may have one or more child
widgets.

### Event handling

golang-ui transform ebitengine input into events which are sent to the widgets using
the HandleWidget() function. Each widget may only be sent events that it should
process. A widget most choose to which of it's child widgets it passes on the
events or not.

Controls are ordered by their layer which is used for drawing the layered
widget correctly.

Particular about golang-ui is that there is no single focused widget, but that each
container decides which child widgets of it is considered as focused
recursively. Container widgets have to decide where to place the focus, but
child widgets can relinquish the focus by calling SetFocus(nil) on their parent.
Child widgets that loose focus will receive an Away event.

### Layout

Each widget has its own layout. Box, Tray and Grid are useful as containers
for other widgets. The layout engine is hierarchical. Each widget lays out it's
child widgets  by calling LayoutWidget() on them recursively lay them out,
and then MoveWidget, to put them in the correct position.

It is normally error for a widget to move itself using MoveWidget, a that is the
responsibility of the parent widget. However, a widget with multiple sub
widgets may move some of them as is needed.

### Drawing

At first, golang-ui used vector graphics, but unfortunately these were too slow to
be usable for complex widgets such as tables. Therefore now golang-ui uses
the sprite atlas resource/theme/ui_atlas.png for drawing the widgets.
The sprite resource resource/icon/icon_atlas.png is also used for icons.

Both sprite resources are loaded with the help of a JSON file that defines
the locations of the sprites and their properties. Widgets are drawn using a
"nine slice" algorithm so they look good even when scaled. Specify the border
thickness of the ui elements by adding a "border" field in the JSON sprite
atlas.

Thanks to the fast bitmap drawing of Ebitengine, which also has drawing cache,
there does't seem to be a need to cache the drawing in golang-ui itself.

### Box model

Unlike CSS, the width and the height of a widget are the real size of the
outer size of the widget. For performance reasons, the background of the
box is always a colored sprite from the ui atlas. The sprite may or may not have
a border, which is drawn using a nine slice algorithm, so the border is
displayed correctly. The margin the styleable distance between the outer size
and the contents of the widget.

Therefore it is neccesary when laying out a widget to take margin
into consideration to make sure the contents fit well and do not touch the
background sprite's border.


    (x, y)
	+--<--------width-------->----------+
	^ 1: margin    1                    |
	|<1>+<----- content--- -------->+<1>|
	|   |  							|   |
	h   |  							|   |
	e   |  							|   |
	i   |  							|   |
	g   |  							|   |
	h   |  							|   |
	t   |  							|   |
	|   +---------------------------+   |
	v                              1    |
	+-----------------------------------+


### Testing

Go unit tests don't work well with ebitengine, so we use test commands in the
test directort. For example to run the dropdown test run
`go run ./test/dropdown`.

### Widgets

golang-ui has the following widgets available:

- Alert
- BasicWidget
- Box
- Button
- Card
- Checkbox
- BasicContainer
- Dialog
- Dropdown
- Entry
- Grid
- Group
- Journal
- Label
- List
- Media
- Menu
- Note
- Pane
- Picture
- Roller
- Scroller
- Slider
- Stack
- Tab
- Table
- TextWidget
- Tray
- Window

### Tables

Tables consist of columns, which do not have widgets embeddded, but self-draw
all row cells for the relevant column. This so we don't have to allocate widgets
for each table cell and row, which could consume considerable memory.

## Design

### Principle

The guiding design principle of golang-ui is to be easy to use for business
applications without having to specify the layout or style manually.

This is achieved by the idea of trust between the parent widget and the child
widget it contains. There are container widgets which
can contain other widgets, and normal or leaf widgets which cannot. These
rely on each other.

Each widget could be used in isolation, but can control its child widgets if
any. The layout that is requested should be correct and only events that a
widget needs should be sent to it.

### Widget Event handling

Widgets will be sent events using the HandleWidget function. Only events
that a widget needs to process should be sent to it. For example, a mouse click
shoudl only be sent to a widget if it shoud process it, even if it may be located
graphically outside of the widget.

It is the responsability of parent widgets to filter the events their child
receive and to decide on focus. This simplifies the implementation of child
widgets susbtantially.

### Layout

An important concept is that with current devices, vertical scrolling is easy
while horizontal scrolling is hard. The UI library Shoes has the same concept of
"only vertical overflow matters." https://shoesrb.com/manual/Slots.html

The above is true for left -to-righ and right-to-left languages.
For top-to bottom or bottom-to top languages the same goes but with the
horizontal and vertical axes swapped.

Another important concept is available space. Container widgets
will lay out their child widgets informing them of the space that is available
for the child. In most cases, the child must either shrink, clip, or provide a
scroll bar so it will fit in the available space.

On the other hand, if additional space is available, then a child widget may
take up that space of the parent container if so configured or styled.
golang-ui calls this feature "filling".

To simplify layout the basic rule is: leaf widgets know their own size.
In case of text, pictures, media, etc, this is the size a widget needs to
display itself without clipping. In cases the wiget is empty, e.g. like an Entry
then the size is determined by the style of the widget using the
active theme.

The initial size of a widget can be set by setting the size in the style of the
theme. However, if the contents are larger than the size, the widget may "fill"
itself as long as there is space available in the parent.

This doesn't mean size needs to be static. For example
a Journal widget gets bigger as text is added. However, the *size* of a widget
does not depend on the size of the parent or other widgets. It is defined as
the size needed for visibility, or in case the leaf widget is dynamically sized,
the size set by the user of the application. The latter can be changed by
clicking, dragging, etc.

### Focus

Each widget now process all events it gets, however, the parent widget is
responsible for filtering the events so the child widgets only get events
that are of interest to them. To help with the "clicking away"
problem, parent widgets must send "Away" events to all widgets that have lost
focus.

Child widgets may in some cases need to call the SetFocus(nil) on their parent.
For example in case they close themselves to ensure the focus is lost correctly.

### Error handling

This library does in general not return any errors. The purpose is to provide a
reliable UI even if something goes wrong. However, the library will panic on
library usage errors.

### User event handling

The library uses callbacks for user event handling. While channels are
nice, ebui tries not to use goroutines as much as is possible to avoid race
conditions. However, almost all ebui callbacks have a signature of
func(*ObjectThatNeedsHandling), which means the callback can be easily
used to send *ObjectThatNeedsHandling over a channel. All ebui calbacks do not
have a return value to influence the caller. In stead, you may make
modifications to the object in question, possibly by sending it events.

## Release Notes

### v0.1.0

* Added a focus system to prevent leaky events.
* Several breaking changes to the API.


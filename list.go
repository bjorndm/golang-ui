package ui

// The list model is used to fetch the rows needed to display,
// and to react to edits.
type ListModel interface {
	NumRows() int
	FetchRow(index int) Row
	UpdateRow(index int, updated Row)
}

// Card is a widget that ideally is used in a List.
type Card struct {
	Box // A card consists of a box with a header tray with the picture and
	//caption, and a second tray with the added items.
	list      *List
	caption   *TextWidget
	picture   *Picture
	header    *Tray
	items     *Tray
	sublist   *List
	index     int
	onClicked func(card *Card, row int)
}

// OnClicked sets a callback function f that is called, if not nil,
// if the card is clicked. This overrides the default behavior.
// If f is nil the default behavior is restored.
func (c *Card) OnClicked(f func(card *Card, row int)) {
	c.onClicked = f
}

func newCard(caption string) *Card {
	card := &Card{}

	card.Box = *NewBox()
	card.header = NewTray()
	card.items = NewTray()

	card.caption = NewTextWidget(caption)
	card.header.Append(card.caption)
	card.Box.Append(card.header)
	card.Box.Append(card.items)

	card.SetStyle(theme.Card)
	card.caption.SetStyle(theme.Card)
	card.width = card.Style().Size.Width.Int()
	card.caption.width = card.width
	return card
}

func NewCard(caption string) *Card {
	card := newCard(caption)
	return card
}

func (c *Card) LayoutWidget(width, height int) {
	c.Box.LayoutWidget(width, height)
	// cards stretch to parent width
	if c.width < width {
		c.width = width
	}
	return
}

func (c *Card) Index() int {
	return c.index
}

func (c *Card) SetPicture(image Image) *Picture {
	item := NewPicture("", image)
	if c.picture == nil {
		c.picture = NewPicture("", image)
		c.header.Append(c.picture)
	} else {
		c.picture.SetImage(image)
	}
	return item
}

func (c *Card) AppendButtonWithIcon(text, icon string) *Button {
	item := NewButtonWithIcon(text, icon)
	c.items.Append(item)
	return item
}

func (c *Card) AppendButton(text string) *Button {
	item := NewButton(text)
	item.SetParent(c)
	c.items.Append(item)
	return item
}

func (c *Card) AppendLabel(text string) *Label {
	item := NewLabel(text)
	item.SetParent(c)
	c.items.Append(item)
	return item
}

func (c *Card) AppendCheckbox(text string) *Checkbox {
	item := NewCheckbox(text)
	c.items.Append(item)
	return item
}

// SetRow copies the values frop the row to the caption or widgets of  the card.
// The value first element of row is used for the caption, the rest of the
// widgets as created by the AppendLabel, AppendButton or AppendCheckbox
// functions.
func (c *Card) SetRow(row Row) {
	for i, value := range row {
		if i > c.items.NumChildren() {
			return
		}
		if i == 0 {
			c.caption.SetText(value.(string))
			return
		}

		control := c.items.Children()[i-1]
		switch widget := control.(type) {
		case *Button:
			widget.SetText(value.(string))
		case *Label:
			widget.SetText(value.(string))
		case *Checkbox:
			widget.SetChecked(value.(bool))
		default:
			panic("Widget not supported in card")
		}
	}
	NeedLayout(c)
}

type List struct {
	Box       // use a box to lay out the cards.
	ListModel // list model for fetching the data.
	cards     []*Card
	template  func(row Row) *Card

	from  int
	shown int
}

func NewList(model ListModel, template func(row Row) *Card) *List {
	g := &List{ListModel: model}
	g.customStyle = theme.List
	if template == nil {
		panic("NewList: template is mandatory")
	}
	g.template = template
	g.Box = *NewBox()
	NeedLayout(g)
	return g
}

// DropCards frops all scards from the list.
func (g *List) DropCards() {
	g.cards = []*Card{}
	g.Box.Destroy()
	g.Box = *NewBox()
	NeedLayout(g)
}

// AppendCard appends a card to the list.
func (g *List) AppendCard(card *Card) {
	idx := len(g.cards)
	card.index = idx
	card.list = g
	g.cards = append(g.cards, card)
	g.Box.Append(card)
	NeedLayout(g)
}

func (t *List) Card(index int) *Card {
	if index < 0 || index > len(t.cards) {
		return nil
	}
	return t.cards[index]
}

// NumCards returns the amount of cards the list has.
func (t List) NumCards() int {
	return len(t.cards)
}

// CreateCards creates cards of the list based on the table model, using
// the template function. Calls DropCards first to clean the card list.
func (t *List) CreateCards() {
	t.DropCards()
	nr := t.ListModel.NumRows()
	for i := 0; i < nr; i++ {
		row := t.ListModel.FetchRow(i)
		card := t.template(row)
		t.AppendCard(card)
	}
}

// Updates the cards of the list based on the table model and SetRow.
func (t *List) UpdateCards() {
	t.cards = []*Card{}
	nr := t.ListModel.NumRows()
	for i := 0; i < nr; i++ {
		row := t.ListModel.FetchRow(i)
		card := t.Box.Get(i).(*Card)
		card.SetRow(row)
	}
	NeedLayout(t)
}

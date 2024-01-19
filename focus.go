package ui

func SetNewFocus(be *BasicEvent, c Control, newFocus Control, children ...Control) {
	if newFocus != c.Focus() {
		for _, child := range children {
			if child.Hidden() {
				continue
			}
			if child != newFocus {
				child.HandleWidget(&AwayEvent{BasicEvent: *be, Focused: newFocus})
			}
		}
	}

	c.SetFocus(newFocus)
}

func EventInsideOneOf(ev Event, children ...Control) (found Control, index int) {
	found = nil
	index = -1

	for i := len(children) - 1; i >= 0; i-- {
		child := children[i]
		if child.Hidden() {
			continue
		}
		if EventInside(ev, child) {
			found = child
			index = i
			return found, index
		}
	}
	return nil, -1
}

func HandleWidgetFocus(ev Event, c Control, children ...Control) (focus Control, index int) {
	// look for clicked control using event handling order
	// loose focus by default
	var newFocus Control
	index = -1

	if !EventIsPress(ev) {
		return nil, -1
	}

	newFocus, index = EventInsideOneOf(ev, children...)
	if newFocus == nil {
		return nil, -1
	}

	SetNewFocus(ev.Event(), c, newFocus, children...)
	return newFocus, index
}

func HandleContainerIfNeeded(e Event, c Container) (used bool) {
	// Check if one of the ordered child widgets has to be focused and focus it if so.
	HandleWidgetFocus(e, c, c.Ordered()...)

	// Now let the focus handle it if needed.
	focus := c.Focus()
	if focus != nil && !focus.Hidden() {
		focus.HandleWidget(e)
		return true
	}

	// Not handled.
	return false
}

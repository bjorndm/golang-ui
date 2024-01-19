package ui

import "time"

import "golang.org/x/exp/slices"
import "github.com/hajimehoshi/ebiten/v2/inpututil"
import "github.com/hajimehoshi/ebiten/v2"

// Event represents the base event interface
type Event interface {
	Event() *BasicEvent
	// Origin is the control that the event orignated in.
	// For most events this will be the top level Window.
	Origin() Control
	Modifiers() EventModifiers
}

// EventModifiers are modifiers of the events
type EventModifiers struct {
	Alt     bool
	Control bool
	Shift   bool
	Meta    bool
}

// EventAt is an event that has X and Y coordinates, such as a mouse click or release,
// or a touch or touch release.
type EventAt interface {
	EventAt() (x, y int)
}

// EventPress is an event that can be considered as a press, such as a mouse click,
// a touch, or a gamepad button press.
type EventPress interface {
	EventPress() bool
}

func EventIsPress(e Event) bool {
	press, ok := e.(EventPress)
	return ok && press.EventPress()
}

// BasicEvent is the basic type of all events.
type BasicEvent struct {
	EventOrigin    Control // EventOrigin is the control that the event orignated in.
	EventModifiers         // EventModifiers are keyboard modifiers of the events.
}

// Event implements the Event interface.
func (be *BasicEvent) Event() *BasicEvent {
	return be
}

// Window implements the Event interface.
func (be BasicEvent) Origin() Control {
	return be.EventOrigin
}

// Modifiers implements the Event interface.
func (be BasicEvent) Modifiers() EventModifiers {
	return be.EventModifiers
}

// UpdateEvent is sent regularly for animated widgets
// or media playback
type UpdateEvent struct {
	BasicEvent
	time.Duration
}

// MouseEvent represents a mouse event
type MouseEvent struct {
	BasicEvent
	X, Y int // Absolute coordinates of mouse event
}

// EventAt implement the EventAt interface
func (e MouseEvent) EventAt() (x, y int) {
	return e.X, e.Y
}

var _ EventAt = MouseEvent{}

func EventInside(ev Event, control Control) bool {
	me, ok := ev.(EventAt)
	if !ok {
		return false
	}

	twx, twy := ControlAbsolute(control)
	ww, wh := control.WidgetSize()
	mx, my := me.EventAt()
	return mx >= twx &&
		mx < twx+ww &&
		my >= twy &&
		my < twy+wh
}

func (me MouseEvent) Inside(control Control) bool {
	twx, twy := ControlAbsolute(control)
	ww, wh := control.WidgetSize()
	return me.X >= twx && me.X < twx+ww &&
		me.Y >= twy && me.Y < twy+wh
}

type MouseMoveEvent struct {
	MouseEvent
	MoveX, MoveY int // Motion of mouse.

}

type MouseClickEvent struct {
	MouseEvent
	Button ebiten.MouseButton
}

// EventPress implements the EventPress interface.
func (me *MouseClickEvent) EventPress() bool {
	return true
}

type MouseButton = ebiten.MouseButton

type MouseReleaseEvent struct {
	MouseEvent
	Button MouseButton
	time.Duration
}

type WheelEvent struct {
	MouseEvent
	WheelX, WheelY float64
}

// KeyEvent represents a key event
type KeyEvent struct {
	BasicEvent
	Key Key // Key associated with key event
}

func (ke KeyEvent) Name() string {
	return ebiten.KeyName(ke.Key)
}

type KeyPressEvent struct {
	KeyEvent
}

// EventPress implements the EventPress interface.
func (me *KeyPressEvent) EventPress() bool {
	return true
}

type KeyReleaseEvent struct {
	KeyEvent
}

// CharEvent represents a character input event
type CharEvent struct {
	BasicEvent
	Runes []rune
}

func (ce CharEvent) Text() string {
	return string(ce.Runes)
}

type gamepadState struct {
	GamepadEvent // reuse event for simplicity
	axes         []float64
}

type inputState struct {
	clicked         []MouseButton
	released        []MouseButton
	mouseX          int
	mouseY          int
	keys            []ebiten.Key
	modifiers       EventModifiers
	buttonsPressed  []GamepadButton
	buttonsReleaded []GamepadButton
	connected       []GamepadID
	gamepadStates   []gamepadState
}

var modifierKeys = []Key{KeyAlt, KeyControl, KeyShift, KeyMeta}

func (in *inputState) convertKeyboardInputToEvents(basic *BasicEvent, window *Window, handle func(Event)) {
	// Check for just pressed and released keys
	pressedKeys := inpututil.AppendJustPressedKeys(nil)
	releasedKeys := inpututil.AppendJustReleasedKeys(nil)

	// Handle just pressed and just released modifier keys
	for _, key := range pressedKeys {
		switch key {
		case KeyAlt:
			in.modifiers.Alt = true
		case KeyControl:
			in.modifiers.Control = true
		case KeyShift:
			in.modifiers.Shift = true
		case KeyMeta:
			in.modifiers.Meta = true
		}
	}

	for _, key := range releasedKeys {
		switch key {
		case KeyAlt:
			in.modifiers.Alt = false
		case KeyControl:
			in.modifiers.Control = false
		case KeyShift:
			in.modifiers.Shift = false
		case KeyMeta:
			in.modifiers.Meta = false
		}
	}

	basic.EventModifiers = in.modifiers

	// Handle just pressed and just released keys
	for _, key := range pressedKeys {
		keyEvent := &KeyPressEvent{KeyEvent: KeyEvent{BasicEvent: *basic, Key: key}}
		handle(keyEvent)
	}

	for _, key := range releasedKeys {
		keyEvent := &KeyReleaseEvent{KeyEvent: KeyEvent{BasicEvent: *basic, Key: key}}
		handle(keyEvent)
	}

	in.keys = pressedKeys

	runes := ebiten.AppendInputChars(nil)
	if len(runes) > 0 {
		charEvent := &CharEvent{BasicEvent: *basic, Runes: runes}
		handle(charEvent)
	}
}

func (in *inputState) convertMouseInputToEvents(basic BasicEvent, window *Window, handle func(Event)) {
	// Check for input events
	mouseX, mouseY := ebiten.CursorPosition()

	basicMouseEvent := MouseEvent{BasicEvent: basic, X: mouseX, Y: mouseY}

	if mouseX != in.mouseX && mouseY != in.mouseY {
		mouseEvent := &MouseMoveEvent{MouseEvent: basicMouseEvent}
		mouseEvent.MoveX = mouseX - in.mouseX
		mouseEvent.MoveY = mouseY - in.mouseY
		handle(mouseEvent)
	}
	in.mouseX = mouseX
	in.mouseY = mouseY

	var justClickedButtons []ebiten.MouseButton
	for btn := ebiten.MouseButton0; btn < ebiten.MouseButtonMax; btn++ {
		if inpututil.IsMouseButtonJustPressed(btn) {
			justClickedButtons = append(justClickedButtons, btn)
		}
	}

	// Check for just released mouse buttons
	var justReleasedButtons []ebiten.MouseButton
	for btn := ebiten.MouseButton0; btn < ebiten.MouseButtonMax; btn++ {
		if inpututil.IsMouseButtonJustReleased(btn) {
			justReleasedButtons = append(justReleasedButtons, btn)
		}
	}

	// Handle just clicked and just released mouse buttons
	for _, btn := range justClickedButtons {
		mouseEvent := &MouseClickEvent{MouseEvent: basicMouseEvent, Button: btn}
		handle(mouseEvent)
	}

	for _, btn := range justReleasedButtons {
		mouseEvent := &MouseReleaseEvent{MouseEvent: basicMouseEvent, Button: btn}
		handle(mouseEvent)
	}

	// Mouse wheel.
	wheelX, wheelY := ebiten.Wheel()
	if wheelX != 0 || wheelY != 0 {
		wheelEvent := &WheelEvent{MouseEvent: basicMouseEvent, WheelX: wheelX, WheelY: wheelY}
		handle(wheelEvent)
	}
}

func (in *inputState) convertTouchInputToEvents(basic BasicEvent, window *Window, handle func(Event)) {
	// Check for just pressed and released touches
	pressedTouches := inpututil.AppendJustPressedTouchIDs(nil)
	releasedTouches := inpututil.AppendJustReleasedTouchIDs(nil)

	for _, id := range pressedTouches {
		te := TouchEvent{BasicEvent: basic}
		te.X, te.Y = ebiten.TouchPosition(id)
		tpe := &TouchPressEvent{TouchEvent: te}
		handle(tpe)
	}

	for _, id := range releasedTouches {
		te := TouchEvent{BasicEvent: basic}
		te.X, te.Y = ebiten.TouchPosition(id)
		tre := &TouchReleaseEvent{TouchEvent: te}
		tre.Duration = convertDuration(inpututil.TouchPressDuration(id))
		handle(tre)
	}
}

func convertGamepad(basic BasicEvent, id GamepadID) GamepadEvent {
	ge := GamepadEvent{BasicEvent: basic}
	ge.Name = ebiten.GamepadName(id)
	ge.SDLID = ebiten.GamepadSDLID(id)
	ge.ButtonCount = ebiten.GamepadButtonCount(id)
	ge.AxisCount = ebiten.GamepadAxisCount(id)
	return ge
}

func (in *inputState) convertGamepadInputToEvents(basic BasicEvent, window *Window, handle func(Event)) {

	disconnected := []int{}
	for _, id := range in.connected {
		if inpututil.IsGamepadJustDisconnected(id) {
			ge := convertGamepad(basic, id)
			gde := &GamepadDisconnectEvent{GamepadEvent: ge}
			handle(gde)
			idx := slices.Index(in.connected, id)
			if idx >= 0 {
				disconnected = append(disconnected, idx)
			}
		}
	}

	for _, idx := range disconnected {
		in.connected = slices.Delete(in.connected, idx, idx+1)
		in.gamepadStates = slices.Delete(in.gamepadStates, idx, idx+1)
	}

	connected := inpututil.AppendJustConnectedGamepadIDs(nil)
	for _, id := range connected {
		state := gamepadState{GamepadEvent: convertGamepad(basic, id)}
		state.axes = make([]float64, state.AxisCount)
		in.gamepadStates = append(in.gamepadStates, state)
		in.connected = append(in.connected, id)
	}

	for _, state := range in.gamepadStates {
		ge := state.GamepadEvent
		id := ge.ID
		gde := &GamepadConnectEvent{GamepadEvent: ge}
		handle(gde)

		pressed := inpututil.AppendJustPressedGamepadButtons(id, nil)
		for _, button := range pressed {
			gbe := GamepadButtonEvent{GamepadEvent: ge}
			gbe.Button = button
			gpe := &GamepadButtonPressEvent{GamepadButtonEvent: gbe}
			handle(gpe)
		}

		released := inpututil.AppendJustReleasedGamepadButtons(id, nil)
		for _, button := range released {
			gbe := GamepadButtonEvent{GamepadEvent: ge}
			gbe.Button = button
			gpe := &GamepadButtonReleaseEvent{GamepadButtonEvent: gbe}
			gpe.Duration = convertDuration(inpututil.GamepadButtonPressDuration(id, button))
			handle(gpe)
		}

		for axis := 0; axis < state.AxisCount; axis++ {
			old := state.axes[axis]
			value := ebiten.GamepadAxis(id, axis)
			if value != old {
				gae := &GamepadAxisEvent{GamepadEvent: ge}
				handle(gae)
			}
		}
	}
}

func convertDuration(ticks int) time.Duration {
	return time.Duration((float64(ticks) / float64(ebiten.TPS())) * float64(time.Second))
}

func (in *inputState) convertInputToEvents(window *Window, handle func(Event)) {
	// Update state
	basic := BasicEvent{EventOrigin: window}

	// Send an update event every time.
	update := &UpdateEvent{BasicEvent: basic}
	update.Duration = convertDuration(1) // one tick
	handle(update)

	in.convertKeyboardInputToEvents(&basic, window, handle)
	in.convertMouseInputToEvents(basic, window, handle)
	in.convertTouchInputToEvents(basic, window, handle)

}

// AwayEvent is sent to notify a child widget that the fous has gone away from
// the child. It should be sent to from a parent widge to a child widget,
// to notify the child that Focused will become the focused widget.
// If Focused is nil, then the focus is still lost but not gained by any other
// widget.
type AwayEvent struct {
	BasicEvent
	Focused Control // Focused is the control that received the focus. If nil, it is unknown.
}

// FocusEvent is sent to notify a child widget that it has received the focus.
// It should be sent to from a parent widget to a child widget.
// If requested is nil, then the focus should be dropped.
type FocusEvent struct {
	BasicEvent
	Focused Control // Focused is the control that received the focus. If nil, it is unknown.
}

type GamepadID = ebiten.GamepadID
type GamepadButton = ebiten.GamepadButton
type StandardGamepadButton = ebiten.StandardGamepadButton
type StandardGamepadAxis = ebiten.StandardGamepadAxis

type GamepadEvent struct {
	BasicEvent
	ID          GamepadID
	ButtonCount int
	Name        string
	SDLID       string
	AxisCount   int
}

type GamepadConnectEvent struct {
	GamepadEvent
}

type GamepadDisconnectEvent struct {
	GamepadEvent
}

type GamepadButtonEvent struct {
	GamepadEvent
	Button         GamepadButton
	StandardButton StandardGamepadButton
}

type GamepadAxisEvent struct {
	GamepadEvent
	Axis         int
	StandardAxis StandardGamepadAxis
}

type GamepadButtonPressEvent struct {
	GamepadButtonEvent
}

// EventPress implements the EventPress interface.
func (me *GamepadButtonPressEvent) EventPress() bool {
	return true
}

type GamepadButtonReleaseEvent struct {
	GamepadButtonEvent
	time.Duration
}

type TouchID = ebiten.TouchID

type TouchEvent struct {
	BasicEvent
	ID   TouchID
	X, Y int
	time.Duration
}

// EventAt implement the EventAt interface
func (e TouchEvent) EventAt() (x, y int) {
	return e.X, e.Y
}

type TouchPressEvent struct {
	TouchEvent
}

// EventPress implements the EventPress interface.
func (me *TouchPressEvent) EventPress() bool {
	return true
}

type TouchReleaseEvent struct {
	TouchEvent
	time.Duration
}

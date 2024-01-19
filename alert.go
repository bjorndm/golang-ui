package ui

type alertKind int

const (
	alertKindMessage alertKind = iota
	alertKindError
	alertKindConfirm
)

// Alert prepares modal dialog boxes that can only be used for warnings or
// simple OK/cancel questions.
type Alert struct {
	*Pane         // An alert is a pane.
	box      *Box // Use a box as the base widget/layout.
	icon     *Picture
	text     *TextWidget
	ok       *Button
	cancel   *Button
	onResult func(a *Alert)
	result   DialogResult
}

const alertMessage = "information"
const alertExclamation = "exclamation"
const alertQuestion = "question"

func (a *Alert) Result() DialogResult {
	return a.result
}

func (a *Alert) SetResult(res DialogResult) {
	a.result = res
}

func (a *Alert) SendResult(res DialogResult) {
	a.SetResult(res)
	if a.onResult != nil {
		a.onResult(a)
	}
	a.closePane()
}

func copyStyleColors(to, from Control) {
	style := to.Style()
	// Take the colors of from, but not the other style parts.
	style.Fill.Color = from.Style().Fill.Color
	style.Color = from.Style().Color
	to.SetStyle(&style)
}

func newAlert(title, message string, kind alertKind) *Alert {
	a := &Alert{}

	a.box = NewBox()
	a.text = NewTextWidget(message)
	// TODO: size

	paneW, paneH := theme.Alert.Size.Width.Int(), theme.Alert.Size.Height.Int()

	a.Pane = NewPane(title, paneW, paneH, false)
	a.Pane.SetChild(a.box)
	a.Pane.OnClosing(func(p *Pane) {
		a.SendResult(DialogResultCancel)
	})
	a.Pane.SetPermanent(false)
	a.Pane.SetPreserved(false)

	a.SetStyle(theme.Alert)

	copyStyleColors(a.box, a)

	icon := alertMessage

	switch kind {
	case alertKindMessage:
		icon = alertMessage
		a.ok = NewButton("OK")
	case alertKindError:
		icon = alertExclamation
		a.ok = NewButton("OK")
	case alertKindConfirm:
		a.cancel = NewButton("Cancel")
		a.ok = NewButton("OK")
		icon = alertQuestion
	}
	a.icon = NewPictureWithIcon("", icon)
	a.icon.SetBorderless(true)
	copyStyleColors(a.icon, a)

	iconText := NewTray()

	iconText.Append(a.icon)
	iconText.Append(a.text)
	a.box.Append(iconText)
	copyStyleColors(iconText, a)

	buttons := NewTray()
	copyStyleColors(buttons, a)

	if a.ok != nil {
		copyStyleColors(a.ok, a)
		a.ok.OnClicked(func(b *Button) {
			a.SendResult(DialogResultOK)
			a.Pane.closePane()
		})
		buttons.Append(a.ok)
	}
	if a.cancel != nil {
		copyStyleColors(a.cancel, a)
		a.cancel.OnClicked(func(b *Button) {
			a.SendResult(DialogResultCancel)
			a.Pane.closePane()
		})
		buttons.Append(a.cancel)
	}
	a.box.Append(buttons)

	return a
}

// Display displays the alert above over, returning by calling callback.
// If the user closes an Alert, it will be destroyed, so, after calling
// this function, the Alert a should not be reused.
func (a *Alert) Display(over Control, callback func(DialogResult)) {
	a.DisplayAlert(over, AlertCallback(callback))
}

// DisplayAlert displays the alert above over, returning by calling callback.
// If the user closes an Alert, it will be destroyed, so, after calling
// this function, the Alert a should not be reused.
func (a *Alert) DisplayAlert(over Control, callback func(a *Alert)) {
	a.Pane.closed = false
	a.onResult = callback
	ControlStartDialog(over, a.Pane, a.title, true)
}

// NewAlert creates an alert dialog without displaying it.
// If the user closes an Alert, it will be destroyed, so after calling Display
// Alert cannot be reused.
func NewAlert(title, message string) *Alert {
	return newAlert(title, message, alertKindMessage)
}

func NewErrorAlert(title, message string) *Alert {
	return newAlert(title, message, alertKindError)
}

func NewConfirmAlert(title, message string) *Alert {
	return newAlert(title, message, alertKindConfirm)
}

// ShowAlert creates an alert dialog and displays it above over.
// If the user closes an Alert, it will be destroyed.
func ShowAlert(over Control, title, message string, res func(DialogResult)) {
	NewAlert(title, message).Display(over, res)
}

func ShowErrorAlert(over Control, title, message string, res func(DialogResult)) {
	NewErrorAlert(title, message).Display(over, res)
}

func ShowConfirmAlert(over Control, title, message string, res func(DialogResult)) {
	NewConfirmAlert(title, message).Display(over, res)
}

// AlertCallback can be used if you don't need the whole dialog and want
// to use a dialog callback with a result only. Pass in your desired callback
// and this function will return an adpted function for use with ShowDialog.
func AlertCallback(rcb func(DialogResult)) func(*Alert) {
	if rcb == nil {
		return nil
	}
	return func(d *Alert) {
		result := d.Result()
		rcb(result)
	}
}

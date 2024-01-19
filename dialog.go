package ui

type dialogKind int

const (
	dialogKindMessage dialogKind = iota
	dialogKindError
	dialogKindConfirm
)

// Dialog is a dialog box that can be used for warnings or simple
// OK/cancel questions. The Pane that a Dialog uses is preserved on closing.
type Dialog struct {
	*Pane            // A dialog is a pane.
	box      *Box    // Box for the layout and because a pane only can have one child.
	child    Control // Child widget with dialog contents.
	tray     *Tray   // Tray with buttons to display on the bottom.
	onResult func(*Dialog)
	result   DialogResult
}

// DialogResult is the result of showing an Dialog or dialog, decided by the
// button the user clicks.
type DialogResult string

const (
	DialogResultNone   DialogResult = ""
	DialogResultCancel DialogResult = "cancel"
	DialogResultOK     DialogResult = "ok"
	DialogResultYes    DialogResult = "yes"
	DialogResultNo     DialogResult = "no"
)

func (r DialogResult) String() string {
	return string(r)
}

const dialogMessage = "information"
const dialogExclamation = "exclamation"
const dialogQuestion = "question"

func (a *Dialog) SetResult(res DialogResult) {
	a.result = res
}

func (a *Dialog) SendResult(res DialogResult) {
	a.SetResult(res)
	if a.onResult != nil {
		a.onResult(a)
	}
}

func (a *Dialog) Result() DialogResult {
	return a.result
}

func newDialog(title string, child Control) *Dialog {
	a := &Dialog{}

	a.box = NewBox()
	a.tray = NewTray()
	a.child = child
	style := theme.Pane

	a.Pane = NewPane(title, style.Size.Width.Int(), style.Size.Height.Int(), false)
	a.Pane.SetChild(a.box)
	a.Pane.SetStyle(style)
	a.Pane.SetPreserved(true)
	a.Pane.OnClosing(func(p *Pane) {
		a.SendResult(DialogResultCancel)
	})

	if a.child != nil {
		a.box.Append(a.child)
	}
	a.box.Append(a.tray)
	return a
}

func (a *Dialog) AddButton(title string, result DialogResult) *Dialog {
	button := NewButton(title)
	button.OnClicked(func(b *Button) {
		a.SendResult(result)
		a.closePane()
	})
	a.tray.Append(button)
	return a
}

func (a *Dialog) AddButtonKeepOpen(title string, result DialogResult) *Dialog {
	button := NewButton(title)
	button.OnClicked(func(b *Button) {
		a.SendResult(result)
	})
	a.tray.Append(button)
	return a
}

func (a *Dialog) DisplayDialog(over Control, cb func(*Dialog)) {
	a.Pane.closed = false
	if a.Pane.Hidden() {
		// We are reusing a hiddden dialog, lay it out again.
		NeedLayout(a)
	}
	a.Pane.Show()
	a.onResult = cb
	ControlStartDialog(over, a.Pane, a.title, true)
}

func (a *Dialog) Display(over Control, cb func(DialogResult)) {
	a.DisplayDialog(over, DialogCallback(cb))
}

// Create a new dialog. The Pane is set to Preserved so it can be reused
// even if closed.
func NewDialog(title string, child Control) *Dialog {
	return newDialog(title, child)
}

func ShowDialog(over Control, title string, child Control, cb func(DialogResult)) *Dialog {
	dialog := NewDialog(title, child)
	dialog.Display(over, cb)
	return dialog
}

func ShowDialogWithDialog(over Control, title string, child Control, cb func(*Dialog)) *Dialog {
	dialog := NewDialog(title, child)
	dialog.DisplayDialog(over, cb)
	return dialog
}

// DialogCallback can be used if you don't need the whole dialog and want
// to use a dialog callback with a result only. Pass in your desired callback
// and this function will return an adpted function for use with ShowDialog.
func DialogCallback(rcb func(DialogResult)) func(*Dialog) {
	if rcb == nil {
		return nil
	}
	return func(d *Dialog) {
		result := d.Result()
		rcb(result)
	}
}

package ui

// Ability are the abilities of a Window or Pane.
// The default Ability is not modal, not plain, not rigid,
// not fixed, not full, not permanent, not preserved.
type Ability struct {
	modal     bool
	plain     bool
	rigid     bool
	fixed     bool
	full      bool
	permanent bool
	preserved bool
	onChanged func(a *Ability)
}

// SetModal sets whether or not the widget is modal
// If true the Window or Pane will not be modal.
// If false Window or Pane may be modal screen if the platform allows it.
func (a *Ability) SetModal(modal bool) {
	a.modal = modal
	a.changed()
}

// Modal returns wheter or not the widget is modal.
// If true the Window or Pane will not be modal.
// If false Window or Pane may be modal if the platform allows it.
func (a Ability) Modal() bool {
	return a.modal
}

// SetFull sets whether or not the widget is full screen
// If true the Window or Pane will not be full screen.
// If false Window or Pane may be full screen if the platform allows it.
func (a *Ability) SetFull(full bool) {
	a.full = full
	a.changed()
}

// Full returns wheter or not the widget is full screen.
// If true the Window or Pane will not be full screen.
// If false Window or Pane may be full screen if the platform allows it.
func (a Ability) Full() bool {
	return a.full
}

// SetPlain sets whether or not the widget is plain, and without decorations.
// If true the Window or Pane will not have decorations.
// If false Window or Pane may have decorations if the platform allows it.
func (a *Ability) SetPlain(plain bool) {
	a.plain = plain
	a.changed()
}

// Plain returns wheter or not the widget is plain and without decorations.
// If true the Window or Pane will not have decorations.
// If false Window or Pane may have decorations if the platform allows it.
func (a Ability) Plain() bool {
	return a.plain
}

// SetRigid sets whether or not the widget is rigid in size.
// If true the Window or Pane cannot be resized by the user.
// If false Window or Pane can be resized by the user.
func (a *Ability) SetRigid(rigid bool) {
	a.rigid = rigid
	a.changed()
}

// Rigid returns wheter or not the widget is rigid in place.
func (a Ability) Rigid() bool {
	return a.rigid
}

// SetFixed sets whether or not the widget is fixed in size.
// If true the Window or Pane will not be movable.
// If false Window or Pane will be movable.
func (a *Ability) SetFixed(fixed bool) {
	a.fixed = fixed
	a.changed()
}

// Fixed returns wheter or not the widget is fixed in place.
func (a Ability) Fixed() bool {
	return a.fixed
}

// SetPermanent sets whether or not the widget is permanent.
// If true the Window or Pane will cannot be closed by the user directly.
// If false the Window or Pane will can be closed by the user directly.
func (a *Ability) SetPermanent(permanent bool) {
	a.permanent = permanent
	a.changed()
}

// Permanent returns wheter or not the widget is permanent.
// If true the Window or Pane will cannot be closed by the user directly.
// If false the Window or Pane will can be closed by the user directly.
func (a Ability) Permanent() bool {
	return a.permanent
}

// SetPreserved sets whether or not the widget will be preserved on closing.
// If true the Window or Pane will not be destroyed on closing so it can be reused.
// If false the Window or Pane will destroyed on closing so it cannot be reused.
func (a *Ability) SetPreserved(preserved bool) {
	a.preserved = preserved
	a.changed()
}

// Preserved returns wheter or not the widget will be preserved on closing.
func (a Ability) Preserved() bool {
	return a.preserved
}

// setOnChanged sets a callback on the change of the ability a.
// this method as well as the callback are for internal use by Pane and Window.
func (a *Ability) setOnChanged(cb func(a *Ability)) {
	a.onChanged = cb
}

// canged calls the change callback if not nil
func (a *Ability) changed() {
	if a.onChanged != nil {
		a.onChanged(a)
	}
}

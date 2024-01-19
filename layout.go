package ui

type Layout interface {
	Layout(parent Control, controls ...Control)
}

func VerticalLayout(parent Control, children ...Control) {

}

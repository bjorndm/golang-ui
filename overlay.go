package ui

import "golang.org/x/exp/slices"

type BasicOverlayer struct {
	overlays []Control
}

// HandleEventForOverlays returns whether or not the event was used by the
// overlay widgets or not.
func (w BasicOverlayer) HandleEventForOverlays(e Event) bool {
	if len(w.overlays) > 0 {
		for _, overlay := range w.overlays {
			// TODO: support touches as well.
			if mc, ok := e.(*MouseClickEvent); ok {
				if mc.Inside(overlay) {
					overlay.HandleWidget(e)
					return true
				}
			} else {
				overlay.HandleWidget(e)
				return true
			}
		}
	}
	return false
}

// DrawOverlays draw the overlays. This should be drawn over the other widgets.
func (w BasicOverlayer) DrawOverlays(screen *Graphic) {
	if len(w.overlays) > 0 {
		for _, overlay := range w.overlays {
			overlay.DrawWidget(screen)
		}
	}
}

// StartOverlay requests that the widget c will become an overlay in the Overlayer.
func (w *BasicOverlayer) StartOverlay(c Control) {
	w.overlays = append(w.overlays, c)
}

// EndOverlay requests that the widget c will not be an overlay anomore.
func (w *BasicOverlayer) EndOverlay(c Control) {
	index := slices.IndexFunc(w.overlays, func(seek Control) bool { return seek == c })
	if index < 0 {
		return
	}
	w.overlays = slices.Delete(w.overlays, index, index+1)
}

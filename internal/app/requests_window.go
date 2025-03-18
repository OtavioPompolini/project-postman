package app

import (
	"log"

	"github.com/awesome-gocui/gocui"

	"github.com/OtavioPompolini/project-postman/internal/ui"
)

// Why not to save a ref to ui.GUI and ui.Window on every Window implementation?????
// Since every F method needs both params...
type RequestsWindow struct {
	isActive     bool
	name         string
	x, y, h, w   int
	stateService StateService
	loadRequests func() error
}

func NewRequestsWindow(GUI *ui.UI, stateService StateService) *ui.Window {
	a, b := GUI.Size()
	return ui.NewWindow(
		&RequestsWindow{
			name:         "RequestsWindow",
			x:            0,
			y:            0,
			h:            b - 1,
			w:            a * 20 / 100,
			stateService: stateService,
			isActive:     true,
		},
		true,
	)
}

func (w RequestsWindow) Name() string {
	return w.name
}

func (w *RequestsWindow) Setup(ui ui.UI, v ui.Window) {
	ui.SelectWindow(&v)
	v.SetTitle("Requests:")
	v.SetSelectedBgColor(gocui.ColorRed)
	v.SetHightlight(true)
	w.ReloadContent(&ui, &v)
}

func (w *RequestsWindow) Update(ui ui.UI, v ui.Window) {
}

func (w *RequestsWindow) Size() (x, y, width, height int) {
	return w.x, w.y, w.x + w.w, w.y + w.h
}

func (w *RequestsWindow) IsActive() bool {
	return w.isActive
}

func (w *RequestsWindow) SetKeybindings(ui *ui.UI, win *ui.Window) error {
	if err := ui.NewKeyBinding(w.Name(), 'j', func(g *gocui.Gui, v *gocui.View) error {
		w.stateService.SelectNext()
		w.ReloadContent(ui, win)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'k', func(g *gocui.Gui, v *gocui.View) error {
		w.stateService.SelectPrev()
		w.ReloadContent(ui, win)
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'p', func(g *gocui.Gui, v *gocui.View) error {
		w.stateService.state.selectedRequest.Execute()
		return nil
	}); err != nil {
		return err
	}

	if err := ui.NewKeyBinding(w.Name(), 'n', func(g *gocui.Gui, v *gocui.View) error {
		win, err := ui.GetWindow("CreateRequestWindow")
		if err != nil {
			return err
		}

		win.OpenWindow()

		return nil
	}); err != nil {
		return err
	}

	// TODO: BUT I STILL HAVEN'T FOUND WHAT I'M LOOKING FOR...
	// Handle change window with a "const" and not a string
	// and need to abstract gocui
	if err := ui.NewKeyBinding(w.Name(), gocui.KeyEnter, func(g *gocui.Gui, v *gocui.View) error {
		if len(w.stateService.state.requests) <= 0 {
			return nil
		}

		ui.SelectWindowByName("RequestDetailsWindow")

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *RequestsWindow) OnDeselect(ui ui.UI, v ui.Window) error {
	return nil
}

func (w *RequestsWindow) OnSelect(ui ui.UI, v ui.Window) error {
	return nil
}

// Make ReloadWindowContent an IWindow Interface func? so other widnows can Reload others
func (w *RequestsWindow) ReloadContent(ui *ui.UI, v *ui.Window) {
	v.ClearWindow()

	cursorPosition := 0
	cursorPositionFound := false
	requests := w.stateService.state.requests

	lines := []string{}

	for _, r := range requests {
		if !cursorPositionFound {
			if r.Id == w.stateService.state.selectedRequest.Id {
				cursorPositionFound = true
			} else {
				cursorPosition += 1
			}
		}
		lines = append(lines, r.Name)
	}

	v.WriteLines(lines)

	err := v.SetCursor(0, cursorPosition)
	if err != nil {
		log.Panic(err)
	}
}

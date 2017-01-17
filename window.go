package meson

import (
	"errors"
	"github.com/koron/go-dproxy"
)

// Rect represents a rectangular region on the screen
type Rect struct {
	Width  int `json:"width"`  // Width in pixels
	Height int `json:"height"` // Height in pixels
	Left   int `json:"left"`   // Left is offset from left in pixel
	Top    int `json:"top"`    // Left is offset from top in pixels
}

type Window struct {
	object
}

func newWindow(id int64) *Window {
	win := &Window{object: newObject(id, objWindow)}
	addObject(id, win)
	// register default handler
	return win
}

// WindowOptions contains options for creating windows
type WindowOptions struct {
	Title            string `json:"title"` // String to display in title bar
	IconPath         string `json:"icon_path"`
	Shape            Rect   `json:"shape"`       // Initial size and position of window
	TitleBar         bool   `json:"titleBar"`    // Whether the window title bar
	Frame            bool   `json:"has_frame"`   // Whether the window has a frame
	Resizable        bool   `json:"resizable"`   // Whether the window border can be dragged to change its shape
	CloseButton      bool   `json:"closeButton"` // Whether the window has a close button
	MinButton        bool   `json:"minButton"`   // Whether the window has a miniaturize button
	FullScreenButton bool   `json:"maxButton"`   // Whether the window has a full screen button
	//	Menu             []MenuEntry
}

// FramedWindowOptions contains options for an "ordinary" window with title bar,
// frame, and min/max/close buttons.
var FramedWindowOptions = WindowOptions{
	Shape:            Rect{Width: 800, Height: 600, Left: 100, Top: 100},
	TitleBar:         true,
	Frame:            true,
	Resizable:        true,
	CloseButton:      true,
	MinButton:        true,
	FullScreenButton: true,
	Title:            "Meson",
}

// NewBrowserWindow Create and control browser windows.
func NewBrowserWindow(opt *WindowOptions) (*Window, error) {
	if !apiReady {
		return nil, errors.New("meson api is not ready yet")
	}
	cmd := makeCreateCommand(objWindow, opt)

	response, err := sendMessage(&cmd)
	if err != nil {
		return nil, err
	}

	id, err := dproxy.New(response).Int64()
	if err != nil {
		return nil, err
	}

	return newWindow(id), nil
}

//LoadURLOptions is optional parameter for Window.LoadURL and WebContents.LoadURL
type LoadURLOptions struct {
	HTTPReferer  string `json:"httpReferrer"` // A HTTP Referrer url.
	UserAgent    string `json:"userAgent"`    // A user agent originating the request.
	ExtraHeaders string `json:"extraHeaders"` // Extra headers separated by “\n”
}

//LoadURL is same as WebContents.LoadURL
func (w *Window) LoadURL(url string) {
	cmd := makeCallCommand(w.objType, w.id, "loadURL", url)
	postMessage(&cmd)
	//_, err := sendMessage(&cmd)
	//return err
}

func (w *Window) LoadURLWithOptions(url string, opt *LoadURLOptions) error {
	cmd := makeCallCommand(objWindow, w.id, "loadURL", url, opt)
	_, err := sendMessage(&cmd)
	return err
}

func (w *Window) Close() {
	cmd := makeCallCommand(objWindow, w.id, "close")
	postMessage(&cmd)
}

func (w *Window) OpenDevTool() {
	// TODO: options??
	cmd := makeCallCommand(objWindow, w.id, "openDevTools")
	postMessage(&cmd)
}

func (w *Window) CloseDevTool() {
	cmd := makeCallCommand(objWindow, w.id, "closeDevTools")
	postMessage(&cmd)
}

func (w *Window) IsDevToolOpened() bool {
	cmd := makeCallCommand(objWindow, w.id, "isDevToolsOpened")
	r, err := sendMessage(&cmd)
	if err != nil {
		return false
	}
	b, _ := dproxy.New(r).Bool()
	return b
}

//------------------------------------------------------------------------
// Callbacks

func (w *Window) OnWindowClose(callback CommonPreventableCallbackHandler) {
	const en = "close"
	w.addCallback(en, commonPreventableCallbackItem{f: callback})
}

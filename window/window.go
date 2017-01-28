package window

import (
	"errors"
	evt "github.com/go-meson/meson/event"
	"github.com/go-meson/meson/internal/binding"
	"github.com/go-meson/meson/internal/command"
	"github.com/go-meson/meson/internal/event"
	"github.com/go-meson/meson/internal/object"
	"github.com/go-meson/meson/util"
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
	object.Object
}

func newWindow(id int64) *Window {
	win := &Window{Object: object.NewObject(id, binding.ObjWindow)}
	object.AddObject(id, win)
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
	Title:            util.ApplicationName,
}

// NewBrowserWindow Create and control browser windows.
func NewBrowserWindow(opt *WindowOptions) (*Window, error) {
	if !command.APIReady {
		return nil, errors.New("meson api is not ready yet")
	}
	cmd := command.MakeCreateCommand(binding.ObjWindow, opt)

	response, err := command.SendMessage(&cmd)
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
func (w *Window) LoadURL(url string) error {
	cmd := command.MakeCallCommand(w.ObjType, w.Id, "loadURL", url)
	return command.PostMessage(&cmd)
}

func (w *Window) LoadURLWithOptions(url string, opt *LoadURLOptions) error {
	cmd := command.MakeCallCommand(w.ObjType, w.Id, "loadURL", opt)
	return command.PostMessage(&cmd)
}

func (w *Window) Close() {
	cmd := command.MakeCallCommand(w.ObjType, w.Id, "close")
	command.PostMessage(&cmd)
}

func (w *Window) OpenDevTool() {
	// TODO: options??
	cmd := command.MakeCallCommand(w.ObjType, w.Id, "OpenDevTools")
	if err := command.PostMessage(&cmd); err != nil {
		panic(err)
	}
}

func (w *Window) CloseDevTool() {
	cmd := command.MakeCallCommand(w.ObjType, w.Id, "closeDevTools")
	command.PostMessage(&cmd)
}

func (w *Window) IsDevToolOpened() bool {
	cmd := command.MakeCallCommand(w.ObjType, w.Id, "isDevToolsOpened")
	r, err := command.SendMessage(&cmd)
	if err != nil {
		return false
	}
	b, _ := dproxy.New(r).Bool()
	return b
}

//------------------------------------------------------------------------
// Callbacks

func (w *Window) OnWindowClose(callback evt.CommonPreventableCallbackHandler) {
	const en = "close"
	event.AddCallback(&w.Object, en, event.CommonPreventableCallbackItem{F: callback})
}

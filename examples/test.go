package main

import (
	"log"
	"os"

	"github.com/go-meson/meson"
	"github.com/go-meson/meson/app"
	"github.com/go-meson/meson/dialog"
	"github.com/go-meson/meson/logger"
	"github.com/go-meson/meson/menu"
	"github.com/go-meson/meson/object"
	"github.com/go-meson/meson/util"
	"github.com/go-meson/meson/window"
)

var (
	counter = 0
)

func onClick(mi *menu.MenuItemTemplate, w *window.Window) {
	log.Printf("clicked: %#v\n", mi)
	str := "This app is running in bundle : "
	dialog.ShowMessageBox(w, str, "Test", dialog.MessageBoxTypeInfo, nil)
}

func onOpenDevTool(mi *menu.MenuItemTemplate, w *window.Window) {
	log.Printf("opendev: %#v, %#v", mi, w)
	if w.IsDevToolOpened() {
		w.CloseDevTool()
	} else {
		w.OpenDevTool()
	}
}

var mainMenu = menu.MenuTemplate{
	{Label: "Test1111111", SubMenu: menu.MenuTemplate{
		{Label: "Test1-1", Click: onClick},
		{Label: "Test1-2"},
		{Label: "Quit", Role: "quit"}}},
	{Label: "Test22222", SubMenu: menu.MenuTemplate{
		{Label: "openDevTool", Click: onOpenDevTool},
		{Label: "Test2-2"}}},
}

type WinUserData struct {
	doClosing bool
}

func main() {
	if err := logger.SetFileLogger("/Users/yoshikawa/.go/src/github.com/go-meson/meson/test.log"); err != nil {
		log.Fatal(err)
		return
	}
	logger.RedirectStdout()
	logger.RedirectStderr()
	log.Printf("bundlePath = %s\n", util.GetApplicationBundlePath())
	meson.MainLoop(os.Args, func(a *app.App) {
		//meson.ShowMessageBox(nil, "This is Menu callback", "Test", meson.MessageBoxTypeInfo, nil)

		m, err := menu.NewMenuWithTemplate(mainMenu)
		log.Printf("menu: %#v, err: %#v\n", m, err)
		if err != nil {
			log.Fatal(err)
			app.Exit(-1)
		}
		a.SetApplicationMenu(m)

		a.OnWindowCloseAll(func(sender object.ObjectRef) {
			log.Println("**** WindowCloseAll ***")
			if a := sender.(*app.App); a != nil {
				app.Exit(0)
			}
		})
		log.Println("Called Init Handler!")
		opt := window.FramedWindowOptions
		opt.Shape.Width = 320
		opt.Shape.Height = 240
		win, err := window.NewBrowserWindow(&opt)
		if err != nil {
			log.Printf("Create window fail: %q", err)
			return
		}
		//win.OpenDevTool()
		win.UserData = &WinUserData{}
		log.Printf("win = %#v\n", win)
		win.OnWindowClose(func(sender object.ObjectRef) bool {
			ud := win.UserData.(*WinUserData)
			if ud.doClosing {
				return false
			}
			dialog.ShowMessageBoxAsync(win, "really quit?", "Quit?", dialog.MessageBoxTypeQuestion, nil, func(buttonId int, err error) {
				if err != nil {
					log.Panicf("err!: %#v", err)
				}
				if buttonId == 0 {
					ud.doClosing = true
					win.Close()
				}
			})
			return true
		})
		//win.LoadURL("http://www.google.co.jp")
		win.LoadURL("file:////Users/yoshikawa/.go/src/github.com/go-meson/meson/test.html")
		/*
			if err != nil {
				log.Printf("LoadURL fail: %q", err)
				return
			}
			log.Printf("***** LoadURL success")
		*/
	})
}

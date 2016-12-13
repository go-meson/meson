package main

import (
	"log"
	"os"

	"github.com/go-meson/meson"
)

var (
	counter = 0
)

func onClick(mi *meson.MenuItemTemplate, w *meson.Window) {
	log.Printf("clicked: %#v\n", mi)
	meson.ShowMessageBox(nil, "This is Menu callback", "Test", meson.MessageBoxTypeInfo, nil)
}

var menu = meson.MenuTemplate{
	{Label: "Test1111111", SubMenu: meson.MenuTemplate{
		{Label: "Test1-1", Click: onClick},
		{Label: "Test1-2"},
		{Label: "Quit", Role: "quit"}}},
	{Label: "Test22222", SubMenu: meson.MenuTemplate{
		{Label: "Test2-1"},
		{Label: "Test2-2"}}},
}

type WinUserData struct {
	doClosing bool
}

func main() {
	meson.MainLoop(os.Args, func(app *meson.App) {
		//meson.ShowMessageBox(nil, "This is Menu callback", "Test", meson.MessageBoxTypeInfo, nil)
		m, err := meson.NewMenuWithTemplate(menu)
		log.Printf("menu: %#v, err: %#v\n", m, err)
		if err != nil {
			log.Fatal(err)
			meson.Exit(-1)
		}
		app.SetApplicationMenu(m)

		app.OnWindowCloseAll(func(sender meson.ObjectRef) {
			log.Println("**** WindowCloseAll ***")
			if app := sender.(*meson.App); app != nil {
				meson.Exit(0)
			}
		})
		log.Println("Called Init Handler!")
		opt := meson.FramedWindowOptions
		opt.Shape.Width = 320
		opt.Shape.Height = 240
		win, err := meson.NewBrowserWindow(&opt)
		if err != nil {
			log.Printf("Create window fail: %q", err)
			return
		}
		win.UserData = &WinUserData{}
		log.Printf("win = %#v\n", win)
		win.OnWindowClose(func(sender meson.ObjectRef) bool {
			ud := win.UserData.(*WinUserData)
			if ud.doClosing {
				return false
			}
			meson.ShowMessageBoxAsync(win, "really quit?", "Quit?", meson.MessageBoxTypeQuestion, nil, func(buttonId int, err error) {
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
		win.LoadURL("http://www.google.co.jp")
		/*
			if err != nil {
				log.Printf("LoadURL fail: %q", err)
				return
			}
			log.Printf("***** LoadURL success")
		*/
	})
}

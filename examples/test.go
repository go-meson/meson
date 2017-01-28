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
	"net/url"
	"os/user"
	"path/filepath"
)

func onClick(mi *menu.ItemTemplate, w *window.Window) {
	log.Printf("clicked: %#v\n", mi)
	str := "This app is running in bundle : "
	dialog.ShowMessageBox(w, str, "Test", dialog.MessageBoxTypeInfo, nil)
}

func onOpenDevTool(mi *menu.ItemTemplate, w *window.Window) {
	log.Printf("opendev: %#v, %#v", mi, w)
	if w.IsDevToolOpened() {
		w.CloseDevTool()
	} else {
		w.OpenDevTool()
	}
}

func onClickDialogTest(mi *menu.ItemTemplate, w *window.Window) {
	opt := dialog.MessageBoxOpt{
		Buttons:   []string{"OK", "Cancel", "Foo1", "Foo2"},
		Detail:    "some details",
		DefaultID: 0,
		CancelID:  1,
		NoLink:    true,
	}
	dialog.ShowMessageBox(w, "This is options test", "Option test", dialog.MessageBoxTypeQuestion, &opt)
}

var mainMenu = menu.Template{
	{Label: "Test1111111", SubMenu: menu.Template{
		{Label: "Test1-1", Click: onClick},
		{Label: "Test1-2", Click: onClickDialogTest},
		{Label: "Quit", Role: "quit"}}},
	{Label: "Test22222", SubMenu: menu.Template{
		{Label: "openDevTool", Click: onOpenDevTool},
		{Label: "Test2-2"}}},
}

type winUserData struct {
	doClosing bool
}

func main() {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
		return
	}
	if err := logger.SetFileLogger(filepath.Join(u.HomeDir, util.ApplicationName+".log")); err != nil {
		log.Fatal(err)
		return
	}
	logger.RedirectStdout()
	logger.RedirectStderr()
	log.Printf("bundlePath = %s\n", util.ApplicationBundlePath)
	meson.MainLoop(os.Args, func(a *app.App) {
		m, err := menu.NewWithTemplate(mainMenu)
		if err != nil {
			log.Fatal(err)
			app.Exit(-1)
		}
		menu.SetApplicationMenu(m)

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
		win.UserData = &winUserData{}
		win.OpenDevTool()

		win.OnWindowClose(func(sender object.ObjectRef) bool {
			ud := win.UserData.(*winUserData)
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
		u := url.URL{
			Scheme: "file",
			Path:   filepath.ToSlash(filepath.Join(util.ApplicationAssetsPath, "test.html")),
		}
		win.LoadURL(u.String())
	})
}

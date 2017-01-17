package meson

import (
	"errors"
	"fmt"
	"github.com/koron/go-dproxy"
	"log"
)

type MenuRole struct {
	Label             string
	Accelerator       string
	WindowMethod      string
	WebContentsMethod string
	AppMethod         string
}

// platform dependents
const (
	MenuLabelAbout = "About {{AppName}}"
	MenuLabelClose = "Close Window"
	MenuLabelQuit  = "Quit {{AppName}}"

	MenuAcceleratorQuit             = "CommandOrControl+Q"
	MenuAcceleratorRedo             = "Shift+CommandOrControl+Z"
	MenuAcceleratorToggleFullscreen = "Control+Command+F"
)

const (
	RoleAbout string = "about" //RoleAbout map to the orderFrontStandardAboutPanel action
	//RoleHide map to the hide action
	RoleHide = "hide"
	//RoleHideOthers map to the hideOtherApplications action
	RoleHideOthers = "hideothers"
	//RoleUnHide map to the unhideAllApplications action
	RoleUnHide = "unhide"
	//RoleStartSpeaking map to the startSpeaking action
	RoleStartSpeaking = "startspeaking"
	//RoleStopSpeaking map to the stopSpeaking action
	RoleStopSpeaking = "stopspeaking"
	//RoleFront map to the arrangeInFront action
	RoleFront = "front"
	//RoleZoom map to the performZoom action
	RoleZoom = "zoom"
	//RoleWindow is the submenu of a “Window” menu
	RoleWindow = "window"
	//RoleHelp is the submenu of a “Help” menu
	RoleHelp = "help"
	//RoleServices is the submenu of a “Services” menu
	RoleServices = "services"
)

// platform dependents
var menuRolePlatform = map[string]MenuRole{
	RoleAbout:         MenuRole{Label: MenuLabelAbout},
	RoleHide:          MenuRole{Label: "Hide {{AppName}}", Accelerator: "Command+H"},
	RoleHideOthers:    MenuRole{Label: "Hide Others", Accelerator: "Command+Alt+H"},
	RoleUnHide:        MenuRole{Label: "Show All"},
	RoleStartSpeaking: MenuRole{Label: "Start Speaking"},
	RoleStopSpeaking:  MenuRole{Label: "Stop Speaking"},
	RoleFront:         MenuRole{Label: "Bring All to Front"},
	RoleZoom:          MenuRole{Label: "Zoom"},
	RoleWindow:        MenuRole{Label: "Window"},
	RoleHelp:          MenuRole{Label: "Help"},
	RoleServices:      MenuRole{Label: "Services"},
}

var menuRoleMap = map[string]MenuRole{
	"close":              MenuRole{Label: MenuLabelClose, Accelerator: "CommandOrControl+W", WindowMethod: "close"},
	"copy":               MenuRole{Label: "Copy", Accelerator: "CommandOrControl+C", WebContentsMethod: "copy"},
	"cut":                MenuRole{Label: "Cut", Accelerator: "CommandOrControl+X", WebContentsMethod: "cut"},
	"delete":             MenuRole{Label: "Delete", WebContentsMethod: "delete"},
	"minimize":           MenuRole{Label: "Minimize", Accelerator: "CommandOrControl+M", WindowMethod: "minimize"},
	"paste":              MenuRole{Label: "Paste", Accelerator: "CommandOrControl+V", WebContentsMethod: "paste"},
	"pasteandmatchstyle": MenuRole{Label: "Paste and Match Style", Accelerator: "Shift+CommandOrControl+V", WebContentsMethod: "pasteAndMatchStyle"},
	"quit":               MenuRole{Label: MenuLabelQuit, Accelerator: MenuAcceleratorQuit, AppMethod: "quit"},
	"redo":               MenuRole{Label: "Redo", Accelerator: MenuAcceleratorRedo, WebContentsMethod: "redo"},
	"resetzoom":          MenuRole{Label: "Actual Size", Accelerator: "CommandOrControl+0", WebContentsMethod: "_menuResetZoom"},
	"selectall":          MenuRole{Label: "Select All", Accelerator: "CommandOrControl+A", WebContentsMethod: "selectAll"},
	"togglefullscreen":   MenuRole{Label: "Toggle Full Screen", Accelerator: MenuAcceleratorToggleFullscreen, WindowMethod: "_menuToggleFullscreen"},
	"undo":               MenuRole{Label: "Undo", Accelerator: "CommandOrControl+Z", WebContentsMethod: "undo"},
	"zoomin":             MenuRole{Label: "Zoom In", Accelerator: "CommandOrControl+Plus", WebContentsMethod: "_menuZoomIn"},
	"zoomout":            MenuRole{Label: "Zoom Out", Accelerator: "CommandOrControl+-", WebContentsMethod: "_MenuZoomOut"},
}

type MenuItemClickHandler func(*MenuItemTemplate, *Window)

type MenuItemTemplate struct {
	Type        MenuType             `json:"type"`
	Role        string               `json:"role,omitempty"`
	Label       string               `json:"label,omitempty"`
	SubLabel    string               `json:"sublabel,omitempty"`
	Accelerator string               `json:"accelerator,omitempty"`
	ID          int                  `json:"id"`
	Disabled    bool                 `json:"disabled"`
	Invisible   bool                 `json:"invisible"`
	Checked     bool                 `json:"checked"`
	SubMenu     MenuTemplate         `json:"-"`
	Click       MenuItemClickHandler `json:"-"`
	//Icon        image.Image // TODO: handle native image?
	// hidden properties
	windowMethod      string
	webContentsMethod string
	appMethod         string
	eventName         string
	subMenuID         int64
}

type menuItemTemplateWrapper struct {
	MenuItemTemplate
	WindowMethod      string `json:"windowMethod,omitempty"`
	WebContentsMethod string `json:"webContentsMethod,omitempty"`
	AppMethod         string `json:"appMethod,omitempty"`
	ClickEventName    string `json:"clickEventName,omitempty"`
	SubMenuID         int64  `json:"subMenuId"`
}

func newMenuItemTemplateWrapper(mi *MenuItemTemplate) *menuItemTemplateWrapper {
	return &menuItemTemplateWrapper{
		MenuItemTemplate:  *mi,
		WindowMethod:      mi.windowMethod,
		WebContentsMethod: mi.webContentsMethod,
		AppMethod:         mi.appMethod,
		ClickEventName:    mi.eventName,
		SubMenuID:         mi.subMenuID,
	}
}

func (mi *MenuItemTemplate) fixMenuType() error {
	if len(mi.SubMenu) > 0 {
		mi.Type = MenuTypeSubmenu
	} else if mi.Type == MenuTypeSubmenu {
		return fmt.Errorf("MenuTemplate type is MenuTypeSubmenu, but not have SubMenu.")
	}
	return nil
}

func (mi *MenuItemTemplate) applyRole() error {
	if mi.Role == "" {
		return nil
	}
	var r MenuRole
	var ok bool
	r, ok = menuRolePlatform[mi.Role]
	if !ok {
		r, ok = menuRoleMap[mi.Role]
	}
	if !ok {
		return fmt.Errorf("unrecognized role %q", mi.Role)
	}
	if mi.Label == "" {
		mi.Label = r.Label
	}
	if mi.Accelerator == "" {
		mi.Accelerator = r.Accelerator
	}
	mi.webContentsMethod = r.WebContentsMethod
	mi.windowMethod = r.WindowMethod
	mi.appMethod = r.AppMethod
	return nil
}

type MenuTemplate []MenuItemTemplate

/*
func (mi *MenuItemTemplate) apply(f func(*MenuItemTemplate) error) error {
	if err := f(mi); err != nil {
		return err
	}
	if mi.SubMenu != nil {
		if err := mi.SubMenu.apply(f); err != nil {
			return err
		}
	}
	return nil
}

func (m *MenuTemplate) apply(f func(*MenuItemTemplate) error) error {
	for i := 0; i < len(*m); i++ {
		mi := &(*m)[i]
		fmt.Printf("mi[%d]: %s\n", i, mi.Label)
		if err := mi.apply(f); err != nil {
			return err
		}

	}
	return nil
}
*/

type menuItemClickItem struct {
	mi *MenuItemTemplate
}

func (p menuItemClickItem) Call(o ObjectRef, arg interface{}) (bool, error) {
	args, ok := arg.([]interface{})
	if !ok || len(args) != 1 {
		log.Panic("Invalid arg type: %%v", arg)
	}
	id, err := dproxy.New(args[0]).Int64()
	if err != nil {
		return false, err
	}
	win := getObject(id).(*Window)
	p.mi.Click(p.mi, win)
	return false, nil
}

type menuCallback map[int]func()

var menuCallbackMap = map[int]menuCallback{}

func (m *MenuTemplate) collectMenuID(idMap map[int]*MenuItemTemplate) error {
	// collect menu ID
	for i := 0; i < len(*m); i++ {
		mi := &(*m)[i]

		if mi.ID == 0 {
			continue
		}
		if _, ok := idMap[mi.ID]; ok {
			return fmt.Errorf("Menu ID conflict: %d", mi.ID)
		}
		idMap[mi.ID] = mi
	}
	return nil
}

// set menu id to empty commands
func (m *MenuTemplate) fillMenuID(idMap map[int]*MenuItemTemplate) {
	menuIndex := 0
	for i := 0; i < len(*m); i++ {
		mi := &(*m)[i]
		if mi.ID != 0 {
			continue
		}
		if mi.Type == MenuTypeSeparator /*|| mi.Type == MenuTypeSubmenu*/ {
			//サブメニューもIDは不要?
			continue
		}
		for {
			menuIndex++
			if _, ok := idMap[menuIndex]; !ok {
				mi.ID = menuIndex
				idMap[menuIndex] = mi
				break
			}
		}
	}
}

type Menu struct {
	object
}

func newMenu(id int64) *Menu {
	menu := &Menu{object: newObject(id, objMenu)}
	addObject(id, menu)
	return menu
}

func NewMenuWithTemplate(template MenuTemplate) (*Menu, error) {
	if !apiReady {
		return nil, errors.New("meson api is not ready yet")
	}
	cmd := makeCreateCommand(objMenu)

	resp, err := sendMessage(&cmd)
	if err != nil {
		return nil, err
	}
	id, err := dproxy.New(resp).Int64()
	if err != nil {
		return nil, err
	}

	menu := newMenu(id)
	if err := menu.LoadTemplate(template); err != nil {
		//TODO: destory object...
		return nil, err
	}

	return menu, nil
}

func (m *Menu) LoadTemplate(template MenuTemplate) error {
	idMap := make(map[int]*MenuItemTemplate)
	if err := template.collectMenuID(idMap); err != nil {
		return err
	}
	template.fillMenuID(idMap)

	ids := make([]int, 0, len(idMap))
	for idx := 0; idx < len(template); idx++ {
		mi := &template[idx]
		mi.fixMenuType()
		mi.applyRole()
		if mi.Click != nil {
			ids = append(ids, mi.ID)
		}
		if mi.Type == MenuTypeSubmenu {
			sm, err := NewMenuWithTemplate(mi.SubMenu)
			if err != nil {
				return err
			}
			mi.subMenuID = sm.id
		}
	}

	tempEvents, err := m.makeTemporaryEvents(len(ids))
	if err != nil {
		return nil
	}

	for idx := 0; idx < len(ids); idx++ {
		id := ids[idx]
		mi := idMap[id]
		eventID := tempEvents[idx].id
		eventName := tempEvents[idx].name
		mi.eventName = eventName
		m.addRegisterdCallback(eventID, menuItemClickItem{mi: mi})
	}

	items := make([]interface{}, len(template))
	for i, t := range template {
		items[i] = newMenuItemTemplateWrapper(&t)
	}
	cmd := makeCallCommand(m.objType, m.id, "loadTemplate", items...)
	_, err = sendMessage(&cmd)
	if err != nil {
		return err
	}
	return nil
}

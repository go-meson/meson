package menu

import (
	"errors"
	"fmt"
	"github.com/go-meson/meson/internal/binding"
	"github.com/go-meson/meson/internal/command"
	"github.com/go-meson/meson/internal/event"
	"github.com/go-meson/meson/internal/object"
	obj "github.com/go-meson/meson/object"
	"github.com/go-meson/meson/window"
	"github.com/koron/go-dproxy"
	"log"
)

type Role struct {
	Label             string
	Accelerator       string
	WindowMethod      string
	WebContentsMethod string
	AppMethod         string
}

// platform dependents
const (
	LabelAbout = "About {{AppName}}"
	LabelClose = "Close Window"
	LabelQuit  = "Quit {{AppName}}"

	AcceleratorQuit             = "CommandOrControl+Q"
	AcceleratorRedo             = "Shift+CommandOrControl+Z"
	AcceleratorToggleFullscreen = "Control+Command+F"
)

type RoleType string

const (
	//RoleAbout map to the orderFrontStandardAboutPanel action
	RoleAbout RoleType = "about"
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
var menuRolePlatform = map[RoleType]Role{
	RoleAbout:         Role{Label: LabelAbout},
	RoleHide:          Role{Label: "Hide {{AppName}}", Accelerator: "Command+H"},
	RoleHideOthers:    Role{Label: "Hide Others", Accelerator: "Command+Alt+H"},
	RoleUnHide:        Role{Label: "Show All"},
	RoleStartSpeaking: Role{Label: "Start Speaking"},
	RoleStopSpeaking:  Role{Label: "Stop Speaking"},
	RoleFront:         Role{Label: "Bring All to Front"},
	RoleZoom:          Role{Label: "Zoom"},
	RoleWindow:        Role{Label: "Window"},
	RoleHelp:          Role{Label: "Help"},
	RoleServices:      Role{Label: "Services"},
}

var menuRoleMap = map[RoleType]Role{
	"close":              Role{Label: LabelClose, Accelerator: "CommandOrControl+W", WindowMethod: "close"},
	"copy":               Role{Label: "Copy", Accelerator: "CommandOrControl+C", WebContentsMethod: "copy"},
	"cut":                Role{Label: "Cut", Accelerator: "CommandOrControl+X", WebContentsMethod: "cut"},
	"delete":             Role{Label: "Delete", WebContentsMethod: "delete"},
	"minimize":           Role{Label: "Minimize", Accelerator: "CommandOrControl+M", WindowMethod: "minimize"},
	"paste":              Role{Label: "Paste", Accelerator: "CommandOrControl+V", WebContentsMethod: "paste"},
	"pasteandmatchstyle": Role{Label: "Paste and Match Style", Accelerator: "Shift+CommandOrControl+V", WebContentsMethod: "pasteAndMatchStyle"},
	"quit":               Role{Label: LabelQuit, Accelerator: AcceleratorQuit, AppMethod: "quit"},
	"redo":               Role{Label: "Redo", Accelerator: AcceleratorRedo, WebContentsMethod: "redo"},
	"resetzoom":          Role{Label: "Actual Size", Accelerator: "CommandOrControl+0", WebContentsMethod: "_menuResetZoom"},
	"selectall":          Role{Label: "Select All", Accelerator: "CommandOrControl+A", WebContentsMethod: "selectAll"},
	"togglefullscreen":   Role{Label: "Toggle Full Screen", Accelerator: AcceleratorToggleFullscreen, WindowMethod: "_menuToggleFullscreen"},
	"undo":               Role{Label: "Undo", Accelerator: "CommandOrControl+Z", WebContentsMethod: "undo"},
	"zoomin":             Role{Label: "Zoom In", Accelerator: "CommandOrControl+Plus", WebContentsMethod: "_menuZoomIn"},
	"zoomout":            Role{Label: "Zoom Out", Accelerator: "CommandOrControl+-", WebContentsMethod: "_menuZoomOut"},
}

type ItemClickHandler func(*ItemTemplate, *window.Window)

type MenuType binding.MenuType

type ItemTemplate struct {
	Type        MenuType         `json:"type"`
	Role        RoleType         `json:"role,omitempty"`
	Label       string           `json:"label,omitempty"`
	SubLabel    string           `json:"sublabel,omitempty"`
	Accelerator string           `json:"accelerator,omitempty"`
	ID          int              `json:"id"`
	Disabled    bool             `json:"disabled"`
	Invisible   bool             `json:"invisible"`
	Checked     bool             `json:"checked"`
	SubMenu     Template         `json:"-"`
	Click       ItemClickHandler `json:"-"`
	//Icon        image.Image // TODO: handle native image?
	// hidden properties
	windowMethod      string
	webContentsMethod string
	appMethod         string
	eventName         string
	subMenuID         int64
}

type menuItemTemplateWrapper struct {
	ItemTemplate
	WindowMethod      string `json:"windowMethod,omitempty"`
	WebContentsMethod string `json:"webContentsMethod,omitempty"`
	AppMethod         string `json:"appMethod,omitempty"`
	ClickEventName    string `json:"clickEventName,omitempty"`
	SubMenuID         int64  `json:"subMenuId"`
}

const (
	MenuTypeNormal    MenuType = MenuType(binding.MenuTypeNormal)
	MenuTypeSeparator          = MenuType(binding.MenuTypeSeparator)
	MenuTypeSubmenu            = MenuType(binding.MenuTypeSubmenu)
	MenuTypeCheckBox           = MenuType(binding.MenuTypeCheckBox)
	MenuTypeRadio              = MenuType(binding.MenuTypeRadio)
)

func newItemTemplateWrapper(mi *ItemTemplate) *menuItemTemplateWrapper {
	return &menuItemTemplateWrapper{
		ItemTemplate:      *mi,
		WindowMethod:      mi.windowMethod,
		WebContentsMethod: mi.webContentsMethod,
		AppMethod:         mi.appMethod,
		ClickEventName:    mi.eventName,
		SubMenuID:         mi.subMenuID,
	}
}

func (mi *ItemTemplate) fixMenuType() error {
	if len(mi.SubMenu) > 0 {
		mi.Type = binding.MenuTypeSubmenu
	} else if mi.Type == MenuTypeSubmenu {
		return fmt.Errorf("Template type is MenuTypeSubmenu, but not have SubMenu.")
	}
	return nil
}

func (mi *ItemTemplate) applyRole() error {
	if mi.Role == "" {
		return nil
	}
	var r Role
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

type Template []ItemTemplate

type menuItemClickItem struct {
	mi *ItemTemplate
}

func (p menuItemClickItem) Call(o obj.ObjectRef, arg interface{}) (bool, error) {
	args, ok := arg.([]interface{})
	if !ok || len(args) != 1 {
		log.Panicf("Invalid arg type: %#v", arg)
	}
	id, err := dproxy.New(args[0]).Int64()
	if err != nil {
		return false, err
	}
	win := object.GetObject(id).(*window.Window)
	p.mi.Click(p.mi, win)
	return false, nil
}

type menuCallback map[int]func()

var menuCallbackMap = map[int]menuCallback{}

func (m *Template) collectMenuID(idMap map[int]*ItemTemplate) error {
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
func (m *Template) fillMenuID(idMap map[int]*ItemTemplate) {
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
	object.Object
}

func newMenu(id int64) *Menu {
	menu := &Menu{Object: object.NewObject(id, binding.ObjMenu)}
	object.AddObject(id, menu)
	return menu
}

func NewWithTemplate(template Template) (*Menu, error) {
	if !command.APIReady {
		return nil, errors.New("meson api is not ready yet")
	}
	cmd := command.MakeCreateCommand(binding.ObjMenu)

	resp, err := command.SendMessage(&cmd)
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

func (m *Menu) LoadTemplate(template Template) error {
	idMap := make(map[int]*ItemTemplate)
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
			sm, err := NewWithTemplate(mi.SubMenu)
			if err != nil {
				return err
			}
			mi.subMenuID = sm.Id
		}
	}

	tempEvents, err := event.MakeTemporaryEvents(&m.Object, len(ids))
	if err != nil {
		return nil
	}

	for idx := 0; idx < len(ids); idx++ {
		id := ids[idx]
		mi := idMap[id]
		eventID := tempEvents[idx].Id
		eventName := tempEvents[idx].Name
		mi.eventName = eventName
		m.AddRegisterdCallback(eventID, menuItemClickItem{mi: mi})
	}

	items := make([]interface{}, len(template))
	for i, t := range template {
		items[i] = newItemTemplateWrapper(&t)
	}
	cmd := command.MakeCallCommand(m.ObjType, m.Id, "loadTemplate", items...)
	_, err = command.SendMessage(&cmd)
	if err != nil {
		return err
	}
	return nil
}

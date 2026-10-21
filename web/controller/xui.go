package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// XUIController is the main controller for the X-UI panel, managing sub-controllers.
type XUIController struct {
	BaseController

	settingController     *SettingController
	xraySettingController *XraySettingController
}

// NewXUIController creates a new XUIController and initializes its routes.
func NewXUIController(g *gin.RouterGroup) *XUIController {
	a := &XUIController{}
	a.initRouter(g)
	return a
}

// initRouter sets up the panel's settings endpoints, behind the login check,
// and sends the old interface's pages to the same screens of the new one, so
// bookmarks keep working.
func (a *XUIController) initRouter(g *gin.RouterGroup) {
	for from, to := range map[string]string{"/panel/": "", "/panel/inbounds": "inbounds", "/panel/settings": "settings", "/panel/xray": "xray"} {
		g.GET(from, func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, c.GetString("base_path")+to)
		})
	}

	g = g.Group("/panel")
	g.Use(a.checkLogin)

	a.settingController = NewSettingController(g)
	a.xraySettingController = NewXraySettingController(g)
}

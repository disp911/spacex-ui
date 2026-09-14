package controller

import (
	"strconv"

	"github.com/disp911/spacex-ui/v2/database/model"
	"github.com/disp911/spacex-ui/v2/web/service"

	"github.com/gin-gonic/gin"
)

// telemtUserForm is the editable part of a Telegram proxy user.
type telemtUserForm struct {
	Username   string `form:"username"`
	Secret     string `form:"secret"`
	Enable     bool   `form:"enable"`
	ExpiryTime int64  `form:"expiryTime"`
	TotalBytes int64  `form:"totalBytes"`
	LimitIp    int    `form:"limitIp"`
	Comment    string `form:"comment"`
}

func (f telemtUserForm) model() *model.TelemtUser {
	return &model.TelemtUser{
		Username:   f.Username,
		Secret:     f.Secret,
		Enable:     f.Enable,
		ExpiryTime: f.ExpiryTime,
		TotalBytes: f.TotalBytes,
		LimitIp:    f.LimitIp,
		Comment:    f.Comment,
	}
}

// TelemtController serves the API of the Telegram proxy page.
type TelemtController struct {
	telemtService service.TelemtService
}

// NewTelemtController creates a TelemtController and registers its routes.
func NewTelemtController(g *gin.RouterGroup) *TelemtController {
	a := &TelemtController{}
	a.initRouter(g)
	return a
}

func (a *TelemtController) initRouter(g *gin.RouterGroup) {
	g.GET("/overview", a.overview)
	g.GET("/logs", a.logs)
	g.POST("/settings", a.saveSettings)
	g.POST("/restart", a.restart)
	g.POST("/users/add", a.addUser)
	g.POST("/users/update/:id", a.updateUser)
	g.POST("/users/enable/:id", a.setUserEnable)
	g.POST("/users/del/:id", a.deleteUser)
	g.POST("/users/rotateSecret/:id", a.rotateSecret)
	g.POST("/users/resetTraffic/:id", a.resetTraffic)
}

func (a *TelemtController) overview(c *gin.Context) {
	ov, err := a.telemtService.GetOverview()
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.telemt.toasts.loadError"), err)
		return
	}
	jsonObj(c, ov, nil)
}

func (a *TelemtController) logs(c *gin.Context) {
	jsonObj(c, a.telemtService.Logs(), nil)
}

func (a *TelemtController) saveSettings(c *gin.Context) {
	var ts service.TelemtSettings
	if err := c.ShouldBind(&ts); err != nil {
		jsonMsg(c, I18nWeb(c, "pages.telemt.toasts.failed"), err)
		return
	}
	err := a.telemtService.SaveSettings(ts)
	telemtResult(c, "pages.telemt.toasts.saved", err)
}

func (a *TelemtController) restart(c *gin.Context) {
	err := a.telemtService.Restart()
	telemtResult(c, "pages.telemt.toasts.restarted", err)
}

func (a *TelemtController) addUser(c *gin.Context) {
	var form telemtUserForm
	if err := c.ShouldBind(&form); err != nil {
		jsonMsg(c, I18nWeb(c, "pages.telemt.toasts.failed"), err)
		return
	}
	err := a.telemtService.AddUser(form.model())
	telemtResult(c, "pages.telemt.toasts.userAdded", err)
}

func (a *TelemtController) updateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.telemt.toasts.failed"), err)
		return
	}
	var form telemtUserForm
	if err := c.ShouldBind(&form); err != nil {
		jsonMsg(c, I18nWeb(c, "pages.telemt.toasts.failed"), err)
		return
	}
	err = a.telemtService.UpdateUser(id, form.model())
	telemtResult(c, "pages.telemt.toasts.userUpdated", err)
}

func (a *TelemtController) setUserEnable(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.telemt.toasts.failed"), err)
		return
	}
	enable, err := strconv.ParseBool(c.PostForm("enable"))
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.telemt.toasts.failed"), err)
		return
	}
	err = a.telemtService.SetUserEnable(id, enable)
	telemtResult(c, "pages.telemt.toasts.userUpdated", err)
}

func (a *TelemtController) deleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.telemt.toasts.failed"), err)
		return
	}
	err = a.telemtService.DeleteUser(id)
	telemtResult(c, "pages.telemt.toasts.userDeleted", err)
}

func (a *TelemtController) rotateSecret(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.telemt.toasts.failed"), err)
		return
	}
	err = a.telemtService.RotateSecret(id)
	telemtResult(c, "pages.telemt.toasts.secretRotated", err)
}

func (a *TelemtController) resetTraffic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.telemt.toasts.failed"), err)
		return
	}
	err = a.telemtService.ResetTraffic(id)
	telemtResult(c, "pages.telemt.toasts.trafficReset", err)
}

// telemtResult reports the outcome of an action: the success message, or a
// generic failure message with the error.
func telemtResult(c *gin.Context, successKey string, err error) {
	if err != nil {
		jsonMsg(c, I18nWeb(c, "pages.telemt.toasts.failed"), err)
		return
	}
	jsonMsg(c, I18nWeb(c, successKey), nil)
}

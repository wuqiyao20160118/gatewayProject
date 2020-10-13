package controller

import (
	"encoding/json"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"src/gatewayProject/dao"
	"src/gatewayProject/dto"
	"src/gatewayProject/golang_common/lib"
	"src/gatewayProject/middleware"
	"src/gatewayProject/public"
	"time"
)

type AdminLoginController struct {
}

/*
	Example:
	func DemoRegister(router *gin.RouterGroup) {
		demo := DemoController{}
		router.GET("/index", demo.Index)
		router.Any("/bind", demo.Bind)
		router.GET("/dao", demo.Dao)
		router.GET("/redis", demo.Redis)
	}

	// ListPage godoc
	// @Summary 测试数据绑定
	// @Description 测试数据绑定
	// @Tags 用户
	// @ID /demo/bind
	// @Accept  json
	// @Produce  json
	// @Param polygon body dto.DemoInput true "body"
	// @Success 200 {object} middleware.Response{data=dto.DemoInput} "success"
	// @Router /demo/bind [post]
	func (demo *DemoController) Bind(c *gin.Context) {
		params := &dto.DemoInput{}
		if err := params.BindingValidParams(c); err != nil {
			middleware.ResponseError(c, 2001, err)
			return
		}
		middleware.ResponseSuccess(c, params)
		return
	}
*/

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	group.GET("/logout", adminLogin.AdminLogOut)
}

// AdminLogin godoc
// @Summary Admin Login
// @Description Admin Login
// @Tags Admin interface
// @ID /admin_login/login
// @Accept  json
// @Param body body dto.AdminLoginInput true "body"
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminLoginOutput} "success"
// @Router /admin_login/login [post]
func (adminLogin *AdminLoginController) AdminLogin(ctx *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindingValidParams(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	// Step 1: params.UserName 取得管理员信息 adminInfo
	// Step 2: adminInfo.salt + params.Password sha256 => saltPassword
	// Step 3: saltPassWord == adminInfo.password
	// Realized in dao/admin.go

	/*
		Example for dao:

		func (demo *DemoController) Dao(c *gin.Context) {
			tx, err := lib.GetGormPool("default")
			if err != nil {
				middleware.ResponseError(c, 2000, err)
				return
			}
			area, err := (&dao.Area{}).Find(c, tx, c.DefaultQuery("id", "1"));
			if err != nil {
				middleware.ResponseError(c, 2001, err)
				return
			}
			middleware.ResponseSuccess(c, area)
			return
		}
	*/
	tx, err := lib.GetGormPool("default") // 使用配置文件中default的数据库连接池
	if err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	admin := &dao.Admin{}
	admin, err = admin.LoginAndCheck(ctx, tx, params)
	if err != nil {
		middleware.ResponseError(ctx, 2002, err)
		return
	}

	// set the user's session
	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo) // binary stream
	if err != nil {
		middleware.ResponseError(ctx, 2003, err)
		return
	}

	sess := sessions.Default(ctx)
	sess.Set(public.AdminSessionInfoKey, string(sessBts))
	sess.Save()

	out := &dto.AdminLoginOutput{Token: admin.UserName}
	middleware.ResponseSuccess(ctx, out)
}

// AdminLogin godoc
// @Summary Admin Logout
// @Description Admin Logout
// @Tags Admin interface
// @ID /admin_login/logout
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin_login/logout [get]
func (adminLogin *AdminLoginController) AdminLogOut(ctx *gin.Context) {
	// get the session and delete it
	sess := sessions.Default(ctx)
	sess.Delete(public.AdminSessionInfoKey)
	sess.Save()

	middleware.ResponseSuccess(ctx, "")
}

package controller

import (
	"encoding/json"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"src/gatewayProject/dao"
	"src/gatewayProject/dto"
	"src/gatewayProject/middleware"
	"src/gatewayProject/public"
)

type AdminController struct {

}

func AdminRegister(group *gin.RouterGroup) {
	admin := &AdminController{}
	group.GET("/admin_info", admin.AdminInfo)
	group.POST("/change_pwd", admin.ChangePwd)
}


// AdminInfo godoc
// @Summary Admin Information
// @Description Admin Information
// @Tags Admin interface
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (admin *AdminController) AdminInfo(ctx *gin.Context) {
	// Step 1: Read the json information in sessionKey and transform it into struct AdminInfoOutput
	// Step 2: Extract the data from AdminInfoOutput, construct and wrap it into the output

	// Step 1
	sess := sessions.Default(ctx)
	sessInfo := sess.Get(public.AdminSessionInfoKey)  // return an interface
	sessionInfoStr := sessInfo.(string)  // or use fmt.Sprint(sessInfo)
	adminSessionInfo := &dto.AdminSessionInfo{}
	// unmarshal the json data
	if err := json.Unmarshal([]byte(sessionInfoStr), adminSessionInfo); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	// Step 2
	out := &dto.AdminInfoOutput{
		ID: adminSessionInfo.ID,
		Name: adminSessionInfo.UserName,
		LoginTime: adminSessionInfo.LoginTime,
		Avatar: "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		Introduction: "Default administrator",
		Roles: []string{"admin"},
	}

	middleware.ResponseSuccess(ctx, out)
}

// ChangePwd godoc
// @Summary Admin Change Password
// @Description Admin Change Password
// @Tags Admin interface
// @ID /admin/change_pwd
// @Accept  json
// @Param body body dto.ChangePwdInput true "body"
// @Produce  json
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (admin *AdminController) ChangePwd(ctx *gin.Context) {
	params := &dto.ChangePwdInput{}
	if err := params.BindingValidParams(ctx); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	// Step 1: Read the user info => sessInfo
	// Step 2: Use sessInfo.ID to search information in the database => adminInfo
	// Step 3: generate the new password and save it into the database

	// Step 1
	sess := sessions.Default(ctx)
	sessInfo := sess.Get(public.AdminSessionInfoKey)
	sessionInfoStr := sessInfo.(string)  // or use fmt.Sprint(sessInfo)
	adminSessionInfo := &dto.AdminSessionInfo{}
	// unmarshal the json data
	if err := json.Unmarshal([]byte(sessionInfoStr), adminSessionInfo); err != nil {
		middleware.ResponseError(ctx, 2000, err)
		return
	}

	// Step 2
	tx, err := lib.GetGormPool("default")  // 使用配置文件中default的数据库连接池
	if err != nil {
		middleware.ResponseError(ctx, 2001, err)
		return
	}
	adminInfo := &dao.Admin{}
	adminInfo, err = adminInfo.Find(ctx, tx, &(dao.Admin{UserName: adminSessionInfo.UserName}))
	if err != nil {
		middleware.ResponseError(ctx, 2002, err)
		return
	}

	// Step 3
	saltPassword := public.GenSaltPassword(adminInfo.Salt, params.Password)
	adminInfo.Password = saltPassword
	// Gin 自动更新 CreatedAt与UpdatedAt
	if err := adminInfo.Save(ctx, tx); err != nil {
		middleware.ResponseError(ctx, 2003, err)
		return
	}

	middleware.ResponseSuccess(ctx, "")
}
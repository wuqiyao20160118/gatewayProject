package dto

import (
	"github.com/gin-gonic/gin"
	"src/gatewayProject/public"
	"time"
)

type AdminSessionInfo struct {
	ID        int         `json:"id"`  // id
	UserName  string      `json:"username"`  // username
	LoginTime time.Time   `json:"login_time"`  // login_time
}

type AdminLoginInput struct {
	UserName  string `json:"username" form:"username" comment:"username" example:"admin" validate:"required,valid_username"`  // username
	Password  string `json:"password" form:"password" comment:"password" example:"123456" validate:"required"`  // password
}

// validate the input parameter with tags defined in struct(AdminLoginInput)
func (param *AdminLoginInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type AdminLoginOutput struct {
	Token string `json:"token" form:"token" comment:"token" example:"token" validate:""`  // token
}
package dto

import (
	"github.com/gin-gonic/gin"
	"src/gatewayProject/public"
	"time"
)

type AdminInfoOutput struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	LoginTime     time.Time   `json:"login_time"`
	Avatar        string      `json:"avatar"`
	Introduction  string      `json:"introduction"`
	Roles         []string    `json:"roles"`
}

type ChangePwdInput struct {
	Password  string `json:"password" form:"password" comment:"password" example:"123456" validate:"required"`  // password
}

func (param *ChangePwdInput) BindingValidParams(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

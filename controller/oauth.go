package controller

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"src/gatewayProject/dao"
	"src/gatewayProject/dto"
	"src/gatewayProject/golang_common/lib"
	"src/gatewayProject/middleware"
	"src/gatewayProject/public"
	"strings"
	"time"
)

type OAuthController struct {
}

func OAuthRegister(group *gin.RouterGroup) {
	oauth := &OAuthController{}

	group.POST("/tokens", oauth.Tokens)
}

// Tokens godoc
// @Summary Get TOKEN
// @Description Get TOKEN
// @Tags OAUTH
// @ID /oauth/tokens
// @Accept  json
// @Produce  json
// @Param body body dto.TokensInput true "body"
// @Success 200 {object} middleware.Response{data=dto.TokensOutput} "success"
// @Router /oauth/tokens [post]
func (oauth *OAuthController) Tokens(c *gin.Context) {
	params := &dto.TokensInput{}
	if err := params.BindingValidParams(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	splits := strings.Split(c.GetHeader("Authorization"), " ")
	if len(splits) != 2 {
		middleware.ResponseError(c, 2001, errors.New("wrong format of username or password"))
		return
	}

	appSecret, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	//  取出 app_id secret
	//  生成 app_list
	//  匹配 app_id
	//  基于 jwt生成token
	//  生成 output
	parts := strings.Split(string(appSecret), ":")
	if len(parts) != 2 {
		middleware.ResponseError(c, 2003, errors.New("wrong format of username or password"))
		return
	}

	appList := dao.AppManagerHandler.GetAppList()
	for _, appInfo := range appList {
		if appInfo.AppID == parts[0] && appInfo.Secret == parts[1] {
			claims := jwt.StandardClaims{
				Issuer:    appInfo.AppID,
				ExpiresAt: time.Now().Add(public.JwtExpires * time.Second).In(lib.TimeLocation).Unix(),
			}
			token, err := public.JwtEncode(claims)
			if err != nil {
				middleware.ResponseError(c, 2004, err)
				return
			}

			out := &dto.TokensOutput{
				AccessToken: token,
				ExpiresIn:   public.JwtExpires,
				TokenType:   "Bearer",
				Scope:       "read_write",
			}

			middleware.ResponseSuccess(c, out)
			return
		}
	}

	middleware.ResponseError(c, 2005, errors.New("no matching App ID"))
}

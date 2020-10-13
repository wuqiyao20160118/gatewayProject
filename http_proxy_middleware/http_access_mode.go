package http_proxy_middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"src/gatewayProject/dao"
	"src/gatewayProject/middleware"
	"src/gatewayProject/public"
)

// 匹配接入方式 基于请求信息 (domain / URL)
func HTTPAccessModeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, err := dao.ServiceManagerHandler.HTTPAccessMode(c)
		if err != nil {
			middleware.ResponseError(c, 1001, err)
			c.Abort()
			return
		}

		fmt.Println("matched service", public.Obj2Json(service))

		c.Set("service", service)
		c.Next()
	}
}

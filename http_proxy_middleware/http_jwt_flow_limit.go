package http_proxy_middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"src/gatewayProject/dao"
	"src/gatewayProject/middleware"
	"src/gatewayProject/public"
)

func HTTPJwtFlowLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("app not found"))
			c.Abort()
			return
		}
		appInfo := appInterface.(*dao.App)

		if appInfo.Qps > 0 {
			appLimiter, err := public.FlowLimiterHandler.GetLimiter(
				public.FlowAppPrefix+appInfo.AppID+"_"+c.ClientIP(),
				float64(appInfo.Qps))
			if err != nil {
				middleware.ResponseError(c, 5001, err)
				c.Abort()
				return
			}
			if !appLimiter.Allow() {
				middleware.ResponseError(c, 5002,
					errors.New(fmt.Sprintf("%v flow limit %v", c.ClientIP(), appInfo.Qps)))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

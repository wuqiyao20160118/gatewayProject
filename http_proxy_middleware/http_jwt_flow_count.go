package http_proxy_middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"src/gatewayProject/dao"
	"src/gatewayProject/middleware"
	"src/gatewayProject/public"
)

func HTTPJwtFlowCountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appInterface, ok := c.Get("app")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("app not found"))
			c.Abort()
			return
		}
		appInfo := appInterface.(*dao.App)

		// 统计项： 租户
		appCounter, err := public.FlowCounterHandler.GetFlowCounter(public.FlowAppPrefix + appInfo.Name)
		if err != nil {
			middleware.ResponseError(c, 4001, err)
			c.Abort()
			return
		}
		appCounter.Increase()

		if appInfo.Qpd > 0 && appCounter.TotalCount > appInfo.Qpd {
			middleware.ResponseError(c, 2003, errors.New(fmt.Sprintf("租户日请求量限流 limit:%v current:%v", appInfo.Qpd, appCounter.TotalCount)))
			c.Abort()
			return
		}

		c.Next()
	}
}

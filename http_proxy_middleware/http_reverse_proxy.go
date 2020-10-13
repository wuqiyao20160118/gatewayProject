package http_proxy_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"src/gatewayProject/dao"
	"src/gatewayProject/middleware"
	"src/gatewayProject/reverse_proxy"
)

func HTTPReverseProxyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		// 接口转型（类型断言，运行期确定）
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// 获取负载均衡器
		lb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			c.Abort()
			return
		}

		// 获取连接池
		trans, err := dao.TransportorHandler.GetTrans(serviceDetail)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			c.Abort()
			return
		}

		// 创建reverse proxy
		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, lb, trans)
		// 使用 reverse proxy中的ServeHttp(c.Writer, c.Request)方法
		proxy.ServeHTTP(c.Writer, c.Request)
		c.Abort()
		return
	}
}

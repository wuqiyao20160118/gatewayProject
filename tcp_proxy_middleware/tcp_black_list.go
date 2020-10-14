package tcp_proxy_middleware

import (
	"fmt"
	"src/gatewayProject/dao"
	"src/gatewayProject/public"
	"strings"
)

func TCPBlackListMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		serverInterface := c.Get("service")
		if serverInterface == nil {
			c.conn.Write([]byte("get service empty"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		splits := strings.Split(c.conn.RemoteAddr().String(), ":")
		clientIP := ""
		if len(splits) == 2 {
			clientIP = splits[0]
		}

		var whiteList []string
		if serviceDetail.AccessControl.WhiteList != "" {
			// strings.Split() 在src字符串为空时会返回[""]，其长度为1，故需要首先判断字符串是否为空
			whiteList = strings.Split(serviceDetail.AccessControl.WhiteList, ",")
		}

		var blackList []string
		if serviceDetail.AccessControl.BlackList != "" {
			// strings.Split() 在src字符串为空时会返回[""]，其长度为1，故需要首先判断字符串是否为空
			blackList = strings.Split(serviceDetail.AccessControl.BlackList, ",")
		}

		if serviceDetail.AccessControl.OpenAuth == 1 && len(whiteList) == 0 && len(blackList) > 0 {
			if public.InStringSlice(blackList, clientIP) {
				c.conn.Write([]byte(fmt.Sprintf("%s in black ip list", clientIP)))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

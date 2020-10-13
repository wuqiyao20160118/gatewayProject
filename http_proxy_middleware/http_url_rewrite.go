package http_proxy_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"regexp"
	"src/gatewayProject/dao"
	"src/gatewayProject/middleware"
	"strings"
)

func HTTPUrlRewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		serverInterface, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("service not found"))
			c.Abort()
			return
		}
		serviceDetail := serverInterface.(*dao.ServiceDetail)

		// ^/gatekeeper/test_service(.*) $1...
		//	Example
		//	Code:
		//		re := regexp.MustCompile(`a(x*)b`)
		//		fmt.Printf("%s\n", re.ReplaceAll([]byte("-ab-axxb-"), []byte("T")))
		//		fmt.Printf("%s\n", re.ReplaceAll([]byte("-ab-axxb-"), []byte("$1")))
		//		fmt.Printf("%s\n", re.ReplaceAll([]byte("-ab-axxb-"), []byte("$1W")))
		//		fmt.Printf("%s\n", re.ReplaceAll([]byte("-ab-axxb-"), []byte("${1}W")))
		//	Output:
		//		-T-T-
		//		--xx-
		//		---
		//		-W-xxW-
		for _, item := range strings.Split(serviceDetail.HTTPRule.UrlRewrite, ",") {
			items := strings.Split(item, " ")
			if len(items) != 2 {
				continue
			}
			re, err := regexp.Compile(items[0])
			if err != nil {
				continue
			}

			// ReplaceAll returns a copy of src, replacing matches of the Regexp with the replacement text repl.
			// Inside repl, $ signs are interpreted as in Expand, so for instance $1 represents the text of the first submatch.
			// regExp.ReplaceAll() returns byte slice
			replacePath := re.ReplaceAll([]byte(c.Request.URL.Path), []byte(items[1]))
			c.Request.URL.Path = string(replacePath)
		}

		c.Next()
	}
}

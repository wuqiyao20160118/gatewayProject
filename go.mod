module src/gatewayProject

go 1.14

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/boj/redistore v0.0.0-20180917114910-cd5dcc76aeff // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/e421083458/go_gateway v0.0.0-20200620084504-d602eb8bc883
	github.com/e421083458/golang_common v1.0.3
	github.com/e421083458/gorm v1.0.1
	github.com/e421083458/grpc-proxy v0.2.0
	github.com/garyburd/redigo v1.6.0
	github.com/gin-gonic/contrib v0.0.0-20191209060500-d6e26eeaa607
	github.com/gin-gonic/gin v1.4.0
	github.com/go-playground/locales v0.12.1
	github.com/go-playground/universal-translator v0.16.0
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/sessions v1.1.3 // indirect
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/viper v1.7.0
	github.com/swaggo/files v0.0.0-20190704085106-630677cd5c14
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.5
	golang.org/x/net v0.0.0-20201010224723-4f7140c49acb
	golang.org/x/sys v0.0.0-20201015000850-e3ed0017c211 // indirect
	golang.org/x/time v0.0.0-20200416051211-89c76fbcd5d1
	google.golang.org/genproto v0.0.0-20201014134559-03b6142f0dc9 // indirect
	google.golang.org/grpc v1.33.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/go-playground/validator.v9 v9.29.0
)

replace github.com/gin-contrib/sse v0.1.0 => github.com/e421083458/sse v0.1.1

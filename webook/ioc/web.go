package ioc

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"geektime/webook/internal/web"
	"geektime/webook/internal/web/middleware"
	"geektime/webook/pkg/ginx/middlewares/ratelimit"
)

func InitGin(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			// AllowOrigins:     []string{"https://localhost:3000"},
			// AllowMethods:     []string{"POST", "GET"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			ExposeHeaders:    []string{"x-jwt-token"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					// 开发环境
					return true
				}
				return origin == "your_company.com"
			},
			MaxAge: 12 * time.Hour,
		}),
		middleware.NewLoginJWTMiddlewareBuilder().IgnorePaths("/users/signup",
			"/users/login", "/users/login_sms/code/send", "/users/login_sms", "/hello").Build(),
		ratelimit.NewBuilder(redisClient, time.Second, 10).Build(),
	}
}

package main

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"geektime/webook/internal/repository"
	"geektime/webook/internal/repository/dao"
	"geektime/webook/internal/service"
	"geektime/webook/internal/web"
	"geektime/webook/internal/web/middleware"
)

func main() {
	server := initWebServer()
	db := initDB()
	u := initUser(db)
	u.RegisterRoutes(server)
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	// 解决跨域问题（CORS）
	server.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"https://localhost:3000"},
		//AllowMethods:     []string{"POST", "GET"},
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
	}))

	// 设置session
	// 步骤1
	//store := cookie.NewStore([]byte("secret"))
	// 单实例部署 memstore 基于内存的实现
	//store := memstore.NewStore([]byte("mQ5>dY9%bZ4,uI6,oF4~aU4(nU0&sK5."),
	//	[]byte("aY3?fW6+kK9~mX7!yQ5|wS7%vR8_lO1`"))
	// 多实例部署 redis 基于redis的实现
	store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("mQ5>dY9%bZ4,uI6,oF4~aU4(nU0&sK5."),
		[]byte("aY3?fW6+kK9~mX7!yQ5|wS7%vR8_lO1`"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("mysession", store))
	// 步骤3
	server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/signup", "/users/login").Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDao(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}
	if err = dao.InitTable(db); err != nil {
		panic(err)
	}
	return db
}

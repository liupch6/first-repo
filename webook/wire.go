//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"geektime/webook/internal/repository"
	"geektime/webook/internal/repository/cache"
	"geektime/webook/internal/repository/dao"
	"geektime/webook/internal/service"
	"geektime/webook/internal/web"
	"geektime/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(ioc.InitDB, ioc.InitRedis,
		dao.NewUserDAO, cache.NewUserCache, cache.NewCodeCache,
		repository.NewUserRepository, repository.NewCodeRepository,
		service.NewUserService, service.NewCodeService, ioc.InitSMSService,
		web.NewUserHandler, ioc.InitWebServer, ioc.InitMiddlewares)
	return new(gin.Engine)
}

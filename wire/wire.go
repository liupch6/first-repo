//go:build wireinject

package wire

import (
	"github.com/google/wire"

	"geektime/wire/repository"
	"geektime/wire/repository/dao"
)

func InitRepository() *repository.UserRepository {
	// 这个方法传入各个组件的初始化方法
	wire.Build(InitDB, dao.NewUserDAO, repository.NewUserRepository)
	return new(repository.UserRepository)
}

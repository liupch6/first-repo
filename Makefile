.PHONY: mock
mock:
	@mockgen -source=webook/internal/service/user.go -package=svcmocks -destination=webook/internal/service/mocks/user.mock.go
	@mockgen -source=webook/internal/service/code.go -package=svcmocks -destination=webook/internal/service/mocks/code.mock.go
	@mockgen -source=webook/internal/repository/user.go -destination=webook/internal/repository/mocks/user.mock.go -package=repomocks
	@mockgen -source=webook/internal/repository/code.go -destination=webook/internal/repository/mocks/code.mock.go -package=repomocks
	@mockgen -source=webook/internal/repository/dao/user.go -destination=webook/internal/repository/dao/mocks/user.mock.go -package=daomocks
	@mockgen -source=webook/internal/repository/cache/user.go -destination=webook/internal/repository/cache/mocks/user.mock.go -package=cachemocks
	@mockgen -destination=webook/internal/repository/cache/redismocks/cmdable.mock.go -package=redismocks github.com/redis/go-redis/v9 Cmdable
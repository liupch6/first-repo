package repository

import (
	"context"

	"geektime/webook/internal/domain"
	"geektime/webook/internal/repository/cache"
	"geektime/webook/internal/repository/dao"
)

var (
	ErrUserEmailDuplicate = dao.ErrUserEmailDuplicate
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	ud    *dao.UserDao
	cache *cache.UserCache
}

func NewUserRepository(ud *dao.UserDao, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		ud:    ud,
		cache: c,
	}
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.ud.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.ud.FindByEmail(ctx, email)
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, err
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	// 缓存命中
	if err == nil {
		return u, err
	}
	user, err := r.ud.FindById(ctx, id)
	// 数据库出错
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}
	// 异步设置缓存
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// log
		}
	}()
	return u, err
}

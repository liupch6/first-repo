package repository

import (
	"context"
	"database/sql"
	"time"

	"geektime/webook/internal/domain"
	"geektime/webook/internal/repository/cache"
	"geektime/webook/internal/repository/dao"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
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
	return r.ud.Insert(ctx, r.domainToEntity(u))
}

func (r *UserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.ud.FindByPhone(ctx, phone)
	return r.entityToDomain(u), err
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.ud.FindByEmail(ctx, email)
	return r.entityToDomain(u), err
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
	u = r.entityToDomain(user)
	// 异步设置缓存
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// log
		}
	}()
	return u, err
}

func (r *UserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Ctime: u.Ctime.UnixMilli(),
	}
}

func (r *UserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}

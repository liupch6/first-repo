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

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	Update(ctx context.Context, u domain.User) error
}

type CacheUserRepository struct {
	ud    dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(ud dao.UserDAO, c cache.UserCache) UserRepository {
	return &CacheUserRepository{
		ud:    ud,
		cache: c,
	}
}

func (r *CacheUserRepository) Update(ctx context.Context, u domain.User) error {
	err := r.ud.UpdateNonZeroFields(ctx, r.domainToEntity(u))
	if err != nil {
		return err
	}
	return r.cache.Delete(ctx, u.Id)
}

func (r *CacheUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.ud.Insert(ctx, r.domainToEntity(u))
}

func (r *CacheUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.ud.FindByPhone(ctx, phone)
	return r.entityToDomain(u), err
}

func (r *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.ud.FindByEmail(ctx, email)
	return r.entityToDomain(u), err
}

func (r *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	// 缓存命中
	if err == nil {
		return u, nil
	}
	user, err := r.ud.FindById(ctx, id)
	// 数据库出错
	if err != nil {
		return domain.User{}, err
	}
	u = r.entityToDomain(user)
	// 异步设置缓存
	// go func() {
	// 	err = r.cache.Set(ctx, u)
	// 	if err != nil {
	// 		// log
	// 	}
	// }()
	go func() {
		_ = r.cache.Set(ctx, u)
	}()

	return u, nil
}

func (r *CacheUserRepository) domainToEntity(u domain.User) dao.User {
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

func (r *CacheUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}

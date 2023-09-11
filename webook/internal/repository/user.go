package repository

import (
	"context"

	"geektime/webook/internal/domain"
	"geektime/webook/internal/repository/dao"
)

var (
	ErrUserEmailDuplicate = dao.ErrUserEmailDuplicate
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	ud *dao.UserDao
}

func NewUserRepository(ud *dao.UserDao) *UserRepository {
	return &UserRepository{
		ud: ud,
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

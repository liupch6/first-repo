package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"geektime/webook/internal/domain"
	"geektime/webook/internal/repository"
)

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("账户或密码错误")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return domain.User{}, ErrInvalidUserOrPassword
		}
		return domain.User{}, err
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {

	// 快路径
	u, err := svc.repo.FindByPhone(ctx, phone)
	// 判断有没有这个用户
	if !errors.Is(err, repository.ErrUserNotFound) {
		return u, err
	}

	// 慢路径
	// 没有这个用户 => 创建用户
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil && !errors.Is(err, ErrUserDuplicate) {
		return domain.User{}, err
	}
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func PathsDownGrage(ctx context.Context, quick, slow func()) {
	quick()
	if ctx.Value("降级") == "true" {
		return
	}
	slow()
}

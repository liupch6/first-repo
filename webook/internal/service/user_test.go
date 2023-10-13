package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"geektime/webook/internal/domain"
	"geektime/webook/internal/repository"
	repomocks "geektime/webook/internal/repository/mocks"
)

func TestUserServiceImpl_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 输入
		// ctx      context.Context
		email    string
		password string

		// 输出
		expectedUser domain.User
		expectedErr  error
	}{
		{
			name: "登录成功", // 邮箱和密码正确
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$10$aQc9gokDobCC5ci4QlHVVOuDKZu7vFsak9w3y/7kwYiLvRbO7w90e",
					Phone:    "15212345678",
					Ctime:    now,
				}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123",
			expectedUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$aQc9gokDobCC5ci4QlHVVOuDKZu7vFsak9w3y/7kwYiLvRbO7w90e",
				Phone:    "15212345678",
				Ctime:    now,
			},
			expectedErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:        "123@qq.com",
			password:     "hello#world123",
			expectedUser: domain.User{},
			expectedErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB 错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), gomock.Any()).Return(domain.User{}, errors.New("mock db 错误"))
				return repo
			},
			email:        "123@qq.com",
			password:     "hello#world123",
			expectedUser: domain.User{},
			expectedErr:  errors.New("mock db 错误"),
		},
		{
			name: "密码不对",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email:    "123@qq.com",
					Password: "$2a$10$aQc9gokDobCC5ci4QlHVVOuDKZu7vFsak9w3y/7kwYiLvRbO7w90e",
					Phone:    "15212345678",
					Ctime:    now,
				}, nil)
				return repo
			},
			email:        "123@qq.com",
			password:     "hello#world12345",
			expectedUser: domain.User{},
			expectedErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedUser, u)
		})
	}
}

func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("hello#world123"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}

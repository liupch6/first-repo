package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"geektime/webook/internal/domain"
	"geektime/webook/internal/repository/cache"
	cachemocks "geektime/webook/internal/repository/cache/mocks"
	"geektime/webook/internal/repository/dao"
	daomocks "geektime/webook/internal/repository/dao/mocks"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	now := time.Now() // 这个 now 带纳秒，而 now.UnixMilli() 已经把纳秒去掉了，只有毫秒
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		ctx context.Context
		id  int64

		expectedUser domain.User
		expectedErr  error
	}{
		{
			name: "缓存未命中，数据库查询成功，并设置缓存",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				ud := daomocks.NewMockUserDAO(ctrl)
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, cache.ErrKeyNotExist)
				ud.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.User{
					Id: 1,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password: "$2a$10$aQc9gokDobCC5ci4QlHVVOuDKZu7vFsak9w3y/7kwYiLvRbO7w90e",
					Phone: sql.NullString{
						String: "15212345678",
						Valid:  true,
					},
					Ctime: now.UnixMilli(),
					Utime: now.UnixMilli(),
				}, nil)
				uc.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil)
				return ud, uc
			},
			ctx: context.Background(),
			id:  1,
			expectedUser: domain.User{
				Id:       1,
				Email:    "123@qq.com",
				Password: "$2a$10$aQc9gokDobCC5ci4QlHVVOuDKZu7vFsak9w3y/7kwYiLvRbO7w90e",
				Phone:    "15212345678",
				Ctime:    now,
			},
			expectedErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				ud := daomocks.NewMockUserDAO(ctrl)
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{
					Id:       1,
					Email:    "123@qq.com",
					Password: "$2a$10$aQc9gokDobCC5ci4QlHVVOuDKZu7vFsak9w3y/7kwYiLvRbO7w90e",
					Phone:    "15212345678",
					Ctime:    now,
				}, nil)
				return ud, uc
			},
			ctx: context.Background(),
			id:  1,
			expectedUser: domain.User{
				Id:       1,
				Email:    "123@qq.com",
				Password: "$2a$10$aQc9gokDobCC5ci4QlHVVOuDKZu7vFsak9w3y/7kwYiLvRbO7w90e",
				Phone:    "15212345678",
				Ctime:    now,
			},
			expectedErr: nil,
		},
		{
			name: "缓存未命中，数据库查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				ud := daomocks.NewMockUserDAO(ctrl)
				uc := cachemocks.NewMockUserCache(ctrl)
				uc.EXPECT().Get(gomock.Any(), int64(1)).Return(domain.User{}, cache.ErrKeyNotExist)
				ud.EXPECT().FindById(gomock.Any(), int64(1)).Return(dao.User{}, errors.New("数据库查询失败"))
				return ud, uc
			},
			ctx:          context.Background(),
			id:           1,
			expectedUser: domain.User{},
			expectedErr:  errors.New("数据库查询失败"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewUserRepository(tc.mock(ctrl))
			u, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedUser, u)
			time.Sleep(time.Second) // 异步设置缓存
		})
	}
}

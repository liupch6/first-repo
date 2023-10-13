package cache

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"geektime/webook/internal/repository/cache/redismocks"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) redis.Cmdable

		ctx   context.Context
		biz   string
		phone string
		code  string

		expectedErr error
	}{
		{
			name: "验证码设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(nil)
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).
					Return(res)
				return cmd
			},
			ctx:         context.Background(),
			biz:         "login",
			phone:       "152",
			code:        "123456",
			expectedErr: nil,
		},
		{
			name: "验证码发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(nil)
				res.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).
					Return(res)
				return cmd
			},
			ctx:         context.Background(),
			biz:         "login",
			phone:       "152",
			code:        "123456",
			expectedErr: ErrCodeSendTooMany,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(nil)
				res.SetVal(int64(-2))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).
					Return(res)
				return cmd
			},
			ctx:         context.Background(),
			biz:         "login",
			phone:       "152",
			code:        "123456",
			expectedErr: errors.New("系统错误"),
		},
		{
			name: "验证码设置失败",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetErr(errors.New("mock error"))
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).
					Return(res)
				return cmd
			},
			ctx:         context.Background(),
			biz:         "login",
			phone:       "152",
			code:        "123456",
			expectedErr: errors.New("mock error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewCodeCache(tc.mock(ctrl))
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

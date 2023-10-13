package web

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"geektime/webook/internal/domain"
	"geektime/webook/internal/service"
	svcmocks "geektime/webook/internal/service/mocks"
)

func TestEncrypt(t *testing.T) {
	password := "hello#world123"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	assert.NoError(t, err)
}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) service.UserService

		reqBody string

		expectedCode int
		expectedBody string
	}{
		{
			name: "OJBK",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(nil) // 注册成功肯定是 nil
				return userSvc
			},
			reqBody:      `{"email": "123@qq.com","password": "hello#world123","confirmPassword": "hello#world123"}`,
			expectedCode: http.StatusOK,
			expectedBody: "注册成功",
		},
		// TODO
		{
			name: "参数不对， bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody:      `{"email": "123@qq.com","password": "hello#world123",}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
		},
		{
			name: "邮箱错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody:      `{"email": "123qq.com","password": "hello#world123","confirmPassword": "hello#world123"}`,
			expectedCode: http.StatusOK,
			expectedBody: "邮箱格式错误",
		},
		{
			name: "密码格式错误",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody:      `{"email": "123@qq.com","password": "hello#world","confirmPassword": "hello#world123"}`,
			expectedCode: http.StatusOK,
			expectedBody: "密码至少8个字符，至少1个字母，1个数字和1个特殊字符",
		},
		{
			name: "密码不一致",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc
			},
			reqBody:      `{"email": "123@qq.com","password": "hello#world1234","confirmPassword": "hello#world123"}`,
			expectedCode: http.StatusOK,
			expectedBody: "两次输入密码不一致",
		},
		{
			name: "注册邮箱重复了",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(service.ErrUserDuplicate)
				return userSvc
			},
			reqBody:      `{"email": "123@qq.com","password": "hello#world123","confirmPassword": "hello#world123"}`,
			expectedCode: http.StatusOK,
			expectedBody: "邮箱重复",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("随便"))
				return userSvc
			},
			reqBody:      `{"email": "123@qq.com","password": "hello#world123","confirmPassword": "hello#world123"}`,
			expectedCode: http.StatusOK,
			expectedBody: "系统异常",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// 准备一个 gin.Engine，并注册路由
			server := gin.Default()
			h := NewUserHandler(tc.mock(ctrl), nil)
			h.RegisterRoutes(server)

			// 准备请求
			req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// 准备记录响应的 Recorder
			resp := httptest.NewRecorder()

			// 直接发起调用，也就是假装收到了 HTTP 请求
			server.ServeHTTP(resp, req)

			// 比较 Recorder 里面记录的响应
			assert.Equal(t, tc.expectedCode, resp.Code)
			assert.Equal(t, tc.expectedBody, resp.Body.String())
		})
	}
}

func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	usersvc := svcmocks.NewMockUserService(ctrl)
	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock error"))

	err := usersvc.SignUp(context.Background(), domain.User{
		Email: "123@qq.com",
	})
	t.Log(err)
}

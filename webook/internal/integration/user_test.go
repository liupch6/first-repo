package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"geektime/webook/internal/web"
	"geektime/webook/ioc"
)

func TestUserHandler_SendLoginSMSCode(t *testing.T) {
	server := InitWebServer()
	rdb := ioc.InitRedis()
	testCases := []struct {
		name string

		// 提前准备数据
		before func(t *testing.T)
		// 验证并且清理数据
		after   func(t *testing.T)
		reqBody string

		expectedCode int
		expectedBody web.Result
	}{
		{
			name: "发送成功",
			before: func(t *testing.T) {
				// 不需要做什么，也就是 redis 里面什么数据都没有
			},
			after: func(t *testing.T) {
				// 清理数据
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				cancel()
				assert.NoError(t, err)
				// 验证码为 6 位
				assert.True(t, len(val) == 6)
			},
			reqBody:      `{"phone": "15212345678"}`,
			expectedCode: http.StatusOK,
			expectedBody: web.Result{
				Code: 0,
				Msg:  "发送成功",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				// 这个手机号已经有一个验证码了
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:15212345678", "123456", time.Minute*9+time.Second*30).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				cancel()
				assert.NoError(t, err)
				// 验证码没有被覆盖，还是 123456
				assert.Equal(t, "123456", val)
			},
			reqBody:      `{"phone": "15212345678"}`,
			expectedCode: http.StatusOK,
			expectedBody: web.Result{
				Code: 0,
				Msg:  "发送太频繁，请稍后再试",
			},
		},
		{
			name: "系统错误",
			before: func(t *testing.T) {
				// 这个手机号已经有一个验证码了，但是没有过期时间
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				_, err := rdb.Set(ctx, "phone_code:login:15212345678", "123456", 0).Result()
				cancel()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
				val, err := rdb.GetDel(ctx, "phone_code:login:15212345678").Result()
				cancel()
				assert.NoError(t, err)
				// 验证码为 6 位
				assert.True(t, len(val) == 6)
			},
			reqBody:      `{"phone": "15212345678"}`,
			expectedCode: http.StatusOK,
			expectedBody: web.Result{
				Code: 5,
				Msg:  "系统错误",
			},
		},
		{
			name: "手机号码为空",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},
			reqBody:      `{"phone": ""}`,
			expectedCode: http.StatusOK,
			expectedBody: web.Result{
				Code: 4,
				Msg:  "输入错误",
			},
		},
		{
			name: "手机号码格式错误， bind 失败",
			before: func(t *testing.T) {

			},
			after: func(t *testing.T) {

			},
			// reqBody:      `{"phone": "",}`,
			reqBody:      `{"phone": ,}`,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			req, err := http.NewRequest(http.MethodPost, "/users/login_sms/code/send", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			assert.Equal(t, tc.expectedCode, resp.Code)
			if resp.Code != 200 {
				// 手机号码格式不对会导致反序列化失败
				return
			}
			var result web.Result
			// err = json.NewDecoder(resp.Body).Decode(&result)
			err = json.Unmarshal(resp.Body.Bytes(), &result)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedBody, result)
			tc.after(t)
		})
	}
}

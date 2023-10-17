package ratelimit

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"geektime/webook/internal/service/sms"
	smsmocks "geektime/webook/internal/service/sms/mocks"
	"geektime/webook/pkg/ratelimit"
	limitmocks "geektime/webook/pkg/ratelimit/mocks"
)

func TestRatelimitSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter)

		// 因为这边测试限流，输入什么不用管

		expectedErr error
	}{
		{
			name: "正常发送",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				svc := smsmocks.NewMockService(ctrl)
				limiter := limitmocks.NewMockLimiter(ctrl)
				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, nil)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return svc, limiter
			},
			expectedErr: nil,
		},
		{
			name: "限流异常",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				limiter := limitmocks.NewMockLimiter(ctrl)
				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(false, errors.New("mock limit err"))
				return nil, limiter
			},
			expectedErr: fmt.Errorf("短信服务判断是否限流异常：%w", errors.New("mock limit err")),
		},
		{
			name: "触发限流",
			mock: func(ctrl *gomock.Controller) (sms.Service, ratelimit.Limiter) {
				limiter := limitmocks.NewMockLimiter(ctrl)
				limiter.EXPECT().Limit(gomock.Any(), gomock.Any()).Return(true, nil)
				return nil, limiter
			},
			expectedErr: errLimited,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			limitSvc := NewRatelimitSMSService(tc.mock(ctrl))
			err := limitSvc.Send(context.Background(), "mytpl", []string{"123"}, "152xxx")
			fmt.Printf("%+v\n", err)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

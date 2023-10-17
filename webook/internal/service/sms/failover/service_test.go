package failover

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"geektime/webook/internal/service/sms"
	smsmocks "geektime/webook/internal/service/sms/mocks"
)

func TestFailoverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) []sms.Service

		expectedErr error
	}{
		{
			name: "重试最终都失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock send err"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock send err again"))
				return []sms.Service{svc0, svc1}
			},
			expectedErr: errors.New("全部服务商都失败了"),
		},
		{
			name: "重试成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock send err"))
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
			expectedErr: nil,
		},
		{
			name: "一次成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc := smsmocks.NewMockService(ctrl)
				svc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc}
			},
			expectedErr: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewFailoverSMSService(tc.mock(ctrl))
			err := svc.Send(context.Background(), "tplId", []string{"args"}, "152XXX")
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

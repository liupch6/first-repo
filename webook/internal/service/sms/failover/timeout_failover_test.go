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

func TestTimeoutFailoverSMSService_Send(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) []sms.Service

		threshold int32
		// 通过控制私有字段的取值，来模拟各种场景
		idx int32
		cnt int32

		expectedErr error
		expectedIdx int32
		expectedCnt int32
	}{
		{
			name: "触发了切换，切换之后成功",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				return []sms.Service{svc0, svc1}
			},
			threshold:   3,
			idx:         0,
			cnt:         3,
			expectedErr: nil,
			expectedIdx: 1,
			expectedCnt: 0,
		},
		{
			name: "触发了切换，切换之后依然超时",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded)
				return []sms.Service{svc0, svc1}
			},
			threshold:   3,
			idx:         0,
			cnt:         3,
			expectedErr: context.DeadlineExceeded,
			expectedIdx: 1,
			expectedCnt: 1,
		},
		{
			name: "触发了切换，切换之后失败",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc1 := smsmocks.NewMockService(ctrl)
				svc1.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("mock CAS error"))
				return []sms.Service{svc0, svc1}
			},
			threshold:   3,
			idx:         0,
			cnt:         3,
			expectedErr: errors.New("mock CAS error"),
			expectedIdx: 1,
			expectedCnt: 0,
		}, {
			name: "超时，但没连续超时",
			mock: func(ctrl *gomock.Controller) []sms.Service {
				svc0 := smsmocks.NewMockService(ctrl)
				svc0.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(context.DeadlineExceeded)
				return []sms.Service{svc0}
			},
			threshold:   3,
			idx:         0,
			cnt:         0,
			expectedErr: context.DeadlineExceeded,
			expectedIdx: 0,
			expectedCnt: 1,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := NewTimeoutFailoverSMSService(tc.mock(ctrl), tc.threshold)
			svc.idx = tc.idx
			svc.cnt = tc.cnt
			err := svc.Send(context.Background(), "tplId", []string{"args"}, "152xxx")

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedIdx, svc.idx)
			assert.Equal(t, tc.expectedCnt, svc.cnt)
		})
	}
}

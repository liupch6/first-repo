package failover

import (
	"context"
	"sync/atomic"

	"geektime/webook/internal/service/sms"
)

type TimeoutFailoverSMSService struct {
	svcs      []sms.Service // 你的服务商
	idx       int32
	cnt       int32 // 连续超时的个数
	threshold int32 // 阈值：连续超时超过这个数字，就要切换
}

func NewTimeoutFailoverSMSService(svcs []sms.Service, threshold int32) *TimeoutFailoverSMSService {
	return &TimeoutFailoverSMSService{
		svcs:      svcs,
		threshold: threshold,
	}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt >= t.threshold {
		// 触发切换，计算新的下标
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			// 说明你切换了
			atomic.StoreInt32(&t.cnt, 0)
		}
		// CAS 操作失败，说明出现并发，有人切换了
		// idx = newIdx
		idx = atomic.LoadInt32(&t.idx)
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tplId, args, numbers...)
	switch err {
	case nil:
		// 没有错误，你的连续状态被打断了，重置计数器
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
		return err
	default:
		// 不知道什么错误，可以考虑切换下一个，语义则是：
		// -超时错误：可能是偶发的，我尽量再试试
		// -非超时错误：我直接下一个
		return err
	}
}

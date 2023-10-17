package ratelimit

import (
	"context"
	"fmt"

	"geektime/webook/internal/service/sms"
	"geektime/webook/pkg/ratelimit"
)

var errLimited = fmt.Errorf("短信服务触发限流")

type RatelimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRatelimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RatelimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s *RatelimitSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 这里加一些代码，新特性
	limit, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		/*
			系统错误，要不要限流呢？
			可以限流：保守策略，你的下游很坑的时候
			可以不限流：你的下有很强，业务可用性要求很高，尽量容错策略
		*/
		return fmt.Errorf("短信服务判断是否限流异常：%w", err)
	}
	if limit {
		return errLimited
	}
	err = s.svc.Send(ctx, tplId, args, numbers...)
	// 这里也可以加一些代码，新特性
	return err
}

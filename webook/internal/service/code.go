package service

import (
	"context"
	"fmt"
	"math/rand"

	"geektime/webook/internal/repository"
	"geektime/webook/internal/service/sms"
)

const codeTplId = "1877556"

type CodeService struct {
	repo   *repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo *repository.CodeRepository, smsSvc sms.Service) *CodeService {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// Send 发送验证码
func (svc *CodeService) Send(ctx context.Context, biz, phone string) error {
	// 生成一个验证码
	code := svc.generateCode()
	// 塞进去 Redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	// 发送出去
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	return err
}

func (svc *CodeService) generateCode() string {
	// 生成 6 位数的验证码(0~999999)
	num := rand.Intn(1000000)
	// 6 位数，不足 6 位前面补 0
	return fmt.Sprintf("%06d", num)
}

func (svc *CodeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

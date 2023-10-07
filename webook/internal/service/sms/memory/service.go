package memory

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	fmt.Printf("%v 验证码是%v\n", time.Now().Format("2006-01-02 15:04:05"), args)
	return nil
}

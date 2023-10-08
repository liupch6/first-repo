package ioc

import (
	"geektime/webook/internal/service/sms"
	"geektime/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewService()
}

package tencent

import (
	"context"
	"testing"

	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func TestService_Send(t *testing.T) {
	type fields struct {
		appId    *string
		signName *string
		client   *sms.Client
	}
	type args struct {
		ctx     context.Context
		tplId   string
		args    []string
		numbers []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				appId:    tt.fields.appId,
				signName: tt.fields.signName,
				client:   tt.fields.client,
			}
			if err := s.Send(tt.args.ctx, tt.args.tplId, tt.args.args, tt.args.numbers...); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

package sms

import (
	"context"
	"fmt"

	"github.com/maxfeizi04-cloud/recruitment-platform/internal/config"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	smsapi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Sender interface {
	SendVerificationCode(ctx context.Context, phone, code string) error
}

type TencentSMS struct {
	client     *smsapi.Client
	sdkAppID   string
	templateID string
	signName   string
}

func NewTencentSMS(cfg config.SMSConfig) (Sender, error) {
	credential := common.NewCredential(cfg.SecretID, cfg.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	client, err := smsapi.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		return nil, fmt.Errorf("create sms client: %w", err)
	}

	return &TencentSMS{
		client:     client,
		sdkAppID:   cfg.SDKAppID,
		templateID: cfg.TemplateID,
		signName:   cfg.SignName,
	}, nil
}

func (s *TencentSMS) SendVerificationCode(ctx context.Context, phone, code string) error {
	req := smsapi.NewSendSmsRequest()
	req.PhoneNumberSet = common.StringPtrs([]string{"+86" + phone})
	req.SmsSdkAppId = common.StringPtr(s.sdkAppID)
	req.TemplateId = common.StringPtr(s.templateID)
	req.SignName = common.StringPtr(s.signName)
	req.TemplateParamSet = common.StringPtrs([]string{code})

	resp, err := s.client.SendSms(req)
	if err != nil {
		return fmt.Errorf("send sms: %w", err)
	}

	for _, status := range resp.Response.SendStatusSet {
		if *status.Code != "Ok" {
			return fmt.Errorf("sms send failed: code=%s, msg=%s", *status.Code, *status.Message)
		}
	}

	return nil
}

type MockSender struct{}

func (m *MockSender) SendVerificationCode(ctx context.Context, phone, code string) error {
	fmt.Printf("[MOCK SMS] To: %s, Code: %s\n", phone, code)
	return nil
}

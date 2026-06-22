package auth

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	pkgauth "github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/auth"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/broker"
	redisclient "github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/redis"
	"github.com/maxfeizi04-cloud/recruitment-platform/internal/pkg/sms"
)

type Service struct {
	repo       *Repository
	jwtManager *pkgauth.JWTManager
	smsSender  sms.Sender
	redis      *redisclient.Client
	broker     broker.MessageBroker
}

func NewService(
	repo *Repository,
	jwtManager *pkgauth.JWTManager,
	smsSender sms.Sender,
	redis *redisclient.Client,
	broker broker.MessageBroker,
) *Service {
	return &Service{
		repo:       repo,
		jwtManager: jwtManager,
		smsSender:  smsSender,
		redis:      redis,
		broker:     broker,
	}
}

func (s *Service) SendVerificationCode(ctx context.Context, phone string) error {
	allowed, err := s.redis.CanSendCode(ctx, phone)
	if err != nil {
		return fmt.Errorf("rate limit check: %w", err)
	}
	if !allowed {
		return fmt.Errorf("verification code sent too recently, please wait 60 seconds")
	}

	code := generateCode()

	if err := s.redis.SetVerificationCode(ctx, phone, code); err != nil {
		return fmt.Errorf("store code: %w", err)
	}

	if err := s.smsSender.SendVerificationCode(ctx, phone, code); err != nil {
		return fmt.Errorf("send sms: %w", err)
	}

	return nil
}

func (s *Service) Login(ctx context.Context, phone, code, role string) (*LoginResult, error) {
	valid, err := s.redis.GetAndDeleteVerificationCode(ctx, phone, code)
	if err != nil {
		return nil, fmt.Errorf("verify code: %w", err)
	}
	if !valid {
		return nil, fmt.Errorf("invalid or expired verification code")
	}

	user, err := s.repo.FindOrCreate(ctx, phone, role)
	if err != nil {
		return nil, fmt.Errorf("find or create user: %w", err)
	}

	token, err := s.jwtManager.Generate(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	s.broker.Publish(ctx, "user.login", []byte(fmt.Sprintf(`{"user_id":"%s","role":"%s"}`, user.ID, user.Role)))

	return &LoginResult{
		Token: token,
		User:  *user,
	}, nil
}

type LoginResult struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

func generateCode() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06d", rng.Intn(1000000))
}

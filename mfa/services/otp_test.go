package services

import (
	"context"
	"mfa/config"
	"testing"
	"time"
)

const email = "test@mail.com"

func TestGenerateVerifyOtp(t *testing.T) {
	srv := getOtpService()

	otp, err := srv.Generate(context.Background(), email)
	if err != nil {
		t.Errorf("Error generating otp: %v", err)
	}

	if otp == "" {
		t.Error("Otp is empty")
	}

	verified, err := srv.Verify(context.Background(), email, otp)
	if err != nil {
		t.Errorf("Error verifying otp: %v", err)
	}

	if !verified {
		t.Error("Otp is not valid")
	}
}

func TestGenerateDifferentOtp(t *testing.T) {
	srv := getOtpService()

	otp1, err := srv.Generate(context.Background(), email)
	if err != nil {
		t.Errorf("Error generating otp: %v", err)
	}

	otp2, err := srv.Generate(context.Background(), email)
	if err != nil {
		t.Errorf("Error generating otp: %v", err)
	}

	if otp1 == otp2 {
		t.Error("Otp is not different")
	}
}

func getOtpService() OtpService {
	return &otpService{
		cfg: config.Config{
			OtpExpiry: 10 * time.Minute,
		},
		repository: &otpRepositoryMock{
			store: make(map[string]string),
		},
	}
}

type otpRepositoryMock struct {
	store map[string]string
}

func (r *otpRepositoryMock) Set(ctx context.Context, key string, otp string, expiry time.Duration) error {
	r.store[key] = otp
	return nil
}

func (r *otpRepositoryMock) Get(ctx context.Context, key string) (string, error) {
	return r.store[key], nil
}

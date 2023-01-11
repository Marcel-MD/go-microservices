package repositories

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/rs/zerolog/log"
)

type IOtpRepository interface {
	Set(ctx context.Context, key string, otp string, expiry time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

type otpRepository struct {
	rdb *redis.Client
}

var (
	otpOnce sync.Once
	otpRepo IOtpRepository
)

func GetOtpRepository() IOtpRepository {

	otpOnce.Do(func() {
		log.Info().Msg("Initializing otp repository")

		otpRepo = &otpRepository{
			rdb: GetRDB(),
		}
	})

	return otpRepo
}

func (r *otpRepository) Set(ctx context.Context, key string, otp string, expiry time.Duration) error {
	return r.rdb.Set(ctx, key, otp, expiry).Err()
}

func (r *otpRepository) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

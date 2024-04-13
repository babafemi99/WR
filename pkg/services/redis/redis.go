package redis

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"github.com/babafemi99/WR/internal/config"
	"github.com/babafemi99/WR/internal/values"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type IRedisService interface {
	SetStaffSessionToken(ctx context.Context, session AuthSession) error
	SetAuthSession(ctx context.Context, session AuthSession, refreshToken string) error
	GetAuthSession(ctx context.Context, userType, userId string) (AuthSession, error)
	GetRefreshSession(ctx context.Context, userType, userId, refreshToken string) (AuthSession, error)
	DeleteAuthSession(ctx context.Context, session []string) error
}

type Service struct {
	Client *redis.Client
}

func (s Service) SetStaffSessionToken(ctx context.Context, session AuthSession) error {
	var b bytes.Buffer

	if encodeErr := gob.NewEncoder(&b).Encode(&session); encodeErr != nil {
		return encodeErr
	}

	err := s.Client.Set(ctx, fmt.Sprintf("%s-%s", session.For, session.UserID), b.Bytes(),
		values.StaffTokenExpiry*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetRefreshSession(ctx context.Context, userType, userId, refreshToken string) (AuthSession, error) {
	key := fmt.Sprintf("ref-%s-%s-%s", userType, userId, refreshToken)
	cmd := s.Client.Get(ctx, key)

	cmdb, err := cmd.Bytes()
	if err != nil {
		return AuthSession{}, err
	}

	b := bytes.NewReader(cmdb)

	var session AuthSession
	if decodeErr := gob.NewDecoder(b).Decode(&session); decodeErr != nil {
		return AuthSession{}, err
	}

	return session, nil
}

func (s Service) DeleteAuthSession(ctx context.Context, session []string) error {
	return s.Client.Del(ctx, session...).Err()
}

func (s Service) GetAuthSession(ctx context.Context, userType, userId string) (AuthSession, error) {

	cmd := s.Client.Get(context.TODO(), fmt.Sprintf("%s-%s", userType, userId))

	cmdb, err := cmd.Bytes()
	if err != nil {
		return AuthSession{}, err
	}

	b := bytes.NewReader(cmdb)

	var session AuthSession
	if decodeErr := gob.NewDecoder(b).Decode(&session); decodeErr != nil {
		return AuthSession{}, err
	}

	return session, nil
}

func (s Service) SetAuthSession(ctx context.Context, session AuthSession, refreshToken string) error {
	var b bytes.Buffer

	if encodeErr := gob.NewEncoder(&b).Encode(&session); encodeErr != nil {
		return encodeErr
	}

	err := s.Client.Set(ctx, fmt.Sprintf("%s-%s", session.For, session.UserID), b.Bytes(),
		values.AccessTokenExpiry*time.Minute).Err()
	if err != nil {
		return err
	}

	err = s.Client.Set(ctx, fmt.Sprintf("ref-%s-%s-%s", session.For, session.UserID, refreshToken), b.Bytes(),
		values.RefreshTokenExpiry*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

// connectRedis creates a connection to redis
func connectRedis(cfg *config.Config) *redis.Client {
	opt, parseErr := redis.ParseURL(cfg.RedisURL)
	if parseErr != nil {
		log.Fatal("[Redis]: failed to initiate connection")
	}

	client := redis.NewClient(opt)
	_, pingErr := client.Ping(context.Background()).Result()
	if pingErr != nil {
		log.Fatalf("[Redis]: unable to ping redis client: %v", pingErr)
	}
	return client
}

func New(cfg *config.Config) IRedisService {
	client := connectRedis(cfg)
	return Service{
		Client: client,
	}
}

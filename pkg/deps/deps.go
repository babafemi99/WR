package deps

import (
	"github.com/babafemi99/WR/internal/config"
	"github.com/babafemi99/WR/pkg/db/cockroachDB"
	"github.com/babafemi99/WR/pkg/services/mail"
	"github.com/babafemi99/WR/pkg/services/redis"
)

type Dependencies struct {
	//Repo Layer
	Repository *cockroachDB.Repository

	// Services
	Redis redis.IRedisService
	//email
	IMailer mail.IMailer
}

func New(cfg *config.Config) *Dependencies {

	deps := Dependencies{
		Repository: cockroachDB.New(cfg.DataBaseUrl),
		Redis:      redis.New(cfg),
		IMailer:    mail.NewMailgun(cfg),
	}
	return &deps
}

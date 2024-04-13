package deps

import (
	"github.com/babafemi99/WR/pkg/db/weddingRepo"
	"github.com/babafemi99/WR/pkg/services/redis"
)

type Dependencies struct {
	//Repo Layer
	Repository *weddingRepo.Repository

	// Services
	Redis redis.IRedisService
	//email
}

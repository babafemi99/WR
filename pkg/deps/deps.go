package deps

import (
	"github.com/babafemi99/WR/pkg/db/weddingRepo"
)

type Dependencies struct {
	//Repo Layer
	*weddingRepo.Repository

	// Services
	//email
}

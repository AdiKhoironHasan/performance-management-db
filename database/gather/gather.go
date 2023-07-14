package gather

import (
	"engine-db/logger"

	"gorm.io/gorm"
)

var (
	logs = logger.NewLogger().Logger.With().Str("pkg", "gather").Logger()
)

type gather struct {
	db *gorm.DB
}

func NewGather(db *gorm.DB) Gather {
	return &gather{
		db: db,
	}
}

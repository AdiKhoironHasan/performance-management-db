package database

import (
	"engine-db/entity"
	"engine-db/logger"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	logs = logger.NewLogger().Logger.With().Str("pkg", "seeder").Logger()
)

func New() *gorm.DB {
	var sqlQueryLogLevel gormLogger.LogLevel = gormLogger.Info

	queryLogging := newSQLLogging(SQLLogConfig{
		Level:                     sqlQueryLogLevel,
		IgnoreRecordNotFoundError: true,
	})
	db, err := gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=project-PA-7 port=5432 sslmode=disable"),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,

			// set default to no log
			Logger:               queryLogging,
			DisableAutomaticPing: true,
		},
	)
	if err != nil {
		logs.Error().Err(err).Msg("failed to connect to database")
		os.Exit(1)
	}

	err = db.Exec("DROP TABLE IF EXISTS users, events, sessions, questions, question_answers;").Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to drop table")
		os.Exit(1)
	}

	err = db.AutoMigrate(
		&entity.Versioning{},
		&entity.User{},
		&entity.Event{},
		&entity.Session{},
		&entity.Question{},
		&entity.QuestionAnswer{},
		&entity.FormAnswerScale{},
		&entity.FormAnswerText{},
	)
	if err != nil {
		logs.Error().Err(err).Msg("failed to migrate")
		os.Exit(1)
	}

	return db
}

package database

import (
	entity "engine-db/entity/new"
	"engine-db/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var (
	logs = logger.NewLogger().Logger.With().Str("pkg", "database").Logger()
)

func New() (*gorm.DB, error) {
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
		return nil, err
	}

	err = db.Exec(`DROP TABLE IF EXISTS
			users,
			enterprises,
			user_versions,
			events,
			sessions,
			questions,
			form_tasks,
			form_scale_answers,
			form_text_answers,
			documents
	`).Error
	if err != nil {
		logs.Error().Err(err).Msg("failed to drop table")
		return nil, err

	}

	err = db.AutoMigrate(
		&entity.User{},
		&entity.Enterprise{},
		&entity.UserVersion{},
		&entity.Event{},
		&entity.Session{},
		&entity.Question{},
		&entity.FormTask{},
		&entity.FormScaleAnswer{},
		&entity.FormTextAnswer{},
		&entity.Document{},
	)
	if err != nil {
		logs.Error().Err(err).Msg("failed to migrate")
		return nil, err

	}

	return db, nil
}

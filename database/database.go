package database

import (
	"engine-db/logger"
	"fmt"
	"net/url"

	_ "github.com/jinzhu/gorm/dialects/postgres"
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
	dsn := url.URL{
		User:     url.UserPassword("USER_NAME", "PASSWORD"),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", "HOST", 5432),
		Path:     "DB_NAME",
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
	x := dsn.String()
	_ = x
	db, err := gorm.Open(postgres.Open(dsn.String()),
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

	// err = db.Exec(`DROP TABLE IF EXISTS
	// 		users,
	// 		enterprises,
	// 		user_versions,
	// 		events,
	// 		sessions,
	// 		questions,
	// 		form_tasks,
	// 		form_scale_answers,
	// 		form_text_answers,
	// 		documents
	// `).Error
	// if err != nil {
	// 	logs.Error().Err(err).Msg("failed to drop table")
	// 	return nil, err

	// }

	// err = db.AutoMigrate(
	// 	&entity.User{},
	// 	&entity.Enterprise{},
	// 	&entity.UserVersion{},
	// 	&entity.Event{},
	// 	&entity.Session{},
	// 	&entity.Question{},
	// 	&entity.FormTask{},
	// 	&entity.FormScaleAnswer{},
	// 	&entity.FormTextAnswer{},
	// 	&entity.Document{},
	// )
	// if err != nil {
	// 	logs.Error().Err(err).Msg("failed to migrate")
	// 	return nil, err

	// }

	return db, nil
}

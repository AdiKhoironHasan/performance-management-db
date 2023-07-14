package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/timemore/foundation/errors"
	"github.com/timemore/foundation/logger"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type sqlLogger struct {
	logger.PkgLogger
	SQLLogConfig
}

var _ gormlogger.Interface = sqlLogger{}

func (l sqlLogger) LogMode(logLevel gormlogger.LogLevel) gormlogger.Interface {
	l.Level = logLevel
	return l
}

func (l sqlLogger) LogRecordNotFoundError(val bool) gormlogger.Interface {
	l.IgnoreRecordNotFoundError = val

	return l
}

func (l sqlLogger) Info(ctx context.Context, msg string, data ...any) {
	if l.Level >= gormlogger.Info {
		l.Logger.Info().Fields(map[string]any{
			"file": utils.FileWithLineNum(),
		}).Msgf(msg, data...)
	}
}
func (l sqlLogger) Warn(ctx context.Context, msg string, data ...any) {
	if l.Level >= gormlogger.Warn {
		l.Logger.Warn().Fields(map[string]any{
			"file": utils.FileWithLineNum(),
		}).Msgf(msg, data...)
	}
}
func (l sqlLogger) Error(ctx context.Context, msg string, data ...any) {
	if l.Level >= gormlogger.Error {
		l.Logger.Warn().Fields(map[string]any{
			"file": utils.FileWithLineNum(),
		}).Msgf(msg, data...)
	}
}
func (l sqlLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	elapsed := time.Since(begin)
	qry, rows := fc()
	switch {
	case err != nil && l.Level >= gormlogger.Error && errors.Is(err, gorm.ErrRecordNotFound) && !l.IgnoreRecordNotFoundError:
		if err == gorm.ErrRecordNotFound || err == sql.ErrNoRows {
			if l.Level >= gormlogger.Info {
				l.Logger.Info().Fields(map[string]any{
					"file":    utils.FileWithLineNum(),
					"type":    "sql",
					"latency": float64(elapsed.Nanoseconds()) / 1e6,
					"rows":    rows,
				}).Msg(qry)
			}
			return
		}
		l.Logger.Warn().Err(err).Fields(map[string]any{
			"file":    utils.FileWithLineNum(),
			"type":    "sql",
			"latency": float64(elapsed.Nanoseconds()) / 1e6,
		}).Msg(qry)
	case l.Level >= gormlogger.Info:
		l.Logger.Info().Fields(map[string]any{
			"file":    utils.FileWithLineNum(),
			"type":    "sql",
			"latency": float64(elapsed.Nanoseconds()) / 1e6,
			"rows":    rows,
		}).Msg(qry)
	}

}

type SQLLogConfig struct {
	Level                     gormlogger.LogLevel
	IgnoreRecordNotFoundError bool
}

func newSQLLogging(cfg SQLLogConfig) sqlLogger {
	return sqlLogger{
		SQLLogConfig: cfg,
		PkgLogger:    logger.NewPkgLogger(),
	}
}

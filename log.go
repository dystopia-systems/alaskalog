package alaskalog

import (
	"github.com/kz/discordrus"
	"github.com/orandin/lumberjackrus"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Logger = logrus.New()

func Configure() {

	Logger.Formatter = &TextFormatter{TimestampFormat: time.RFC822, FullTimestamp: true}
	Logger.Level = logrus.DebugLevel

	Logger.AddHook(discordrus.NewHook(
		os.Getenv("DISCORDRUS_WEBHOOK_URL"),
		logrus.WarnLevel,
		&discordrus.Opts{
			Username: os.Getenv("LOG_USERNAME"),
			TimestampFormat: defaultTimestampFormat,
			DisableTimestamp: false,
			EnableCustomColors: true,
			CustomLevelColors: &discordrus.LevelColors{
				Trace: 3092790,
				Debug: 10170623,
				Info:  3581519,
				Warn:  14327864,
				Error: 13631488,
				Panic: 13631488,
				Fatal: 13631488,
			},
		}))

	if isProd() {
		lumberjackrusHook, err := lumberjackrus.NewHook(
			&lumberjackrus.LogFile{
				Filename:   "logs/general.json",
				MaxSize:    20,
				MaxBackups: 1,
				MaxAge:     1,
				Compress:   false,
				LocalTime:  false,
			},
			logrus.InfoLevel,
			&logrus.JSONFormatter{},
			&lumberjackrus.LogFileOpts{
				logrus.WarnLevel: &lumberjackrus.LogFile{
					Filename: "logs/warn.json",
					MaxSize: 20,
				},
				logrus.FatalLevel: &lumberjackrus.LogFile{
					Filename: "logs/error.json",
					MaxSize: 50,
				},
			})

		if err != nil {
			return
		}

		Logger.AddHook(lumberjackrusHook)
	}

}

func isProd() bool {
	for _, arg := range os.Args {
		switch arg {
		case "-p":
			return true
		}
	}

	return false
}

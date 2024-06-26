package bot

import (
	"fmt"
	"ratatoskr/internal/config"
	"ratatoskr/internal/db"
	"ratatoskr/internal/logger"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func Run(
	db db.DB,
	logger *logger.Logger, config *config.BotConfig) error {
	logger.Info("initializing bot...")
	bot, err := gotgbot.NewBot(config.Token, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to initialize new bot, error: %v", err))
		return err
	}
	logger.Info("bot initialized")

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
		Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
			return ext.DispatcherActionNoop
		},
		MaxRoutines: ext.DefaultMaxRoutines,
	})

	updater := ext.NewUpdater(dispatcher, nil)

	addHandlers(db, dispatcher, logger, config)

	logger.Info("staring polling...")
	err = updater.StartPolling(bot, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: &gotgbot.GetUpdatesOpts{
			Timeout: 9,
			RequestOpts: &gotgbot.RequestOpts{
				Timeout: time.Second * 10,
			},
		},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("failed to start polling, error: %v", err))
		return err
	}
	logger.Info("polling started")
	logger.Info(fmt.Sprintf("WebApp url - %s", config.WebAppUrl))
	logger.Info(fmt.Sprintf("%s is live", bot.FirstName))
	updater.Idle()
	return nil
}

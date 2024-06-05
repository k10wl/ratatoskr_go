package bot

import (
	"fmt"
	"ratatoskr/internal/config"
	"ratatoskr/internal/logger"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

func addHandlers(dispatcher *ext.Dispatcher, logger *logger.Logger, config *config.Config) {
	dispatcher.AddHandler(
		handlers.NewMessage(message.Text,
			adminOnly(logger, config, handleMessage(logger))),
	)
}

func handleMessage(logger *logger.Logger) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
		if err != nil {
			return logger.Error(fmt.Sprintf("failed to echo message: %v", err))
		}
		return nil
	}
}

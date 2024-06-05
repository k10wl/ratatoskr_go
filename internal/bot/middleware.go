package bot

import (
	"ratatoskr/internal/config"
	"ratatoskr/internal/logger"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

func adminOnly(
	logger *logger.Logger,
	config *config.Config,
	next handlers.Response,
) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		if slices.Index(config.AdminIDs, ctx.EffectiveSender.User.Id) == -1 {
			return logger.Error("unauthorized sender")
		}
		return next(b, ctx)
	}
}

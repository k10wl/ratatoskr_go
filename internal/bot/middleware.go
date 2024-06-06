package bot

import (
	"fmt"
	"ratatoskr/internal/config"
	"ratatoskr/internal/logger"
	"slices"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

type middleware struct {
	logger *logger.Logger
	config *config.Config
}

func newMidlleware(logger *logger.Logger, config *config.Config) *middleware {
	return &middleware{
		logger: logger,
		config: config,
	}
}

func (m middleware) adminOnly(
	next handlers.Response,
) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		if slices.Index(m.config.AdminIDs, ctx.EffectiveSender.User.Id) == -1 {
			return m.logger.Error(fmt.Sprintf("unauthorized sender %+v", ctx.EffectiveSender))
		}
		return next(b, ctx)
	}
}

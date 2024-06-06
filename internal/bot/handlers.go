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

type handler struct {
	logger *logger.Logger
}

func newHandler(
	logger *logger.Logger,
) *handler {
	return &handler{
		logger: logger,
	}
}

func addHandlers(dispatcher *ext.Dispatcher, logger *logger.Logger, config *config.Config) {
	handler := newHandler(logger)
	middleware := newMidlleware(logger, config)
	dispatcher.AddHandler(
		handlers.NewMessage(message.Text,
			middleware.adminOnly(handler.handleEchoMessage(handler.removeOriginal()))),
	)
	dispatcher.AddHandler(
		handlers.NewMessage(message.MediaGroup,
			middleware.adminOnly(handler.handleMediaGroup(handler.removeOriginal()))),
	)
	dispatcher.AddHandler(
		handlers.NewMessage(message.Photo,
			middleware.adminOnly(handler.handlePhoto(handler.removeOriginal()))),
	)
	dispatcher.AddHandler(
		handlers.NewMessage(message.Video,
			middleware.adminOnly(handler.handleVideo(handler.removeOriginal()))),
	)
	dispatcher.AddHandler(
		handlers.NewMessage(message.Animation,
			middleware.adminOnly(handler.handleAnimation(handler.removeOriginal()))),
	)
}

func (h handler) handleMediaGroup(next handlers.Response) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		chatID := ctx.EffectiveMessage.GetSender().Id()
		h.logger.Info(fmt.Sprintf("%+v", ctx.EffectiveMessage))
		b.SendMessage(
			chatID,
			"I don't know how to process group yet",
			&gotgbot.SendMessageOpts{},
		)
		return next(b, ctx)
	}
}

func (h handler) handlePhoto(next handlers.Response) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		chatID := ctx.EffectiveMessage.GetSender().Id()
		_, err := b.SendPhoto(
			chatID,
			ctx.EffectiveMessage.Photo[0].FileId,
			&gotgbot.SendPhotoOpts{},
		)
		if err != nil {
			return h.logger.Error(
				fmt.Sprintf("failed to reply with photo, error: %v", err),
			)
		}
		h.logger.Info(fmt.Sprintf("photo message reply success"))
		return next(b, ctx)
	}
}

func (h handler) handleVideo(next handlers.Response) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		_, err := b.SendVideo(
			ctx.EffectiveMessage.GetSender().Id(),
			ctx.EffectiveMessage.Video.FileId,
			&gotgbot.SendVideoOpts{},
		)
		if err != nil {
			return h.logger.Error(
				fmt.Sprintf("failed to reply with video, error: %v", err),
			)
		}
		h.logger.Info(fmt.Sprintf("video message reply success"))
		return next(b, ctx)
	}
}

func (h handler) handleAnimation(next handlers.Response) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		chatID := ctx.EffectiveMessage.GetSender().Id()
		_, err := b.SendAnimation(
			chatID,
			ctx.EffectiveMessage.Animation.FileId,
			&gotgbot.SendAnimationOpts{},
		)
		if err != nil {
			return h.logger.Error(
				fmt.Sprintf("failed to reply with animation, error: %v", err),
			)
		}
		h.logger.Info(fmt.Sprintf("animation message reply success"))
		return next(b, ctx)
	}
}

func (h handler) handleEchoMessage(next handlers.Response) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		h.logger.Info(fmt.Sprintf("%+v", ctx.EffectiveMessage.GetChat()))
		_, err := ctx.EffectiveMessage.Reply(b, ctx.EffectiveMessage.Text, nil)
		if err != nil {
			return h.logger.Error(fmt.Sprintf("failed to echo message: %v", err))
		}
		h.logger.Info("echo message reply success")
		return next(b, ctx)
	}
}

func (h handler) removeOriginal() handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		ok, err := b.DeleteMessage(
			ctx.EffectiveMessage.GetSender().Id(),
			ctx.EffectiveMessage.MessageId,
			&gotgbot.DeleteMessageOpts{},
		)
		if !ok {
			h.logger.Warning("failed to delete video reply message")
		}
		if err != nil {
			h.logger.Warning(fmt.Sprintf("failed to delete video reply message, error: %v", err))
		}
		h.logger.Info("original message removed successfully")
		return nil
	}
}

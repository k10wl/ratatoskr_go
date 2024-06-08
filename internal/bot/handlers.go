package bot

import (
	"fmt"
	"ratatoskr/internal/config"
	"ratatoskr/internal/logger"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

type handler struct {
	logger        *logger.Logger
	mediaGroupMap map[string][]string
}

func newHandler(
	logger *logger.Logger,
) *handler {
	return &handler{
		logger:        logger,
		mediaGroupMap: map[string][]string{},
	}
}

func addHandlers(
	dispatcher *ext.Dispatcher,
	logger *logger.Logger,
	config *config.Config,
) {
	handler := newHandler(logger)
	middleware := newMidlleware(logger, config)

	dispatcher.AddHandler(
		handlers.NewMessage(
			message.Text,
			middleware.adminOnly(
				handler.handleEchoMessage(
					handler.removeOneEffectiveMessage(),
				),
			),
		),
	)

	dispatcher.AddHandler(
		handlers.NewMessage(
			message.MediaGroup,
			middleware.adminOnly(
				handler.receiveGroup(
					time.Millisecond*500,
					handler.respondWithMediaGroup(),
				),
			),
		),
	)

	dispatcher.AddHandler(
		handlers.NewMessage(
			message.Photo,
			middleware.adminOnly(
				handler.handlePhoto(
					handler.removeOneEffectiveMessage(),
				),
			),
		),
	)

	dispatcher.AddHandler(
		handlers.NewMessage(
			message.Video,
			middleware.adminOnly(
				handler.handleVideo(
					handler.removeOneEffectiveMessage(),
				),
			),
		),
	)

	dispatcher.AddHandler(
		handlers.NewMessage(
			message.Animation,
			middleware.adminOnly(
				handler.handleAnimation(
					handler.removeOneEffectiveMessage(),
				),
			),
		),
	)
}

func (h handler) handlePhoto(next handlers.Response) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		chatID := ctx.EffectiveMessage.GetSender().Id()
		_, err := sendPhoto(
			b,
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
		_, err := sendVideo(
			b,
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
		_, err := sendAnimation(
			b,
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

func (h handler) removeOneEffectiveMessage() handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		ok, err := deleteMessage(
			b,
			ctx.EffectiveMessage.GetSender().Id(),
			ctx.EffectiveMessage.MessageId,
		)
		if ok {
			h.logger.Info("original message removed successfully")
		}
		if err != nil {
			h.logger.Warning(fmt.Sprintf("failed to delete video reply message, error: %v", err))
		}
		return nil
	}
}

func (h *handler) receiveGroup(
	interval time.Duration,
	next handlers.Response,
) handlers.Response {
	lastMessageID := map[string]int64{}
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		val, _ := h.mediaGroupMap[ctx.EffectiveMessage.MediaGroupId]
		var mediaFileID string
		if ctx.EffectiveMessage.Video != nil {
			mediaFileID = ctx.EffectiveMessage.Video.FileId
		} else {
			mediaFileID = ctx.EffectiveMessage.Photo[0].FileId
		}
		h.mediaGroupMap[ctx.EffectiveMessage.MediaGroupId] = append(
			val,
			mediaFileID,
		)
		lastMessageID[ctx.EffectiveMessage.MediaGroupId] = ctx.EffectiveMessage.MessageId
		go func() {
			time.Sleep(interval)
			if lastMessageID[ctx.EffectiveMessage.MediaGroupId] != ctx.EffectiveMessage.MessageId {
				return
			}
			// I WANT TO BREAK FREE
			delete(lastMessageID, ctx.EffectiveMessage.MediaGroupId)
			next(b, ctx)
		}()
		return nil
	}
}

func (h handler) respondWithMediaGroup() handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		group := []gotgbot.InputMedia{}
		for _, v := range h.mediaGroupMap {
			for _, id := range v {
				group = append(group, gotgbot.InputMediaPhoto{Media: id})
			}
		}
		b.SendMediaGroup(ctx.EffectiveSender.ChatId, group, &gotgbot.SendMediaGroupOpts{})
		return nil
	}
}

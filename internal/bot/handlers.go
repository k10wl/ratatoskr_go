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
	mediaGroupMap *mediaGroupMap
}

func newHandler(
	logger *logger.Logger,
) *handler {
	return &handler{
		logger:        logger,
		mediaGroupMap: newMediaGroupMap(),
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
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		var mediaFileID string
		var mediaType string
		if ctx.EffectiveMessage.Video != nil {
			mediaFileID = ctx.EffectiveMessage.Video.FileId
			mediaType = "video"
		} else {
			mediaFileID = ctx.EffectiveMessage.Photo[0].FileId
			mediaType = "photo"
		}
		h.mediaGroupMap.add(ctx.EffectiveMessage.MediaGroupId, item{
			fileID:    mediaFileID,
			mediaType: mediaType,
			messageID: ctx.EffectiveMessage.MessageId,
		})
		go func() {
			time.Sleep(interval)
			related := h.mediaGroupMap.get(ctx.EffectiveMessage.MediaGroupId)
			if len(related) == 0 {
				h.logger.Error("map is empty for " + ctx.EffectiveMessage.MediaGroupId)
				return
			}
			if related[0].messageID != ctx.EffectiveMessage.MessageId {
				return
			}
			next(b, ctx)
		}()
		return nil
	}
}

func (h handler) respondWithMediaGroup() handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		group := []gotgbot.InputMedia{}
		for _, item := range h.mediaGroupMap.get(ctx.EffectiveMessage.MediaGroupId) {
			switch item.mediaType {
			case "photo":
				group = append(group, gotgbot.InputMediaPhoto{Media: item.fileID})
			case "video":
				group = append(group, gotgbot.InputMediaVideo{Media: item.fileID})
			default:
				h.logger.Error(fmt.Sprintf("unhandler media type in %+v", item))
			}
		}
		sendMediaGroup(b, ctx.EffectiveSender.ChatId, group, &gotgbot.SendMediaGroupOpts{})
		h.mediaGroupMap.remove(ctx.EffectiveMessage.MediaGroupId)
		return nil
	}
}

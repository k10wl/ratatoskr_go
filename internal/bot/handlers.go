package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"ratatoskr/internal/config"
	"ratatoskr/internal/db"
	"ratatoskr/internal/logger"
	"ratatoskr/internal/models"
	"ratatoskr/internal/utils"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

type handler struct {
	logger        *logger.Logger
	mediaGroupMap *mediaGroupMap
	config        *config.BotConfig
	db            db.DB
}

func newHandler(
	db db.DB,
	logger *logger.Logger,
	config *config.BotConfig,
) *handler {
	return &handler{
		config:        config,
		logger:        logger,
		mediaGroupMap: newMediaGroupMap(),
		db:            db,
	}
}

func addHandlers(
	db db.DB,
	dispatcher *ext.Dispatcher,
	logger *logger.Logger,
	config *config.BotConfig,
) {
	handler := newHandler(db, logger, config)
	middleware := newMidlleware(logger, config)

	dispatcher.AddHandler(
		handlers.NewCommand("ping",
			middleware.adminOnly(
				handler.handlePing()),
		),
	)

	dispatcher.AddHandler(
		handlers.NewMessage(isTagsMessage, middleware.adminOnly(handler.handleUpdateTags())),
	)

	dispatcher.AddHandler(
		handlers.NewMessage(func(msg *gotgbot.Message) bool {
			return msg.WebAppData != nil
		},
			middleware.adminOnly(
				handler.handleWebAppData(time.Now)),
		),
	)

	dispatcher.AddHandler(
		handlers.NewMessage(
			message.MediaGroup,
			middleware.adminOnly(
				handler.receiveGroup(
					time.Millisecond*500,
					handler.respondWithMediaGroup(
						handler.removeEffectiveMediaGroup(),
					),
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
		h.logger.Info(fmt.Sprintf("received photo %d", ctx.EffectiveMessage.MessageId))
		m, err := sendPhoto(
			b,
			ctx.EffectiveChat.Id,
			ctx.EffectiveMessage.Photo[0].FileId,
			&gotgbot.SendPhotoOpts{},
		)
		if err != nil {
			return h.logger.Error(
				fmt.Sprintf("failed to reply with photo, error: %v", err),
			)
		}
		err = h.sendWebAppMarkup(b, ctx.EffectiveChat.Id, []int64{m.MessageId})
		if err != nil {
			return h.logger.Error(
				fmt.Sprintf("failed to reply with webapp, error: %v", err),
			)
		}
		h.logger.Info(
			fmt.Sprintf("photo message reply success %d", m.MessageId),
		)
		return next(b, ctx)
	}
}

func (h handler) handleVideo(next handlers.Response) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		h.logger.Info(fmt.Sprintf("received video %d", ctx.EffectiveMessage.MessageId))
		m, err := sendVideo(
			b,
			ctx.EffectiveChat.Id,
			ctx.EffectiveMessage.Video.FileId,
			&gotgbot.SendVideoOpts{},
		)
		if err != nil {
			return h.logger.Error(
				fmt.Sprintf("failed to reply with video, error: %v", err),
			)
		}
		err = h.sendWebAppMarkup(b, ctx.EffectiveChat.Id, []int64{m.MessageId})
		h.logger.Info(
			fmt.Sprintf("video message reply success %d", m.MessageId),
		)
		if err != nil {
			return h.logger.Error(
				fmt.Sprintf("failed to reply with video, error: %v", err),
			)
		}
		return next(b, ctx)
	}
}

func (h handler) handleAnimation(next handlers.Response) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		h.logger.Info(fmt.Sprintf("received animation %d", ctx.EffectiveMessage.MessageId))
		m, err := sendAnimation(
			b,
			ctx.EffectiveChat.Id,
			ctx.EffectiveMessage.Animation.FileId,
			&gotgbot.SendAnimationOpts{},
		)
		if err != nil {
			return h.logger.Error(
				fmt.Sprintf("failed to reply with animation, error: %v", err),
			)
		}
		err = h.sendWebAppMarkup(b, ctx.EffectiveChat.Id, []int64{m.MessageId})
		if err != nil {
			return h.logger.Error(
				fmt.Sprintf("failed to reply with animation, error: %v", err),
			)
		}
		h.logger.Info(fmt.Sprintf("animation message reply success %d", m.MessageId))
		return next(b, ctx)
	}
}

func (h handler) removeOneEffectiveMessage() handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		h.logger.Info(fmt.Sprintf("removing message %d", ctx.EffectiveMessage.MessageId))
		ok, err := deleteMessage(
			b,
			ctx.EffectiveMessage.GetSender().Id(),
			ctx.EffectiveMessage.MessageId,
		)
		if ok {
			h.logger.Info(
				fmt.Sprintf("message successfully removed %d", ctx.EffectiveMessage.MessageId),
			)
		}
		if err != nil {
			h.logger.Warning(fmt.Sprintf("failed to delete reply message, error: %v", err))
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
		h.logger.Info(
			fmt.Sprintf(
				"receiving group %s, current file %s",
				ctx.EffectiveMessage.MediaGroupId,
				mediaFileID,
			),
		)
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
			h.logger.Info(
				fmt.Sprintf(
					"processing group %s",
					ctx.EffectiveMessage.MediaGroupId,
				),
			)
			next(b, ctx)
		}()
		return nil
	}
}

func (h handler) respondWithMediaGroup(next handlers.Response) handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		h.logger.Info(
			fmt.Sprintf(
				"responding with media group %d",
				ctx.EffectiveMessage.MessageId,
			),
		)
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
		messages, err := sendMediaGroup(
			b,
			ctx.EffectiveChat.Id,
			group,
			&gotgbot.SendMediaGroupOpts{},
		)
		if err != nil {
			return h.logger.Error(
				fmt.Sprintf("failed to reply with animation, error: %v", err),
			)
		}
		messageIDs := []int64{}
		for _, message := range messages {
			messageIDs = append(messageIDs, message.MessageId)
		}
		h.logger.Info(
			fmt.Sprintf(
				"media group files send %d",
				ctx.EffectiveMessage.MessageId,
			),
		)
		err = h.sendWebAppMarkup(b, ctx.EffectiveChat.Id, messageIDs)
		if err != nil {
			return h.logger.Error(
				fmt.Sprintf(
					"failed to send web app markup  %d - %v",
					ctx.EffectiveMessage.MessageId,
					err,
				),
			)
		}
		return next(b, ctx)
	}
}

func (h handler) removeEffectiveMediaGroup() handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		h.logger.Info(
			fmt.Sprintf(
				"removing effective media group - %s",
				ctx.EffectiveMessage.MediaGroupId,
			),
		)
		toDelete := []int64{}
		for _, v := range h.mediaGroupMap.get(ctx.EffectiveMessage.MediaGroupId) {
			toDelete = append(toDelete, v.messageID)
		}
		_, err := deleteMessages(b, ctx.EffectiveChat.Id, toDelete)
		if err != nil {
			h.logger.Error(fmt.Sprintf("did not remove group media, reason: %+v", err))
		}
		h.mediaGroupMap.remove(ctx.EffectiveMessage.MediaGroupId)
		h.logger.Info(
			fmt.Sprintf(
				"successfully removed media group id - %s",
				ctx.EffectiveMessage.MediaGroupId,
			),
		)
		return nil
	}
}

func (h handler) sendWebAppMarkup(b bot, chatID int64, messageID []int64) error {
	h.logger.Info(
		fmt.Sprintf(
			"sending web app markup %d",
			messageID,
		),
	)
	str := []string{}
	for _, id := range messageID {
		str = append(str, fmt.Sprint(id))
	}
	m, err := sendMessage(b, chatID, "* * *", &gotgbot.SendMessageOpts{})
	if err != nil {
		return h.logger.Error(err.Error())
	}
	_, err = sendMessage(b, chatID, "* * *", &gotgbot.SendMessageOpts{
		ReplyMarkup: gotgbot.ReplyKeyboardMarkup{
			ResizeKeyboard: true,
			IsPersistent:   true,
			Keyboard: [][]gotgbot.KeyboardButton{
				{{
					Text: "#tag",
					WebApp: &gotgbot.WebAppInfo{Url: fmt.Sprintf(
						"%s/%s?message-id=%v&media-id=%v",
						h.config.WebAppUrl,
						h.config.Token,
						m.MessageId+1,
						strings.Join(str, ","),
					)},
				}},
			},
		},
	})
	if err != nil {
		return h.logger.Error(err.Error())
	}
	_, err = deleteMessage(b, chatID, m.MessageId)
	if err != nil && err != gotgbot.ErrNilBotClient {
		return h.logger.Error(err.Error())
	}
	h.logger.Info(
		fmt.Sprintf("webapp markup send - %d", messageID),
	)
	return nil
}

func (h handler) handleWebAppData(now func() time.Time) handlers.Response {
	type data struct {
		MediaIDs  string     `json:"mediaIds,required"`
		MessageID string     `json:"messageId,required"`
		Data      [][]string `json:"data,required"`
	}
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		h.logger.Info(
			fmt.Sprintf(
				"received webapp markup %d %+v",
				ctx.EffectiveMessage.MessageId,
				ctx.EffectiveMessage.WebAppData.Data,
			),
		)
		var d data
		err := json.Unmarshal([]byte(ctx.EffectiveMessage.WebAppData.Data), &d)
		if err != nil {
			return h.logger.Error(err.Error())
		}
		mediaIDs, err := utils.StringToIntSlice(d.MediaIDs)
		if err != nil {
			return h.logger.Error(err.Error())
		}
		tags := []string{}
		analytics := []models.Analytics{}
		for _, v := range d.Data {
			if len(v) != 2 {
				return h.logger.Error(
					fmt.Sprintf("Failed to parse data from web app, wrong format: %+v", v),
				)
			}
			group := v[0]
			tag := v[1]
			tags = append(tags, tag)
			analytics = append(
				analytics,
				models.Analytics{
					Group: group,
					Tag:   tag,
					Date:  now(),
				},
			)
		}
		_, _, err = editMessageCaption(b, &gotgbot.EditMessageCaptionOpts{
			ChatId:    ctx.EffectiveChat.Id,
			MessageId: mediaIDs[0],
			Caption:   strings.Join(tags, "\n"),
		})
		if err != nil &&
			!strings.Contains(err.Error(), "are exactly the same as a current content") {
			return h.logger.Error(err.Error())
		}
		_, err = copyMessages(b, h.config.ReceiverID, ctx.EffectiveChat.Id, mediaIDs, nil)
		if err != nil {
			return h.logger.Error(err.Error())
		}
		messageId, err := strconv.Atoi(d.MessageID)
		if err != nil {
			return h.logger.Error(err.Error())
		}
		_, err = deleteMessages(
			b,
			ctx.EffectiveChat.Id,
			[]int64{int64(messageId), ctx.EffectiveMessage.MessageId},
		)
		if err != nil {
			return h.logger.Error(err.Error())
		}
		c, cancel := context.WithTimeout(context.TODO(), time.Second*5)
		defer cancel()
		err = h.db.InsertAnalytics(c, &analytics)
		if err != nil {
			h.logger.Error(err.Error())
		}
		h.logger.Info(
			fmt.Sprintf(
				"webapp data processed succussfully %d",
				ctx.EffectiveMessage.MessageId,
			),
		)
		return nil
	}
}

func (h handler) handlePing() handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		h.logger.Info(fmt.Sprintf("received ping command %d", ctx.EffectiveMessage.MessageId))
		_, err := sendMessage(
			b,
			ctx.EffectiveChat.Id,
			fmt.Sprintf("pong (v%s)", h.config.Version),
			nil,
		)
		if err != nil {
			return h.logger.Error(err.Error())
		}
		h.logger.Info(fmt.Sprintf("successfully send pong %d", ctx.EffectiveMessage.MessageId))
		return nil
	}
}

var tagsRegexp = regexp.MustCompile(`(•.*:\n((#.*)(\n?))*)`)

func isTagsMessage(msg *gotgbot.Message) bool {
	return tagsRegexp.Match([]byte(msg.Text))
}

func (h handler) handleUpdateTags() handlers.Response {
	return func(b *gotgbot.Bot, ctx *ext.Context) error {
		h.logger.Info(
			fmt.Sprintf("received update tags request %d", ctx.EffectiveMessage.MessageId),
		)
		all := ctx.EffectiveMessage.Text[strings.IndexRune(ctx.EffectiveMessage.Text, '•'):]
		groupStrings := strings.Split(all, "\n\n")
		g := []models.Group{}
		ok := true
		for i, group := range groupStrings {
			data := strings.Split(group, "\n")
			r := regexp.MustCompile("• (.*):")
			matched := r.FindStringSubmatch(data[0])
			name := matched[1]
			tags := data[1:]
			if len(tags) == 0 || len(matched) != 2 {
				ok = false
				break
			}
			t := []models.Tag{}
			for _, v := range tags {
				t = append(t, models.Tag{Name: v})
			}
			g = append(g, models.Group{Name: name, Tags: t, OriginalIndex: i})
		}
		if !ok {
			sendMessage(b, ctx.EffectiveChat.Id, "error", nil)
			return h.logger.Error("failed to parse tags")
		}
		err := h.db.UpdateTags(context.Background(), &g)
		if err != nil {
			sendMessage(b, ctx.EffectiveChat.Id, "error", nil)
			return h.logger.Error(err.Error())
		}
		sendMessage(b, ctx.EffectiveChat.Id, "👍", nil)
		h.logger.Info(
			fmt.Sprintf("updated tags %d", ctx.EffectiveMessage.MessageId),
		)
		return err
	}
}

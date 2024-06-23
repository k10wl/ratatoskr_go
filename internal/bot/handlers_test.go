package bot

import (
	"fmt"
	"ratatoskr/internal/config"
	"ratatoskr/internal/logger"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func fakeLogger() *logger.Logger {
	return logger.NewLogger(
		"test logger",
		&strings.Builder{},
		&strings.Builder{},
	)
}

func TestReceiveGroupMedia(t *testing.T) {
	fakeHandler := newHandler(fakeLogger(), &config.BotConfig{Token: "TOKEN", WebAppUrl: webAppUrl})

	calls := 0
	res := fakeHandler.receiveGroup(
		time.Millisecond*500,
		func(b *gotgbot.Bot, ctx *ext.Context) error {
			calls++
			return nil
		},
	)
	res(&gotgbot.Bot{}, &ext.Context{
		EffectiveMessage: &gotgbot.Message{
			MessageId:    1,
			MediaGroupId: "1",
			Photo:        []gotgbot.PhotoSize{{FileId: "1"}},
		},
	})
	res(&gotgbot.Bot{}, &ext.Context{
		EffectiveMessage: &gotgbot.Message{
			MessageId:    2,
			MediaGroupId: "1",
			Photo:        []gotgbot.PhotoSize{{FileId: "2"}},
		},
	})
	res(&gotgbot.Bot{}, &ext.Context{
		EffectiveMessage: &gotgbot.Message{
			MessageId:    3,
			MediaGroupId: "1",
			Photo:        []gotgbot.PhotoSize{{FileId: "3"}},
		},
	})
	res(&gotgbot.Bot{}, &ext.Context{
		EffectiveMessage: &gotgbot.Message{
			MessageId:    4,
			MediaGroupId: "1",
			Photo:        []gotgbot.PhotoSize{{FileId: "4"}},
		},
	})

	time.Sleep(time.Second)

	expected := map[string][]item{
		"1": {
			{messageID: 1, mediaType: "photo", fileID: "1"},
			{messageID: 2, mediaType: "photo", fileID: "2"},
			{messageID: 3, mediaType: "photo", fileID: "3"},
			{messageID: 4, mediaType: "photo", fileID: "4"},
		},
	}
	if !reflect.DeepEqual(fakeHandler.mediaGroupMap.hashMap, expected) {
		t.Errorf(
			"Failed to receive media group\nexpected: %+v\nactual:   %+v",
			expected,
			fakeHandler.mediaGroupMap.hashMap,
		)
	}
	if calls != 1 {
		t.Errorf(
			"Wrong calls amount\nexpected: %+v\nactual:   %+v",
			1,
			calls,
		)
	}

}

func TestRemoveOneEffectiveMessage(t *testing.T) {
	type arg struct {
		messageId int64
		chatId    int64
	}
	var removed arg
	calls := 0
	fakeHandler := newHandler(fakeLogger(), &config.BotConfig{WebAppUrl: webAppUrl})
	original := deleteMessage
	defer func() {
		deleteMessage = original
	}()
	deleteMessage = func(b bot, chatId int64, messageId int64) (bool, error) {
		calls++
		removed = arg{
			chatId:    chatId,
			messageId: messageId,
		}
		return true, nil
	}

	err := fakeHandler.removeOneEffectiveMessage()(&gotgbot.Bot{}, &ext.Context{
		EffectiveMessage: &gotgbot.Message{
			MessageId:  1,
			SenderChat: &gotgbot.Chat{Id: 1},
		},
	})

	if err != nil {
		t.Errorf("Unexpected error in removeOriginal")
	}
	if calls != 1 {
		t.Errorf("Wrong amount of deletion calls (%d)", calls)
	}
	if !reflect.DeepEqual(removed, arg{messageId: 1, chatId: 1}) {
		t.Errorf("Did not remove correct original (%+v)", removed)
	}
}

func TestSendPhoto(t *testing.T) {
	type photoArg struct {
		fileID gotgbot.InputFile
		chatID int64
	}
	var sendPhotoArg photoArg
	var sendWebAppUrl string
	nextCalled := false
	createdMessageID := 0
	fakeHandler := newHandler(fakeLogger(), &config.BotConfig{WebAppUrl: webAppUrl})
	originalSendPhoto := sendPhoto
	originalSendMessage := sendMessage
	originalEditMessageReplyMarkup := editMessageReplyMarkup
	defer func() {
		sendPhoto = originalSendPhoto
		sendMessage = originalSendMessage
		editMessageReplyMarkup = originalEditMessageReplyMarkup
	}()
	sendPhoto = func(
		b bot,
		chatId int64,
		fileID gotgbot.InputFile,
		opts *gotgbot.SendPhotoOpts,
	) (*gotgbot.Message, error) {
		createdMessageID++
		sendPhotoArg = photoArg{
			chatID: chatId,
			fileID: fileID,
		}
		return &gotgbot.Message{
			MessageId: int64(createdMessageID),
		}, nil
	}
	sendMessage = func(
		b bot,
		chatId int64,
		message string,
		opts *gotgbot.SendMessageOpts,
	) (*gotgbot.Message, error) {
		createdMessageID++
		if opts.ReplyMarkup != nil {
			sendWebAppUrl = opts.ReplyMarkup.(gotgbot.ReplyKeyboardMarkup).Keyboard[0][0].WebApp.Url
		}
		return &gotgbot.Message{MessageId: int64(createdMessageID)}, nil
	}

	mockNext := func(b *gotgbot.Bot, ctx *ext.Context) error {
		nextCalled = true
		return nil
	}

	err := fakeHandler.handlePhoto(mockNext)(
		&gotgbot.Bot{},
		&ext.Context{
			EffectiveMessage: &gotgbot.Message{
				MessageId: 1,
				Photo:     []gotgbot.PhotoSize{{FileId: "unique file id"}},
			},
			EffectiveChat: &gotgbot.Chat{Id: 1},
		},
	)

	if err != nil {
		t.Errorf("Unexpected error in handlePhoto")
	}
	if !nextCalled {
		t.Errorf("Next was not called after handlePhoto")
	}
	expectedWebAppUrl := fmt.Sprintf(
		"%s/%s?message-id=3&media-id=1",
		webAppUrl,
		fakeHandler.config.Token,
	)
	if sendWebAppUrl != expectedWebAppUrl {
		t.Errorf(
			"Did not send correct webApp message-id query params\nexpected: %v\nactual:   %v",
			expectedWebAppUrl,
			sendWebAppUrl,
		)
	}
	if !reflect.DeepEqual(sendPhotoArg, photoArg{fileID: "unique file id", chatID: 1}) {
		t.Errorf("Did not send correct photo (%+v)", sendPhotoArg)
	}
}

func TestSendVideo(t *testing.T) {
	type arg struct {
		fileID gotgbot.InputFile
		chatID int64
	}
	var send arg
	calls := 0
	createdMessageID := 0
	var sendWebAppUrl string
	nextCalled := false
	fakeHandler := newHandler(fakeLogger(), &config.BotConfig{WebAppUrl: webAppUrl})
	originalSendVideo := sendVideo
	originalEditMessageReplyMarkup := editMessageReplyMarkup
	defer func() {
		sendVideo = originalSendVideo
		editMessageReplyMarkup = originalEditMessageReplyMarkup
	}()
	sendVideo = func(
		b bot,
		chatId int64,
		fileID gotgbot.InputFile,
		opts *gotgbot.SendVideoOpts,
	) (*gotgbot.Message, error) {
		calls++
		createdMessageID++
		send = arg{
			chatID: chatId,
			fileID: fileID,
		}
		return &gotgbot.Message{
			MessageId: int64(createdMessageID),
		}, nil
	}
	sendMessage = func(
		b bot,
		chatId int64,
		message string,
		opts *gotgbot.SendMessageOpts,
	) (*gotgbot.Message, error) {
		createdMessageID++
		if opts.ReplyMarkup != nil {
			sendWebAppUrl = opts.ReplyMarkup.(gotgbot.ReplyKeyboardMarkup).Keyboard[0][0].WebApp.Url
		}
		return &gotgbot.Message{MessageId: int64(createdMessageID)}, nil
	}

	mockNext := func(b *gotgbot.Bot, ctx *ext.Context) error {
		nextCalled = true
		return nil
	}

	err := fakeHandler.handleVideo(mockNext)(
		&gotgbot.Bot{},
		&ext.Context{
			EffectiveChat: &gotgbot.Chat{Id: 1},
			EffectiveMessage: &gotgbot.Message{
				MessageId: 1,
				Video:     &gotgbot.Video{FileId: "unique file id"},
			},
		},
	)

	if err != nil {
		t.Errorf("Unexpected error in handleVideo")
	}
	if !nextCalled {
		t.Errorf("Next was not called after handleVideo")
	}
	expectedWebAppUrl := fmt.Sprintf(
		"%s/%s?message-id=3&media-id=1",
		webAppUrl,
		fakeHandler.config.Token,
	)
	if sendWebAppUrl != expectedWebAppUrl {
		t.Errorf(
			"Did not send correct webApp message-id query params\nexpected: %v\nactual:   %v",
			expectedWebAppUrl,
			sendWebAppUrl,
		)
	}
	if !reflect.DeepEqual(send, arg{fileID: "unique file id", chatID: 1}) {
		t.Errorf("Did not send correct video (%+v)", send)
	}
}

func TestSendAnimation(t *testing.T) {
	type arg struct {
		fileID gotgbot.InputFile
		chatID int64
	}
	var send arg
	calls := 0
	createdMessageID := 0
	var sendWebAppUrl string
	nextCalled := false
	fakeHandler := newHandler(fakeLogger(), &config.BotConfig{WebAppUrl: webAppUrl})
	originalSendAnimation := sendAnimation
	originalEditMessageReplyMarkup := editMessageReplyMarkup
	defer func() {
		sendAnimation = originalSendAnimation
		editMessageReplyMarkup = originalEditMessageReplyMarkup
	}()
	sendAnimation = func(
		b bot,
		chatId int64,
		fileID gotgbot.InputFile,
		opts *gotgbot.SendAnimationOpts,
	) (*gotgbot.Message, error) {
		calls++
		createdMessageID++
		send = arg{
			chatID: chatId,
			fileID: fileID,
		}
		return &gotgbot.Message{MessageId: int64(createdMessageID)}, nil
	}

	sendMessage = func(
		b bot,
		chatId int64,
		message string,
		opts *gotgbot.SendMessageOpts,
	) (*gotgbot.Message, error) {
		createdMessageID++
		if opts.ReplyMarkup != nil {
			sendWebAppUrl = opts.ReplyMarkup.(gotgbot.ReplyKeyboardMarkup).Keyboard[0][0].WebApp.Url
		}
		return &gotgbot.Message{MessageId: int64(createdMessageID)}, nil
	}
	editMessageReplyMarkup = func(b bot, opts *gotgbot.EditMessageReplyMarkupOpts) (*gotgbot.Message, bool, error) {
		sendWebAppUrl = opts.ReplyMarkup.InlineKeyboard[0][0].WebApp.Url
		return nil, true, nil
	}

	mockNext := func(b *gotgbot.Bot, ctx *ext.Context) error {
		nextCalled = true
		return nil
	}

	err := fakeHandler.handleAnimation(mockNext)(
		&gotgbot.Bot{},
		&ext.Context{
			EffectiveChat: &gotgbot.Chat{Id: 1},
			EffectiveMessage: &gotgbot.Message{
				MessageId: 1,
				Animation: &gotgbot.Animation{FileId: "unique file id"},
			},
		},
	)

	if err != nil {
		t.Errorf("Unexpected error in handleAnimation")
	}
	if !nextCalled {
		t.Errorf("Next was not called after handleAnimation")
	}
	expectedWebAppUrl := fmt.Sprintf("%s/?message-id=3&media-id=1", webAppUrl)
	if sendWebAppUrl != expectedWebAppUrl {
		t.Errorf(
			"Did not send correct webApp message-id query params\nexpected: %v\nactual:   %v",
			expectedWebAppUrl,
			sendWebAppUrl,
		)
	}
	if !reflect.DeepEqual(send, arg{fileID: "unique file id", chatID: 1}) {
		t.Errorf("Did not send correct animation (%+v)", send)
	}
}

func TestRespondWithMediaGroup(t *testing.T) {
	type arg struct {
		inputMedia []gotgbot.InputMedia
		chatID     int64
	}
	var send arg
	var sendWebAppUrl string
	sendMessageID := 0
	fakeHandler := newHandler(fakeLogger(), &config.BotConfig{Token: "TOKEN", WebAppUrl: webAppUrl})
	fakeHandler.mediaGroupMap = newMediaGroupMap()
	fakeHandler.mediaGroupMap.hashMap = map[string][]item{
		"1": {
			{
				mediaType: "photo",
				messageID: 1,
				fileID:    "file 1",
			},
			{
				mediaType: "photo",
				messageID: 2,
				fileID:    "file 2",
			},
			{
				mediaType: "video",
				messageID: 3,
				fileID:    "file 3",
			},
		},
	}
	originalSendMediaGroup := sendMediaGroup
	originalSendMessage := sendMessage
	originalEditMessageReplyMarkup := editMessageReplyMarkup
	nextCalled := false
	defer func() {
		sendMediaGroup = originalSendMediaGroup
		sendMessage = originalSendMessage
		editMessageReplyMarkup = originalEditMessageReplyMarkup
	}()
	sendMediaGroup = func(
		b bot,
		chatId int64,
		inputMedia []gotgbot.InputMedia,
		opts *gotgbot.SendMediaGroupOpts,
	) ([]gotgbot.Message, error) {
		send = arg{
			chatID:     chatId,
			inputMedia: inputMedia,
		}
		return []gotgbot.Message{
			{MessageId: 1},
			{MessageId: 2},
			{MessageId: 3},
		}, nil
	}
	sendMessage = func(
		b bot,
		chatId int64,
		message string,
		opts *gotgbot.SendMessageOpts,
	) (*gotgbot.Message, error) {
		sendMessageID++
		if opts.ReplyMarkup != nil {
			sendWebAppUrl = opts.ReplyMarkup.(gotgbot.ReplyKeyboardMarkup).Keyboard[0][0].WebApp.Url
		}

		return &gotgbot.Message{MessageId: int64(sendMessageID)}, nil
	}

	err := fakeHandler.respondWithMediaGroup(func(b *gotgbot.Bot, ctx *ext.Context) error {
		nextCalled = true
		return nil
	})(
		&gotgbot.Bot{},
		&ext.Context{
			EffectiveChat: &gotgbot.Chat{
				Id: 1,
			},
			EffectiveMessage: &gotgbot.Message{
				MediaGroupId: "1",
			},
		},
	)

	if err != nil {
		t.Errorf("Unexpected error in handleMediaGroup")
	}
	expected := arg{inputMedia: []gotgbot.InputMedia{
		gotgbot.InputMediaPhoto{Media: "file 1"},
		gotgbot.InputMediaPhoto{Media: "file 2"},
		gotgbot.InputMediaVideo{Media: "file 3"},
	}, chatID: 1}
	if !reflect.DeepEqual(
		send,
		expected,
	) {
		t.Errorf("Did not send correct media group:\nexpected: %+v\nactual:   %+v", expected, send)
	}
	expectedWebAppUrl := fmt.Sprintf(
		"%s/%s?message-id=2&media-id=1,2,3",
		webAppUrl,
		fakeHandler.config.Token,
	)
	if sendWebAppUrl != expectedWebAppUrl {
		t.Errorf(
			"Did not send correct webApp message-id query params\nexpected: %v\nactual:   %v",
			expectedWebAppUrl,
			sendWebAppUrl,
		)
	}
	if !nextCalled {
		t.Error("did not call next after clearing messages")
	}
}

var webAppUrl = "https://webapp.url"

func TestRemoveEffectiveMediaGroup(t *testing.T) {
	type arg struct {
		messageId []int64
		chatId    int64
	}
	removed := arg{}
	calls := 0
	fakeHandler := newHandler(fakeLogger(), &config.BotConfig{Token: "TOKEN", WebAppUrl: webAppUrl})
	original := deleteMessages
	defer func() {
		deleteMessages = original
	}()
	deleteMessages = func(b bot, chatId int64, messageId []int64) (bool, error) {
		calls++
		removed = arg{
			chatId:    chatId,
			messageId: messageId,
		}
		return true, nil
	}
	fakeHandler.mediaGroupMap.hashMap = map[string][]item{
		"1": {
			{
				messageID: 1,
				mediaType: "photo",
				fileID:    "2",
			},
			{
				messageID: 2,
				mediaType: "video",
				fileID:    "2",
			},
		},
	}

	err := fakeHandler.removeEffectiveMediaGroup()(&gotgbot.Bot{}, &ext.Context{
		EffectiveChat: &gotgbot.Chat{
			Id: 1,
		},
		EffectiveMessage: &gotgbot.Message{
			MessageId:    1,
			MediaGroupId: "1",
			SenderChat:   &gotgbot.Chat{Id: 1},
		},
	})

	if err != nil {
		t.Errorf("Unexpected error in removeOriginal")
	}
	if len(fakeHandler.mediaGroupMap.hashMap) != 0 {
		t.Errorf("Did not remove related media group from hash map")
	}
	if !reflect.DeepEqual(removed, arg{
		chatId:    1,
		messageId: []int64{1, 2},
	}) {
		t.Errorf("Did not remove correct media group, removed: %+v", removed)
	}
}

func TestSendWebAppMarkup(t *testing.T) {
	type tc struct {
		name     string
		expected string
		run      func() string
	}

	originalSendMessage := sendMessage

	defer func() {
		sendMessage = originalSendMessage
	}()

	table := []tc{
		{
			name: "should create webapp with url for one photo",
			run: func() string {
				var url string
				sendMessageCalls := 0
				sendMessage = func(b bot, chatId int64, message string, opts *gotgbot.SendMessageOpts) (*gotgbot.Message, error) {
					sendMessageCalls++
					if opts.ReplyMarkup != nil {
						url = opts.ReplyMarkup.(gotgbot.ReplyKeyboardMarkup).Keyboard[0][0].WebApp.Url
					}
					return &gotgbot.Message{MessageId: int64(sendMessageCalls)}, nil
				}
				fakeHandler := newHandler(
					fakeLogger(),
					&config.BotConfig{Token: "TOKEN", WebAppUrl: webAppUrl},
				)
				fakeHandler.sendWebAppMarkup(&gotgbot.Bot{}, int64(sendMessageCalls), []int64{1234})
				return url
			},
			expected: webAppUrl + "/TOKEN?message-id=2&media-id=1234",
		},

		{
			name: "should create webapp with url for one animation",
			run: func() string {
				var url string
				sendMessageCalls := 0
				sendMessage = func(b bot, chatId int64, message string, opts *gotgbot.SendMessageOpts) (*gotgbot.Message, error) {
					sendMessageCalls++
					if opts.ReplyMarkup != nil {
						url = opts.ReplyMarkup.(gotgbot.ReplyKeyboardMarkup).Keyboard[0][0].WebApp.Url
					}
					return &gotgbot.Message{MessageId: int64(sendMessageCalls)}, nil
				}
				fakeHandler := newHandler(
					fakeLogger(),
					&config.BotConfig{Token: "TOKEN", WebAppUrl: webAppUrl},
				)
				fakeHandler.sendWebAppMarkup(&gotgbot.Bot{}, 1, []int64{1234})
				return url
			},
			expected: webAppUrl + "/TOKEN?message-id=2&media-id=1234",
		},

		{
			name: "should create webapp with url for one video",
			run: func() string {
				var url string
				sendMessageCalls := 0
				sendMessage = func(b bot, chatId int64, message string, opts *gotgbot.SendMessageOpts) (*gotgbot.Message, error) {
					sendMessageCalls++
					if opts.ReplyMarkup != nil {
						url = opts.ReplyMarkup.(gotgbot.ReplyKeyboardMarkup).Keyboard[0][0].WebApp.Url
					}
					return &gotgbot.Message{MessageId: int64(sendMessageCalls)}, nil
				}
				fakeHandler := newHandler(
					fakeLogger(),
					&config.BotConfig{Token: "TOKEN", WebAppUrl: webAppUrl},
				)
				fakeHandler.sendWebAppMarkup(&gotgbot.Bot{}, 1, []int64{1234})
				return url
			},
			expected: webAppUrl + "/TOKEN?message-id=2&media-id=1234",
		},

		{
			name: "should create webapp with url for group",
			run: func() string {
				var url string

				sendMessage = func(b bot, chatId int64, message string, opts *gotgbot.SendMessageOpts) (*gotgbot.Message, error) {
					if opts.ReplyMarkup != nil {
						url = opts.ReplyMarkup.(gotgbot.ReplyKeyboardMarkup).Keyboard[0][0].WebApp.Url
					}
					return &gotgbot.Message{MessageId: 1}, nil
				}

				fakeHandler := newHandler(
					fakeLogger(),
					&config.BotConfig{Token: "TOKEN", WebAppUrl: webAppUrl},
				)
				fakeHandler.sendWebAppMarkup(&gotgbot.Bot{}, 1, []int64{1234, 1235, 1236})
				return url
			},
			expected: webAppUrl + "/TOKEN?message-id=2&media-id=1234,1235,1236",
		},
	}

	for _, test := range table {
		actual := test.run()
		if test.expected != actual {
			t.Errorf(
				"%s - unexpected webapp url\nexpected: %s\nactual:   %s",
				test.name,
				test.expected,
				actual,
			)
		}
	}
}

func TestHandleWebAppData(t *testing.T) {
	type editCaptionsResult struct {
		chatID    int64
		messageID int64
		caption   string
	}
	type copyMessagesResult struct {
		from       int64
		receiver   int64
		messageIDs []int64
	}
	type deleteMessagesResult struct {
		chatID     int64
		massageIDs []int64
	}
	type result struct {
		editCaption    editCaptionsResult
		copyMessages   copyMessagesResult
		deleteMessages deleteMessagesResult
	}
	type tc struct {
		name     string
		ctx      *ext.Context
		call     func(ctx *ext.Context) result
		expected result
	}

	originalEditMessageCaptionOpts := editMessageCaption
	originalDeleteMessages := deleteMessages
	originalCopyMessages := copyMessages
	defer func() {
		editMessageCaption = originalEditMessageCaptionOpts
		deleteMessages = originalDeleteMessages
		copyMessages = originalCopyMessages
	}()
	var res result
	deleteMessages = func(b bot, chatId int64, messageIds []int64) (bool, error) {
		res.deleteMessages.chatID = chatId
		res.deleteMessages.massageIDs = messageIds
		return true, nil
	}
	copyMessages = func(b bot, chatId, fromChatId int64, messageIds []int64, opts *gotgbot.CopyMessagesOpts) ([]gotgbot.MessageId, error) {
		res.copyMessages.from = fromChatId
		res.copyMessages.receiver = chatId
		res.copyMessages.messageIDs = messageIds
		return []gotgbot.MessageId{}, nil
	}
	editMessageCaption = func(b bot, opts *gotgbot.EditMessageCaptionOpts) (*gotgbot.Message, bool, error) {
		res.editCaption.chatID = opts.ChatId
		res.editCaption.messageID = opts.MessageId
		res.editCaption.caption = opts.Caption
		return nil, true, nil
	}
	fh := newHandler(fakeLogger(), &config.BotConfig{
		ReceiverID: 7890,
	})

	table := []tc{
		{
			name: "edit one caption and copy one message",
			ctx: &ext.Context{
				EffectiveChat: &gotgbot.Chat{
					Id: 1,
				},
				EffectiveMessage: &gotgbot.Message{
					MessageId: 3,
					WebAppData: &gotgbot.WebAppData{
						Data: `{"tags": ["#tag1", "#tag2", "#tag3"], "mediaIds": "1", "messageId": "2"}`,
					}},
			},
			call: func(ctx *ext.Context) result {
				res = result{}
				fh.handleWebAppData()(&gotgbot.Bot{}, ctx)
				return res
			},
			expected: result{
				editCaption: editCaptionsResult{
					chatID:    1,
					messageID: 1,
					caption:   "#tag1\n#tag2\n#tag3",
				},
				copyMessages: copyMessagesResult{
					from:       1,
					receiver:   7890,
					messageIDs: []int64{1},
				},
				deleteMessages: deleteMessagesResult{
					chatID:     1,
					massageIDs: []int64{2, 3},
				},
			},
		},

		{
			name: "edit one caption and copy multiple messages",
			ctx: &ext.Context{
				EffectiveChat: &gotgbot.Chat{
					Id: 1,
				},
				EffectiveMessage: &gotgbot.Message{
					MessageId: 3,
					WebAppData: &gotgbot.WebAppData{
						Data: `{"tags": ["#tag1", "#tag2", "#tag3"], "mediaIds": "1,2,3", "messageId": "2"}`,
					}},
			},
			call: func(ctx *ext.Context) result {
				res = result{}
				fh.handleWebAppData()(&gotgbot.Bot{}, ctx)
				return res
			},
			expected: result{
				editCaption: editCaptionsResult{
					chatID:    1,
					messageID: 1,
					caption:   "#tag1\n#tag2\n#tag3",
				},
				copyMessages: copyMessagesResult{
					from:       1,
					receiver:   7890,
					messageIDs: []int64{1, 2, 3},
				},
				deleteMessages: deleteMessagesResult{
					chatID:     1,
					massageIDs: []int64{2, 3},
				},
			},
		},
	}

	for _, test := range table {
		actual := test.call(test.ctx)
		if !reflect.DeepEqual(test.expected, actual) {
			t.Errorf(
				"Unexpected result - %s\nexpected: %+v\nactual:   %+v",
				test.name,
				test.expected,
				actual,
			)
		}
	}
}

func TestPingHandler(t *testing.T) {
	fakeHandler := newHandler(fakeLogger(), &config.BotConfig{Token: "TOKEN", WebAppUrl: webAppUrl})
	originalSendMessage := sendMessage
	defer func() {
		sendMessage = originalSendMessage
	}()
	messageText := ""
	sendInChat := -1
	sendMessage = func(b bot, chatId int64, message string, opts *gotgbot.SendMessageOpts) (*gotgbot.Message, error) {
		messageText = message
		sendInChat = int(chatId)
		return nil, nil
	}
	fakeHandler.handlePing()(&gotgbot.Bot{}, &ext.Context{

		EffectiveChat: &gotgbot.Chat{Id: 1},
	})
	if messageText != "pong" {
		t.Errorf("did not respond with `ping`. Actual: %s", messageText)
	}
	if sendInChat != 1 {
		t.Errorf("did not respond in same chat. Actual: %d", sendInChat)
	}
}

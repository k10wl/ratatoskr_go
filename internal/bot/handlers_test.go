package bot

import (
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
	fakeHandler := newHandler(fakeLogger())

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

	expected := map[string][]string{"1": {"1", "2", "3", "4"}}
	if !reflect.DeepEqual(fakeHandler.mediaGroupMap, expected) {
		t.Errorf(
			"Failed to receive media group\nexpected: %+v\nactual:   %+v",
			expected,
			fakeHandler.mediaGroupMap,
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
	fakeHandler := newHandler(fakeLogger())
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
	type arg struct {
		fileID gotgbot.InputFile
		chatID int64
	}
	var send arg
	calls := 0
	nextCalled := false
	fakeHandler := newHandler(fakeLogger())
	original := sendPhoto
	defer func() {
		sendPhoto = original
	}()
	sendPhoto = func(
		b bot,
		chatId int64,
		fileID gotgbot.InputFile,
		opts *gotgbot.SendPhotoOpts,
	) (*gotgbot.Message, error) {
		calls++
		send = arg{
			chatID: chatId,
			fileID: fileID,
		}
		return &gotgbot.Message{}, nil
	}

	mockNext := func(b *gotgbot.Bot, ctx *ext.Context) error {
		nextCalled = true
		return nil
	}

	err := fakeHandler.handlePhoto(mockNext)(
		&gotgbot.Bot{},
		&ext.Context{
			EffectiveMessage: &gotgbot.Message{
				MessageId:  1,
				SenderChat: &gotgbot.Chat{Id: 1},
				Photo:      []gotgbot.PhotoSize{{FileId: "unique file id"}},
			},
		},
	)

	if err != nil {
		t.Errorf("Unexpected error in handlePhoto")
	}
	if !nextCalled {
		t.Errorf("Next was not called after handlePhoto")
	}
	if calls != 1 {
		t.Errorf("Wrong amount of send image calls (%d)", calls)
	}
	if !reflect.DeepEqual(send, arg{fileID: "unique file id", chatID: 1}) {
		t.Errorf("Did not send correct photo (%+v)", send)
	}
}

func TestSendVideo(t *testing.T) {
	type arg struct {
		fileID gotgbot.InputFile
		chatID int64
	}
	var send arg
	calls := 0
	nextCalled := false
	fakeHandler := newHandler(fakeLogger())
	original := sendVideo
	defer func() {
		sendVideo = original
	}()
	sendVideo = func(
		b bot,
		chatId int64,
		fileID gotgbot.InputFile,
		opts *gotgbot.SendVideoOpts,
	) (*gotgbot.Message, error) {
		calls++
		send = arg{
			chatID: chatId,
			fileID: fileID,
		}
		return &gotgbot.Message{}, nil
	}

	mockNext := func(b *gotgbot.Bot, ctx *ext.Context) error {
		nextCalled = true
		return nil
	}

	err := fakeHandler.handleVideo(mockNext)(
		&gotgbot.Bot{},
		&ext.Context{
			EffectiveMessage: &gotgbot.Message{
				MessageId:  1,
				SenderChat: &gotgbot.Chat{Id: 1},
				Video:      &gotgbot.Video{FileId: "unique file id"},
			},
		},
	)

	if err != nil {
		t.Errorf("Unexpected error in handleVideo")
	}
	if !nextCalled {
		t.Errorf("Next was not called after handleVideo")
	}
	if calls != 1 {
		t.Errorf("Wrong amount of send video calls (%d)", calls)
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
	nextCalled := false
	fakeHandler := newHandler(fakeLogger())
	original := sendAnimation
	defer func() {
		sendAnimation = original
	}()
	sendAnimation = func(
		b bot,
		chatId int64,
		fileID gotgbot.InputFile,
		opts *gotgbot.SendAnimationOpts,
	) (*gotgbot.Message, error) {
		calls++
		send = arg{
			chatID: chatId,
			fileID: fileID,
		}
		return &gotgbot.Message{}, nil
	}

	mockNext := func(b *gotgbot.Bot, ctx *ext.Context) error {
		nextCalled = true
		return nil
	}

	err := fakeHandler.handleAnimation(mockNext)(
		&gotgbot.Bot{},
		&ext.Context{
			EffectiveMessage: &gotgbot.Message{
				MessageId:  1,
				SenderChat: &gotgbot.Chat{Id: 1},
				Animation:  &gotgbot.Animation{FileId: "unique file id"},
			},
		},
	)

	if err != nil {
		t.Errorf("Unexpected error in handleAnimation")
	}
	if !nextCalled {
		t.Errorf("Next was not called after handleAnimation")
	}
	if calls != 1 {
		t.Errorf("Wrong amount of send animation calls (%d)", calls)
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
	calls := 0
	fakeHandler := newHandler(fakeLogger())
	fakeHandler.mediaGroupMap = map[string][]string{
		"1": {"file 1", "file 2", "file 3"},
	}
	original := sendMediaGroup
	defer func() {
		sendMediaGroup = original
	}()
	sendMediaGroup = func(
		b bot,
		chatId int64,
		inputMedia []gotgbot.InputMedia,
		opts *gotgbot.SendMediaGroupOpts,
	) ([]gotgbot.Message, error) {
		calls++
		send = arg{
			chatID:     chatId,
			inputMedia: inputMedia,
		}
		return []gotgbot.Message{}, nil
	}

	err := fakeHandler.respondWithMediaGroup()(
		&gotgbot.Bot{},
		&ext.Context{
			EffectiveSender: &gotgbot.Sender{
				ChatId: 1,
			},
		},
	)

	if err != nil {
		t.Errorf("Unexpected error in handleMediaGroup")
	}
	if calls != 1 {
		t.Errorf("Wrong amount of send media group calls (%d)", calls)
	}
	expected := arg{inputMedia: []gotgbot.InputMedia{
		gotgbot.InputMediaPhoto{Media: "file 1"},
		gotgbot.InputMediaPhoto{Media: "file 2"},
		gotgbot.InputMediaPhoto{Media: "file 3"},
	}, chatID: 1}
	if !reflect.DeepEqual(
		send,
		expected,
	) {
		t.Errorf("Did not send correct media group:\nexpected: %+v\nactual:   %+v", expected, send)
	}
}

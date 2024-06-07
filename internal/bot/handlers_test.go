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
	original := DeleteMessage
	defer func() {
		DeleteMessage = original
	}()
	DeleteMessage = func(chatId, messageId int64, opts *gotgbot.DeleteMessageOpts) (bool, error) {
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
		t.Errorf("Wrong amount of deletion calls")
	}
	if !reflect.DeepEqual(removed, arg{messageId: 1, chatId: 1}) {
		t.Errorf("Did not remove correct original")
	}
}

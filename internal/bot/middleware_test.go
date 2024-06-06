package bot

import (
	"ratatoskr/internal/config"
	"ratatoskr/internal/logger"
	"strings"
	"testing"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func TestAdminOnly(t *testing.T) {
	called := false
	fakeHandler := func(b *gotgbot.Bot, ctx *ext.Context) error {
		called = true
		return nil
	}

	middleware := newMidlleware(
		logger.NewLogger("test", &strings.Builder{}, &strings.Builder{}),
		&config.Config{AdminIDs: []int64{1234}},
	)

	middleware.adminOnly(fakeHandler)(
		nil,
		&ext.Context{EffectiveSender: &gotgbot.Sender{User: &gotgbot.User{Id: 9876}}},
	)

	if called {
		t.Errorf("middleware did not block unauthorized request")
		called = false
	}

	middleware.adminOnly(fakeHandler)(
		nil,
		&ext.Context{EffectiveSender: &gotgbot.Sender{User: &gotgbot.User{Id: 1234}}},
	)

	if !called {
		t.Errorf("middleware blocked authorized request")
		called = false
	}

	middleware.adminOnly(fakeHandler)(
		nil,
		&ext.Context{EffectiveSender: &gotgbot.Sender{User: &gotgbot.User{Id: 1234}}},
	)

	if !called {
		t.Errorf("middleware did not allow authorized request")
	}
}

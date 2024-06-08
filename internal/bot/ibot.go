package bot

import "github.com/PaulSonOfLars/gotgbot/v2"

type bot interface {
	DeleteMessage(int64, int64, *gotgbot.DeleteMessageOpts) (bool, error)
	SendPhoto(int64, gotgbot.InputFile, *gotgbot.SendPhotoOpts) (*gotgbot.Message, error)
}

var (
	deleteMessage = botDeleteMessage
	sendPhoto     = botSendPhoto
)

func botSendPhoto(
	b bot,
	chatId int64,
	fileID gotgbot.InputFile,
	opts *gotgbot.SendPhotoOpts,
) (*gotgbot.Message, error) {
	return b.SendPhoto(
		chatId,
		fileID,
		opts,
	)
}

func botDeleteMessage(b bot, chatId int64, messageId int64) (bool, error) {
	return b.DeleteMessage(
		chatId, messageId, &gotgbot.DeleteMessageOpts{},
	)
}

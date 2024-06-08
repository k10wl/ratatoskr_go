package bot

import "github.com/PaulSonOfLars/gotgbot/v2"

type bot interface {
	DeleteMessage(int64, int64, *gotgbot.DeleteMessageOpts) (bool, error)
	SendPhoto(int64, gotgbot.InputFile, *gotgbot.SendPhotoOpts) (*gotgbot.Message, error)
	SendVideo(int64, gotgbot.InputFile, *gotgbot.SendVideoOpts) (*gotgbot.Message, error)
	SendAnimation(int64, gotgbot.InputFile, *gotgbot.SendAnimationOpts) (*gotgbot.Message, error)
	SendMediaGroup(
		int64,
		[]gotgbot.InputMedia,
		*gotgbot.SendMediaGroupOpts,
	) ([]gotgbot.Message, error)
}

var (
	deleteMessage  = botDeleteMessage
	sendPhoto      = botSendPhoto
	sendVideo      = botSendVideo
	sendAnimation  = botSendAnimation
	sendMediaGroup = botSendMediaGroup
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

func botSendVideo(
	b bot,
	chatId int64,
	fileID gotgbot.InputFile,
	opts *gotgbot.SendVideoOpts,
) (*gotgbot.Message, error) {
	return b.SendVideo(
		chatId,
		fileID,
		opts,
	)
}

func botSendAnimation(
	b bot,
	chatId int64,
	fileID gotgbot.InputFile,
	opts *gotgbot.SendAnimationOpts,
) (*gotgbot.Message, error) {
	return b.SendAnimation(
		chatId,
		fileID,
		opts,
	)
}

func botSendMediaGroup(
	b bot,
	chatId int64,
	inputMedia []gotgbot.InputMedia,
	opts *gotgbot.SendMediaGroupOpts,
) ([]gotgbot.Message, error) {
	return b.SendMediaGroup(
		chatId,
		inputMedia,
		opts,
	)
}

func botDeleteMessage(b bot, chatId int64, messageId int64) (bool, error) {
	return b.DeleteMessage(
		chatId, messageId, &gotgbot.DeleteMessageOpts{},
	)
}

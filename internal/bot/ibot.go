package bot

import "github.com/PaulSonOfLars/gotgbot/v2"

type bot interface {
	DeleteMessage(int64, int64, *gotgbot.DeleteMessageOpts) (bool, error)
	DeleteMessages(int64, []int64, *gotgbot.DeleteMessagesOpts) (bool, error)
	SendPhoto(int64, gotgbot.InputFile, *gotgbot.SendPhotoOpts) (*gotgbot.Message, error)
	SendVideo(int64, gotgbot.InputFile, *gotgbot.SendVideoOpts) (*gotgbot.Message, error)
	SendAnimation(int64, gotgbot.InputFile, *gotgbot.SendAnimationOpts) (*gotgbot.Message, error)
	SendMediaGroup(
		int64,
		[]gotgbot.InputMedia,
		*gotgbot.SendMediaGroupOpts,
	) ([]gotgbot.Message, error)
	SendMessage(int64, string, *gotgbot.SendMessageOpts) (*gotgbot.Message, error)
	EditMessageReplyMarkup(*gotgbot.EditMessageReplyMarkupOpts) (*gotgbot.Message, bool, error)
}

var (
	deleteMessage          = botDeleteMessage
	deleteMessages         = botDeleteMessages
	sendPhoto              = botSendPhoto
	sendVideo              = botSendVideo
	sendAnimation          = botSendAnimation
	sendMediaGroup         = botSendMediaGroup
	sendMessage            = botSendMessage
	editMessageReplyMarkup = botEditMessageReplyMarkup
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

func botDeleteMessages(b bot, chatId int64, messageIds []int64) (bool, error) {
	return b.DeleteMessages(
		chatId, messageIds, &gotgbot.DeleteMessagesOpts{},
	)
}

func botSendMessage(
	b bot,
	chatId int64,
	message string,
	opts *gotgbot.SendMessageOpts,
) (*gotgbot.Message, error) {
	return b.SendMessage(chatId, message, opts)
}

func botEditMessageReplyMarkup(
	b bot,
	opts *gotgbot.EditMessageReplyMarkupOpts,
) (*gotgbot.Message, bool, error) {
	return b.EditMessageReplyMarkup(opts)
}

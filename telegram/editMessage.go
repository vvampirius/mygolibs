package telegram

type EditMessageIntInlineKeyboardMarkup struct {
	ChatId          int                  `json:"chat_id"`
	MessageId       int                  `json:"message_id"`
	InlineMessageId string               `json:"inline_message_id"`
	ReplyMarkup     InlineKeyboardMarkup `json:"reply_markup"`
}

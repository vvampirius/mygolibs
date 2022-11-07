package telegram

type SendMessageText struct {
	Text string `json:"text"`
}

type SendMessageInt struct {
	SendMessageText
	ChatId int `json:"chat_id"`
}

type SendMessageIntWithoutReplyMarkup struct {
	SendMessageInt
	ParseMode string `json:"parse_mode"`
	//Entities []Entities //TODO: check this
	DisableWebPagePreview bool `json:"disable_web_page_preview"`
	DisableNotification bool `json:"disable_notification"`
	ProtectContent bool `json:"protect_content"`
	ReplyToMessageId int `json:"reply_to_message_id"`
	AllowSendingWithoutReply bool `json:"allow_sending_without_reply"`
}

type SendMessageIntWithInlineKeyboardMarkup struct {
	SendMessageIntWithoutReplyMarkup
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
}

type SendMessageIntWithReplyKeyboardMarkup struct {
	SendMessageIntWithoutReplyMarkup
	ReplyMarkup ReplyKeyboardMarkup `json:"reply_markup"`
}


type ForceReply struct {
	ForceReply bool `json:"force_reply"`
	InputFieldPlaceholder string `json:"input_field_placeholder"`
	Selective bool `json:"selective"`
}

type SendMessageIntWithForceReply struct {
	SendMessageIntWithoutReplyMarkup
	ReplyMarkup ForceReply `json:"reply_markup"`
}

package telegram

type InlineKeyboardButton struct {
	Text string `json:"text"`
	Url string `json:"url"`
	CallbackData string `json:"callback_data"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

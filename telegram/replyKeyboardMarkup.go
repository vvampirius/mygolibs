package telegram

type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton		`json:"keyboard"`
	ResizeKeyboard bool				`json:"resize_keyboard"`
	OneTimeKeyboard bool			`json:"one_time_keyboard"`
	InputFieldPlaceholder string	`json:"input_field_placeholder"`
	Selective bool					`json:"selective"`
}

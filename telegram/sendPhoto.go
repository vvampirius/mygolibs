package telegram

type SendPhotoUrl struct {
	ChatId int		`json:"chat_id"`
	Photo string	`json:"photo"`
	Caption string 	`json:"caption"`
}

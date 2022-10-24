package telegram

type KeyboardButton struct {
	Text string 			`json:"text"`
	RequestContact bool 	`json:"request_contact"`
	RequestLocation bool	`json:"request_location"`
}

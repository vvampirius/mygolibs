package telegram

type RequestResponse struct {
	Ok bool
	Result struct {
		Chat Chat
		Date int
		From User
		MessageId int `json:"message_id"`
		Text string
	}
	ErrorCode int `json:"error_code"`
	Description string
}

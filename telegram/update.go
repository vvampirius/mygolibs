package telegram

type Update struct {
	Id      int     `json:"update_id"`
	Message Message `json:"message"`
}

type Message struct {
	Id   int    `json:"message_id"`
	Date int    `json:"date"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	Id                          int    `json:"id"`
	Type                        string `json:"type"`
	Title                       string `json:"title"`
	Username                    string `json:"username"`
	FirstName                   string `json:"first_name"`
	LastName                    string `json:"last_name"`
	AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
}

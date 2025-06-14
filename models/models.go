package models

type Message struct {
	Message_id     string `json:"id"`
	Content        string `json:"content"`
	Sent_timestamp int64  `json:"timestamp"`
	Status         string `json:"status"`
}

type ModeratorResponse struct {
	Mod_id        string
	Message_id    string
	Status        string
	Response_time int
}

type APIResponse struct {
	Message_id     string `json:"id"`
	Sent_timestamp int64  `json:"timestamp"`
	Status         string `json:"status"`
}

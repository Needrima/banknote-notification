package models

type Notification struct {
	Reference string `json:"reference"`
	Channel   string `json:"channel" binding:"required"`
	Type      string `json:"type" binding:"required,eq=SCHEDULED|eq=INSTANT"`
	Subject   string `json:"subject" binding:"required"`
	From      string `json:"from" binding:"required"`
	To        string `json:"to" binding:"required"`
	Message   string `json:"message" binding:"required"`
	SendAt    string `json:"send_at" binding:"required"`
	SentAt    string `json:"sent_at"`
	Status    string `json:"status"`
}

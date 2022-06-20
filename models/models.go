package models

type NotificationObject struct {
	Reference string `json:"reference" binding:"required"`
	Channel   string `json:"channel" binding:"required"`
	Type      string `json:"type" binding:"required,eq=SCHEDULED|INSTANT"`
	Subject   string `json:"subject" binding:"required"`
	From      string `json:"from" binding:"required"`
	To        string `json:"to" binding:"required"`
	Message   string `json:"message" binding:"required"`
	SendAt    string `json:"send_at" binding:"required"`
	SentAt    string `json:"sent_at" binding:"required"`
	Status    string `json:"status" binding:"required,eq=PENDING|SENT"`
}

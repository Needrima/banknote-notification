package models

type Notification struct {
	Reference string `json:"reference" bson:"reference"`
	Channel   string `json:"channel" bson:"channel" binding:"required,eq=EMAIL"`
	Type      string `json:"type" bson:"type" binding:"required,eq=SCHEDULED|eq=INSTANT"`
	Subject   string `json:"subject" bson:"subject" binding:"required"`
	From      string `json:"from" bson:"from" binding:"required"`
	To        string `json:"to" bson:"to" binding:"required"`
	Message   string `json:"message" bson:"message" binding:"required"`
	SendAt    string `json:"send_at" bson:"send_at" binding:"required"`
	SentAt    string `json:"sent_at" bson:"sent_at"`
	Status    string `json:"status" bson:"status"`
}

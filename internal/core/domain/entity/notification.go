package entity

type Notification struct {
	Reference string `json:"reference" bson:"reference"`
	Channel   string `json:"channel" bson:"channel" binding:"required,eq=EMAIL"`
	Type      string `json:"type" bson:"type" binding:"required,eq=SCHEDULED|eq=INSTANT"`
	Subject   string `json:"subject" bson:"subject" binding:"required"`
	To        string `json:"to" bson:"to" binding:"required"`
	From        string `json:"from" bson:"from" binding:"required"`
	Message   string `json:"message" bson:"message" binding:"required"`
	SendAt    string `json:"send_at" bson:"send_at" binding:"required"` // comes in RFC3339 format E.G 2022-11-02T23:47:00
	SentAt    string `json:"sent_at" bson:"sent_at"`                    // comes in RFC3339 format E.G 2022-11-02T23:47:00
	Status    string `json:"status" bson:"status"`
}

package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subscriber struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Email        string             `json:"email" bson:"email" binding:"email"`
	SubscribedAt string             `json:"subscribed_at,omitempty" bson:"subscribed_at"` // comes in RFC3339 format E.G 2022-11-02T23:47:00
}

type ContactMessage struct {
	Name    string `json:"name" binding:"required"`
	Email   string `json:"email" binding:"email"` // email from contact message
	Message string `json:"message" binding:"required"`
	To      string `json:"to,omitempty" bson:"to"`           // Receiver's email address
	From    string `json:"from,omitempty" bson:"from"`       // Sender's email address
	SentAt  string `json:"sent_at,omitempty" bson:"sent_at"` // comes in RFC3339 format E.G 2022-11-02T23:47:00
}

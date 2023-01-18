package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ProductID   primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Name        string             `json:"name,omitempty" bson:"name" binding:"required"`
	ImageNames  []string           `json:"img_names,omitempty" bson:"img_names" binding:"required"`
	Price       float64            `json:"price,omitempty" bson:"price" binding:"required"`
	Description string             `json:"desc,omitempty" bson:"desc" binding:"required"`
	Sizes       []string           `json:"sizes,omitempty" bson:"sizes" binding:"required"` //S,M,L,XL,XXL
	Rating      int                `json:"rating,omitempty" bson:"rating" binding:"required"`
	Brand       string             `json:"brand,omitempty" bson:"brand" binding:"required"` // default is DONA
	Colours     []string           `json:"colours,omitempty" bson:"colours" binding:"required"`
	Category    string             `json:"category,omitempty" bson:"category" binding:"required"`
	CreatedAt   string             `json:"created_at,omitempty" bson:"created_at"` // comes in RFC3339 format E.G 2022-11-02T23:47:00
}

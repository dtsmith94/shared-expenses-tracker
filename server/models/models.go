package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ExpenseList struct {
	ID     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"task,omitempty"`
	Status bool               `json:"status,omitempty"`
}

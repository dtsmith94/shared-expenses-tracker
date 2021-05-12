package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	Database struct {
		ConnectionString string `yaml:"connectionString"`
	}
}

type Expense struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name"`
	Amount   float64            `json:"amount"`
	Settled  bool               `json:"settled"`
	Created  time.Time          `json:"created"`
	Modified time.Time          `json:"modified"`
}

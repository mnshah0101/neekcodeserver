package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Solution struct {
	Username    string             `json:"username"`
	Email       string             `json:"email"`
	File_url    string             `json:"file_url"`
	Likes       int                `json:"likes"`
	Question_id string             `json:"question_id"`
	Date        string             `json:"date"`
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
}

package user

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id primitive.ObjectID `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Username string `json:"username" bson:"username"`
	Email string `json:"email"`
	Handle string `json:"handle" bson:"handle"`
	DOB int64 `json:"dob" bson:"dob"`
	Password string `json:"password" bson:"password"`
	CreatedAt int64 `json:"created_at" bson:"created_at"`
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
	DeletedAt int64 `json:"deleted_at" bson:"deleted_at, omitempty"`
}
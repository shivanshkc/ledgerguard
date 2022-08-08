package models

// UserDoc is the schema of the user's document as stored in the database.
type UserDoc struct {
	UserID    string `json:"user_id" bson:"user_id"`
	Email     string `json:"email" bson:"email"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}

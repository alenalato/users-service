package storage

import "time"

type UserDetails struct {
	ID           string    `bson:"_id,omitempty"`
	FirstName    string    `bson:"first_name,omitempty"`
	LastName     string    `bson:"last_name,omitempty"`
	Nickname     string    `bson:"nickname,omitempty"`
	Email        string    `bson:"email,omitempty"`
	PasswordHash string    `bson:"password,omitempty"`
	Country      string    `bson:"country,omitempty"`
	CreatedAt    time.Time `bson:"created_at,omitempty"`
	UpdatedAt    time.Time `bson:"updated_at,omitempty"`
}

type User struct {
	ID        string    `bson:"_id,omitempty"`
	FirstName string    `bson:"first_name,omitempty"`
	LastName  string    `bson:"last_name,omitempty"`
	Nickname  string    `bson:"nickname,omitempty"`
	Email     string    `bson:"email,omitempty"`
	Country   string    `bson:"country,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
}

type UserUpdate struct {
	FirstName *string    `bson:"first_name,omitempty"`
	LastName  *string    `bson:"last_name,omitempty"`
	Country   *string    `bson:"country,omitempty"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
}

type UserFilter struct {
	FirstName *string `json:"first_name" bson:"first_name,omitempty"`
	LastName  *string `json:"last_name" bson:"last_name,omitempty"`
	Country   *string `json:"country" bson:"country,omitempty"`
}

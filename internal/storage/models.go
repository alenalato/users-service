package storage

import "time"

// UserDetails represents the input details of a user to be created
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

// User represents a user output from the storage
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

// UserUpdate represents the input details of a user to be updated
// If a field is nil, it will not be updated
type UserUpdate struct {
	FirstName *string    `bson:"first_name,omitempty"`
	LastName  *string    `bson:"last_name,omitempty"`
	Nickname  *string    `bson:"nickname,omitempty"`
	Email     *string    `bson:"email,omitempty"`
	Country   *string    `bson:"country,omitempty"`
	UpdatedAt *time.Time `bson:"updated_at,omitempty"`
}

// UserFilter represents the input filter criteria for listing users
// If a field is nil, it will not be used in the filter
type UserFilter struct {
	FirstName *string `json:"first_name" bson:"first_name,omitempty"`
	LastName  *string `json:"last_name" bson:"last_name,omitempty"`
	Country   *string `json:"country" bson:"country,omitempty"`
}

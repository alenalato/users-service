package businesslogic

import "time"

// PasswordDetails represents the input details of a password
// Text is required as input
// Hash is generated and then handled to be write-only
type PasswordDetails struct {
	Text string `validate:"required"`
	Hash string
}

// UserDetails represents the input details of a user to be created
// Nickname, email and password are required fields
type UserDetails struct {
	FirstName string
	LastName  string
	Nickname  string          `validate:"required"`
	Email     string          `validate:"required,email"`
	Password  PasswordDetails `validate:"required"`
	Country   string
}

// User represents a user
type User struct {
	ID        string
	FirstName string
	LastName  string
	Nickname  string
	Email     string
	Country   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserUpdate represents the input details of a user to be updated
// UpdateMask is a required list of fields that will be considered for the update
type UserUpdate struct {
	FirstName  string
	LastName   string
	Nickname   string
	Email      string
	Country    string
	UpdateMask []string `validate:"required"`
}

// UserFilter represents the input filter criteria for listing users
// If a field is nil, it will not be used in the filter
type UserFilter struct {
	FirstName *string
	LastName  *string
	Nickname  *string
	Email     *string
	Country   *string
}

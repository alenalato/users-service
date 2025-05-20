package businesslogic

import "time"

type PasswordDetails struct {
	Text string `validate:"required,min=8"`
	Hash string
}

type UserDetails struct {
	FirstName string
	LastName  string
	Nickname  string          `validate:"required"`
	Email     string          `validate:"required,email"`
	Password  PasswordDetails `validate:"required"`
	Country   string
}

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

type UserUpdate struct {
	FirstName  string
	LastName   string
	Country    string
	UpdateMask []string
}

type UserFilter struct {
	FirstName *string
	LastName  *string
	Country   *string
}

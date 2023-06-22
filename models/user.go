package models

// User defines the structure for an API users
type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

// Users  is a collection of User
type Users []*User

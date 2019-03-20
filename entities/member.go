package entities

import "time"

// Member represents a member of a test
type Member struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createAt"`
}

// NewMember creates a new member
func NewMember(name, email, role string) (Member, error) {
	return Member{
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: now(),
	}, nil
}

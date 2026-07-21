package dto

type UserResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number,omitempty"`
	RoleID      string `json:"role_id"`
}

type UpdateProfileRequest struct {
	Name        string `json:"name" validate:"omitempty,min=2,max=255"`
	PhoneNumber string `json:"phone_number" validate:"omitempty,max=20"`
}

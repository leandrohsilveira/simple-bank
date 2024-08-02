package user_domain

import database "github.com/leandrohsilveira/simple-bank/server/database"

type UserDTO struct {
	Id    int64             `json:"id"`
	Email string            `json:"email"`
	Role  database.UserRole `json:"role"`
	Name  string            `json:"name"`
}

func FromUser(user database.User) UserDTO {
	return UserDTO{
		Id:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		Name:  user.Name,
	}
}

func FromUsers(users []database.User) []UserDTO {
	result := make([]UserDTO, 0)

	for _, user := range users {
		result = append(result, FromUser(user))
	}

	return result
}

func IsUserDTOEmpty(user UserDTO) bool {
	return user == UserDTO{}
}

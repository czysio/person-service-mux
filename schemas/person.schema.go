package schemas

type CreatePerson struct {
	FirstName string `json:"first_name" binding:"required"`
	Surname   string `json:"surname" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Nickname  string `json:"nickname" binding:"required"`
}

type UpdatePerson struct {
	FirstName string `json:"first_name"`
	Surname   string `json:"surname"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
}

package user

type createUserRequest struct {
	Name string `json:"name"`
}

type userResponse struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
}

package domain

type UserData struct {
	Id       string `json: "id"`
	User     string `json: "user"`
	Password string `json: "password"`
	Admin    bool   `json: "admin"`
}

type LoginData struct {
	Token  string `json: "token"`
	IdU    string `json: "idu"`
	AdminU bool   `json:"adminu"`
}

type UsersData []UserData

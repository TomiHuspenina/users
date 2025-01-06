package users

type User struct {
	Id       string `bson:"_id,omitempty"`
	User     string `bson:"user"`
	Password string `bson:"password"`
	Admin    bool   `bson:"admin"`
}

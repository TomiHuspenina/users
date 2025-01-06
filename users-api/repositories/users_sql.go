package users

import (
	"fmt"
	dao "users-api/dao"
	"users-api/domain"
	e "users-api/errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type SQLConfig struct {
	Name string
	User string
	Pass string
	Host string
}

type SQL struct {
	db       *gorm.DB
	Database string
}

// DeleteReserva implements reservas.Repository.

func NewSql(config SQLConfig) SQL {
	db, err := gorm.Open("mysql", config.User+":"+config.Pass+"@tcp("+config.Host+":3306)/"+config.Name+"?charset=utf8&parseTime=True")
	if err != nil {
		log.Println("Connection Failed to Open")
		log.Fatal(err)
	} else {
		log.Println("Connection Established gg")
	}
	db.AutoMigrate(&dao.User{})
	return SQL{
		db:       db,
		Database: config.Name,
	}

}

func (repository SQL) GetUserById(id int64) (dao.User, error) {
	var buscado dao.User
	log.Println("ID: ", id)
	result := repository.db.Where("ID = ?", id).First(&buscado)
	log.Println("resultado: ", result)
	if result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return dao.User{}, fmt.Errorf("no user found with ID: %d", id)
		}
		return dao.User{}, fmt.Errorf("error finding document: %v", result.Error)
	}
	return buscado, nil
}

func (repository SQL) InsertUser(user dao.User) (dao.User, error) {
	result := repository.db.Create(&user)
	if result.Error != nil {
		log.Panic("Error creating the user")
		return user, fmt.Errorf("error inserting document:")
	}
	return user, nil
}

func (repository SQL) GetUserByName(user dao.User) (dao.User, error) {
	var usuario dao.User

	result := repository.db.Where("User = ?", user.User).First(&usuario)
	log.Debug("User: ", usuario)
	if result.Error != nil {
		log.Error("Error al buscar el usuario")
		log.Error(result.Error)
		return user, e.NewBadRequestApiError("Error al buscar el usuario")
	}

	return usuario, nil
}

func (repository SQL) Login(user dao.User) (domain.LoginData, error) {
	var usuario dao.User

	result := repository.db.Where("User = ?", user.User).First(&usuario)
	log.Debug("User: ", usuario)
	if result.Error != nil {
		log.Error("Error al buscar el usuario")
		log.Error(result.Error)
		return domain.LoginData{}, e.NewBadRequestApiError("Error al buscar el usuario")
	}

	return domain.LoginData{}, nil
}

package users

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
	dao "users-api/dao"
	domain "users-api/domain"
	e "users-api/errors"

	"github.com/golang-jwt/jwt"
)

type Repository interface {
	GetUserById(id string) (dao.User, error)
	InsertUser(usuario dao.User) (dao.User, error)
	Login(User dao.User) (domain.LoginData, error)
	GetUserByName(User dao.User) (dao.User, error)
}

type Service struct {
	mainRepo Repository
}

func NewService(mainRepo Repository) Service {
	return Service{
		mainRepo: mainRepo,
	}
}

func (service Service) GetUserById(ctx context.Context, id string) (domain.UserData, error) {
	userDAO, err := service.mainRepo.GetUserById(id)
	if err != nil {
		return domain.UserData{}, fmt.Errorf("error getting user from main repository: %v", err)
	}

	return domain.UserData{
		Id:       userDAO.Id,
		User:     userDAO.User,
		Password: userDAO.Password,
		Admin:    userDAO.Admin,
	}, nil
}

func (service Service) InsertUser(ctx context.Context, usuarioDomain domain.UserData) (domain.UserData, error) {
	userDao := dao.User{
		User:     usuarioDomain.User,
		Password: usuarioDomain.Password,
		Admin:    usuarioDomain.Admin,
	}
	hash := md5.New()
	hash.Write([]byte(usuarioDomain.Password))
	userDao.Password = hex.EncodeToString(hash.Sum(nil))

	_, err := service.mainRepo.InsertUser(userDao)
	if err != nil {
		return domain.UserData{}, fmt.Errorf("error inserting user into main repository: %v", err)
	}

	return usuarioDomain, nil
}

func (service Service) Login(ctx context.Context, User domain.UserData) (domain.LoginData, error) {
	userDAO := dao.User{
		User:     User.User,
		Password: User.Password,
	}

	var tokenDomain domain.LoginData
	user, err := service.mainRepo.GetUserByName(userDAO)
	if err != nil {
		return tokenDomain, e.NewBadRequestApiError("usuario no encontrado")
	}

	var Logpsw = md5.Sum([]byte(User.Password))
	psw := hex.EncodeToString(Logpsw[:])

	if psw == user.Password {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"idU":    user.Id,
			"Adminu": user.Admin,
			"exp":    time.Now().Add(time.Hour * 72).Unix(),
		})
		t, _ := token.SignedString([]byte("frantomi"))
		tokenDomain.Token = t
		tokenDomain.IdU = user.Id
		tokenDomain.AdminU = user.Admin
		return tokenDomain, nil
	} else {
		return tokenDomain, e.NewBadRequestApiError("Contrasenia incorrecta")
	}
}

func (service Service) GetUserByName(ctx context.Context, user domain.UserData) (domain.UserData, error) {
	userDAO := dao.User{
		User:     user.User,
		Password: user.Password,
	}

	userDAO, err := service.mainRepo.GetUserByName(userDAO)
	if err != nil {
		return domain.UserData{}, fmt.Errorf("error getting user from main repository: %v", err)
	}

	return domain.UserData{
		Id:       userDAO.Id,
		User:     userDAO.User,
		Password: userDAO.Password,
		Admin:    userDAO.Admin,
	}, nil
}

/*package users

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
	dao "users-api/dao"
	domain "users-api/domain"
	e "users-api/errors"

	"github.com/golang-jwt/jwt"
)

type Repository interface {
	GetUserById(id int64) (dao.User, error)
	InsertUser(usuario dao.User) (dao.User, error)
	Login(User dao.User) (domain.LoginData, error)
	GetUserByName(User dao.User) (dao.User, error)
}
type Service struct {
	mainRepo        Repository
	//cacheRepository Repository
}

func NewService(mainRepo Repository, cacheRepository Repository) Service {
	return Service{
		mainRepo:        mainRepo,
		//cacheRepository: cacheRepository,
	}
}

func (service Service) GetUserById(ctx context.Context, id int64) (domain.UserData, error) {
	userDAO, err := service.cacheRepository.GetUserById(id)
	if err != nil {
		userDAO, err = service.mainRepo.GetUserById(id)
		if err != nil {
			return domain.UserData{}, fmt.Errorf("error getting hotel from main repository: %v", err)
		}
	}

	return domain.UserData{
		Id:       userDAO.Id,
		User:     userDAO.User,
		Password: userDAO.Password,
		Admin:    userDAO.Admin,
	}, nil

}

/*
func (service Service) InsertUser(ctx context.Context, usuarioDomain domain.UserData) (domain.UserData, error) {

	var usuario dao.User
	var result, er = service.mainRepo.GetUserByName(usuario)

	if er != nil {

		usuario.User = usuarioDomain.User
		usuario.Password = usuarioDomain.Password
		usuario.Admin = usuarioDomain.Admin

		hash := md5.New()
		hash.Write([]byte(usuarioDomain.Password))
		usuario.Password = hex.EncodeToString(hash.Sum(nil))

		var usuario2, err = service.mainRepo.InsertUser(usuario)

		if err != nil {
			return usuarioDomain, e.NewBadRequestApiError("usuario no inseertado")
		}

		usuarioDomain.Id = usuario2.Id
		usuarioDomain.Admin = usuario.Admin

		return usuarioDomain, nil
	}

	return domain.UserData(result), e.NewBadRequestApiError("Nombre de usuario existente")

}*/

/*

func (service Service) InsertUser(ctx context.Context, usuarioDomain domain.UserData) (domain.UserData, error) {
	userDao := dao.User{
		User:     usuarioDomain.User,
		Password: usuarioDomain.Password,
		Admin:    usuarioDomain.Admin,
	}
	hash := md5.New()
	hash.Write([]byte(usuarioDomain.Password))
	userDao.Password = hex.EncodeToString(hash.Sum(nil))
	_, err := service.cacheRepository.InsertUser(userDao)
	if err != nil {
		return domain.UserData{}, fmt.Errorf("Error inserting user into cache repository: %v", err)
	}
	_, err = service.mainRepo.InsertUser(userDao)
	if err != nil {
		return domain.UserData{}, fmt.Errorf("Error inserting user into main repository: %v", err)
	}

	return usuarioDomain, nil
}

func (service Service) Login(ctx context.Context, User domain.UserData) (domain.LoginData, error) {

	userDAO := dao.User{
		User:     User.User,
		Password: User.Password,
	}
	var tokenDomain domain.LoginData
	user, err := service.cacheRepository.GetUserByName(userDAO)
	if err != nil {
		user, err = service.mainRepo.GetUserByName(userDAO)
		if err != nil {
			return tokenDomain, e.NewBadRequestApiError("usuario no encontrado")
		}

		if _, err := service.cacheRepository.InsertUser(user); err != nil {
			return tokenDomain, fmt.Errorf("error inserting user into cache: %v", err)
		}
	}

	var Logpsw = md5.Sum([]byte(User.Password))
	psw := hex.EncodeToString(Logpsw[:])

	if psw == user.Password {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"idU":    user.Id,
			"Adminu": user.Admin,
			"exp":    time.Now().Add(time.Hour * 72).Unix(),
		})
		t, _ := token.SignedString([]byte("frantomi"))
		tokenDomain.Token = t
		tokenDomain.IdU = user.Id
		tokenDomain.AdminU = user.Admin
		return tokenDomain, nil
	} else {
		return tokenDomain, e.NewBadRequestApiError("Contrasenia incorrecta")
	}
}

func (service Service) GetUserByName(ctx context.Context, user domain.UserData) (domain.UserData, error) {
	userDAO := dao.User{
		User:     user.User,
		Password: user.Password,
	}
	userDAO, err := service.cacheRepository.GetUserByName(userDAO)
	if err != nil {
		userDAO, err = service.mainRepo.GetUserByName(userDAO)
		if err != nil {
			return domain.UserData{}, fmt.Errorf("error getting user from main repository: %v", err)
		}
	}

	return domain.UserData{
		Id:       userDAO.Id,
		User:     userDAO.User,
		Password: userDAO.Password,
		Admin:    userDAO.Admin,
	}, nil
}
*/

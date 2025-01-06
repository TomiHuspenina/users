package users

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	Domain "users-api/domain"
	middle "users-api/middleware"

	"github.com/gin-gonic/gin"
)

type Service interface {
	InsertUser(ctx context.Context, user Domain.UserData) (Domain.UserData, error)
	GetUserById(ctx context.Context, id string) (Domain.UserData, error)
	GetUserByName(ctx context.Context, user Domain.UserData) (Domain.UserData, error)
	Login(ctx context.Context, user Domain.UserData) (Domain.LoginData, error)
}

type Controller struct {
	service Service
}

func NewController(service Service) Controller {
	return Controller{
		service: service,
	}
}

func (controller Controller) Login(c *gin.Context) {
	var userData Domain.UserData
	c.BindJSON(&userData)

	loginResponse, err := controller.service.Login(c.Request.Context(), userData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, loginResponse)
}
func (controller Controller) Extrac(c *gin.Context) {
	data := strings.TrimSpace(c.GetHeader("Authorization"))
	log.Println("token buscado: ", data)
	response, err := middle.ExtractClaims(data)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (controller Controller) GetUserByName(c *gin.Context) {

	var userDomain Domain.UserData
	c.BindJSON(&userDomain)

	userDomain, err := controller.service.GetUserByName(c.Request.Context(), userDomain)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, userDomain)

}

func (controller Controller) InsertUser(c *gin.Context) {
	var userDomain Domain.UserData
	err := c.BindJSON(&userDomain)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userDomain.Admin {
		log.Println("creating admin user")
	} else {
		log.Println("creating regular user")
	}

	userDomain, er := controller.service.InsertUser(c.Request.Context(), userDomain)

	if er != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": er.Error()})
		return
	}

	c.JSON(http.StatusCreated, userDomain)

}

func (controller Controller) GetUserById(ctx *gin.Context) {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id cannot be empty",
		})
		return
	}

	user, err := controller.service.GetUserById(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error getting user: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

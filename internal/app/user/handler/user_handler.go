package handler

import (
	"exercise/internal/app/domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		db: db,
	}
}

func (uh UserHandler) Register(c *gin.Context) {
	var userRegister domain.UserRegister
	if err := c.ShouldBind(&userRegister); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid body",
		})
	}

	// setelah line ini adalah usecase
	user, err := domain.NewUser(userRegister.Email, userRegister.Name, userRegister.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	if err := uh.db.Create(user).Error; err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}
	token, err := user.GenerateJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func (uh UserHandler) Login(c *gin.Context) {
	var userRequest domain.UserLogin
	if err := c.ShouldBind(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
		return
	}

	var user domain.User
	err := uh.db.Where("email = ?", userRequest.Email).Take(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid username or password",
		})
		return
	}

	if correctPassword := user.CorrectPassword(userRequest.Password); !correctPassword {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "invalid username or password",
		})
		return
	}

	token, err := user.GenerateJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "error when generate token",
		})
		return
	}
	c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

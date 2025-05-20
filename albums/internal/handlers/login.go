package handlers

import (
	"albums/internal/models"
	"albums/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	LoginService service.LoginService
}

// Login godoc
//
//	@Summary		Login and get a JWT Token
//	@Description	Authenticates a user and returns a signed JWT token.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body	models.LoginPayload	true	"Login Credentials"
//	@Success		200	{object}	map[string]interface{}	"Login successful"
//	@Failure		400	{object}	map[string]string		"Bad request"
//	@Failure		404	{object}	map[string]string		"User not found"
//	@Failure		500	{object}	map[string]string		"Internal server error"
//	@Router			/login [post]
func (handler *LoginHandler) Login(c *gin.Context) {
	userID, exists := c.Get("UserID")
	if !exists || len(userID.(string)) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	token, err := handler.LoginService.JWTTokenGenerator(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not able to generate token : ", "messsage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Loggedin", "token": token, "ID": userID})
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user account and returns user details.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.UserSwagger	true	"New user details"
//	@Success		200		{object}	map[string]interface{}	"User created successfully"
//	@Failure		400		{object}	map[string]string		"Invalid JSON"
//	@Failure		500		{object}	map[string]string		"Failed to create user"
//	@Router			/register [post]
func (handler *LoginHandler) Register(c *gin.Context) {
	var newUser models.Users
	if err := c.ShouldBindBodyWithJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	user, err := handler.LoginService.RegisterUser(newUser.Username, newUser.Password)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "user created", "ID": user.ID, "username": user.Username, "password": user.Password})
}

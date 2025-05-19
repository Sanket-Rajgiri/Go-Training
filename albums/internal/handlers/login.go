package handlers

import (
	"albums/internal/models"
	"albums/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	userID, exists := c.Get("UserID")
	if !exists || len(userID.(string)) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	token, err := service.JWTTokenGenerator(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not able to generate token : ", "messsage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Loggedin", "token": token, "ID": userID})
}

func RegisterHandler(c *gin.Context) {
	var newUser models.Users
	if err := c.ShouldBindBodyWithJSON(&newUser); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	user, err := service.RegisterUser(newUser.Username, newUser.Password)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "user created", "ID": user.ID, "username": user.Username, "password": user.Password})
}

package handlers

import (
	"albums/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	token, err := service.JWTTokenGenerator(c.Request.URL.User.Username())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not able to generate token : ", "messsage": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Loggedin", "token": token})
}

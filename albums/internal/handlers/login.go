package handlers

import (
	"albums/internal/models"
	"albums/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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
	ctx := c.Request.Context()
	tracer := otel.Tracer("Login-Tracer")
	ctx, span := tracer.Start(ctx, "Login-Handler")
	defer span.End()
	userID, exists := c.Get("UserID")
	if !exists || len(userID.(string)) == 0 {
		span.SetStatus(codes.Error, "user not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	span.SetAttributes(attribute.String("userId", userID.(string)))
	token, err := handler.LoginService.JWTTokenGenerator(ctx, userID.(string))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not able to generate token : ", "messsage": err.Error()})
		return
	}
	span.SetStatus(codes.Ok, "token generated")
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
	ctx := c.Request.Context()
	tracer := otel.Tracer("Register-Tracer")
	ctx, span := tracer.Start(ctx, "Register-Handler")
	defer span.End()
	var newUser models.Users
	if err := c.ShouldBindBodyWithJSON(&newUser); err != nil {
		span.SetStatus(codes.Error, "Invalid JSON")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	user, err := handler.LoginService.RegisterUser(ctx, newUser.Username, newUser.Password)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	span.SetAttributes(attribute.Int64("userID", int64(user.ID)), attribute.String("username", user.Username))
	span.SetStatus(codes.Ok, "user created")
	c.IndentedJSON(http.StatusOK, gin.H{"message": "user created", "ID": user.ID, "username": user.Username, "password": user.Password})
}

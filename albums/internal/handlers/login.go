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
	ctx, span := tracer.Start(ctx, "LoginHandler")
	defer span.End()
	userID, exists := c.Get("UserID")
	if !exists || len(userID.(string)) == 0 {
		span.SetStatus(codes.Error, "user not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	role, exists := c.Get("Role")
	if !exists || len(role.(string)) == 0 {
		span.SetStatus(codes.Error, "role not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}
	span.SetAttributes(attribute.String("userId", userID.(string)))

	acces_token, refresh_token, err := handler.LoginService.Login(ctx, userID.(string), role.(string))
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	span.SetStatus(codes.Ok, "token generated")
	c.JSON(http.StatusOK, gin.H{"message": "Loggedin", "access_token": acces_token, "refresh_token": refresh_token, "ID": userID})
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
//	@Router			/login/register [post]
func (handler *LoginHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("Register-Tracer")
	ctx, span := tracer.Start(ctx, "RegisterHandler")
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

// Refresh godoc
//
//	@Summary		Refresh Access Token
//	@Description	Creates and returns new accessToken.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			tokens	body		RefreshHandlerBody	true	"username and refresh token"
//	@Success		200		{object}	map[string]interface{}	"token refreshed"
//	@Failure		400		{object}	map[string]string		"Invalid JSON"
//	@Failure		500		{object}	map[string]string		"Failed to refresh token"
//	@Router			/login/refresh [post]

type RefreshHandlerBody struct {
	Username string `json:"username" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

func (handler *LoginHandler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("Handler-Tracer")
	ctx, span := tracer.Start(ctx, "RefreshHandler")
	defer span.End()
	var requestBody RefreshHandlerBody
	if err := c.ShouldBindBodyWithJSON(&requestBody); err != nil {
		span.SetStatus(codes.Error, "Invalid JSON")
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	newToken, err := handler.LoginService.RefreshToken(ctx, requestBody.Username, requestBody.Token)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	span.SetStatus(codes.Ok, "token refreshed")
	c.JSON(http.StatusOK, gin.H{"user": requestBody.Username, "message": "refreshed", "token": newToken})

}

// Logout godoc
//
//	@Summary		Refresh Access Token
//	@Description	Creates and returns new accessToken.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			username	body	struct{ username string }	true				"username"
//	@Success		200		{object}	map[string]interface{}	"token refreshed"
//	@Failure		400		{object}	map[string]string		"Invalid JSON"
//	@Failure		500		{object}	map[string]string		"Failed to refresh token"
//	@Router			/logout [post]
//	@Security		BearerAuth

func (handler *LoginHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	tracer := otel.Tracer("Handler-Tracer")
	ctx, span := tracer.Start(ctx, "LogoutHandler")
	defer span.End()
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		span.SetStatus(codes.Error, "invalid json")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	username, ok := body["username"]
	if !ok {
		span.SetStatus(codes.Error, "'username' is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "'username' is required"})
		return
	}
	usernameStr, ok := username.(string)
	if !ok || usernameStr == "" {
		span.SetStatus(codes.Error, "'username' must be a non-empty string")
		c.JSON(http.StatusBadRequest, gin.H{"error": "'username' must be a non-empty string"})
		return
	}

	if err := handler.LoginService.Logout(ctx, username.(string)); err != nil {
		span.SetStatus(codes.Error, err.Error())
		c.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
	span.SetStatus(codes.Ok, "Logged out")
	c.JSON(http.StatusOK, gin.H{"user": username, "message": "Logged out"})
}

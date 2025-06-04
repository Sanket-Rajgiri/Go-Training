package service

import (
	"albums/internal/models"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	getUserQuery     = "SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?"
	addUserQuery     = "INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`username`,`password`,`secret_key`,`role`,`revoked`) VALUES (?,?,?,?,?,?,?,?)"
	getUserbyIdQuery = "SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL"
	LoginUserQuery   = "UPDATE `users` SET `id`=?,`created_at`=?,`updated_at`=?,`deleted_at`=?,`username`=?,`password`=?,`secret_key`=?,`role`=?,`revoked`=? WHERE `users`.`deleted_at` IS NULL AND `id` = ?"
	LogoutUserQuery  = "UPDATE `users` SET `id`=?,`created_at`=?,`updated_at`=?,`deleted_at`=?,`username`=?,`password`=?,`secret_key`=?,`role`=?,`revoked`=? WHERE `users`.`deleted_at` IS NULL AND `id` = ?"
)

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock, username, password string)
		username    string
		password    string
		expectError bool
	}{
		{
			name:     "successful registration",
			username: "newuser",
			password: "password123",
			setupMock: func(mock sqlmock.Sqlmock, username, password string) {
				mock.ExpectQuery(regexp.QuoteMeta(getUserQuery)).
					WithArgs(username, 1).
					WillReturnError(gorm.ErrRecordNotFound)
				// Mock Create to succeed
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(addUserQuery)).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						username,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						"user",
						sqlmock.AnyArg(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name:     "username already exists",
			username: "existinguser",
			password: "password123",
			setupMock: func(mock sqlmock.Sqlmock, username, password string) {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "secret_key", "role"}).
					AddRow(1, username, password, "secretkey", "user")
				mock.ExpectQuery(regexp.QuoteMeta(getUserQuery)).
					WithArgs(username, 1).
					WillReturnRows(rows)
			},
			expectError: true,
		},
		{
			name:     "database error on user creation",
			username: "newuser",
			password: "password123",
			setupMock: func(mock sqlmock.Sqlmock, username, password string) {
				mock.ExpectQuery(regexp.QuoteMeta(getUserQuery)).
					WithArgs(username, 1).
					WillReturnError(gorm.ErrRecordNotFound)
				// Mock Create to return error
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(addUserQuery)).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						username,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						"user",
						sqlmock.AnyArg(),
					).
					WillReturnError(errors.New("insert error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			tt.setupMock(mock, tt.username, tt.password)

			svc := &LoginServiceImpl{DB: db}
			ctx := context.Background()

			user, err := svc.RegisterUser(ctx, tt.username, tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, models.Users{}, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.username, user.Username)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock, userID string)
		userID      string
		expectError bool
	}{
		{
			name:   "successful login",
			userID: "1",
			setupMock: func(mock sqlmock.Sqlmock, userID string) {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "secret_key", "role"}).
					AddRow(1, "validuser", "hashedpassword", "secretkey", "user")
				mock.ExpectQuery(regexp.QuoteMeta(getUserbyIdQuery)).
					WithArgs(userID).
					WillReturnRows(rows)
				mock.ExpectBegin()
				id, _ := strconv.Atoi(userID)
				mock.ExpectExec(regexp.QuoteMeta(LoginUserQuery)).
					WithArgs(
						id,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						"validuser",
						"hashedpassword",
						"secretkey",
						"user",
						false,
						id,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name:   "user not found",
			userID: "99",
			setupMock: func(mock sqlmock.Sqlmock, userID string) {
				mock.ExpectQuery(regexp.QuoteMeta(getUserbyIdQuery)).
					WithArgs(userID).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name:        "db error on save",
			expectError: true,
			userID:      "2",
			setupMock: func(mock sqlmock.Sqlmock, userID string) {
				id, _ := strconv.Atoi(userID)
				rows := sqlmock.NewRows([]string{"id", "username", "password", "secret_key", "role"}).
					AddRow(id, "validuser", "hashedpassword", "secretkey", "user")
				mock.ExpectQuery(regexp.QuoteMeta(getUserbyIdQuery)).
					WithArgs(userID).
					WillReturnRows(rows)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(LoginUserQuery)).
					WithArgs(
						id,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						"validuser",
						"hashedpassword",
						"secretkey",
						"user",
						false,
						id,
					).
					WillReturnError(fmt.Errorf("db failed to save"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			tt.setupMock(mock, tt.userID)

			svc := &LoginServiceImpl{DB: db}
			ctx := context.Background()

			_, _, err := svc.Login(ctx, tt.userID, "user")

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock, username string)
		username    string
		expectError bool
	}{
		{
			name:        "successful logout",
			username:    "validusername",
			expectError: false,
			setupMock: func(mock sqlmock.Sqlmock, username string) {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "secret_key", "role"}).
					AddRow(1, username, "hashedpassword", "secretkey", "user")
				mock.ExpectQuery(regexp.QuoteMeta(getUserQuery)).
					WithArgs(username, 1).
					WillReturnRows(rows)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(LogoutUserQuery)).
					WithArgs(
						1,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						username,
						"hashedpassword",
						"secretkey",
						"user",
						true,
						1,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:        "user not found",
			username:    "validusername",
			expectError: true,
			setupMock: func(mock sqlmock.Sqlmock, username string) {
				mock.ExpectQuery(regexp.QuoteMeta(getUserQuery)).
					WithArgs(username, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
		},
		{
			name:        "db failed to save",
			username:    "validusername",
			expectError: true,
			setupMock: func(mock sqlmock.Sqlmock, username string) {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "secret_key", "role"}).
					AddRow(1, username, "hashedpassword", "secretkey", "user")
				mock.ExpectQuery(regexp.QuoteMeta(getUserQuery)).
					WithArgs(username, 1).
					WillReturnRows(rows)
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(LogoutUserQuery)).
					WithArgs(
						1,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						username,
						"hashedpassword",
						"secretkey",
						"user",
						true,
						1,
					).
					WillReturnError(fmt.Errorf("db failed to save"))
				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			tt.setupMock(mock, tt.username)

			svc := &LoginServiceImpl{DB: db}
			ctx := context.Background()

			err := svc.Logout(ctx, tt.username)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRefreshToken(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		token       string
		setupMock   func(mock sqlmock.Sqlmock, username string)
		expectError bool
	}{
		{
			name:     "successful refresh token",
			username: "validusername",
			token:    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkwMjE0MjcsImV4cCI6MzA0Mjg2MTU1MCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJ1c2VyIn0.-aQx2LRZuqSBLYVDOBAq_1TNL2bKoorzHaEZEJ0tNx4",
			setupMock: func(mock sqlmock.Sqlmock, username string) {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "secret_key", "role"}).
					AddRow(1, username, "hashedpassword", "secretkey", "user")
				mock.ExpectQuery(regexp.QuoteMeta(getUserQuery)).
					WithArgs(username, 1).
					WillReturnRows(rows)
			},
			expectError: false,
		},
		{
			name:     "user not found",
			username: "username",
			token:    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkwMjE0MjcsImV4cCI6MzA0Mjg2MTU1MCwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJ1c2VyIn0.-aQx2LRZuqSBLYVDOBAq_1TNL2bKoorzHaEZEJ0tNx4",
			setupMock: func(mock sqlmock.Sqlmock, username string) {
				mock.ExpectQuery(regexp.QuoteMeta(getUserQuery)).
					WithArgs(username, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name:     "invalid Token",
			username: "validusername",
			token:    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkwMjE0MjcsImV4cCI6MTc0OTAyMTczMiwiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJ1c2VyIn0.dHpj5hZykc_SJ12M_Q5TXPSnoWjgmew_7BlCfRCXXWw",
			setupMock: func(mock sqlmock.Sqlmock, username string) {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "secret_key", "role"}).
					AddRow(1, username, "hashedpassword", "secretkey", "user")
				mock.ExpectQuery(regexp.QuoteMeta(getUserQuery)).
					WithArgs(username, 1).
					WillReturnRows(rows)
			},
			expectError: true,
		},
		{
			name:     "invalid Token claims",
			username: "validusername",
			token:    "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE3NDkwMjE0MjcsImV4cCI6MTc4MDU1Nzg3NywiYXVkIjoiIiwic3ViIjoiMSIsIlJvbGUiOiJhZG1pbiJ9.hnAg84vNSz9xBQYSCsXVlcErx6Qv3nx8620q7sWXGRg",
			setupMock: func(mock sqlmock.Sqlmock, username string) {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "secret_key", "role"}).
					AddRow(1, username, "hashedpassword", "secretkey", "user")
				mock.ExpectQuery(regexp.QuoteMeta(getUserQuery)).
					WithArgs(username, 1).
					WillReturnRows(rows)
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			tt.setupMock(mock, tt.username)

			svc := &LoginServiceImpl{DB: db}
			ctx := context.Background()

			_, err := svc.RefreshToken(ctx, tt.username, tt.token)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestValidateCredentials(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		password    string
		setupMock   func(mock sqlmock.Sqlmock, username, password string)
		expectError bool
	}{
		{
			name:     "Valid Credentials",
			username: "validUsername",
			password: "validPassword",
			setupMock: func(mock sqlmock.Sqlmock, username, password string) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				rows := sqlmock.NewRows([]string{"id", "username", "password", "secret_key", "role"}).
					AddRow(1, username, hashedPassword, "secretkey", "user")
				mock.ExpectQuery(regexp.QuoteMeta(getUserQuery)).
					WithArgs(username, 1).
					WillReturnRows(rows)
			},
		},
		{
			name:     "User Not Found",
			username: "InvalidUsername",
			password: "InvalidPassword",
			setupMock: func(mock sqlmock.Sqlmock, username, password string) {
				mock.ExpectQuery(regexp.QuoteMeta(getUserQuery)).
					WithArgs(username, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name:     "Invalid Credentials",
			username: "InvalidUsername",
			password: "InvalidPassword",
			setupMock: func(mock sqlmock.Sqlmock, username, password string) {
				rows := sqlmock.NewRows([]string{"id", "username", "password", "secret_key", "role"}).
					AddRow(1, username, password, "secretkey", "user")
				mock.ExpectQuery(regexp.QuoteMeta(getUserQuery)).
					WithArgs(username, 1).
					WillReturnRows(rows)
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			tt.setupMock(mock, tt.username, tt.password)

			svc := &LoginServiceImpl{DB: db}
			ctx := context.Background()

			_, _, err := svc.ValidateCredentials(ctx, tt.username, tt.password)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

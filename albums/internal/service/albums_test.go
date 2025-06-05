package service_test

import (
	"albums/internal/customlogs"
	definedMetrics "albums/internal/metrics"
	"albums/internal/models"
	"albums/internal/service"

	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.opentelemetry.io/otel"
	noopMetric "go.opentelemetry.io/otel/metric/noop"
	noopTrace "go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	// Step 1: Create mock SQL DB
	otel.SetTracerProvider(noopTrace.NewTracerProvider())
	customlogs.OtelLogger = otelzap.New(zap.NewNop())
	definedMetrics.AlbumsAdded = noopMetric.Int64Counter{}
	sqlDB, mock, err := sqlmock.New()
	assert.NoError(t, err)

	// Step 2: Configure GORM to use sqlmock DB with mysql dialector
	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true, // <--- required to avoid version query
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	// Step 3: Cleanup function to close mock DB
	cleanup := func() {
		sqlDB.Close()
	}

	return db, mock, cleanup
}

func TestAddAlbums(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(mock sqlmock.Sqlmock, album models.Album)
		album       models.Album
		expectError bool
	}{
		{
			name: "album added ",
			album: models.Album{
				Title:  "Test Album",
				Artist: "Test Artist",
				Price:  9.99,
			},
			setupMock: func(mock sqlmock.Sqlmock, album models.Album) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `albums` (`created_at`,`updated_at`,`deleted_at`,`title`,`artist`,`price`) VALUES (?,?,?,?,?,?)")).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						nil,
						album.Title,
						album.Artist,
						album.Price,
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "fails to begin transaction",
			album: models.Album{
				Title:  "Album X",
				Artist: "Artist Y",
				Price:  12.99,
			},
			setupMock: func(mock sqlmock.Sqlmock, album models.Album) {
				mock.ExpectBegin().WillReturnError(fmt.Errorf("begin error"))
			},
			expectError: true,
		},
		{
			name: "fails to insert album",
			album: models.Album{
				Title:  "Album Y",
				Artist: "Artist Z",
				Price:  13.49,
			},
			setupMock: func(mock sqlmock.Sqlmock, album models.Album) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO `albums` (`created_at`,`updated_at`,`deleted_at`,`title`,`artist`,`price`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, album.Title, album.Artist, album.Price).
					WillReturnError(fmt.Errorf("insert error"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
		{
			name: "fails to commit transaction",
			album: models.Album{
				Title:  "Album Z",
				Artist: "Artist Q",
				Price:  7.50,
			},
			setupMock: func(mock sqlmock.Sqlmock, album models.Album) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"INSERT INTO `albums` (`created_at`,`updated_at`,`deleted_at`,`title`,`artist`,`price`) VALUES (?,?,?,?,?,?)")).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, album.Title, album.Artist, album.Price).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(fmt.Errorf("commit error"))
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()
			svc := &service.AlbumServiceImpl{DB: db}
			ctx := context.Background()
			tt.setupMock(mock, tt.album)
			id, err := svc.AddAlbums(ctx, tt.album)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, uint(1), id)
			}
		})
	}
}

func TestGetAlbums(t *testing.T) {
	tests := []struct {
		name               string
		setupMock          func(mock sqlmock.Sqlmock)
		expectedLen        int
		expectedFirstTitle string
		expectError        bool
	}{
		{
			name: "albums found",
			setupMock: func(mock sqlmock.Sqlmock) {
				mockAlbums := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
					AddRow(1, "Album A", "Artist A", 10.0).
					AddRow(2, "Album B", "Artist B", 15.5)

				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `albums` WHERE `albums`.`deleted_at` IS NULL",
				)).WillReturnRows(mockAlbums)
			},
			expectedLen:        2,
			expectedFirstTitle: "Album A",
			expectError:        false,
		},
		{
			name: "no albums found",
			setupMock: func(mock sqlmock.Sqlmock) {
				mockAlbums := sqlmock.NewRows([]string{"id", "title", "artist", "price"})
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `albums` WHERE `albums`.`deleted_at` IS NULL",
				)).WillReturnRows(mockAlbums)
			},
			expectedLen:        0,
			expectedFirstTitle: "",
			expectError:        false,
		},
		{
			name: "error from db",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `albums` WHERE `albums`.`deleted_at` IS NULL",
				)).WillReturnError(fmt.Errorf("error from db"))
			},
			expectedLen:        0,
			expectedFirstTitle: "",
			expectError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			svc := &service.AlbumServiceImpl{DB: db}
			tt.setupMock(mock)

			ctx := context.Background()
			albums, err := svc.GetAlbums(ctx)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, albums, tt.expectedLen)
				if tt.expectedLen > 0 {
					assert.Equal(t, tt.expectedFirstTitle, albums[0].Title)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAlbumByID(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock, id uint)
		inputID       uint
		expectedTitle string
		expectError   bool
	}{
		{
			name:    "valid album ID",
			inputID: 2,
			setupMock: func(mock sqlmock.Sqlmock, id uint) {
				mockAlbums := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
					AddRow(2, "Album D", "Artist B", 15.5)
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `albums` WHERE `albums`.`id` = ? AND `albums`.`deleted_at` IS NULL ORDER BY `albums`.`id` LIMIT ?",
				)).WithArgs(id, 1).WillReturnRows(mockAlbums)
			},
			expectedTitle: "Album D",
			expectError:   false,
		},
		{
			name:    "album not found",
			inputID: 99,
			setupMock: func(mock sqlmock.Sqlmock, id uint) {
				mockAlbums := sqlmock.NewRows([]string{"id", "title", "artist", "price"})
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `albums` WHERE `albums`.`id` = ? AND `albums`.`deleted_at` IS NULL ORDER BY `albums`.`id` LIMIT ?",
				)).WithArgs(id, 1).WillReturnRows(mockAlbums)
			},
			expectedTitle: "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			svc := &service.AlbumServiceImpl{DB: db}

			tt.setupMock(mock, tt.inputID)

			ctx := context.Background()
			album, err := svc.GetAlbumByID(ctx, tt.inputID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTitle, album.Title)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdatePrice(t *testing.T) {
	tests := []struct {
		name        string
		album       models.Album
		setupMock   func(mock sqlmock.Sqlmock, album models.Album)
		expectError bool
	}{
		{
			name: "successful update",
			album: models.Album{
				Model: gorm.Model{ID: 1},
				Price: 20.0,
			},
			setupMock: func(mock sqlmock.Sqlmock, album models.Album) {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `albums` WHERE `albums`.`id` = ? AND `albums`.`deleted_at` IS NULL ORDER BY `albums`.`id` LIMIT ?")).
					WithArgs(album.ID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
						AddRow(album.ID, "Old Title", "Old Artist", 15.0))

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"UPDATE `albums` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`title`=?,`artist`=?,`price`=? WHERE `albums`.`deleted_at` IS NULL AND `id` = ?")).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						nil,
						"Old Title",
						"Old Artist",
						album.Price,
						album.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 2))
				mock.ExpectCommit()
			},
			expectError: false,
		},
		{
			name: "album not found",
			album: models.Album{
				Model: gorm.Model{ID: 999},
				Price: 99.99,
			},
			setupMock: func(mock sqlmock.Sqlmock, album models.Album) {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `albums` WHERE `albums`.`id` = ? AND `albums`.`deleted_at` IS NULL ORDER BY `albums`.`id` LIMIT ?")).
					WithArgs(album.ID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name: "db error on save",
			album: models.Album{
				Model: gorm.Model{ID: 2},
				Price: 45.0,
			},
			setupMock: func(mock sqlmock.Sqlmock, album models.Album) {

				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `albums` WHERE `albums`.`id` = ? AND `albums`.`deleted_at` IS NULL ORDER BY `albums`.`id` LIMIT ?")).
					WithArgs(album.ID, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
						AddRow(album.ID, "Some Title", "Some Artist", 30.0))

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"UPDATE `albums` SET `created_at`=?,`updated_at`=?,`deleted_at`=?,`title`=?,`artist`=?,`price`=? WHERE `albums`.`deleted_at` IS NULL AND `id` = ?")).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						nil,
						"Some Title",
						"Some Artist",
						album.Price,
						album.ID,
					).
					WillReturnError(fmt.Errorf("db update failed"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()
			svc := &service.AlbumServiceImpl{DB: db}
			ctx := context.Background()

			tt.setupMock(mock, tt.album)

			id, err := svc.UpdatePrice(ctx, tt.album)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, uint(0), id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.album.ID, id)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteAlbum(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock, id uint)
		inputID       uint
		expectedTitle string
		expectError   bool
	}{
		{
			name:    "successful delete",
			inputID: 1,
			setupMock: func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `albums` WHERE `albums`.`id` = ? AND `albums`.`deleted_at` IS NULL ORDER BY `albums`.`id` LIMIT ?")).
					WithArgs(id, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
						AddRow(id, "Test Title", "Test Artist", 10.99))

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"UPDATE `albums` SET `deleted_at`=? WHERE `albums`.`id` = ? AND `albums`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), id).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "album not found",
			inputID: 2,
			setupMock: func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `albums` WHERE `albums`.`id` = ? AND `albums`.`deleted_at` IS NULL ORDER BY `albums`.`id` LIMIT ?")).
					WithArgs(id, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name:    "delete failed",
			inputID: 3,
			setupMock: func(mock sqlmock.Sqlmock, id uint) {
				mock.ExpectQuery(regexp.QuoteMeta(
					"SELECT * FROM `albums` WHERE `albums`.`id` = ? AND `albums`.`deleted_at` IS NULL ORDER BY `albums`.`id` LIMIT ?")).
					WithArgs(id, 1).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
						AddRow(id, "Bad Title", "Bad Artist", 5.55))

				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(
					"UPDATE `albums` SET `deleted_at`=? WHERE `albums`.`id` = ? AND `albums`.`deleted_at` IS NULL")).
					WithArgs(sqlmock.AnyArg(), id).
					WillReturnError(errors.New("delete failed"))
				mock.ExpectRollback()
			},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()
			svc := &service.AlbumServiceImpl{DB: db}
			ctx := context.Background()

			tt.setupMock(mock, tt.inputID)

			id, err := svc.DeleteAlbum(ctx, tt.inputID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, uint(0), id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.inputID, id)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})

	}

}

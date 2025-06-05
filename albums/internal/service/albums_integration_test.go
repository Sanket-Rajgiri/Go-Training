package service_test

// import (
// 	"albums/internal/models"
// 	"albums/internal/service"
// 	"context"
// 	"log"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/testcontainers/testcontainers-go"
// 	testMysql "github.com/testcontainers/testcontainers-go/modules/mysql"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// type TestMySQLContainer struct {
// 	Container *testMysql.MySQLContainer
// 	DSN       string
// }

// func StartMySQLContainer(ctx context.Context) (*TestMySQLContainer, error) {
// 	container, err := testMysql.RunContainer(ctx,
// 		testcontainers.WithImage("mysql:8"),
// 		testMysql.WithDatabase("testdb"),
// 		testMysql.WithUsername("testuser"),
// 		testMysql.WithPassword("testpass"),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	dsn, err := container.ConnectionString(ctx, "test")
// 	if err != nil {
// 		return nil, err
// 	}

// 	log.Printf("Started MySQL container at DSN: %s", dsn)

// 	return &TestMySQLContainer{
// 		Container: container,
// 		DSN:       dsn,
// 	}, nil
// }

// func InitTestDB(ctx context.Context) (*gorm.DB, func(), error) {
// 	mysqlContainer, err := StartMySQLContainer(ctx)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	db, err := gorm.Open(mysql.Open(mysqlContainer.DSN), &gorm.Config{})
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	// Apply schema
// 	err = db.AutoMigrate(&models.Album{})
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	cleanup := func() {
// 		mysqlContainer.Container.Terminate(ctx)
// 	}

// 	return db, cleanup, nil
// }

// func TestAlbumsServiceWithRealDB(t *testing.T) {
// 	ctx := context.Background()

// 	db, cleanup, err := InitTestDB(ctx)
// 	if err != nil {
// 		t.Fatalf("failed to init test DB: %v", err)
// 	}
// 	defer cleanup()

// 	svc := service.AlbumServiceImpl{
// 		DB: db,
// 	}

// 	// Create an album
// 	id, err := svc.AddAlbums(ctx, models.Album{
// 		Title:  "Test Album",
// 		Artist: "Navdeep",
// 		Price:  10.0,
// 	})
// 	assert.NoError(t, err)
// 	assert.IsType(t, uint(0), id)

// 	// Fetch it back
// 	albums, err := svc.GetAlbums(ctx)
// 	assert.NoError(t, err)
// 	assert.Len(t, albums, 1)
// 	assert.Equal(t, "Test Album", albums[0].Title)
// }

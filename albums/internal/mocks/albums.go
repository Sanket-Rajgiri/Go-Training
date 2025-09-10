package mocks

import (
	models "albums/internal/models"
	"context"

	mock "github.com/stretchr/testify/mock"
)

type AlbumService struct {
	mock.Mock
}

func (m *AlbumService) AddAlbums(ctx context.Context, album models.Album) (uint, error) {
	args := m.Called(ctx, album)
	return args.Get(0).(uint), args.Error(1)
}

func (m *AlbumService) DeleteAlbum(ctx context.Context, id uint) (uint, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(uint), args.Error(1)
}

func (m *AlbumService) GetAlbumByID(ctx context.Context, id uint) (models.Album, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Album), args.Error(1)
}

func (m *AlbumService) GetAlbums(ctx context.Context) ([]models.Album, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Album), args.Error(1)
}

func (m *AlbumService) UpdatePrice(ctx context.Context, album models.Album) (uint, error) {
	args := m.Called(ctx, album)
	return args.Get(0).(uint), args.Error(1)
}

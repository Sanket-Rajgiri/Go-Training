package mocks

import (
	models "albums/internal/models"

	mock "github.com/stretchr/testify/mock"
)

type AlbumService struct {
	mock.Mock
}

func (m *AlbumService) AddAlbums(album models.Album) (uint, error) {
	args := m.Called(album)
	return args.Get(0).(uint), args.Error(1)
}

func (m *AlbumService) DeleteAlbum(id uint) (uint, error) {
	args := m.Called(id)
	return args.Get(0).(uint), args.Error(1)
}

func (m *AlbumService) GetAlbumByID(id uint) (models.Album, error) {
	args := m.Called(id)
	return args.Get(0).(models.Album), args.Error(1)
}

func (m *AlbumService) GetAlbums() ([]models.Album, error) {
	args := m.Called()
	return args.Get(0).([]models.Album), args.Error(1)
}

func (m *AlbumService) UpdatePrice(album models.Album) (uint, error) {
	args := m.Called(album)
	return args.Get(0).(uint), args.Error(1)
}

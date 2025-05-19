package service

import (
	"albums/internal/models"

	"gorm.io/gorm"
)

type AlbumsService interface {
	GetAlbums() ([]models.Album, error)
	GetAlbumByID(id uint) (models.Album, error)
	AddAlbums(album models.Album) (uint, error)
	UpdatePrice(album models.Album) (uint, error)
	DeleteAlbum(id uint) (uint, error)
}

type AlbumServiceImpl struct {
	DB *gorm.DB
}

func (service *AlbumServiceImpl) GetAlbums() ([]models.Album, error) {
	var albums []models.Album
	if err := service.DB.Find(&albums).Error; err != nil {
		return nil, err
	}
	return albums, nil
}

func (service *AlbumServiceImpl) GetAlbumByID(id uint) (models.Album, error) {
	var album models.Album
	if err := service.DB.First(&album, id).Error; err != nil {
		return album, err
	}
	return album, nil
}

func (service *AlbumServiceImpl) AddAlbums(album models.Album) (uint, error) {
	if err := service.DB.Create(&album).Error; err != nil {
		return 0, err
	}
	return album.ID, nil
}

func (service *AlbumServiceImpl) UpdatePrice(album models.Album) (uint, error) {
	existingAlbum, err := service.GetAlbumByID(album.ID)
	if err != nil {
		return 0, err
	}
	existingAlbum.Price = album.Price
	if err := service.DB.Save(&existingAlbum).Error; err != nil {
		return 0, err
	}
	return album.ID, nil
}

func (service *AlbumServiceImpl) DeleteAlbum(id uint) (uint, error) {
	album, err := service.GetAlbumByID(id)
	if err != nil {
		return 0, err
	}
	if err := service.DB.Delete(&album).Error; err != nil {
		return 0, err
	}
	return id, nil
}

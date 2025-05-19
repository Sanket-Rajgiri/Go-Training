package service

import (
	"albums/internal/database"
	"albums/internal/models"
)

func GetAlbums() ([]models.Album, error) {
	var albums []models.Album
	if err := database.DB.Find(&albums).Error; err != nil {
		return nil, err
	}
	return albums, nil
}

func GetAlbumByID(id uint) (models.Album, error) {
	var album models.Album
	if err := database.DB.First(&album, id).Error; err != nil {
		return album, err
	}
	return album, nil
}

func AddAlbums(album models.Album) (uint, error) {
	if err := database.DB.Create(&album).Error; err != nil {
		return 0, err
	}
	return album.ID, nil
}

func UpdatePrice(album models.Album) (uint, error) {
	existingAlbum, err := GetAlbumByID(album.ID)
	if err != nil {
		return 0, err
	}
	existingAlbum.Price = album.Price
	if err := database.DB.Save(&existingAlbum).Error; err != nil {
		return 0, err
	}
	return album.ID, nil
}

func DeleteAlbum(id uint) (uint, error) {
	album, err := GetAlbumByID(id)
	if err != nil {
		return 0, err
	}
	if err := database.DB.Delete(&album).Error; err != nil {
		return 0, err
	}
	return id, nil
}

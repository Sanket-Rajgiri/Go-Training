package service

import (
	"albums/internal/customlogs"
	"albums/internal/models"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type AlbumsService interface {
	GetAlbums(ctx context.Context) ([]models.Album, error)
	GetAlbumByID(ctx context.Context, id uint) (models.Album, error)
	AddAlbums(ctx context.Context, album models.Album) (uint, error)
	UpdatePrice(ctx context.Context, album models.Album) (uint, error)
	DeleteAlbum(ctx context.Context, id uint) (uint, error)
}

type AlbumServiceImpl struct {
	DB *gorm.DB
}

func (service *AlbumServiceImpl) GetAlbums(ctx context.Context) ([]models.Album, error) {
	var albums []models.Album
	if err := service.DB.WithContext(ctx).Find(&albums).Error; err != nil {
		customlogs.OtelLogger.Error(fmt.Sprintf("Error fetching albums: %v", err))
		return nil, err
	}
	customlogs.OtelLogger.Info("Albums fetched Successfully")
	return albums, nil
}

func (service *AlbumServiceImpl) GetAlbumByID(ctx context.Context, id uint) (models.Album, error) {
	var album models.Album
	if err := service.DB.WithContext(ctx).First(&album, id).Error; err != nil {
		customlogs.OtelLogger.Error(fmt.Sprintf("Error fetching album: %v", err))
		return album, err
	}
	customlogs.OtelLogger.Info(fmt.Sprintf("Album: %d fetched Successfully", album.ID))
	return album, nil
}

func (service *AlbumServiceImpl) AddAlbums(ctx context.Context, album models.Album) (uint, error) {
	if err := service.DB.WithContext(ctx).Create(&album).Error; err != nil {
		customlogs.OtelLogger.Error(fmt.Sprintf("Error adding album: %v", err))
		return 0, err
	}
	customlogs.OtelLogger.Info(fmt.Sprintf("Album: %d added Successfully", album.ID))
	return album.ID, nil
}

func (service *AlbumServiceImpl) UpdatePrice(ctx context.Context, album models.Album) (uint, error) {
	existingAlbum, err := service.GetAlbumByID(ctx, album.ID)
	if err != nil {
		return 0, err
	}
	existingAlbum.Price = album.Price
	if err := service.DB.WithContext(ctx).Save(&existingAlbum).Error; err != nil {
		customlogs.OtelLogger.Error(fmt.Sprintf("Error updating album: %v", err))
		return 0, err
	}
	customlogs.OtelLogger.Info(fmt.Sprintf("Album: %d updated Successfully", album.ID))
	return album.ID, nil
}

func (service *AlbumServiceImpl) DeleteAlbum(ctx context.Context, id uint) (uint, error) {
	album, err := service.GetAlbumByID(ctx, id)
	if err != nil {
		return 0, err
	}
	if err := service.DB.WithContext(ctx).Delete(&album).Error; err != nil {
		customlogs.OtelLogger.Error(fmt.Sprintf("Error deleting album: %v", err))
		return 0, err
	}
	customlogs.OtelLogger.Info(fmt.Sprintf("Album: %d deleted Successfully", album.ID))
	return id, nil
}

package service

import (
	"albums/internal/customlogs"
	definedMetrics "albums/internal/metrics"
	"albums/internal/models"
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
)

type AlbumServiceImpl struct {
	DB *gorm.DB
}

func addedAlbumsMetric(ctx context.Context) {
	definedMetrics.AlbumsAdded.Add(ctx, 1)
}

func (service *AlbumServiceImpl) GetAlbums(ctx context.Context) ([]models.Album, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "GetAlbums")
	defer span.End()
	var albums []models.Album
	if err := service.DB.WithContext(ctx).Find(&albums).Error; err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("Error fetching albums: %v", err))
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	customlogs.OtelLogger.Ctx(ctx).Info("Albums fetched Successfully")
	span.SetStatus(codes.Ok, "albums fetched")
	return albums, nil
}

func (service *AlbumServiceImpl) GetAlbumByID(ctx context.Context, id uint) (models.Album, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "GetAlbumByID")
	defer span.End()
	var album models.Album
	if err := service.DB.WithContext(ctx).First(&album, id).Error; err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("Error fetching album: %v", err))
		span.SetStatus(codes.Error, err.Error())
		return album, err
	}
	customlogs.OtelLogger.Ctx(ctx).Info(fmt.Sprintf("Album: %d fetched Successfully", album.ID))
	span.SetStatus(codes.Ok, "album fetched")
	span.SetAttributes(attribute.Int64("albumID", int64(album.ID)))
	return album, nil
}

func (service *AlbumServiceImpl) AddAlbums(ctx context.Context, album models.Album) (uint, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "AddAlbums")
	defer span.End()
	if err := service.DB.WithContext(ctx).Create(&album).Error; err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("Error adding album: %v", err))
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}
	customlogs.OtelLogger.Ctx(ctx).Info(fmt.Sprintf("Album: %d added Successfully", album.ID))
	addedAlbumsMetric(ctx)
	span.SetStatus(codes.Ok, "album added")
	span.SetAttributes(attribute.Int64("albumID", int64(album.ID)))
	return album.ID, nil
}

func (service *AlbumServiceImpl) UpdatePrice(ctx context.Context, album models.Album) (uint, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "UpdatePrice")
	defer span.End()
	existingAlbum, err := service.GetAlbumByID(ctx, album.ID)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}
	existingAlbum.Price = album.Price
	if err := service.DB.WithContext(ctx).Save(&existingAlbum).Error; err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("Error updating album: %v", err))
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}
	customlogs.OtelLogger.Ctx(ctx).Info(fmt.Sprintf("Album: %d updated Successfully", album.ID))
	span.SetStatus(codes.Ok, "album updated")
	span.SetAttributes(attribute.Int64("albumID", int64(album.ID)))
	return album.ID, nil
}

func (service *AlbumServiceImpl) DeleteAlbum(ctx context.Context, id uint) (uint, error) {
	tracer := otel.Tracer("Service-Tracer")
	ctx, span := tracer.Start(ctx, "DeleteAlbum")
	defer span.End()
	album, err := service.GetAlbumByID(ctx, id)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}
	if err := service.DB.WithContext(ctx).Delete(&album).Error; err != nil {
		customlogs.OtelLogger.Ctx(ctx).Error(fmt.Sprintf("Error deleting album: %v", err))
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}
	customlogs.OtelLogger.Ctx(ctx).Info(fmt.Sprintf("Album: %d deleted Successfully", album.ID))
	span.SetStatus(codes.Ok, "album deleted")
	span.SetAttributes(attribute.Int64("albumID", int64(album.ID)))
	return id, nil
}

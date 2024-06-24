package db

import (
	"context"
	"ratatoskr/internal/models"
)

type DB interface {
	GetAllGroupsWithTags(context.Context) (*[]models.Group, error)
	UpdateTags(context.Context, *[]models.Group) error
	InsertAnalytics(context.Context, *[]models.Analytics) error
}

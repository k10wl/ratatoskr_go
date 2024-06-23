package db

import (
	"context"
	"ratatoskr/internal/models"
)

type DB interface {
	GetAllGroupsWithTags(context.Context) (*[]models.Group, error)
}

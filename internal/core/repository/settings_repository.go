package repository

import (
	"context"

	"gohac/internal/core/domain"
)

// SettingsRepository defines the interface for settings data access
type SettingsRepository interface {
	// GetGlobalSettings retrieves global site settings
	GetGlobalSettings(ctx context.Context) (*domain.GlobalSettings, error)

	// UpdateGlobalSettings updates global site settings
	UpdateGlobalSettings(ctx context.Context, settings *domain.GlobalSettings) error
}

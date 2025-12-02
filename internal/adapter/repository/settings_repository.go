package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"gohac/internal/core/domain"
	"gohac/internal/core/repository"

	"gorm.io/gorm"
)

// settingsRepository implements the SettingsRepository interface using GORM
type settingsRepository struct {
	db *gorm.DB
}

// NewSettingsRepository creates a new settings repository instance
func NewSettingsRepository(db *gorm.DB) repository.SettingsRepository {
	return &settingsRepository{db: db}
}

// GetGlobalSettings retrieves global site settings from system_configs table
func (r *settingsRepository) GetGlobalSettings(ctx context.Context) (*domain.GlobalSettings, error) {
	var config domain.SystemConfig
	err := r.db.WithContext(ctx).
		Where("key = ? AND tenant_id = ?", "global_settings", "").
		First(&config).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default settings if not found
			return &domain.GlobalSettings{
				SiteName:     "Gohac CMS",
				Logo:         "",
				Favicon:      "",
				ContactEmail: "",
			}, nil
		}
		return nil, fmt.Errorf("failed to get global settings: %w", err)
	}

	var settings domain.GlobalSettings
	if len(config.Value) > 0 {
		if err := json.Unmarshal(config.Value, &settings); err != nil {
			return nil, fmt.Errorf("failed to unmarshal global settings: %w", err)
		}
	} else {
		// Return default settings if value is empty
		settings = domain.GlobalSettings{
			SiteName:     "Gohac CMS",
			Logo:         "",
			Favicon:      "",
			ContactEmail: "",
		}
	}

	return &settings, nil
}

// UpdateGlobalSettings updates global site settings in system_configs table
func (r *settingsRepository) UpdateGlobalSettings(ctx context.Context, settings *domain.GlobalSettings) error {
	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal global settings: %w", err)
	}

	// Check if config exists
	var existingConfig domain.SystemConfig
	err = r.db.WithContext(ctx).
		Where("key = ? AND tenant_id = ?", "global_settings", "").
		First(&existingConfig).Error

	if err == gorm.ErrRecordNotFound {
		// Create new config
		config := &domain.SystemConfig{
			TenantID: "",
			Key:      "global_settings",
			Value:    settingsJSON,
		}
		if err := r.db.WithContext(ctx).Create(config).Error; err != nil {
			return fmt.Errorf("failed to create global settings: %w", err)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check existing global settings: %w", err)
	}

	// Update existing config
	existingConfig.Value = settingsJSON
	if err := r.db.WithContext(ctx).Save(&existingConfig).Error; err != nil {
		return fmt.Errorf("failed to update global settings: %w", err)
	}

	return nil
}

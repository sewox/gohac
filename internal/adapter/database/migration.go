package database

import (
	"fmt"
	"log"

	"gohac/internal/core/domain"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

// Migrate runs all database migrations using gormigrate
func Migrate(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "20240101_init",
			Migrate: func(tx *gorm.DB) error {
				log.Println("Running migration 20240101_init: Creating Page table")
				return tx.AutoMigrate(&domain.Page{})
			},
			Rollback: func(tx *gorm.DB) error {
				log.Println("Rolling back migration 20240101_init")
				return tx.Migrator().DropTable(&domain.Page{})
			},
		},
		{
			ID: "20240102_settings",
			Migrate: func(tx *gorm.DB) error {
				log.Println("Running migration 20240102_settings: Creating SystemConfig and Menu tables")
				return tx.AutoMigrate(&domain.SystemConfig{}, &domain.Menu{})
			},
			Rollback: func(tx *gorm.DB) error {
				log.Println("Rolling back migration 20240102_settings")
				return tx.Migrator().DropTable(&domain.SystemConfig{}, &domain.Menu{})
			},
		},
		{
			ID: "20240103_remove_menu_position_column",
			Migrate: func(tx *gorm.DB) error {
				log.Println("Running migration 20240103_remove_menu_position_column: Refactoring Menu table")

				// Check if menus table has old 'position' column
				if tx.Migrator().HasTable(&domain.Menu{}) && tx.Migrator().HasColumn(&domain.Menu{}, "position") {
					log.Println("Old 'position' column found in 'menus' table. Performing data migration.")

					// SQLite migration: recreate table
					tempTable := "menus_new"

					// Use transaction to ensure atomicity
					return tx.Transaction(func(tx2 *gorm.DB) error {
						// Step 1: Create temporary table with new schema
						if err := tx2.Exec(`
							CREATE TABLE IF NOT EXISTS ` + tempTable + ` (
								id TEXT PRIMARY KEY,
								tenant_id TEXT,
								name TEXT NOT NULL,
								description TEXT,
								items TEXT,
								created_at DATETIME,
								updated_at DATETIME
							)
						`).Error; err != nil {
							return fmt.Errorf("failed to create temp table: %w", err)
						}

						// Step 2: Copy data (if any) - use position value as name
						if err := tx2.Exec(`
							INSERT INTO ` + tempTable + ` (id, tenant_id, name, description, items, created_at, updated_at)
							SELECT id, tenant_id, COALESCE(position, 'Unnamed Menu') as name, '' as description, items, created_at, updated_at
							FROM menus
						`).Error; err != nil {
							return fmt.Errorf("failed to copy data: %w", err)
						}

						// Step 3: Drop old table
						if err := tx2.Exec("DROP TABLE menus").Error; err != nil {
							return fmt.Errorf("failed to drop old table: %w", err)
						}

						// Step 4: Rename new table
						if err := tx2.Exec("ALTER TABLE " + tempTable + " RENAME TO menus").Error; err != nil {
							return fmt.Errorf("failed to rename table: %w", err)
						}

						// Step 5: Recreate indexes
						if err := tx2.Exec("CREATE INDEX IF NOT EXISTS idx_menus_tenant_id ON menus(tenant_id)").Error; err != nil {
							return fmt.Errorf("failed to create index: %w", err)
						}

						log.Println("✅ Menu table migration completed successfully")
						return nil
					})
				} else {
					log.Println("No old 'position' column found or 'menus' table does not exist. Skipping specific column migration.")
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				log.Println("Rollback for 20240103_remove_menu_position_column is not fully implemented for complex schema changes.")
				return nil // Complex rollback for table recreation is often manual or not fully automated
			},
		},
		{
			ID: "20240104_users",
			Migrate: func(tx *gorm.DB) error {
				log.Println("Running migration 20240104_users: Creating Users table")

				// Create users table
				if err := tx.AutoMigrate(&domain.User{}); err != nil {
					return fmt.Errorf("failed to create users table: %w", err)
				}

				// Check if table is empty, then seed default admin
				var count int64
				if err := tx.Model(&domain.User{}).Count(&count).Error; err != nil {
					return fmt.Errorf("failed to count users: %w", err)
				}

				if count == 0 {
					log.Println("Users table is empty. Seeding default admin user...")
					adminUser := &domain.User{
						Name:     "Admin",
						Email:    "admin@example.com",
						Password: "password", // Will be hashed
						Role:     domain.UserRoleAdmin,
					}

					// Hash password before saving
					if err := adminUser.HashPassword(); err != nil {
						return fmt.Errorf("failed to hash admin password: %w", err)
					}

					if err := tx.Create(adminUser).Error; err != nil {
						return fmt.Errorf("failed to seed admin user: %w", err)
					}

					log.Println("✅ Default admin user created: admin@example.com / password")
				} else {
					log.Println("Users table already has data. Skipping seed.")
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				log.Println("Rolling back migration 20240104_users")
				return tx.Migrator().DropTable(&domain.User{})
			},
		},
		{
			ID: "20240105_blog",
			Migrate: func(tx *gorm.DB) error {
				log.Println("Running migration 20240105_blog: Creating Posts and Categories tables")

				// Create posts table
				if err := tx.AutoMigrate(&domain.Post{}); err != nil {
					return fmt.Errorf("failed to create posts table: %w", err)
				}

				// Create categories table
				if err := tx.AutoMigrate(&domain.Category{}); err != nil {
					return fmt.Errorf("failed to create categories table: %w", err)
				}

				// Create post_categories join table (many-to-many)
				// GORM will create this automatically, but we ensure it exists
				if !tx.Migrator().HasTable("post_categories") {
					if err := tx.Exec(`
						CREATE TABLE IF NOT EXISTS post_categories (
							post_id TEXT NOT NULL,
							category_id TEXT NOT NULL,
							PRIMARY KEY (post_id, category_id),
							FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
							FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
						)
					`).Error; err != nil {
						return fmt.Errorf("failed to create post_categories table: %w", err)
					}
				}

				log.Println("✅ Blog tables created successfully")
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				log.Println("Rolling back migration 20240105_blog")
				return tx.Migrator().DropTable("post_categories", &domain.Post{}, &domain.Category{})
			},
		},
	})

	if err := m.Migrate(); err != nil {
		return fmt.Errorf("could not migrate: %v", err)
	}

	log.Println("✅ All migrations completed successfully")
	return nil
}

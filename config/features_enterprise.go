//go:build enterprise

package config

// IsEnterprise returns true for enterprise builds
func IsEnterprise() bool {
	return true
}

// IsCommunity returns false for enterprise builds
func IsCommunity() bool {
	return false
}

// GetDatabaseDriver returns the database driver name based on build tags
func GetDatabaseDriver() string {
	return "postgres"
}

// SupportsMultiTenancy returns true if multi-tenancy is supported
func SupportsMultiTenancy() bool {
	return true
}

// SupportsS3Storage returns true if S3 storage is supported
func SupportsS3Storage() bool {
	return true
}

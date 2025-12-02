//go:build community

package config

// IsEnterprise returns false for community builds
func IsEnterprise() bool {
	return false
}

// IsCommunity returns true for community builds
func IsCommunity() bool {
	return true
}

// GetDatabaseDriver returns the database driver name based on build tags
func GetDatabaseDriver() string {
	return "sqlite"
}

// SupportsMultiTenancy returns true if multi-tenancy is supported
func SupportsMultiTenancy() bool {
	return false
}

// SupportsS3Storage returns true if S3 storage is supported
func SupportsS3Storage() bool {
	return false
}

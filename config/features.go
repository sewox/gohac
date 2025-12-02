//go:build !community && !enterprise

package config

// Default implementations when no build tag is specified
// These default to community edition behavior

// IsEnterprise returns false by default (community edition)
func IsEnterprise() bool {
	return false
}

// IsCommunity returns true by default (community edition)
func IsCommunity() bool {
	return true
}

// GetDatabaseDriver returns the database driver name
func GetDatabaseDriver() string {
	return "sqlite"
}

// SupportsMultiTenancy returns false by default
func SupportsMultiTenancy() bool {
	return false
}

// SupportsS3Storage returns false by default
func SupportsS3Storage() bool {
	return false
}

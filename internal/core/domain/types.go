package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONB is a custom type for handling JSONB columns in GORM
// Works with both SQLite (JSON) and PostgreSQL (JSONB)
type JSONB json.RawMessage

// Value implements the driver.Valuer interface for database storage
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for database retrieval
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("failed to unmarshal JSONB value")
	}

	result := json.RawMessage{}
	err := json.Unmarshal(bytes, &result)
	*j = JSONB(result)
	return err
}

// MarshalJSON implements json.Marshaler
func (j JSONB) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler
func (j *JSONB) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("jsonb: UnmarshalJSON on nil pointer")
	}
	*j = JSONB(json.RawMessage(data))
	return nil
}

// ToBlocks converts JSONB to a slice of Block structs
func (j JSONB) ToBlocks() ([]Block, error) {
	if j == nil || len(j) == 0 {
		return []Block{}, nil
	}

	var blocks []Block
	err := json.Unmarshal(j, &blocks)
	return blocks, err
}

// FromBlocks converts a slice of Block structs to JSONB
func FromBlocks(blocks []Block) (JSONB, error) {
	if blocks == nil {
		return JSONB("[]"), nil
	}

	data, err := json.Marshal(blocks)
	if err != nil {
		return nil, err
	}
	return JSONB(data), nil
}

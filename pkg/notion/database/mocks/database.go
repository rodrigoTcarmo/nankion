package mocks

import (
	"errors"

	"github.com/dstotijn/go-notion"
)

// MockDatabaseClient implements notion.DatabaseClient for testing
type MockDatabaseClient struct {
	GetDatabaseFunc            func(databaseID string) (*notion.Database, error)
	ListDatabasePropertiesFunc func(databaseID string) (*notion.DatabaseProperties, error)
}

func (m *MockDatabaseClient) GetDatabase(databaseID string) (*notion.Database, error) {
	if m.GetDatabaseFunc != nil {
		return m.GetDatabaseFunc(databaseID)
	}
	return nil, errors.New("GetDatabaseFunc not set")
}

func (m *MockDatabaseClient) ListDatabaseProperties(databaseID string) (*notion.DatabaseProperties, error) {
	if m.ListDatabasePropertiesFunc != nil {
		return m.ListDatabasePropertiesFunc(databaseID)
	}
	return nil, errors.New("ListDatabasePropertiesFunc not set")
}

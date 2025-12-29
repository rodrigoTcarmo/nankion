package mocks

import (
	"errors"

	"github.com/rodrigoTcarmo/nankion/pkg/notion/page"
	"github.com/rodrigoTcarmo/nankion/pkg/statement"
)

// MockPageClient implements page.Page for testing
type MockPageClient struct {
	BuildPageFunc  func(databaseID string, statement statement.Statement) page.PageData
	CreatePageFunc func(pageData page.PageData) error
}

func (m *MockPageClient) BuildPage(databaseID string, stmt statement.Statement) page.PageData {
	if m.BuildPageFunc != nil {
		return m.BuildPageFunc(databaseID, stmt)
	}
	return page.PageData{}
}

func (m *MockPageClient) CreatePage(pageData page.PageData) error {
	if m.CreatePageFunc != nil {
		return m.CreatePageFunc(pageData)
	}
	return errors.New("CreatePageFunc not set")
}

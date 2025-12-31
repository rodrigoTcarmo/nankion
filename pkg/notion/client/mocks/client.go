package mocks

import (
	"context"
	"fmt"

	"github.com/dstotijn/go-notion"
)

type MockNotionClient struct {
	FindDatabaseByIDFunc func(context.Context, string) (notion.Database, error)
	UpdateDatabaseFunc   func(context.Context, string, notion.UpdateDatabaseParams) (notion.Database, error)
	CreatePageFunc       func(context.Context, notion.CreatePageParams) (notion.Page, error)
}

func (m *MockNotionClient) FindDatabaseByID(ctx context.Context, databaseId string) (notion.Database, error) {
	if m.FindDatabaseByIDFunc != nil {
		return m.FindDatabaseByIDFunc(ctx, databaseId)
	}
	return notion.Database{}, fmt.Errorf("FindDatabaseByIDFunc not set")
}

func (m *MockNotionClient) UpdateDatabase(ctx context.Context, databaseId string, params notion.UpdateDatabaseParams) (notion.Database, error) {
	if m.UpdateDatabaseFunc != nil {
		return m.UpdateDatabaseFunc(ctx, databaseId, params)
	}
	return notion.Database{}, fmt.Errorf("UpdateDatabaseFunc not set")
}

func (m *MockNotionClient) CreatePage(ctx context.Context, params notion.CreatePageParams) (notion.Page, error) {
	if m.CreatePageFunc != nil {
		return m.CreatePageFunc(ctx, params)
	}
	return notion.Page{}, fmt.Errorf("CreatePageFunc not set")
}

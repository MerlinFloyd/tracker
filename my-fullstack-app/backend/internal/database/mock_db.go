package database

import (
	"database/sql"
	"my-fullstack-app/backend/internal/models"
)

// MockDB represents a mock database implementation for testing
type MockDB struct {
	users []models.User
}

// NewMockDB creates a new mock database
func NewMockDB() *MockDB {
	return &MockDB{
		users: []models.User{
			{ID: 1, Username: "testuser", Email: "test@example.com", Password: "hashedpassword"},
		},
	}
}

// ConnectMock returns a mock implementation for testing
func ConnectMock() (*sql.DB, *MockDB, error) {
	mockDB := NewMockDB()
	return nil, mockDB, nil
}

// GetUserByID returns a user by ID
func (m *MockDB) GetUserByID(id int) (models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, sql.ErrNoRows
}

package examples

import (
	"testing"
	"time"
)



func TestUser_Struct(t *testing.T) {
	user := &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		CreateAt: time.Now(),
	}

	if user.ID != 1 {
		t.Errorf("Expected ID 1, got: %d", user.ID)
	}

	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got: %s", user.Username)
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got: %s", user.Email)
	}
}

func TestNewUserService(t *testing.T) {
	mockMapper := NewMockUserMapper()
	service := NewUserService(mockMapper)

	if service == nil {
		t.Error("Expected non-nil service")
	}

	if service.userMapper != mockMapper {
		t.Error("Expected userMapper to be set")
	}
}

func TestUserService_GetUser(t *testing.T) {
	mockMapper := NewMockUserMapper()
	service := NewUserService(mockMapper)

	// Test user not found
	user, err := service.GetUser(999)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if user != nil {
		t.Error("Expected nil user")
	}

	// Add a user to mock
	testUser := &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		CreateAt: time.Now(),
	}
	mockMapper.users[1] = testUser

	// Test user found
	user, err = service.GetUser(1)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if user == nil {
		t.Error("Expected non-nil user")
	}
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got: %s", user.Username)
	}

	// Test database error
	mockMapper.SetError(true)
	user, err = service.GetUser(1)
	if err == nil {
		t.Error("Expected database error")
	}
	if user != nil {
		t.Error("Expected nil user on error")
	}
}

func TestUserService_CreateUser(t *testing.T) {
	mockMapper := NewMockUserMapper()
	service := NewUserService(mockMapper)

	// Test successful creation
	user, err := service.CreateUser("newuser", "new@example.com")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if user == nil {
		t.Error("Expected non-nil user")
	}
	if user.Username != "newuser" {
		t.Errorf("Expected username 'newuser', got: %s", user.Username)
	}
	if user.Email != "new@example.com" {
		t.Errorf("Expected email 'new@example.com', got: %s", user.Email)
	}
	if user.ID == 0 {
		t.Error("Expected non-zero ID")
	}

	// Test database error
	mockMapper.SetError(true)
	user, err = service.CreateUser("erroruser", "error@example.com")
	if err == nil {
		t.Error("Expected database error")
	}
	if user != nil {
		t.Error("Expected nil user on error")
	}
}

func TestUserService_UpdateUserEmail(t *testing.T) {
	mockMapper := NewMockUserMapper()
	service := NewUserService(mockMapper)

	// Add a user to mock
	testUser := &User{
		ID:       1,
		Username: "testuser",
		Email:    "old@example.com",
		CreateAt: time.Now(),
	}
	mockMapper.users[1] = testUser

	// Test successful update
	err := service.UpdateUserEmail(1, "new@example.com")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify email was updated
	updatedUser := mockMapper.users[1]
	if updatedUser.Email != "new@example.com" {
		t.Errorf("Expected email 'new@example.com', got: %s", updatedUser.Email)
	}

	// Test user not found
	err = service.UpdateUserEmail(999, "notfound@example.com")
	if err != nil {
		t.Errorf("Expected no error for non-existent user, got: %v", err)
	}

	// Test database error on get
	mockMapper.SetError(true)
	err = service.UpdateUserEmail(1, "error@example.com")
	if err == nil {
		t.Error("Expected database error")
	}

	// Test database error on update
	mockMapper.SetError(false)
	mockMapper.users[2] = &User{ID: 2, Username: "user2", Email: "user2@example.com"}
	mockMapper.SetError(true)
	err = service.UpdateUserEmail(2, "newemail@example.com")
	if err == nil {
		t.Error("Expected database error on update")
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	mockMapper := NewMockUserMapper()
	service := NewUserService(mockMapper)

	// Add a user to mock
	testUser := &User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		CreateAt: time.Now(),
	}
	mockMapper.users[1] = testUser

	// Test successful deletion
	err := service.DeleteUser(1)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify user was deleted
	if _, exists := mockMapper.users[1]; exists {
		t.Error("Expected user to be deleted")
	}

	// Test deleting non-existent user
	err = service.DeleteUser(999)
	if err != nil {
		t.Errorf("Expected no error for non-existent user, got: %v", err)
	}

	// Test database error
	mockMapper.SetError(true)
	err = service.DeleteUser(1)
	if err == nil {
		t.Error("Expected database error")
	}
}

func TestUserService_SearchUsers(t *testing.T) {
	mockMapper := NewMockUserMapper()
	service := NewUserService(mockMapper)

	// Add users to mock
	user1 := &User{ID: 1, Username: "john", Email: "john@example.com"}
	user2 := &User{ID: 2, Username: "jane", Email: "jane@example.com"}
	user3 := &User{ID: 3, Username: "johnny", Email: "johnny@example.com"}
	mockMapper.users[1] = user1
	mockMapper.users[2] = user2
	mockMapper.users[3] = user3

	// Test search (note: this is a simplified mock that looks for exact matches)
	users, err := service.SearchUsers("john")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// The mock implementation looks for exact username matches
	// In real implementation, this would use LIKE with wildcards
	expectedCount := 1 // Only "john" matches exactly
	if len(users) != expectedCount {
		t.Errorf("Expected %d users, got: %d", expectedCount, len(users))
	}

	// Test database error
	mockMapper.SetError(true)
	users, err = service.SearchUsers("test")
	if err == nil {
		t.Error("Expected database error")
	}
	if users != nil {
		t.Error("Expected nil users on error")
	}
}

func TestUserService_GetUserCount(t *testing.T) {
	mockMapper := NewMockUserMapper()
	service := NewUserService(mockMapper)

	// Test empty count
	count, err := service.GetUserCount()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected count 0, got: %d", count)
	}

	// Add users
	user1 := &User{ID: 1, Username: "user1", Email: "user1@example.com"}
	user2 := &User{ID: 2, Username: "user2", Email: "user2@example.com"}
	mockMapper.users[1] = user1
	mockMapper.users[2] = user2

	// Test count with users
	count, err = service.GetUserCount()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if count != 2 {
		t.Errorf("Expected count 2, got: %d", count)
	}

	// Test database error
	mockMapper.SetError(true)
	count, err = service.GetUserCount()
	if err == nil {
		t.Error("Expected database error")
	}
	if count != 0 {
		t.Errorf("Expected count 0 on error, got: %d", count)
	}
}

func TestMockUserMapper_GetAllUsers(t *testing.T) {
	mockMapper := NewMockUserMapper()

	// Test empty result
	users, err := mockMapper.GetAllUsers()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(users) != 0 {
		t.Errorf("Expected 0 users, got: %d", len(users))
	}

	// Add users
	user1 := &User{ID: 1, Username: "user1", Email: "user1@example.com"}
	user2 := &User{ID: 2, Username: "user2", Email: "user2@example.com"}
	mockMapper.users[1] = user1
	mockMapper.users[2] = user2

	// Test with users
	users, err = mockMapper.GetAllUsers()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got: %d", len(users))
	}

	// Test database error
	mockMapper.SetError(true)
	users, err = mockMapper.GetAllUsers()
	if err == nil {
		t.Error("Expected database error")
	}
	if users != nil {
		t.Error("Expected nil users on error")
	}
}
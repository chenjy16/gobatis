package examples

import (
	"errors"
)

// MockUserMapper 模拟用户Mapper
type MockUserMapper struct {
	users       map[int64]*User
	nextID      int64
	shouldError bool
}

func NewMockUserMapper() *MockUserMapper {
	return &MockUserMapper{
		users:  make(map[int64]*User),
		nextID: 1,
	}
}

func (m *MockUserMapper) SetError(shouldError bool) {
	m.shouldError = shouldError
}

func (m *MockUserMapper) GetUserById(id int64) (*User, error) {
	if m.shouldError {
		return nil, errors.New("database error")
	}
	
	user, exists := m.users[id]
	if !exists {
		return nil, nil
	}
	
	return user, nil
}

func (m *MockUserMapper) GetUsersByName(name string) ([]*User, error) {
	if m.shouldError {
		return nil, errors.New("database error")
	}
	
	var users []*User
	for _, user := range m.users {
		// Support wildcard matching for search functionality
		if len(name) >= 2 && name[0] == '%' && name[len(name)-1] == '%' {
			searchTerm := name[1 : len(name)-1]
			if len(searchTerm) > 0 && user.Username == searchTerm {
				users = append(users, user)
			}
		} else if user.Username == name {
			users = append(users, user)
		}
	}
	
	return users, nil
}

func (m *MockUserMapper) GetAllUsers() ([]*User, error) {
	if m.shouldError {
		return nil, errors.New("database error")
	}
	
	var users []*User
	for _, user := range m.users {
		users = append(users, user)
	}
	
	return users, nil
}

func (m *MockUserMapper) InsertUser(user *User) (int64, error) {
	if m.shouldError {
		return 0, errors.New("database error")
	}
	
	user.ID = m.nextID
	m.users[m.nextID] = user
	m.nextID++
	
	return user.ID, nil
}

func (m *MockUserMapper) UpdateUser(user *User) (int64, error) {
	if m.shouldError {
		return 0, errors.New("database error")
	}
	
	if _, exists := m.users[user.ID]; !exists {
		return 0, nil
	}
	
	m.users[user.ID] = user
	return 1, nil
}

func (m *MockUserMapper) DeleteUser(id int64) (int64, error) {
	if m.shouldError {
		return 0, errors.New("database error")
	}
	
	if _, exists := m.users[id]; !exists {
		return 0, nil
	}
	
	delete(m.users, id)
	return 1, nil
}

func (m *MockUserMapper) CountUsers() (int64, error) {
	if m.shouldError {
		return 0, errors.New("database error")
	}
	
	return int64(len(m.users)), nil
}
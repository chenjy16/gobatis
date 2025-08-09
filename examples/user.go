package examples

import (
	"time"
)

// User 用户实体
type User struct {
	ID       int64     `db:"id"`
	Username string    `db:"username"`
	Email    string    `db:"email"`
	CreateAt time.Time `db:"create_at"`
}

// UserMapper 用户 Mapper 接口
type UserMapper interface {
	GetUserById(id int64) (*User, error)
	GetUsersByName(name string) ([]*User, error)
	GetAllUsers() ([]*User, error)
	InsertUser(user *User) (int64, error)
	UpdateUser(user *User) (int64, error)
	DeleteUser(id int64) (int64, error)
	CountUsers() (int64, error)
}

// UserService 用户服务
type UserService struct {
	userMapper UserMapper
}

// NewUserService 创建用户服务
func NewUserService(userMapper UserMapper) *UserService {
	return &UserService{
		userMapper: userMapper,
	}
}

// GetUser 获取用户
func (s *UserService) GetUser(id int64) (*User, error) {
	return s.userMapper.GetUserById(id)
}

// CreateUser 创建用户
func (s *UserService) CreateUser(username, email string) (*User, error) {
	user := &User{
		Username: username,
		Email:    email,
		CreateAt: time.Now(),
	}

	id, err := s.userMapper.InsertUser(user)
	if err != nil {
		return nil, err
	}

	user.ID = id
	return user, nil
}

// UpdateUserEmail 更新用户邮箱
func (s *UserService) UpdateUserEmail(id int64, email string) error {
	user, err := s.userMapper.GetUserById(id)
	if err != nil {
		return err
	}

	if user == nil {
		return nil
	}

	user.Email = email
	_, err = s.userMapper.UpdateUser(user)
	return err
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int64) error {
	_, err := s.userMapper.DeleteUser(id)
	return err
}

// SearchUsers 搜索用户
func (s *UserService) SearchUsers(name string) ([]*User, error) {
	return s.userMapper.GetUsersByName("%" + name + "%")
}

// GetUserCount 获取用户总数
func (s *UserService) GetUserCount() (int64, error) {
	return s.userMapper.CountUsers()
}
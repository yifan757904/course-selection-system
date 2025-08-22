package service

import (
	"errors"
	"testing"

	"github.com/liuyifan1996/course-selection-system/api/model"
)

type mockAuthRepo struct {
	users map[string]*model.User
}

func (m *mockAuthRepo) FindByIDCard(idCard string) (*model.User, error) {
	user, ok := m.users[idCard]
	if !ok {
		return nil, errors.New("not found")
	}
	return user, nil
}
func (m *mockAuthRepo) CreateUser(user *model.User) error {
	if _, ok := m.users[user.IDCard]; ok {
		return errors.New("exists")
	}
	m.users[user.IDCard] = user
	return nil
}
func (m *mockAuthRepo) DeleteUser(idCard, password string) error  { return nil }
func (m *mockAuthRepo) UpdateUser(user *model.User) error         { return nil }
func (m *mockAuthRepo) GetUserByID(id int64) (*model.User, error) { return nil, nil }

func TestRegister_PasswordComplexity(t *testing.T) {
	repo := &mockAuthRepo{users: map[string]*model.User{}}
	service := NewAuthService(repo)
	input := RegisterInput{IDCard: "123", Name: "A", Password: "abc", Role: "student"}
	_, err := service.Register(input)
	if err == nil {
		t.Error("should fail for weak password")
	}
	input.Password = "Abc123"
	_, err = service.Register(input)
	if err != nil {
		t.Errorf("should pass for strong password, got %v", err)
	}
}

func TestRegister_Uniqueness(t *testing.T) {
	repo := &mockAuthRepo{users: map[string]*model.User{"123": {IDCard: "123"}}}
	service := NewAuthService(repo)
	input := RegisterInput{IDCard: "123", Name: "A", Password: "Abc123", Role: "student"}
	_, err := service.Register(input)
	if err == nil {
		t.Error("should fail for duplicate user")
	}
}

func TestUpdateUser_PasswordCheck(t *testing.T) {
	repo := &mockAuthRepo{users: map[string]*model.User{"123": {IDCard: "123", Password: "Abc123", Name: "A"}}}
	service := NewAuthService(repo)
	input := UpdateUserInput{IDCard: "123", Password: "wrong", Name: nil, NewPassword: nil}
	_, err := service.UpdateUser(input)
	if err == nil {
		t.Error("should pass for no new password")
	}
	input.Password = "Abc123"
	newName := "B"
	input.Name = &newName
	_, err = service.UpdateUser(input)
	if err != nil {
		t.Errorf("should pass for correct password, got %v", err)
	}
}

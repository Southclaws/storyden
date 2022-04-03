package user

import (
	"context"
	"errors"
	"time"

	"4d63.com/optional"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/api/src/utils"
)

type mockRepo struct {
	m map[UserID]User
	r MockRepository
}

func NewMock() Repository {
	return &mockRepo{m: map[UserID]User{}}
}

func (m *mockRepo) CreateUser(ctx context.Context, email string, username string) (*User, error) {
	id := uuid.New()

	if _, ok := lo.Find(lo.Values(m.m), func(t User) bool { return email == t.Email }); ok {
		return nil, errors.New("email already exists")
	}

	u := User{
		ID:    UserID(id),
		Email: email,
		Name:  username,
	}

	m.m[UserID(id)] = u

	return &u, nil
}

func (m *mockRepo) GetUser(ctx context.Context, userId UserID, public bool) (*User, error) {
	u, ok := m.m[userId]
	if !ok {
		return nil, nil
	}
	return utils.Ref(u), nil
}

func (m *mockRepo) GetUserByEmail(ctx context.Context, email string, public bool) (*User, error) {
	u, ok := lo.Find(lo.Values(m.m), func(t User) bool { return email == t.Email })
	if !ok {
		return nil, nil
	}

	return &u, nil
}

func (m *mockRepo) GetUsers(ctx context.Context, sort string, max, skip int, public bool) ([]User, error) {
	return lo.Values(m.m), nil
}

func (m *mockRepo) UpdateUser(ctx context.Context, userId UserID, email, name, bio *string) (*User, error) {
	update := m.m[userId]

	if email != nil {
		update.Email = *email
	}
	if name != nil {
		update.Email = *name
	}
	if bio != nil {
		update.Email = *bio
	}

	m.m[userId] = update

	return &update, nil
}

func (m *mockRepo) SetAdmin(ctx context.Context, userId UserID, status bool) error {
	update := m.m[userId]
	update.Admin = status
	m.m[userId] = update

	return nil
}

func (m *mockRepo) Ban(ctx context.Context, userId UserID) (*User, error) {
	update := m.m[userId]
	update.DeletedAt = optional.Of(time.Now())
	m.m[userId] = update

	return &update, nil
}

func (m *mockRepo) Unban(ctx context.Context, userId UserID) (*User, error) {
	update := m.m[userId]
	update.DeletedAt = nil
	m.m[userId] = update

	return &update, nil
}

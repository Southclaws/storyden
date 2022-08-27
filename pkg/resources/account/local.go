package account

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/Southclaws/storyden/internal/utils"
)

type local struct {
	m map[AccountID]Account
}

func NewLocal() Repository {
	return &local{m: map[AccountID]Account{}}
}

func (m *local) GetResourceName() string { return "account" }

func (m *local) Create(ctx context.Context, email string, username string, opts ...option) (*Account, error) {
	id := uuid.New()

	if _, ok := lo.Find(lo.Values(m.m), func(t Account) bool { return email == t.Email }); ok {
		return nil, errors.New("email already exists")
	}

	u := Account{
		ID:    AccountID(id),
		Email: email,
		Name:  username,
	}

	m.m[AccountID(id)] = u

	return &u, nil
}

func (m *local) GetByID(ctx context.Context, userId AccountID) (*Account, error) {
	u, ok := m.m[userId]
	if !ok {
		return nil, nil
	}

	return utils.Ref(u), nil
}

func (m *local) LookupByEmail(ctx context.Context, email string) (*Account, bool, error) {
	u, ok := lo.Find(lo.Values(m.m), func(t Account) bool { return email == t.Email })
	if !ok {
		return nil, false, nil
	}

	return &u, true, nil
}

func (m *local) List(ctx context.Context, sort string, max, skip int) ([]Account, error) {
	return lo.Values(m.m), nil
}

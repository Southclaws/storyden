package authentication

import (
	"context"

	"github.com/google/uuid"

	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model"
	model_account "github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/account"
	"github.com/Southclaws/storyden/backend/internal/infrastructure/db/model/authentication"
	"github.com/Southclaws/storyden/backend/pkg/resources/account"
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) Create(ctx context.Context,
	id account.AccountID,
	service Service,
	identifier string,
	token string,
	metadata map[string]any,
) (*Authentication, error) {
	r, err := d.db.Authentication.Create().
		SetAccountID(uuid.UUID(id)).
		SetService(string(service)).
		SetIdentifier(identifier).
		SetToken(token).
		SetMetadata(metadata).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	r, err = d.db.Authentication.
		Query().
		Where(authentication.ID(r.ID)).
		WithAccount().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return FromModel(r), nil
}

func (d *database) GetByIdentifier(ctx context.Context, service Service, identifier string) (*Authentication, error) {
	r, err := d.db.Authentication.
		Query().
		Where(
			authentication.IdentifierEQ(identifier),
			authentication.ServiceEQ(string(service)),
		).
		WithAccount().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return FromModel(r), nil
}

func (d *database) GetAuthMethods(ctx context.Context, id account.AccountID) ([]Authentication, error) {
	r, err := d.db.Authentication.
		Query().
		Where(authentication.HasAccountWith(model_account.IDEQ(uuid.UUID(id)))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return FromModelMany(r), nil
}

func (d *database) IsEqual(ctx context.Context, id account.AccountID, identifier string, token string) (bool, error) {
	r, err := d.db.Authentication.
		Query().
		Where(
			authentication.HasAccountWith(model_account.IDEQ(uuid.UUID(id))),
			authentication.IdentifierEQ(identifier),
		).
		Only(ctx)
	if err != nil {
		return false, err
	}

	return r.Token == token, nil
}

package authentication

import (
	"context"

	"github.com/google/uuid"

	"github.com/Southclaws/storyden/api/src/infra/db/model"
	"github.com/Southclaws/storyden/api/src/infra/db/model/authentication"
	model_user "github.com/Southclaws/storyden/api/src/infra/db/model/user"
	"github.com/Southclaws/storyden/api/src/resources/user"
)

type database struct {
	db *model.Client
}

func New(db *model.Client) Repository {
	return &database{db}
}

func (d *database) Create(ctx context.Context,
	userID user.UserID,
	service Service,
	identifier string,
	token string,
	metadata map[string]any,
) (*Authentication, error) {
	r, err := d.db.Authentication.Create().
		SetUserID(uuid.UUID(userID)).
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
		WithUser().
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
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return FromModel(r), nil
}

func (d *database) GetAuthMethods(ctx context.Context, userID user.UserID) ([]Authentication, error) {
	r, err := d.db.Authentication.
		Query().
		Where(authentication.HasUserWith(model_user.IDEQ(uuid.UUID(userID)))).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return FromModelMany(r), nil
}

func (d *database) IsEqual(ctx context.Context, userID user.UserID, identifier string, token string) (bool, error) {
	r, err := d.db.Authentication.
		Query().
		Where(
			authentication.HasUserWith(model_user.IDEQ(uuid.UUID(userID))),
			authentication.IdentifierEQ(identifier),
		).
		Only(ctx)
	if err != nil {
		return false, err
	}

	return r.Token == token, nil
}

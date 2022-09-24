// Code generated by ent, DO NOT EDIT.

package model

import (
	"time"

	"github.com/Southclaws/storyden/internal/infrastructure/db/model/account"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/authentication"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/category"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/notification"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/post"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/react"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/role"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/subscription"
	"github.com/Southclaws/storyden/internal/infrastructure/db/model/tag"
	"github.com/Southclaws/storyden/internal/infrastructure/db/schema"
	"github.com/rs/xid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	accountMixin := schema.Account{}.Mixin()
	accountMixinFields0 := accountMixin[0].Fields()
	_ = accountMixinFields0
	accountMixinFields1 := accountMixin[1].Fields()
	_ = accountMixinFields1
	accountMixinFields2 := accountMixin[2].Fields()
	_ = accountMixinFields2
	accountFields := schema.Account{}.Fields()
	_ = accountFields
	// accountDescCreatedAt is the schema descriptor for created_at field.
	accountDescCreatedAt := accountMixinFields1[0].Descriptor()
	// account.DefaultCreatedAt holds the default value on creation for the created_at field.
	account.DefaultCreatedAt = accountDescCreatedAt.Default.(func() time.Time)
	// accountDescUpdatedAt is the schema descriptor for updated_at field.
	accountDescUpdatedAt := accountMixinFields2[0].Descriptor()
	// account.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	account.DefaultUpdatedAt = accountDescUpdatedAt.Default.(func() time.Time)
	// account.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	account.UpdateDefaultUpdatedAt = accountDescUpdatedAt.UpdateDefault.(func() time.Time)
	// accountDescHandle is the schema descriptor for handle field.
	accountDescHandle := accountFields[1].Descriptor()
	// account.HandleValidator is a validator for the "handle" field. It is called by the builders before save.
	account.HandleValidator = accountDescHandle.Validators[0].(func(string) error)
	// accountDescName is the schema descriptor for name field.
	accountDescName := accountFields[2].Descriptor()
	// account.NameValidator is a validator for the "name" field. It is called by the builders before save.
	account.NameValidator = accountDescName.Validators[0].(func(string) error)
	// accountDescAdmin is the schema descriptor for admin field.
	accountDescAdmin := accountFields[4].Descriptor()
	// account.DefaultAdmin holds the default value on creation for the admin field.
	account.DefaultAdmin = accountDescAdmin.Default.(bool)
	// accountDescID is the schema descriptor for id field.
	accountDescID := accountMixinFields0[0].Descriptor()
	// account.DefaultID holds the default value on creation for the id field.
	account.DefaultID = accountDescID.Default.(func() xid.ID)
	// account.IDValidator is a validator for the "id" field. It is called by the builders before save.
	account.IDValidator = func() func(string) error {
		validators := accountDescID.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(id string) error {
			for _, fn := range fns {
				if err := fn(id); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	authenticationMixin := schema.Authentication{}.Mixin()
	authenticationMixinFields0 := authenticationMixin[0].Fields()
	_ = authenticationMixinFields0
	authenticationMixinFields1 := authenticationMixin[1].Fields()
	_ = authenticationMixinFields1
	authenticationFields := schema.Authentication{}.Fields()
	_ = authenticationFields
	// authenticationDescCreatedAt is the schema descriptor for created_at field.
	authenticationDescCreatedAt := authenticationMixinFields1[0].Descriptor()
	// authentication.DefaultCreatedAt holds the default value on creation for the created_at field.
	authentication.DefaultCreatedAt = authenticationDescCreatedAt.Default.(func() time.Time)
	// authenticationDescService is the schema descriptor for service field.
	authenticationDescService := authenticationFields[0].Descriptor()
	// authentication.ServiceValidator is a validator for the "service" field. It is called by the builders before save.
	authentication.ServiceValidator = authenticationDescService.Validators[0].(func(string) error)
	// authenticationDescToken is the schema descriptor for token field.
	authenticationDescToken := authenticationFields[2].Descriptor()
	// authentication.TokenValidator is a validator for the "token" field. It is called by the builders before save.
	authentication.TokenValidator = authenticationDescToken.Validators[0].(func(string) error)
	// authenticationDescID is the schema descriptor for id field.
	authenticationDescID := authenticationMixinFields0[0].Descriptor()
	// authentication.DefaultID holds the default value on creation for the id field.
	authentication.DefaultID = authenticationDescID.Default.(func() xid.ID)
	// authentication.IDValidator is a validator for the "id" field. It is called by the builders before save.
	authentication.IDValidator = func() func(string) error {
		validators := authenticationDescID.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(id string) error {
			for _, fn := range fns {
				if err := fn(id); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	categoryMixin := schema.Category{}.Mixin()
	categoryMixinFields0 := categoryMixin[0].Fields()
	_ = categoryMixinFields0
	categoryMixinFields1 := categoryMixin[1].Fields()
	_ = categoryMixinFields1
	categoryMixinFields2 := categoryMixin[2].Fields()
	_ = categoryMixinFields2
	categoryFields := schema.Category{}.Fields()
	_ = categoryFields
	// categoryDescCreatedAt is the schema descriptor for created_at field.
	categoryDescCreatedAt := categoryMixinFields1[0].Descriptor()
	// category.DefaultCreatedAt holds the default value on creation for the created_at field.
	category.DefaultCreatedAt = categoryDescCreatedAt.Default.(func() time.Time)
	// categoryDescUpdatedAt is the schema descriptor for updated_at field.
	categoryDescUpdatedAt := categoryMixinFields2[0].Descriptor()
	// category.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	category.DefaultUpdatedAt = categoryDescUpdatedAt.Default.(func() time.Time)
	// category.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	category.UpdateDefaultUpdatedAt = categoryDescUpdatedAt.UpdateDefault.(func() time.Time)
	// categoryDescDescription is the schema descriptor for description field.
	categoryDescDescription := categoryFields[1].Descriptor()
	// category.DefaultDescription holds the default value on creation for the description field.
	category.DefaultDescription = categoryDescDescription.Default.(string)
	// categoryDescColour is the schema descriptor for colour field.
	categoryDescColour := categoryFields[2].Descriptor()
	// category.DefaultColour holds the default value on creation for the colour field.
	category.DefaultColour = categoryDescColour.Default.(string)
	// categoryDescSort is the schema descriptor for sort field.
	categoryDescSort := categoryFields[3].Descriptor()
	// category.DefaultSort holds the default value on creation for the sort field.
	category.DefaultSort = categoryDescSort.Default.(int)
	// categoryDescAdmin is the schema descriptor for admin field.
	categoryDescAdmin := categoryFields[4].Descriptor()
	// category.DefaultAdmin holds the default value on creation for the admin field.
	category.DefaultAdmin = categoryDescAdmin.Default.(bool)
	// categoryDescID is the schema descriptor for id field.
	categoryDescID := categoryMixinFields0[0].Descriptor()
	// category.DefaultID holds the default value on creation for the id field.
	category.DefaultID = categoryDescID.Default.(func() xid.ID)
	// category.IDValidator is a validator for the "id" field. It is called by the builders before save.
	category.IDValidator = func() func(string) error {
		validators := categoryDescID.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(id string) error {
			for _, fn := range fns {
				if err := fn(id); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	notificationMixin := schema.Notification{}.Mixin()
	notificationMixinFields0 := notificationMixin[0].Fields()
	_ = notificationMixinFields0
	notificationMixinFields1 := notificationMixin[1].Fields()
	_ = notificationMixinFields1
	notificationFields := schema.Notification{}.Fields()
	_ = notificationFields
	// notificationDescCreatedAt is the schema descriptor for created_at field.
	notificationDescCreatedAt := notificationMixinFields1[0].Descriptor()
	// notification.DefaultCreatedAt holds the default value on creation for the created_at field.
	notification.DefaultCreatedAt = notificationDescCreatedAt.Default.(func() time.Time)
	// notificationDescID is the schema descriptor for id field.
	notificationDescID := notificationMixinFields0[0].Descriptor()
	// notification.DefaultID holds the default value on creation for the id field.
	notification.DefaultID = notificationDescID.Default.(func() xid.ID)
	// notification.IDValidator is a validator for the "id" field. It is called by the builders before save.
	notification.IDValidator = func() func(string) error {
		validators := notificationDescID.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(id string) error {
			for _, fn := range fns {
				if err := fn(id); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	postMixin := schema.Post{}.Mixin()
	postMixinFields0 := postMixin[0].Fields()
	_ = postMixinFields0
	postMixinFields1 := postMixin[1].Fields()
	_ = postMixinFields1
	postMixinFields2 := postMixin[2].Fields()
	_ = postMixinFields2
	postFields := schema.Post{}.Fields()
	_ = postFields
	// postDescCreatedAt is the schema descriptor for created_at field.
	postDescCreatedAt := postMixinFields1[0].Descriptor()
	// post.DefaultCreatedAt holds the default value on creation for the created_at field.
	post.DefaultCreatedAt = postDescCreatedAt.Default.(func() time.Time)
	// postDescUpdatedAt is the schema descriptor for updated_at field.
	postDescUpdatedAt := postMixinFields2[0].Descriptor()
	// post.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	post.DefaultUpdatedAt = postDescUpdatedAt.Default.(func() time.Time)
	// post.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	post.UpdateDefaultUpdatedAt = postDescUpdatedAt.UpdateDefault.(func() time.Time)
	// postDescPinned is the schema descriptor for pinned field.
	postDescPinned := postFields[3].Descriptor()
	// post.DefaultPinned holds the default value on creation for the pinned field.
	post.DefaultPinned = postDescPinned.Default.(bool)
	// postDescID is the schema descriptor for id field.
	postDescID := postMixinFields0[0].Descriptor()
	// post.DefaultID holds the default value on creation for the id field.
	post.DefaultID = postDescID.Default.(func() xid.ID)
	// post.IDValidator is a validator for the "id" field. It is called by the builders before save.
	post.IDValidator = func() func(string) error {
		validators := postDescID.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(id string) error {
			for _, fn := range fns {
				if err := fn(id); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	reactMixin := schema.React{}.Mixin()
	reactMixinFields0 := reactMixin[0].Fields()
	_ = reactMixinFields0
	reactMixinFields1 := reactMixin[1].Fields()
	_ = reactMixinFields1
	reactFields := schema.React{}.Fields()
	_ = reactFields
	// reactDescCreatedAt is the schema descriptor for created_at field.
	reactDescCreatedAt := reactMixinFields1[0].Descriptor()
	// react.DefaultCreatedAt holds the default value on creation for the created_at field.
	react.DefaultCreatedAt = reactDescCreatedAt.Default.(func() time.Time)
	// reactDescID is the schema descriptor for id field.
	reactDescID := reactMixinFields0[0].Descriptor()
	// react.DefaultID holds the default value on creation for the id field.
	react.DefaultID = reactDescID.Default.(func() xid.ID)
	// react.IDValidator is a validator for the "id" field. It is called by the builders before save.
	react.IDValidator = func() func(string) error {
		validators := reactDescID.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(id string) error {
			for _, fn := range fns {
				if err := fn(id); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	roleMixin := schema.Role{}.Mixin()
	roleMixinFields0 := roleMixin[0].Fields()
	_ = roleMixinFields0
	roleMixinFields1 := roleMixin[1].Fields()
	_ = roleMixinFields1
	roleMixinFields2 := roleMixin[2].Fields()
	_ = roleMixinFields2
	roleFields := schema.Role{}.Fields()
	_ = roleFields
	// roleDescCreatedAt is the schema descriptor for created_at field.
	roleDescCreatedAt := roleMixinFields1[0].Descriptor()
	// role.DefaultCreatedAt holds the default value on creation for the created_at field.
	role.DefaultCreatedAt = roleDescCreatedAt.Default.(func() time.Time)
	// roleDescUpdatedAt is the schema descriptor for updated_at field.
	roleDescUpdatedAt := roleMixinFields2[0].Descriptor()
	// role.DefaultUpdatedAt holds the default value on creation for the updated_at field.
	role.DefaultUpdatedAt = roleDescUpdatedAt.Default.(func() time.Time)
	// role.UpdateDefaultUpdatedAt holds the default value on update for the updated_at field.
	role.UpdateDefaultUpdatedAt = roleDescUpdatedAt.UpdateDefault.(func() time.Time)
	// roleDescID is the schema descriptor for id field.
	roleDescID := roleMixinFields0[0].Descriptor()
	// role.DefaultID holds the default value on creation for the id field.
	role.DefaultID = roleDescID.Default.(func() xid.ID)
	// role.IDValidator is a validator for the "id" field. It is called by the builders before save.
	role.IDValidator = func() func(string) error {
		validators := roleDescID.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(id string) error {
			for _, fn := range fns {
				if err := fn(id); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	subscriptionMixin := schema.Subscription{}.Mixin()
	subscriptionMixinFields0 := subscriptionMixin[0].Fields()
	_ = subscriptionMixinFields0
	subscriptionMixinFields1 := subscriptionMixin[1].Fields()
	_ = subscriptionMixinFields1
	subscriptionFields := schema.Subscription{}.Fields()
	_ = subscriptionFields
	// subscriptionDescCreatedAt is the schema descriptor for created_at field.
	subscriptionDescCreatedAt := subscriptionMixinFields1[0].Descriptor()
	// subscription.DefaultCreatedAt holds the default value on creation for the created_at field.
	subscription.DefaultCreatedAt = subscriptionDescCreatedAt.Default.(func() time.Time)
	// subscriptionDescRefersType is the schema descriptor for refers_type field.
	subscriptionDescRefersType := subscriptionFields[0].Descriptor()
	// subscription.RefersTypeValidator is a validator for the "refers_type" field. It is called by the builders before save.
	subscription.RefersTypeValidator = subscriptionDescRefersType.Validators[0].(func(string) error)
	// subscriptionDescRefersTo is the schema descriptor for refers_to field.
	subscriptionDescRefersTo := subscriptionFields[1].Descriptor()
	// subscription.RefersToValidator is a validator for the "refers_to" field. It is called by the builders before save.
	subscription.RefersToValidator = subscriptionDescRefersTo.Validators[0].(func(string) error)
	// subscriptionDescID is the schema descriptor for id field.
	subscriptionDescID := subscriptionMixinFields0[0].Descriptor()
	// subscription.DefaultID holds the default value on creation for the id field.
	subscription.DefaultID = subscriptionDescID.Default.(func() xid.ID)
	// subscription.IDValidator is a validator for the "id" field. It is called by the builders before save.
	subscription.IDValidator = func() func(string) error {
		validators := subscriptionDescID.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(id string) error {
			for _, fn := range fns {
				if err := fn(id); err != nil {
					return err
				}
			}
			return nil
		}
	}()
	tagMixin := schema.Tag{}.Mixin()
	tagMixinFields0 := tagMixin[0].Fields()
	_ = tagMixinFields0
	tagMixinFields1 := tagMixin[1].Fields()
	_ = tagMixinFields1
	tagFields := schema.Tag{}.Fields()
	_ = tagFields
	// tagDescCreatedAt is the schema descriptor for created_at field.
	tagDescCreatedAt := tagMixinFields1[0].Descriptor()
	// tag.DefaultCreatedAt holds the default value on creation for the created_at field.
	tag.DefaultCreatedAt = tagDescCreatedAt.Default.(func() time.Time)
	// tagDescID is the schema descriptor for id field.
	tagDescID := tagMixinFields0[0].Descriptor()
	// tag.DefaultID holds the default value on creation for the id field.
	tag.DefaultID = tagDescID.Default.(func() xid.ID)
	// tag.IDValidator is a validator for the "id" field. It is called by the builders before save.
	tag.IDValidator = func() func(string) error {
		validators := tagDescID.Validators
		fns := [...]func(string) error{
			validators[0].(func(string) error),
			validators[1].(func(string) error),
		}
		return func(id string) error {
			for _, fn := range fns {
				if err := fn(id); err != nil {
					return err
				}
			}
			return nil
		}
	}()
}

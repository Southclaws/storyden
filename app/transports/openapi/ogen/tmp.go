package ogen

import (
	"context"

	"github.com/Southclaws/storyden/internal/openapi/ogen"
)

type Tmp struct{}

func (t Tmp) AccountsGet(ctx context.Context) (ogen.AccountsGetRes, error) {
	return nil, nil
}

func (t Tmp) AccountsGetAvatar(ctx context.Context, params ogen.AccountsGetAvatarParams) (ogen.AccountsGetAvatarRes, error) {
	return nil, nil
}

func (t Tmp) AccountsSetAvatar(ctx context.Context, req ogen.AccountsSetAvatarReq) (ogen.AccountsSetAvatarRes, error) {
	return nil, nil
}

func (t Tmp) AccountsUpdate(ctx context.Context, req ogen.OptAccountsUpdateBody) (ogen.AccountsUpdateRes, error) {
	return nil, nil
}

func (t Tmp) AuthOAuthProviderCallback(ctx context.Context, req ogen.OptAuthOAuthProviderCallbackBody, params ogen.AuthOAuthProviderCallbackParams) (ogen.AuthOAuthProviderCallbackRes, error) {
	return nil, nil
}

func (t Tmp) AuthOAuthProviderList(ctx context.Context) (ogen.AuthOAuthProviderListRes, error) {
	return nil, nil
}

func (t Tmp) AuthPasswordSignin(ctx context.Context, req ogen.AuthPasswordSigninReq) (ogen.AuthPasswordSigninRes, error) {
	return nil, nil
}

func (t Tmp) AuthPasswordSignup(ctx context.Context, req ogen.AuthPasswordSignupReq) (ogen.AuthPasswordSignupRes, error) {
	return nil, nil
}

func (t Tmp) GetSpec(ctx context.Context) (ogen.GetSpecOK, error) {
	return ogen.GetSpecOK{}, nil
}

func (t Tmp) GetVersion(ctx context.Context) (ogen.GetVersionOK, error) {
	return ogen.GetVersionOK{}, nil
}

func (t Tmp) PostsCreate(ctx context.Context, req ogen.OptPost, params ogen.PostsCreateParams) (ogen.PostsCreateRes, error) {
	return nil, nil
}

func (t Tmp) ProfilesGet(ctx context.Context, params ogen.ProfilesGetParams) (ogen.ProfilesGetRes, error) {
	return nil, nil
}

func (t Tmp) ThreadsList(ctx context.Context, params ogen.ThreadsListParams) (ogen.ThreadsListRes, error) {
	return nil, nil
}

func (t Tmp) WebAuthnGetAssertion(ctx context.Context, req ogen.WebAuthnGetAssertionReq, params ogen.WebAuthnGetAssertionParams) (ogen.WebAuthnGetAssertionRes, error) {
	return nil, nil
}

func (t Tmp) WebAuthnMakeAssertion(ctx context.Context, req ogen.WebAuthnMakeAssertionReq) (ogen.WebAuthnMakeAssertionRes, error) {
	return nil, nil
}

func (t Tmp) WebAuthnMakeCredential(ctx context.Context, req *ogen.WebAuthnMakeCredentialReq) (ogen.WebAuthnMakeCredentialRes, error) {
	return nil, nil
}

func (t Tmp) WebAuthnRequestCredential(ctx context.Context, params ogen.WebAuthnRequestCredentialParams) (ogen.WebAuthnRequestCredentialRes, error) {
	return nil, nil
}

// Code generated by ogen, DO NOT EDIT.

package ogen

import (
	"context"

	ht "github.com/ogen-go/ogen/http"
)

// UnimplementedHandler is no-op Handler which returns http.ErrNotImplemented.
type UnimplementedHandler struct{}

var _ Handler = UnimplementedHandler{}

// AccountsGet implements AccountsGet operation.
//
// Get the information for the currently authenticated account.
//
// GET /v1/accounts
func (UnimplementedHandler) AccountsGet(ctx context.Context) (r AccountsGetRes, _ error) {
	return r, ht.ErrNotImplemented
}

// AccountsGetAvatar implements AccountsGetAvatar operation.
//
// Get an avatar for the specified account.
//
// GET /v1/accounts/{account_handle}/avatar
func (UnimplementedHandler) AccountsGetAvatar(ctx context.Context, params AccountsGetAvatarParams) (r AccountsGetAvatarRes, _ error) {
	return r, ht.ErrNotImplemented
}

// AccountsSetAvatar implements AccountsSetAvatar operation.
//
// Upload an avatar for the authenticated account.
//
// POST /v1/accounts/self/avatar
func (UnimplementedHandler) AccountsSetAvatar(ctx context.Context, req AccountsSetAvatarReq) (r AccountsSetAvatarRes, _ error) {
	return r, ht.ErrNotImplemented
}

// AccountsUpdate implements AccountsUpdate operation.
//
// Update the information for the currently authenticated account.
//
// PATCH /v1/accounts
func (UnimplementedHandler) AccountsUpdate(ctx context.Context, req OptAccountsUpdateBody) (r AccountsUpdateRes, _ error) {
	return r, ht.ErrNotImplemented
}

// AuthOAuthProviderCallback implements AuthOAuthProviderCallback operation.
//
// Sign in to an existing account with a username and password.
//
// POST /v1/auth/oauth/{oauth_provider}/callback
func (UnimplementedHandler) AuthOAuthProviderCallback(ctx context.Context, req OptAuthOAuthProviderCallbackBody, params AuthOAuthProviderCallbackParams) (r AuthOAuthProviderCallbackRes, _ error) {
	return r, ht.ErrNotImplemented
}

// AuthOAuthProviderList implements AuthOAuthProviderList operation.
//
// Retrieve a list of OAuth2 providers and their links.
//
// GET /v1/auth/oauth
func (UnimplementedHandler) AuthOAuthProviderList(ctx context.Context) (r AuthOAuthProviderListRes, _ error) {
	return r, ht.ErrNotImplemented
}

// AuthPasswordSignin implements AuthPasswordSignin operation.
//
// Sign in to an existing account with a username and password.
//
// POST /v1/auth/password/signin
func (UnimplementedHandler) AuthPasswordSignin(ctx context.Context, req AuthPasswordSigninReq) (r AuthPasswordSigninRes, _ error) {
	return r, ht.ErrNotImplemented
}

// AuthPasswordSignup implements AuthPasswordSignup operation.
//
// Register a new account with a username and password.
//
// POST /v1/auth/password/signup
func (UnimplementedHandler) AuthPasswordSignup(ctx context.Context, req AuthPasswordSignupReq) (r AuthPasswordSignupRes, _ error) {
	return r, ht.ErrNotImplemented
}

// GetSpec implements GetSpec operation.
//
// Note: the generator creates a `map[string]interface{}` if this is set to
// `application/json`... so I'm just using plain text for now.
//
// GET /openapi.json
func (UnimplementedHandler) GetSpec(ctx context.Context) (r GetSpecOK, _ error) {
	return r, ht.ErrNotImplemented
}

// GetVersion implements GetVersion operation.
//
// The version number includes the date and time of the release build as
// well as a short representation of the Git commit hash.
//
// GET /version
func (UnimplementedHandler) GetVersion(ctx context.Context) (r GetVersionOK, _ error) {
	return r, ht.ErrNotImplemented
}

// PostsCreate implements PostsCreate operation.
//
// Create a new post within a thread.
//
// POST /v1/threads/{thread_id}/posts
func (UnimplementedHandler) PostsCreate(ctx context.Context, req OptPost, params PostsCreateParams) (r PostsCreateRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ProfilesGet implements ProfilesGet operation.
//
// Get a public profile by ID.
//
// GET /v1/profiles/{account_handle}
func (UnimplementedHandler) ProfilesGet(ctx context.Context, params ProfilesGetParams) (r ProfilesGetRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ThreadsCreate implements ThreadsCreate operation.
//
// Create a new thread within the specified category.
//
// POST /v1/threads
func (UnimplementedHandler) ThreadsCreate(ctx context.Context, req OptThreadsCreateBody) (r ThreadsCreateRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ThreadsGet implements ThreadsGet operation.
//
// Get information about a thread such as its title, author, when it was
// created as well as a list of the posts within the thread.
//
// GET /v1/threads/{thread_id}
func (UnimplementedHandler) ThreadsGet(ctx context.Context, params ThreadsGetParams) (r ThreadsGetRes, _ error) {
	return r, ht.ErrNotImplemented
}

// ThreadsList implements ThreadsList operation.
//
// Get a list of all threads.
//
// GET /v1/threads
func (UnimplementedHandler) ThreadsList(ctx context.Context, params ThreadsListParams) (r ThreadsListRes, _ error) {
	return r, ht.ErrNotImplemented
}

// WebAuthnGetAssertion implements WebAuthnGetAssertion operation.
//
// Start the WebAuthn assertion for an existing account.
//
// POST /v1/auth/webauthn/assert/{account_handle}
func (UnimplementedHandler) WebAuthnGetAssertion(ctx context.Context, req WebAuthnGetAssertionReq, params WebAuthnGetAssertionParams) (r WebAuthnGetAssertionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// WebAuthnMakeAssertion implements WebAuthnMakeAssertion operation.
//
// Complete the credential assertion and sign in to an account.
//
// GET /v1/auth/webauthn/assert
func (UnimplementedHandler) WebAuthnMakeAssertion(ctx context.Context, req WebAuthnMakeAssertionReq) (r WebAuthnMakeAssertionRes, _ error) {
	return r, ht.ErrNotImplemented
}

// WebAuthnMakeCredential implements WebAuthnMakeCredential operation.
//
// Complete WebAuthn registration by creating a new credential.
//
// GET /v1/auth/webauthn/make
func (UnimplementedHandler) WebAuthnMakeCredential(ctx context.Context, req *WebAuthnMakeCredentialReq) (r WebAuthnMakeCredentialRes, _ error) {
	return r, ht.ErrNotImplemented
}

// WebAuthnRequestCredential implements WebAuthnRequestCredential operation.
//
// Start the WebAuthn registration process by requesting a credential.
//
// POST /v1/auth/webauthn/make/{account_handle}
func (UnimplementedHandler) WebAuthnRequestCredential(ctx context.Context, params WebAuthnRequestCredentialParams) (r WebAuthnRequestCredentialRes, _ error) {
	return r, ht.ErrNotImplemented
}

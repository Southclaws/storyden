// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the Go project's LICENSE file.
//
// Adapted from the unmerged golang.org/x/oauth2 Dynamic Client Registration
// proposal by sunsingerus:
// https://github.com/golang/oauth2/pull/417

package oauthremote

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type dynamicClientMetadata struct {
	RedirectURIs            []string `json:"redirect_uris,omitempty"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method,omitempty"`
	GrantTypes              []string `json:"grant_types,omitempty"`
	ResponseTypes           []string `json:"response_types,omitempty"`
	ClientName              string   `json:"client_name,omitempty"`
	ClientURI               string   `json:"client_uri,omitempty"`
	Scope                   string   `json:"scope,omitempty"`
}

type dynamicClientRegistrationResponse struct {
	ClientID                string   `json:"client_id"`
	ClientSecret            string   `json:"client_secret,omitempty"`
	ClientSecretExpiresAt   int64    `json:"client_secret_expires_at,omitempty"`
	ClientIDIssuedAt        int64    `json:"client_id_issued_at,omitempty"`
	RedirectURIs            []string `json:"redirect_uris,omitempty"`
	TokenEndpointAuthMethod string   `json:"token_endpoint_auth_method,omitempty"`
	GrantTypes              []string `json:"grant_types,omitempty"`
	ResponseTypes           []string `json:"response_types,omitempty"`
	ClientName              string   `json:"client_name,omitempty"`
	ClientURI               string   `json:"client_uri,omitempty"`
	Scope                   string   `json:"scope,omitempty"`
}

func registerDynamicClient(ctx context.Context, client *http.Client, endpoint string, metadata dynamicClientMetadata) (dynamicClientRegistrationResponse, error) {
	body, err := json.Marshal(metadata)
	if err != nil {
		return dynamicClientRegistrationResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return dynamicClientRegistrationResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return dynamicClientRegistrationResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return dynamicClientRegistrationResponse{}, fmt.Errorf("dynamic client registration failed: %s", res.Status)
	}

	var out dynamicClientRegistrationResponse
	if err := json.NewDecoder(io.LimitReader(res.Body, maxMetadataBytes)).Decode(&out); err != nil {
		return dynamicClientRegistrationResponse{}, err
	}
	if out.ClientID == "" {
		return dynamicClientRegistrationResponse{}, fmt.Errorf("dynamic client registration response missing client_id")
	}

	return out, nil
}

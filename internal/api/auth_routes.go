package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"

	"github.com/Southclaws/storyden/app/services/authentication/provider/oauth/keycloak"
)

type AuthRoutes struct {
	logger           *slog.Logger
	keycloakProvider *keycloak.Provider
}

func NewAuthRoutes(
	logger *slog.Logger,
	keycloakProvider *keycloak.Provider,
) *AuthRoutes {
	return &AuthRoutes{
		logger:           logger,
		keycloakProvider: keycloakProvider,
	}
}

type bootstrapRequest struct {
	RedirectPath string `json:"redirect_path"`
}

func (ar *AuthRoutes) HandleBootstrap(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req bootstrapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ar.logger.WarnContext(ctx, "invalid bootstrap request", slog.String("error", err.Error()))
		http.Error(w, `{"error": "invalid request"}`, http.StatusBadRequest)
		return
	}

	// Default to "/" if not provided.
	redirect := req.RedirectPath
	if redirect == "" {
		redirect = "/"
	}

	// Call the Bootstrap method on Keycloak provider.
	result, err := ar.keycloakProvider.Bootstrap(redirect)
	if err != nil {
		ar.logger.ErrorContext(ctx, "failed to generate bootstrap", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, fault.Wrap(err,
			fctx.With(ctx),
			fmsg.With("failed to generate bootstrap"),
		)), http.StatusInternalServerError)
		return
	}

	ar.logger.InfoContext(ctx, "generated OIDC bootstrap", slog.String("redirect", redirect))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Mount registers the auth routes with the provided http.ServeMux.
func (ar *AuthRoutes) Mount(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/auth/oidc/bootstrap", ar.HandleBootstrap)
}

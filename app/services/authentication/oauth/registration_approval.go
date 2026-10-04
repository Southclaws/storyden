package oauth

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/ftag"
	"github.com/Southclaws/opt"

	"github.com/Southclaws/storyden/app/resources/account"
	oauthresource "github.com/Southclaws/storyden/app/resources/oauth"
	"github.com/Southclaws/storyden/app/resources/oauth/oauth_writer"
	"github.com/Southclaws/storyden/app/resources/pagination"
	"github.com/Southclaws/storyden/app/resources/resolve"
	"github.com/Southclaws/storyden/app/resources/settings"
)

const registrationPollInterval = 5

type RegistrationApprovalChallenge struct {
	RegistrationCode        string
	VerificationCode        string
	VerificationURI         string
	VerificationURIComplete string
	ExpiresIn               int
	Interval                int
}

type RegistrationApprovalReview struct {
	VerificationCode string
	CreatedAt        time.Time
	ExpiresAt        time.Time
	Metadata         DynamicClientRegistration
}

func (s *Service) startRegistrationApproval(ctx context.Context, configuration settings.OAuthServiceSettings, input DynamicClientRegistration) (*DynamicClientRegistrationResult, *Error, error) {
	ttl := configuration.RegistrationApprovalTTL.OrZero()
	if ttl < time.Second {
		return nil, oauthError("temporarily_unavailable", "Registration approval lifetime is not configured"), nil
	}

	approvalURL, err := s.registrationApprovalURL(configuration.RegistrationApprovalURL.OrZero())
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	input.RegistrationModes = nil
	metadata, err := json.Marshal(input)
	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	var registrationCode, verificationCode string
	for range userCodeGenerationAttempts {
		registrationCode, err = randomToken(32)
		if err != nil {
			return nil, nil, err
		}

		verificationCode, err = generateUserCode()
		if err != nil {
			return nil, nil, err
		}

		now := time.Now()
		_, err = s.tokens.CreateRegistrationApproval(ctx, oauth_writer.RegistrationApprovalCreate{
			RegistrationCodeHash: hashString(registrationCode),
			VerificationCodeHash: hashString(normalizeCode(verificationCode)),
			VerificationCode:     verificationCode,
			Metadata:             metadata,
			ExpiresAt:            now.Add(ttl),
			NextPollAt:           now.Add(registrationPollInterval * time.Second),
			PollIntervalSeconds:  registrationPollInterval,
		})
		if err == nil {
			break
		}

		if ftag.Get(err) != ftag.AlreadyExists {
			return nil, nil, err
		}
	}

	if err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	verificationURI := approvalURL.String()
	query := approvalURL.Query()
	query.Set("verification_code", verificationCode)
	approvalURL.RawQuery = query.Encode()

	return &DynamicClientRegistrationResult{Approval: &RegistrationApprovalChallenge{
		RegistrationCode:        registrationCode,
		VerificationCode:        verificationCode,
		VerificationURI:         verificationURI,
		VerificationURIComplete: approvalURL.String(),
		ExpiresIn:               int(ttl.Seconds()),
		Interval:                registrationPollInterval,
	}}, nil, nil
}

func (s *Service) registrationApprovalURL(address string) (*url.URL, error) {
	if address == "" {
		return resolve.URL(s.cfg.PublicWebAddress, "admin", "oauth-dcr-approval"), nil
	}

	return url.Parse(address)
}

func (s *Service) getRegistrationApproval(ctx context.Context, configuration settings.OAuthServiceSettings, code string) (*oauthresource.RegistrationApproval, *Error, error) {
	if !s.Enabled() || !configuration.DynamicRegistrationEnabled.OrZero() {
		return nil, oauthError("temporarily_unavailable", "Dynamic client registration is not enabled"), nil
	}

	rec, err := s.clients.GetRegistrationApprovalByRegistrationCodeHash(ctx, hashString(code))
	if err != nil {
		if ftag.Get(err) == ftag.NotFound {
			return nil, invalidRegistrationError(), nil
		}

		return nil, nil, err
	}

	if rec.ConsumedAt.Ok() || rec.CancelledAt.Ok() {
		return nil, invalidRegistrationError(), nil
	}

	if !rec.ExpiresAt.After(time.Now()) {
		return nil, oauthError("expired_token", "Registration has expired"), nil
	}

	if rec.DeniedAt.Ok() {
		return nil, oauthError("access_denied", "Registration was denied"), nil
	}

	return rec, nil, nil
}

func (s *Service) PollRegistrationApproval(ctx context.Context, code string) (*DynamicClientRegistrationResult, *Error, error) {
	configuration, err := s.registrationSettings(ctx)
	if err != nil {
		return nil, nil, err
	}

	rec, oauthErr, err := s.getRegistrationApproval(ctx, configuration, code)
	if oauthErr != nil || err != nil {
		return nil, oauthErr, err
	}

	if configuration.AutonomousRegistrationMode.OrZero() != settings.OAuthAutonomousRegistrationModeApproval && configuration.AutonomousRegistrationMode.OrZero() != settings.OAuthAutonomousRegistrationModeOpen {
		return nil, oauthError("access_denied", "Autonomous registration approval is no longer enabled"), nil
	}

	now := time.Now()
	interval := rec.PollIntervalSeconds
	early := now.Before(rec.NextPollAt)
	if early {
		interval += registrationPollInterval
	}

	updated, err := s.tokens.RecordRegistrationApprovalPoll(ctx, rec, now.Add(time.Duration(interval)*time.Second), interval)
	if err != nil {
		return nil, nil, err
	}

	if !updated {
		current, oauthErr, err := s.getRegistrationApproval(ctx, configuration, code)
		if oauthErr != nil || err != nil {
			return nil, oauthErr, err
		}

		return &DynamicClientRegistrationResult{RetryAfter: current.PollIntervalSeconds}, nil, nil
	}

	if early {
		return &DynamicClientRegistrationResult{RetryAfter: interval}, nil, nil
	}

	if !rec.ApprovedAt.Ok() {
		return &DynamicClientRegistrationResult{Approval: &RegistrationApprovalChallenge{}}, nil, nil
	}

	var input DynamicClientRegistration
	if err := json.Unmarshal(rec.Metadata, &input); err != nil {
		return nil, nil, fault.Wrap(err, fctx.With(ctx))
	}

	return s.provisionClientRegistration(ctx, configuration, input, opt.NewEmpty[oauthresource.DynamicRegistrationAccessTokenID](), rec)
}

func (s *Service) CancelRegistrationApproval(ctx context.Context, code string) (*Error, error) {
	configuration, err := s.registrationSettings(ctx)
	if err != nil {
		return nil, err
	}

	rec, oauthErr, err := s.getRegistrationApproval(ctx, configuration, code)
	if oauthErr != nil || err != nil {
		return oauthErr, err
	}

	cancelled, err := s.tokens.CancelRegistrationApproval(ctx, rec.ID, time.Now())
	if err != nil {
		return nil, err
	}

	if !cancelled {
		return invalidRegistrationError(), nil
	}

	return nil, nil
}

func (s *Service) getPendingRegistrationApproval(ctx context.Context, code string) (*oauthresource.RegistrationApproval, *Error, error) {
	configuration, err := s.registrationSettings(ctx)
	if err != nil {
		return nil, nil, err
	}

	if !s.Enabled() || !configuration.DynamicRegistrationEnabled.OrZero() || configuration.AutonomousRegistrationMode.OrZero() != settings.OAuthAutonomousRegistrationModeApproval {
		return nil, oauthError("access_denied", "Registration approval is not enabled"), nil
	}

	normalized, ok := parseUserCode(code)
	if !ok {
		return nil, invalidRegistrationError(), nil
	}

	rec, err := s.clients.GetRegistrationApprovalByVerificationCodeHash(ctx, hashString(normalized))
	if err != nil {
		if ftag.Get(err) == ftag.NotFound {
			return nil, invalidRegistrationError(), nil
		}

		return nil, nil, err
	}

	if !rec.ExpiresAt.After(time.Now()) || rec.ConsumedAt.Ok() || rec.CancelledAt.Ok() || rec.ApprovedAt.Ok() || rec.DeniedAt.Ok() {
		return nil, invalidRegistrationError(), nil
	}

	return rec, nil, nil
}

func (s *Service) GetRegistrationApproval(ctx context.Context, code string) (*RegistrationApprovalReview, *Error, error) {
	rec, oauthErr, err := s.getPendingRegistrationApproval(ctx, code)
	if oauthErr != nil || err != nil {
		return nil, oauthErr, err
	}

	review, err := mapRegistrationApprovalReview(rec)
	return review, nil, fault.Wrap(err, fctx.With(ctx))
}

func (s *Service) DecideRegistrationApproval(ctx context.Context, approver account.AccountID, code string, approved bool) (*Error, error) {
	rec, oauthErr, err := s.getPendingRegistrationApproval(ctx, code)
	if oauthErr != nil || err != nil {
		return oauthErr, err
	}

	updated, err := s.tokens.DecideRegistrationApproval(ctx, rec.ID, approver, approved, time.Now())
	if err != nil {
		return nil, err
	}

	if !updated {
		return invalidRegistrationError(), nil
	}

	return nil, nil
}

func invalidRegistrationError() *Error {
	return oauthError("invalid_registration", "Registration is invalid or unavailable")
}

func (s *Service) ListRegistrationApprovals(ctx context.Context, params pagination.Parameters) (*pagination.Result[*RegistrationApprovalReview], error) {
	records, err := s.clients.ListPendingRegistrationApprovals(ctx, params)
	if err != nil {
		return nil, err
	}

	reviews := make([]*RegistrationApprovalReview, 0, len(records.Items))
	for _, rec := range records.Items {
		review, err := mapRegistrationApprovalReview(rec)
		if err != nil {
			return nil, fault.Wrap(err, fctx.With(ctx))
		}

		reviews = append(reviews, review)
	}

	result := pagination.ConvertPageResult(*records, reviews)
	return &result, nil
}

func mapRegistrationApprovalReview(rec *oauthresource.RegistrationApproval) (*RegistrationApprovalReview, error) {
	var metadata DynamicClientRegistration
	if err := json.Unmarshal(rec.Metadata, &metadata); err != nil {
		return nil, err
	}

	return &RegistrationApprovalReview{
		VerificationCode: rec.VerificationCode, CreatedAt: rec.CreatedAt,
		ExpiresAt: rec.ExpiresAt, Metadata: metadata,
	}, nil
}

func (s *Service) DecideAllRegistrationApprovals(ctx context.Context, approver account.AccountID, approved bool, before time.Time) (int, error) {
	now := time.Now()
	if before.IsZero() || before.After(now) {
		return 0, fault.New("created_before must be a past snapshot time", ftag.With(ftag.InvalidArgument))
	}

	return s.tokens.DecideAllRegistrationApprovals(ctx, approver, approved, before.In(now.Location()), now)
}

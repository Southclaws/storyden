package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Southclaws/fault"
)

type userProfile struct {
	ID    int64   `json:"id"`
	Login *string `json:"login"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

func fetchProfile(ctx context.Context, client *http.Client, token string) (*userProfile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return nil, fault.Wrap(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "Storyden")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fault.Wrap(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fault.New(fmt.Sprintf("GitHub profile request returned HTTP %d", resp.StatusCode))
	}
	const maxProfileBytes = 1 << 20
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxProfileBytes+1))
	if err != nil {
		return nil, fault.Wrap(err)
	}
	if len(body) > maxProfileBytes {
		return nil, fault.New("GitHub profile response exceeds size limit")
	}
	var profile userProfile
	if err := json.Unmarshal(body, &profile); err != nil {
		return nil, fault.Wrap(err)
	}
	if profile.ID <= 0 || profile.Login == nil || strings.TrimSpace(*profile.Login) == "" {
		return nil, fault.New("GitHub profile is missing a valid ID or login")
	}
	return &profile, nil
}

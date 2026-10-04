package config

import "time"

type Registration struct {
	State                   string    `yaml:"state" json:"state"`
	Handle                  string    `yaml:"handle" json:"handle"`
	Endpoint                string    `yaml:"endpoint" json:"-"`
	Code                    string    `yaml:"code,omitempty" json:"-"`
	VerificationCode        string    `yaml:"verification_code,omitempty" json:"verification_code,omitempty"`
	VerificationURI         string    `yaml:"verification_uri,omitempty" json:"verification_uri,omitempty"`
	VerificationURIComplete string    `yaml:"verification_uri_complete,omitempty" json:"verification_uri_complete,omitempty"`
	ExpiresAt               time.Time `yaml:"expires_at,omitempty" json:"expires_at,omitzero"`
	NextPollAt              time.Time `yaml:"next_poll_at,omitempty" json:"next_poll_at,omitzero"`
	Interval                int       `yaml:"interval,omitempty" json:"interval,omitempty"`
}

package password_reset

import (
	"net/url"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"
)

var ErrResetURLHostMismatch = fault.New("password reset URL host does not match the configured web address", ftag.With(ftag.InvalidArgument))

type LinkTemplate struct {
	u  url.URL
	qp string
}

func (r *LinkTemplate) GetURL(token string) string {
	q := r.u.Query()

	q.Add(r.qp, token)

	r.u.RawQuery = q.Encode()

	return r.u.String()
}

func NewLinkTemplate(urlString string, tokenQueryParam string, allowedHost string) (*LinkTemplate, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	// reset token must only ever land on our own host, never an attacker's
	if u.Host == "" || !strings.EqualFold(u.Host, allowedHost) {
		return nil, fault.Wrap(ErrResetURLHostMismatch,
			fmsg.WithDesc("host mismatch", "The password reset link must point to this site."))
	}

	return &LinkTemplate{
		u:  *u,
		qp: tokenQueryParam,
	}, nil
}
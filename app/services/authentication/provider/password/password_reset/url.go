package password_reset

import "net/url"

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

func NewLinkTemplate(urlString string, tokenQueryParam string) (*LinkTemplate, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	return &LinkTemplate{
		u:  *u,
		qp: tokenQueryParam,
	}, nil
}

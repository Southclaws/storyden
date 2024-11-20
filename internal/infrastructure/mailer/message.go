package mailer

import (
	"net/mail"

	"github.com/Southclaws/fault"
)

var (
	ErrNoSubject = fault.New("no subject provided")
	ErrNoContent = fault.New("no content provided")
)

type Content struct {
	HTML  string
	Plain string
}

func NewContent(html string, plain string) (*Content, error) {
	if html == "" && plain == "" {
		return nil, ErrNoContent
	}

	return &Content{
		HTML:  html,
		Plain: plain,
	}, nil
}

type Message struct {
	Address mail.Address
	Name    string
	Subject string
	Content Content
}

func NewMessage(
	address mail.Address,
	name string,
	subject string,
	content Content,
) (*Message, error) {
	if name == "" {
		name = address.Name
	}

	if subject == "" {
		return nil, ErrNoSubject
	}

	return &Message{
		Address: address,
		Name:    name,
		Subject: subject,
		Content: content,
	}, nil
}

package mail

import (
	"context"
	"encoding/json"

	"github.com/cenkalti/backoff"
	"github.com/pkg/errors"
	"go.uber.org/fx"

	"github.com/Southclaws/storyden/api/src/infra/mailer"
	"github.com/Southclaws/storyden/api/src/infra/pubsub"
	"github.com/Southclaws/storyden/api/src/resources/mail/template"
)

type Worker struct {
	t pubsub.Topic
	b pubsub.Bus
	m mailer.Mailer
	r template.Registry
}

func New(lc fx.Lifecycle, b pubsub.Bus, m mailer.Mailer, r template.Registry) *Worker {
	w := &Worker{b.Declare("system.email"), b, m, r}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := backoff.Retry(w.run, backoff.NewExponentialBackOff()); err != nil {
					panic(err)
				}
			}()
			return nil
		},
	})

	return w
}

type Message struct {
	Name     string
	Addr     string
	Subj     string
	Template template.ID
	Data     interface{}
}

func (w *Worker) Enqueue(name, addr, subj string, template template.ID, data interface{}) error {
	body, err := json.Marshal(Message{
		Name:     name,
		Addr:     addr,
		Subj:     subj,
		Template: template,
		Data:     data,
	})
	if err != nil {
		return err
	}
	return w.b.Publish(w.t, body)
}

func (w *Worker) run() error {
	return w.b.Subscribe(w.t, func(body []byte) (bool, error) {
		var message Message
		if err := json.Unmarshal(body, &message); err != nil {
			return true, errors.Wrap(err, "unexpected message in mailer topic")
		}

		t, err := w.r.Get(message.Template, message.Data)
		if err != nil {
			return false, errors.Wrap(err, "mailworker failed to format email")
		}

		if err := w.m.Mail(message.Name, message.Addr, message.Subj, t.Rich, t.Text); err != nil {
			return false, errors.Wrap(err, "mailworker failed to send email")
		}

		return true, nil
	})
}

package openai

import (
	"errors"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/openai/openai-go/v3"
)

var errOpenAIAPI = fault.New("openai api error")

func mapError(err error) error {
	var oe *openai.Error
	if errors.As(err, &oe) {
		return fault.Wrap(errOpenAIAPI,
			fmsg.WithDesc(oe.Type, oe.Message),
		)
	}

	return err
}

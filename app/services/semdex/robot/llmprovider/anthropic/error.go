package anthropic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/anthropics/anthropic-sdk-go"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
)

var errAnthropicAPI = fault.New("anthropic api error")

func mapError(err error) error {
	var ae *anthropic.Error
	if errors.As(err, &ae) {
		var res anthropic.ErrorResponse
		if json.Unmarshal([]byte(ae.RawJSON()), &res) == nil && res.Error.Message != "" {
			return fault.Wrap(errAnthropicAPI,
				fmsg.WithDesc(res.Error.Type, res.Error.Message),
			)
		}

		title := string(ae.Type())
		if title == "" {
			title = "api_error"
		}

		message := http.StatusText(ae.StatusCode)
		if message == "" {
			message = "request failed"
		}

		return fault.Wrap(errAnthropicAPI,
			fmsg.WithDesc(title, fmt.Sprintf("%s (status %d)", message, ae.StatusCode)),
		)
	}

	return err
}

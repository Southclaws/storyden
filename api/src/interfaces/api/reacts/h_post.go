package reacts

import (
	"errors"
	"net/http"

	"github.com/Southclaws/storyden/api/src/infra/web"
	"github.com/Southclaws/storyden/api/src/resources/react"
	"github.com/Southclaws/storyden/api/src/services/authentication"
)

type PostBody struct {
	PostID string `json:"postId"` // ID of the post to react to
	Emoji  string `json:"emoji"`  // A string containing a single emoji character
}

// @Summary  React to a post
// @Tags     reacts
// @Accept   json
// @Produce  json
// @Param    payload  body      PostBody  true  "Post ID and emoji to react with"
// @Success  200      {object}  react.React
// @Failure  400      {object}  web.Error
// @Failure  500      {object}  web.Error
// @Failure  404      {object}  web.Error
// @Router   /reacts [post]
func (c *controller) post(w http.ResponseWriter, r *http.Request) {
	info, ok := authentication.GetAuthenticationInfo(w, r)
	if !ok {
		return
	}

	var b PostBody
	if err := web.DecodeBody(r, &b); err != nil {
		web.StatusBadRequest(w, err)
		return
	}

	reaction, err := c.reacts.Add(r.Context(), info.Cookie.UserID, b.PostID, b.Emoji)
	if err != nil {
		if errors.Is(err, react.ErrAlreadyReacted) {
			web.Write(w, b)
			return
		} else if errors.Is(err, react.ErrInvalidEmoji) {
			web.StatusBadRequest(w, web.WithSuggestion(err,
				"The emoji ID sent to the server was invalid.",
				"This may be an issue with the emoji picker menu, please report this to a site administrator and include the exact emoji you tried to react with."))
			return
		} else {
			web.StatusInternalServerError(w, err)
		}
		return
	}
	if reaction == nil {
		web.StatusNotFound(w, nil)
		return
	}

	web.Write(w, reaction)
}

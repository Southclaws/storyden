package subscription

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/puzpuzpuz/xsync"
	"go.uber.org/zap"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/storyden/app/resources/account"
	"github.com/Southclaws/storyden/app/services/authentication/session"
	"github.com/Southclaws/storyden/internal/config"
)

type Handler struct {
	logger         *zap.Logger
	originPatterns []string
	subscriptions  *xsync.MapOf[account.AccountID, []Channel]
}

func New(cfg config.Config, log *zap.Logger) *Handler {
	return &Handler{
		logger:         log,
		originPatterns: []string{cfg.PublicWebAddress.Host},
		subscriptions: xsync.NewTypedMapOf[account.AccountID, []Channel](func(ai account.AccountID) uint64 {
			return xsync.StrHash64(ai.String())
		}),
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	err := h.handle(r.Context(), w, r)
	if err != nil {
		h.logger.Error("failed to handle websocket request", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	accountID, err := session.GetAccountID(r.Context())
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx))
	}

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: h.originPatterns,
	})
	if err != nil {
		return fault.Wrap(err, fctx.With(ctx), fmsg.With("failed to accept websocket connection"))
	}
	defer func() {
		if err := c.CloseNow(); err != nil {
			h.logger.Error("failed to close websocket connection", zap.Error(err))
		}
	}()

	channels, ok := h.subscriptions.Load(accountID)
	if !ok {
		channels = []Channel{}
	}

	// Set the context as needed. Use of r.Context() is not recommended
	// to avoid surprising behavior (see http.Hijacker).
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	fmt.Println(channels)

	for {

		// err = wsjson.Write(ctx, c, v)
		// if err != nil {
		// 	return err
		// }

	}
}

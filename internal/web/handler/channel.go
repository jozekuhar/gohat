package handler

import (
	"log/slog"
	"net/http"

	"mimokocke/internal/channel"
	"mimokocke/internal/tenant"
	"mimokocke/internal/web/view"
)

type channelHandler struct {
	logger      *slog.Logger
	channelSrv  *channel.Service
	channelView *view.Channel
}

func NewChanneHandler(logger *slog.Logger, channelSrv *channel.Service) *channelHandler {
	return &channelHandler{
		logger:      logger,
		channelSrv:  channelSrv,
		channelView: view.NewChannel(),
	}
}

func (h *channelHandler) GetChannels(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	// channels, err := h.channelSrv.GetChannels(r.Context(), identity)
	// if err != nil {
	// 	h.logger.Error("getting channels", "err", err)
	// 	return
	// }

	render(w, h.channelView.ChannelsPage(identity, nil), http.StatusOK)
}

func (h *channelHandler) PostCreateChannel(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

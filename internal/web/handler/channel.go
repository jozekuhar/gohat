package handler

import (
	"log/slog"
	"net/http"

	"mimokocke/internal/channel"
	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/identity"
	"mimokocke/internal/web/view"

	"github.com/go-playground/form"
)

type channelHandler struct {
	logger      *slog.Logger
	formDecoder *form.Decoder
	channelSrv  *channel.Service
	view        *view.View
}

func NewChanneHandler(
	logger *slog.Logger,
	channelSrv *channel.Service,
	formDecoder *form.Decoder,
) *channelHandler {
	return &channelHandler{
		logger:      logger,
		channelSrv:  channelSrv,
		formDecoder: formDecoder,
		view:        view.NewView(),
	}
}

func (h *channelHandler) GetChannels(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustFromContext(r.Context())

	channels, err := h.channelSrv.GetChannelsOverview(r.Context(), ident)
	if err != nil {
		h.logger.Error("getting channel overview", "err", err)
		render(w, h.view.Toast.Fragment("Something went wrong"), http.StatusInternalServerError)
		return
	}

	render(w, h.view.Channel.ChannelsPage(ident, channels), http.StatusOK)
}

func (h *channelHandler) GetCreateChannelFormModal(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustFromContext(r.Context())

	render(w, h.view.Channel.CreateChannelFormModal(ident), http.StatusOK)
}

func (h *channelHandler) PostCreateChannel(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustFromContext(r.Context())

	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parsing form", "err", err)
		render(w, h.view.Toast.Fragment("Something went wrong"), http.StatusInternalServerError)
		return
	}

	provider := r.PostFormValue("Provider")

	switch provider {
	case db.ChannelProviderWooCommerce:
		params := channel.SaveWooCommerceChannelParams{}

		err := h.formDecoder.Decode(&params, r.Form)
		if err != nil {
			return
		}

		err = h.channelSrv.SaveWooCommerceChannel(r.Context(), ident, params)
		if err != nil {
			return
		}

		return
	case db.ChannelProviderShopify:
		params := channel.SaveShopifyChannelParams{}

		err := h.formDecoder.Decode(&params, r.Form)
		if err != nil {
			return
		}

		err = h.channelSrv.SaveShopifyChannel(r.Context(), ident, params)
		if err != nil {
			return
		}

		return
	}
}

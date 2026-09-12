package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"mimokocke/internal/channel"
	"mimokocke/internal/tenant"
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
	identity := tenant.MustIdentityFromContext(r.Context())

	render(w, h.view.Channel.ChannelsPage(identity), http.StatusOK)
}

func (h *channelHandler) GetCreateChannelFormModal(w http.ResponseWriter, r *http.Request) {
	identity := tenant.MustIdentityFromContext(r.Context())

	render(w, h.view.Channel.CreateChannelFormModal(identity), http.StatusOK)
}

func (h *channelHandler) PostCreateChannel(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		render(w, h.view.Toast.Fragment("Something went wrong"), http.StatusInternalServerError)
		return
	}

	form := struct{}{}

	err = h.formDecoder.Decode(form, r.Form)
	if err != nil {
		render(w, h.view.Toast.Fragment("Something went wrong"), http.StatusInternalServerError)
		return
	}

	fmt.Println("create channel")
}

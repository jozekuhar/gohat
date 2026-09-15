package channel

import (
	"errors"
	"log/slog"
	"net/http"

	"mimokocke/internal/provider/db"
	"mimokocke/internal/shared/identity"
	"mimokocke/internal/shared/routes"
	"mimokocke/internal/web/components"
	"mimokocke/internal/web/request"
	"mimokocke/internal/web/response"

	"github.com/go-playground/form"
)

type Handler struct {
	logger  *slog.Logger
	decoder *form.Decoder
	srv     *Service
}

func NewHandler(logger *slog.Logger, decoder *form.Decoder, srv *Service) *Handler {
	return &Handler{
		logger:  logger,
		decoder: decoder,
		srv:     srv,
	}
}

func (h *Handler) GetChannels(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustFromContext(r.Context())

	channels, err := h.srv.GetChannelsOverview(r.Context(), ident)
	if err != nil {
		h.logger.Error("getting channel overview", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	response.RenderStatus(w, channelsPage(ident, channels), http.StatusOK)
}

func (h *Handler) GetCreateChannelFormModal(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustFromContext(r.Context())

	response.RenderStatus(w, createChannelFormModal(ident), http.StatusOK)
}

func (h *Handler) PostCreateChannel(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) PostCreateWooCommerceChannel(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustFromContext(r.Context())

	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parsing woo form", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	var form struct {
		Name           string
		StoreURL       string
		ConsumerKey    string
		ConsumerSecret string
	}

	err = h.decoder.Decode(&form, r.Form)
	if err != nil {
		h.logger.Error("decode woo form", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	chn, err := h.srv.SaveWooCommerceChannel(
		r.Context(),
		ident,
		SaveWooCommerceChannelParams{
			Name: form.Name,
			Credentials: WooCommerceCredentials{
				StoreURL:       form.StoreURL,
				ConsumerKey:    form.ConsumerKey,
				ConsumerSecret: form.ConsumerSecret,
			},
		},
	)
	if errors.Is(err, db.ErrAlreadyExists) {
		response.RenderStatus(
			w,
			components.ToastFragment("Channel with this name already exists."),
			http.StatusConflict,
		)
		return
	}
	if err != nil {
		h.logger.Error("save woo channel", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	response.HXTrigger(w, components.EventModalClose)
	response.Render(w, channelItem(ident, chn))
	response.RenderStatus(w, components.ToastFragment("Channel created"), http.StatusCreated)
}

func (h *Handler) PostTestWooCommerceChannel(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		h.logger.Error("parse woo form", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	var form struct {
		StoreURL       string
		ConsumerKey    string
		ConsumerSecret string
	}

	err = h.decoder.Decode(&form, r.Form)
	if err != nil {
		h.logger.Error("decode woo form", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	err = h.srv.VerifyWooCommerceCredentials(r.Context(), WooCommerceCredentials{
		StoreURL:       form.StoreURL,
		ConsumerKey:    form.ConsumerKey,
		ConsumerSecret: form.ConsumerSecret,
	})
	if err != nil {
		response.RenderStatus(
			w,
			components.ToastFragment("Test unsuccescfull"),
			http.StatusBadRequest,
		)
		return
	}

	response.RenderStatus(w, components.ToastFragment("Connection test succesfull"), http.StatusOK)
}

func (h *Handler) GetUpdateChannelModalForm(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustFromContext(r.Context())

	chnID, err := request.PathValueUUID(r, routes.PathChannelID)
	if err != nil {
		h.logger.Error("path value channelID", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	chn, err := h.srv.GetChannelDetails(r.Context(), ident, chnID)
	if err != nil {
		h.logger.Error("get channel details", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	response.RenderStatus(w, updateChannelFormModal(ident, chn), http.StatusOK)
}

func (h *Handler) PostDeactivateChannel(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func (h *Handler) DeleteRemoveChannel(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustFromContext(r.Context())

	channelID, err := request.PathValueUUID(r, routes.PathChannelID)
	if err != nil {
		h.logger.Error("path value channelID", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	err = h.srv.RemoveChannel(r.Context(), ident, channelID)
	if err != nil {
		h.logger.Error("get channel details", "err", err)
		response.RenderStatus(
			w,
			components.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	response.HXTrigger(w, components.EventModalClose)
	response.RenderStatus(w, components.ToastFragment("Channel deleted succesfully"), http.StatusOK)
}

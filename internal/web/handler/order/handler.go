package orderweb

import (
	"fmt"
	"log/slog"
	"net/http"

	"mimokocke/internal/domain/channel"
	"mimokocke/internal/domain/courier"
	"mimokocke/internal/shared/identity"
	orderuc "mimokocke/internal/usecase/order"
	"mimokocke/internal/web/response"
	"mimokocke/internal/web/view"
)

type handler struct {
	listOrdersUseCase *orderuc.ListOrdersUseCase
	logger            *slog.Logger
}

func NewHandler(
	listOrdersUseCase *orderuc.ListOrdersUseCase,
	channelSrv *channel.Service,
	courierSrv *courier.Service,
	logger *slog.Logger,
) *handler {
	return &handler{
		listOrdersUseCase: listOrdersUseCase,
		logger:            logger,
	}
}

func (h *handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustIdentityFromContext(r.Context())

	orders, err := h.listOrdersUseCase.Execute(r.Context(), ident.OrgID)
	if err != nil {
		h.logger.Error("list orders")
		response.RenderStatus(
			w,
			view.ToastFragment("Something went wrong"),
			http.StatusInternalServerError,
		)
		return
	}

	response.RenderStatus(w, ordersPage(ident, orders), http.StatusOK)
}

func (h *handler) GetCreateOrderFormModal(w http.ResponseWriter, r *http.Request) {
	response.RenderStatus(w, createOrderFormModal(), http.StatusOK)
}

func (h *handler) PostCreateOrder(w http.ResponseWriter, r *http.Request) {
	ident := identity.MustIdentityFromContext(r.Context())
	_ = ident

	fmt.Println("order created")
}

package main

import (
	"errors"
	"net/http"

	"github.com/juxue97/common"
	pb "github.com/juxue97/common/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	// gateway - service discovery
	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *handler {
	return &handler{client: client}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.HandleCreateOrder)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return errors.New("items cannot be empty")
	}

	for _, i := range items {
		if i.ID == "" {
			return errors.New("item id is required")
		}
		if i.Quantity <= 0 {
			return errors.New("item quantity must be greater than 0")
		}
	}

	return nil
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(w, r, &items); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	if err := validateItems(items); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	o, err := h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})
	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.InvalidArgument {
			common.BadRequestResponse(w, r, errors.New(rStatus.Message()))
			return
		}
	}
	if err != nil {
		common.InternalServerError(w, r, err)
		return
	}
	if err := common.WriteJSON(w, http.StatusCreated, o); err != nil {
		common.InternalServerError(w, r, err)
	}
}

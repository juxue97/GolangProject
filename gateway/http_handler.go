package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/juxue97/common"
	pb "github.com/juxue97/common/api"
	"github.com/juxue97/gateway/gateway"
	"go.opentelemetry.io/otel"
	otelCodes "go.opentelemetry.io/otel/codes"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	// gateway - service discovery
	gateway gateway.OrdersGateway
}

func NewHandler(gateway gateway.OrdersGateway) *handler {
	return &handler{gateway: gateway}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.Handle("/", http.FileServer(http.Dir("public")))

	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.handleCreateOrder)
	mux.HandleFunc("GET /api/customers/{customerID}/orders/{orderID}", h.handleGetOrder)
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

func (h *handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(w, r, &items); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	if err := validateItems(items); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}

	o, err := h.gateway.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})
	rStatus := status.Convert(err)
	if rStatus != nil {
		span.SetStatus(otelCodes.Error, err.Error())

		if rStatus.Code() != codes.InvalidArgument {
			common.BadRequestResponse(w, r, errors.New(rStatus.Message()))
			return
		}
		common.InternalServerError(w, r, err)
		return
	}

	res := createOrderResponse{
		Order:         o,
		RedirectToUrl: fmt.Sprintf("http://localhost:8080/success.html?customerID=%s&orderID=%s", o.CustomerID, o.ID),
	}

	if err := common.WriteJSON(w, http.StatusCreated, res); err != nil {
		common.InternalServerError(w, r, err)
	}
}

func (h *handler) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	orderID := r.PathValue("orderID")

	tr := otel.Tracer("http")
	ctx, span := tr.Start(r.Context(), fmt.Sprintf("%s %s", r.Method, r.RequestURI))
	defer span.End()

	o, err := h.gateway.GetOrder(ctx, orderID, customerID)
	rStatus := status.Convert(err)
	if rStatus != nil {
		span.SetStatus(otelCodes.Error, err.Error())

		if rStatus.Code() != codes.InvalidArgument {
			common.BadRequestResponse(w, r, errors.New(rStatus.Message()))
			return
		}
		common.InternalServerError(w, r, err)
		return
	}

	if err := common.WriteJSON(w, http.StatusOK, o); err != nil {
		common.InternalServerError(w, r, err)
	}
}

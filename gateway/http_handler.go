package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/juxue97/common"
	pb "github.com/juxue97/common/api"
	"github.com/juxue97/gateway/gateway"
	"go.opentelemetry.io/otel"
	otelCodes "go.opentelemetry.io/otel/codes"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var Validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

type handler struct {
	// gateway - service discovery
	ordersGateway gateway.OrdersGateway
	stocksGateway gateway.StocksGateway
}

func NewHandler(ordersGateway gateway.OrdersGateway, stocksGateway gateway.StocksGateway) *handler {
	return &handler{
		ordersGateway: ordersGateway,
		stocksGateway: stocksGateway,
	}
}

func (h *handler) registerRoutes(mux *http.ServeMux) {
	mux.Handle("/", http.FileServer(http.Dir("public")))

	mux.HandleFunc("POST /api/customers/{customerID}/orders", h.handleCreateOrder)
	mux.HandleFunc("GET /api/customers/{customerID}/orders/{orderID}", h.handleGetOrder)

	mux.HandleFunc("POST /stocks", h.handleCreateItem)
	mux.HandleFunc("GET /stocks", h.handleGetItems)
	mux.HandleFunc("GET /stocks/{id}", h.handleGetItem)
	mux.HandleFunc("PUT /stocks/{id}", h.handleUpdateItem)
	mux.HandleFunc("PUT /stocks/{id}/{quantity}", h.handleUpdateStock)
	mux.HandleFunc("DELETE /stocks/{id}", h.handleDeleteItem)
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

	o, err := h.ordersGateway.CreateOrder(ctx, &pb.CreateOrderRequest{
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

	o, err := h.ordersGateway.GetOrder(ctx, orderID, customerID)
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

func (h *handler) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	// call the service layer
	var payload *pb.CreateItemRequest

	if err := common.ReadJSON(w, r, &payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		common.UnprocessableEntityResponse(w, r, err)
		return
	}

	ctx := r.Context()

	oID, err := h.stocksGateway.CreateItem(ctx, payload)
	if err != nil {
		switch err {
		case common.ErrNoQuantity:
			common.BadRequestResponse(w, r, err)

		case common.ErrConvertID:
			common.InternalServerError(w, r, err)

		default:
			common.InternalServerError(w, r, err)
		}
		return
	}

	if err := common.WriteJSON(w, http.StatusCreated, oID); err != nil {
		common.InternalServerError(w, r, err)
		return
	}
}

func (h *handler) handleGetItems(w http.ResponseWriter, r *http.Request) {
	// no payload, extract all records from db
	ctx := r.Context()

	items, err := h.stocksGateway.GetItems(ctx)
	if err != nil {
		switch err {
		case common.ErrNoDoc:
			common.NotFoundError(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}

	if err := common.WriteJSON(w, http.StatusOK, items); err != nil {
		common.InternalServerError(w, r, err)
		return
	}
}

func (h *handler) handleGetItem(w http.ResponseWriter, r *http.Request) {
	oID := r.PathValue("id")
	ctx := r.Context()

	item, err := h.stocksGateway.GetItem(ctx, oID)
	if err != nil {
		switch err {
		case common.ErrNoDoc:
			common.NotFoundError(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}

	if err := common.WriteJSON(w, http.StatusOK, item); err != nil {
		common.InternalServerError(w, r, err)
		return
	}
}

func (h *handler) handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	oID := r.PathValue("id")
	ctx := r.Context()

	var payload *pb.UpdateStockItemRequest
	if err := common.ReadJSON(w, r, &payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		common.UnprocessableEntityResponse(w, r, err)
		return
	}

	item, err := h.stocksGateway.UpdateItem(ctx, oID, payload)
	if err != nil {
		switch err {
		case common.ErrNoDoc:
			common.NotFoundError(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}

	if err := common.WriteJSON(w, http.StatusOK, item); err != nil {
		common.InternalServerError(w, r, err)
		return
	}
}

func (h *handler) handleUpdateStock(w http.ResponseWriter, r *http.Request) {
	oID := r.PathValue("id")
	quantity, _ := strconv.Atoi(r.PathValue("quantity"))

	ctx := r.Context()

	item, err := h.stocksGateway.UpdateStockQuantity(ctx, oID, quantity)
	if err != nil {
		switch err {
		case common.ErrNoDoc:
			common.NotFoundError(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}
	if err := common.WriteJSON(w, http.StatusOK, item); err != nil {
		common.InternalServerError(w, r, err)
	}
}

func (h *handler) handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	oID := r.PathValue("id")
	ctx := r.Context()

	if err := h.stocksGateway.DeleteItem(ctx, oID); err != nil {
		switch err {
		case common.ErrNoDoc:
			common.NotFoundError(w, r, err)
		default:
			common.InternalServerError(w, r, err)
		}
		return
	}

	if err := common.WriteJSON(w, http.StatusNoContent, ""); err != nil {
		common.InternalServerError(w, r, err)
		return
	}
}

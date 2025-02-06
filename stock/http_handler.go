package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/juxue97/common"
)

var Validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

type stockHandler struct {
	// service add here
	service *stockService
}

func NewStockHandler(service *stockService) *stockHandler {
	return &stockHandler{service: service}
}

func (h *stockHandler) registerRouters(mux *http.ServeMux) {
	mux.HandleFunc("POST /stocks", h.handleCreateItem)
	mux.HandleFunc("GET /stocks", h.handleGetItems)
	mux.HandleFunc("GET /stocks/{id}", h.handleGetItem)
	mux.HandleFunc("PUT /stocks/{id}", h.handleUpdateItem)
	mux.HandleFunc("DELETE /stocks/{id}", h.handleDeleteItem)
}

func (h *stockHandler) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	// call the service layer
	var payload CreateItemRequest

	if err := common.ReadJSON(w, r, &payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		common.UnprocessableEntityResponse(w, r, err)
		return
	}

	ctx := r.Context()

	oID, err := h.service.CreateItem(ctx, &payload)
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

	if err := common.WriteJSON(w, http.StatusCreated, &CreateItemResponse{
		ObjectID: oID,
	}); err != nil {
		common.InternalServerError(w, r, err)
		return
	}
}

func (h *stockHandler) handleGetItems(w http.ResponseWriter, r *http.Request) {
}

func (h *stockHandler) handleGetItem(w http.ResponseWriter, r *http.Request) {
}

func (h *stockHandler) handleUpdateItem(w http.ResponseWriter, r *http.Request) {
}

func (h *stockHandler) handleDeleteItem(w http.ResponseWriter, r *http.Request) {
}

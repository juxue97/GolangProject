package main

import (
	"net/http"
	"strconv"

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
	mux.HandleFunc("PUT /stocks/{id}/{quantity}", h.handleUpdateStock)
	mux.HandleFunc("DELETE /stocks/{id}", h.handleDeleteItem)
}

func (h *stockHandler) handleCreateItem(w http.ResponseWriter, r *http.Request) {
	// call the service layer
	var payload *CreateItemRequest

	if err := common.ReadJSON(w, r, &payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		common.UnprocessableEntityResponse(w, r, err)
		return
	}

	ctx := r.Context()

	oID, err := h.service.CreateItem(ctx, payload)
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
	// no payload, extract all records from db
	ctx := r.Context()

	items, err := h.service.GetItems(ctx)
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

func (h *stockHandler) handleGetItem(w http.ResponseWriter, r *http.Request) {
	oID := r.PathValue("id")
	ctx := r.Context()

	item, err := h.service.GetItem(ctx, oID)
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

func (h *stockHandler) handleUpdateItem(w http.ResponseWriter, r *http.Request) {
	oID := r.PathValue("id")
	ctx := r.Context()

	var payload *UpdateItemRequest
	if err := common.ReadJSON(w, r, &payload); err != nil {
		common.BadRequestResponse(w, r, err)
		return
	}
	if err := Validate.Struct(payload); err != nil {
		common.UnprocessableEntityResponse(w, r, err)
		return
	}

	item, err := h.service.UpdateItem(ctx, oID, payload)
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

func (h *stockHandler) handleUpdateStock(w http.ResponseWriter, r *http.Request) {
	oID := r.PathValue("id")
	quantity, _ := strconv.Atoi(r.PathValue("quantity"))

	ctx := r.Context()

	item, err := h.service.UpdateStock(ctx, oID, quantity)
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

func (h *stockHandler) handleDeleteItem(w http.ResponseWriter, r *http.Request) {
	oID := r.PathValue("id")
	ctx := r.Context()

	if err := h.service.DeleteItem(ctx, oID); err != nil {
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

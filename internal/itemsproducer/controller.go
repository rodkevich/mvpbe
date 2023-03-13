package itemsproducer

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/rodkevich/mvpbe/internal/itemsproducer/model"
	"github.com/rodkevich/mvpbe/pkg/validate"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

// Handler ...
type Handler struct {
	items    ItemsSampleUsage
	validate *validator.Validate
}

// NewItemsHandler ...
func NewItemsHandler(cmd ItemsSampleUsage) *Handler {
	return &Handler{
		items:    cmd,
		validate: validate.New(),
	}
}

// GetItemHandler render an item by id
func (h *Handler) GetItemHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "itemID")
		if id == "" {
			api.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		itemID, err := strconv.Atoi(id)
		if err != nil {
			api.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			log.Println("strconv.Atoi: ", err)
			return
		}

		item, err := h.items.GetOne(r.Context(), itemID)
		if err != nil {
			if errors.Is(err, errDeletedItem) {
				api.Error(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
				log.Printf("deleted item requested")
				return
			}

			api.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			log.Println("usecase.GetOne error:", err)
			return
		}

		resp := api.ResponseBase{
			Data: map[string]interface{}{"item": item},
			Meta: api.MetaData{Size: 1, Total: 1},
		}
		api.RenderJSON(w, http.StatusOK, resp)
	}
}

// CreateItemHandler creates new model.SampleItem
func (h *Handler) CreateItemHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		itemRequest := &api.SampleItemRequest{}

		err := itemRequest.Bind(r.Body)
		if err != nil {
			api.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			log.Println("itemRequest.Bind error: ", err)
			return
		}

		item := &model.SampleItem{
			Status:     itemRequest.Status,
			ManualProc: itemRequest.ManualDelivery,
		}

		err = h.items.AddOne(r.Context(), item)
		if err != nil {
			api.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			log.Println("usecase.AddOne error:", err)
			return
		}

		resp := api.ResponseBase{
			Data: map[string]interface{}{"item": item},
			Meta: api.MetaData{Size: 1, Total: 1},
		}
		api.RenderJSON(w, http.StatusOK, resp)
	}
}

// UpdateItemHandler updates model.SampleItem
func (h *Handler) UpdateItemHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		itemRequest := &api.SampleItemRequest{}

		err := itemRequest.Bind(r.Body)
		if err != nil {
			api.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			log.Println("itemRequest.Bind error: ", err)
			return
		}

		paramID := chi.URLParam(r, "itemID")
		if paramID == "" {
			api.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		itemRequest.ID, err = strconv.Atoi(paramID)
		if err != nil {
			api.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			log.Println("strconv.Atoi error: ", err)
			return
		}

		err = h.validate.Struct(itemRequest)
		if err != nil {
			api.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			log.Println("validate.Struct error: ", err)
			return
		}

		item := &model.SampleItem{
			ID:         itemRequest.ID,
			Status:     itemRequest.Status,
			ManualProc: itemRequest.ManualDelivery,
		}

		err = h.items.UpdateOne(r.Context(), item)
		if err != nil {
			api.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			log.Println("usecase.UpdateOne error: ", err)
			return
		}

		api.Status(w, http.StatusOK)
	}
}

// AllDatabases sample handler to get with all db names
func (h *Handler) AllDatabases() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h.items.AllDatabases(r.Context())
		if err != nil {
			api.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			log.Println("usecase.AllDatabases error: ", err)
			return
		}

		resp := api.ResponseBase{
			Data: data,
			Meta: api.MetaData{Size: len(data), Total: len(data)},
		}

		api.RenderJSON(w, http.StatusOK, resp)
	}
}

// LivenessHandler to check api response
func (h *Handler) LivenessHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		api.Status(w, http.StatusOK)
	}
}

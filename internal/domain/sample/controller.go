package sample

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"

	"github.com/rodkevich/mvpbe/internal/domain/sample/model"
	"github.com/rodkevich/mvpbe/pkg/validate"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

// Handler ...
type Handler struct {
	usecase  ItemsSampleUsage
	validate *validator.Validate
}

// NewHandler ...
func NewHandler(cmd ItemsSampleUsage) *Handler {
	return &Handler{
		usecase:  cmd,
		validate: validate.New(),
	}
}

// GetItemHandler render an item by id
func (h *Handler) GetItemHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			api.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		if _, err := strconv.Atoi(id); err != nil {
			api.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		}

		item, err := h.usecase.GetItem(r.Context(), id)
		if err != nil {
			api.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		resp := api.ResponseBase{
			Data: map[string]interface{}{"item": item},
			Meta: api.MetaData{
				Size:  1,
				Total: 1,
			},
		}
		api.RenderJSON(w, http.StatusOK, resp)
	}
}

// CreateItemHandler creates new model.SampleItem
func (h *Handler) CreateItemHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		item := &model.SampleItem{}
		err := h.usecase.AddItem(r.Context(), item)
		if err != nil {
			api.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		resp := api.ResponseBase{
			Data: map[string]interface{}{"item": item},
			Meta: api.MetaData{
				Size:  1,
				Total: 1,
			},
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
			return
		}

		err = h.validate.Struct(itemRequest)
		if err != nil {
			api.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		item := &model.SampleItem{}
		item.ID = itemRequest.ID
		item.Status = itemRequest.Status
		err = h.usecase.UpdateItem(r.Context(), item)
		if err != nil {
			api.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		api.Status(w, http.StatusOK)
	}
}

// AllDatabases sample handler to get with all db names
func (h *Handler) AllDatabases() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := h.usecase.AllDatabases(r.Context())
		if err != nil {
			api.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
		resp := api.ResponseBase{
			Data: data,
			Meta: api.MetaData{
				Size:  len(data),
				Total: len(data),
			},
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

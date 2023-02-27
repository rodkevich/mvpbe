package sample

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/rodkevich/mvpbe/internal/domain/sample/model"
	"github.com/rodkevich/mvpbe/pkg/validate"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

// Handler ...
type Handler struct {
	usecase  UseCase
	validate *validator.Validate
}

// NewHandler ...
func NewHandler(cmd UseCase) *Handler {
	return &Handler{
		usecase:  cmd,
		validate: validate.New(),
	}
}

// GetItemHandler render an item by id
func (h *Handler) GetItemHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := r.URL.Query().Get("id")
		if id == "" {
			api.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}
		fmt.Println("==========" + id + "===========")

		data, err := h.usecase.GetItem(ctx, id)
		if err != nil {
			api.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		resp := api.ResponseBase{
			Data: data,
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
		ctx := r.Context()
		item := &model.SampleItem{}
		err := h.usecase.CreateItem(ctx, item)
		if err != nil {
			api.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}

		data := map[string]interface{}{"item": item}
		resp := api.ResponseBase{
			Data: data,
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
		ctx := r.Context()
		req := &api.SampleItemRequest{}

		err := req.Bind(r.Body)
		if err != nil {
			api.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		item := &model.SampleItem{}
		item.Status = req.Status

		err = h.validate.Struct(req)
		if err != nil {
			api.Error(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
			return
		}

		err = h.usecase.UpdateItem(ctx, item)
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
		ctx := r.Context()
		data, _ := h.usecase.AllDatabases(ctx)

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

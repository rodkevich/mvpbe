package sample

import (
	"net/http"

	"github.com/go-playground/validator/v10"

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

// LivenessHandler to check api response
func (h *Handler) LivenessHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

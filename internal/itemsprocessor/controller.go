package itemsprocessor

import (
	"net/http"

	"github.com/go-playground/validator/v10"

	"github.com/rodkevich/mvpbe/pkg/validate"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

// Handler ...
type Handler struct {
	processor Processor
	validate  *validator.Validate
}

// LivenessHandler to check api response
func (h *Handler) LivenessHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		api.Status(w, http.StatusOK)
	}
}

// NewItemsHandler ...
func NewItemsHandler(p Processor) *Handler {
	return &Handler{
		processor: p,
		validate:  validate.New(),
	}
}

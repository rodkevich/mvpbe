package sample

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rodkevich/mvpbe/internal/dev"
	"github.com/rodkevich/mvpbe/internal/domain/sample/mocks"
	"github.com/rodkevich/mvpbe/internal/domain/sample/model"
)

func TestHandler_CreateItemHandler_no_error(t *testing.T) {
	t.Parallel()
	// item := &model.SampleItem{}
	// var b bytes.Buffer
	// err := json.NewEncoder(&b).Encode(item)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// req := httptest.NewRequest(http.MethodPost, "/api/v1/sample", &b)
	// w := httptest.NewRecorder()

	useCase := mocks.NewUseCase(t)
	useCase.On("CreateItem", context.Background(), &model.SampleItem{}).Return(nil)
	h := NewHandler(useCase)
	// h.CreateItemHandler()(w, req)
	// uses httptest
	assert.HTTPSuccess(t, h.CreateItemHandler(), "POST", "/api/v1/sample", nil)
	assert.HTTPBodyContains(t, h.CreateItemHandler(), "POST", "/api/v1/sample", nil, "data")
	assert.HTTPBodyContains(t, h.CreateItemHandler(), "POST", "/api/v1/sample", nil, "item")
}

func TestHandler_CreateItemHandler_create_error(t *testing.T) {
	t.Parallel()
	ctx := dev.TestContext(t)

	// item := &model.SampleItem{}
	// var b bytes.Buffer
	// err := json.NewEncoder(&b).Encode(item)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//
	// req := httptest.NewRequest(http.MethodPost, "/api/v1/sample", &b)
	// w := httptest.NewRecorder()

	useCase := mocks.NewUseCase(t)
	useCase.On("CreateItem", ctx, &model.SampleItem{}).Return(errors.New("stub"))
	h := NewHandler(useCase)
	// h.CreateItemHandler()(w, req)
	assert.HTTPError(t, h.CreateItemHandler(), "POST", "/api/v1/sample", nil)
	assert.HTTPBodyNotContains(t, h.CreateItemHandler(), "POST", "/api/v1/sample", nil, "data")
	assert.HTTPStatusCode(t, h.CreateItemHandler(), "POST", "/api/v1/sample", nil, http.StatusInternalServerError)
}

func TestHandler_GetItemHandler_id(t *testing.T) {
	t.Parallel()
	ctx := dev.TestContext(t)

	useCase := mocks.NewUseCase(t)
	useCase.On("GetItem", ctx, "1").Return(nil, nil)
	h := NewHandler(useCase)

	assert.HTTPStatusCode(t, h.GetItemHandler(), "GET", "/api/v1/sample", url.Values{"id": {""}}, http.StatusBadRequest)
	assert.HTTPStatusCode(t, h.GetItemHandler(), "GET", "/api/v1/sample", url.Values{"id": {"1"}}, http.StatusOK)
}

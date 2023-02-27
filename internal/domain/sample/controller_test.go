package sample

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rodkevich/mvpbe/internal/dev"
	"github.com/rodkevich/mvpbe/internal/domain/sample/mocks"
	"github.com/rodkevich/mvpbe/internal/domain/sample/model"

	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

func TestHandler_UpdateItemHandler_use_httptest_example(t *testing.T) {
	t.Parallel()
	data := &api.SampleItemRequest{
		ID:     1,
		Status: model.ItemPending,
	}

	item := &model.SampleItem{
		ID:     1,
		Status: model.ItemPending,
	}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(data)
	if err != nil {
		t.Fatal(err)
	}

	useCase := mocks.NewUseCase(t)
	useCase.On("UpdateItem", dev.TestContext(t), item).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/sample", &b)
	w := httptest.NewRecorder()
	h := NewHandler(useCase)
	h.UpdateItemHandler()(w, req)
	assert.Equal(t, w.Code, 200)
}

func TestHandler_CreateItemHandler_positive(t *testing.T) {
	t.Parallel()

	mockUC := mocks.NewUseCase(t)
	mockUC.On("CreateItem", dev.TestContext(t), &model.SampleItem{}).Return(nil)
	h := NewHandler(mockUC)

	// uses httptest under the hood
	t.Run("no error", func(t *testing.T) {
		assert.HTTPSuccess(t, h.CreateItemHandler(), "POST", "/api/v1/sample", nil)
		assert.HTTPBodyContains(t, h.CreateItemHandler(), "POST", "/api/v1/sample", nil, "data")
		assert.HTTPBodyContains(t, h.CreateItemHandler(), "POST", "/api/v1/sample", nil, "item")
	})
}

func TestHandler_CreateItemHandler_failures(t *testing.T) {
	t.Parallel()

	mockUC := mocks.NewUseCase(t)
	mockUC.On("CreateItem", dev.TestContext(t), &model.SampleItem{}).Return(errors.New("stub"))
	h := NewHandler(mockUC)

	assert.HTTPError(t, h.CreateItemHandler(), "POST", "/api/v1/sample", nil)
	assert.HTTPBodyNotContains(t, h.CreateItemHandler(), "POST", "/api/v1/sample", nil, "data")
	assert.HTTPStatusCode(t, h.CreateItemHandler(), "POST", "/api/v1/sample", nil, http.StatusInternalServerError)
}

func TestHandler_GetItemHandler(t *testing.T) {
	t.Parallel()
	ctx := dev.TestContext(t)

	mockUC := mocks.NewUseCase(t)
	mockUC.On("GetItem", ctx, "1").Return(&model.SampleItem{ID: 99999}, nil)
	mockUC.On("GetItem", ctx, "1O1").Return(nil, nil)
	h := NewHandler(mockUC)

	t.Run("valid_id", func(t *testing.T) {
		t.Parallel()

		assert.HTTPStatusCode(t, h.GetItemHandler(), "GET", "/api/v1/sample", url.Values{"id": {"1"}}, http.StatusOK)
		assert.HTTPBodyContains(t, h.GetItemHandler(), "GET", "/api/v1/sample", url.Values{"id": {"1"}}, 99999)
	})

	t.Run("illegal_id", func(t *testing.T) {
		t.Parallel()

		assert.HTTPStatusCode(t, h.GetItemHandler(), "GET", "/api/v1/sample", url.Values{"id": {""}}, http.StatusBadRequest)
		assert.HTTPStatusCode(t, h.GetItemHandler(), "GET", "/api/v1/sample", url.Values{"id": {"1O1"}}, http.StatusBadRequest)
		assert.HTTPBodyNotContains(t, h.GetItemHandler(), "GET", "/api/v1/sample", url.Values{"id": {""}}, "data")
	})
}

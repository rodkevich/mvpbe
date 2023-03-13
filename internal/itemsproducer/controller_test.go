package itemsproducer

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rodkevich/mvpbe/internal/dev"
	"github.com/rodkevich/mvpbe/internal/itemsproducer/mocks"
	"github.com/rodkevich/mvpbe/internal/itemsproducer/model"
	api "github.com/rodkevich/mvpbe/pkg/api/v1"
)

func TestHandler_UpdateItemHandler_Httptest_Usage_Example(t *testing.T) {
	t.Parallel()

	t.Run("require status 200", func(t *testing.T) {
		t.Parallel()

		data := &api.SampleItemRequest{ID: 777, Status: model.ItemPending}
		item := &model.SampleItem{ID: 777, Status: model.ItemPending}

		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(data)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/items", &b)

		ctx := dev.TestContext(t)
		chiCtx := chi.NewRouteContext()
		// for chi.URLParam reads chi RouteContext
		chiCtx.URLParams.Add("itemID", "777")
		ctx = context.WithValue(ctx, chi.RouteCtxKey, chiCtx)
		r = r.WithContext(ctx)

		useCase := mocks.NewItemsSampleUsage(t)
		useCase.On("UpdateOne", ctx, item).Return(nil)

		h := NewItemsHandler(useCase)
		h.UpdateItemHandler()(w, r)

		assert.Equal(t, 200, w.Code)
	})

	t.Run("fail on empty body", func(t *testing.T) {
		t.Parallel()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/v1/items", nil)

		ctx := dev.TestContext(t)
		chiCtx := chi.NewRouteContext()
		chiCtx.URLParams.Add("itemID", "777")
		ctx = context.WithValue(ctx, chi.RouteCtxKey, chiCtx)
		r = r.WithContext(ctx)

		useCase := mocks.NewItemsSampleUsage(t)

		h := NewItemsHandler(useCase)
		h.UpdateItemHandler()(w, r)

		body, err := io.ReadAll(w.Body)
		require.NoError(t, err, "failed to read HTTP body")

		assert.Equal(t, 400, w.Code)
		assert.Equal(t, `{"error":"Bad Request"}`, string(body))
	})

	t.Run("fail on params validation", func(t *testing.T) {
		t.Parallel()

		tests := []struct{ itemID string }{
			{itemID: "0"},
			{itemID: "O1"}, // contains letter
			{itemID: ""},
		}

		for _, tt := range tests {
			test := tt
			t.Run(tt.itemID, func(t *testing.T) {
				t.Parallel()

				w := httptest.NewRecorder()
				r := httptest.NewRequest(http.MethodPost, "/api/v1/items", nil)

				ctx := dev.TestContext(t)
				chiCtx := chi.NewRouteContext()

				chiCtx.URLParams.Add("itemID", test.itemID)
				ctx = context.WithValue(ctx, chi.RouteCtxKey, chiCtx)
				r = r.WithContext(ctx)

				useCase := mocks.NewItemsSampleUsage(t)

				h := NewItemsHandler(useCase)
				h.UpdateItemHandler()(w, r)
				assert.Equal(t, 400, w.Code)

				body, err := io.ReadAll(w.Body)
				require.NoError(t, err, "failed to read HTTP body")
				assert.Equal(t, `{"error":"Bad Request"}`, string(body))
			})
		}
	})
}

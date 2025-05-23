package order_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"l0/internal/http-server/handlers/order"
	"l0/internal/http-server/handlers/order/mocks"
	"l0/internal/lib/storage"
	"l0/internal/model"

	"log/slog"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestGetOrder(t *testing.T) {
	tests := []struct {
		name            string
		id              string
		acceptHeader    string
		mockReturnOrder model.Order
		mockReturnError error
		wantStatus      int
		wantBody        string
	}{
		{
			name:            "OK",
			id:              "order123",
			acceptHeader:    "application/json",
			mockReturnOrder: model.Order{OrderUID: "order123"},
			mockReturnError: nil,
			wantStatus:      http.StatusOK,
			wantBody:        `"order_uid":"order123"`,
		},
		{
			name:            "Not Found",
			id:              "missing123",
			acceptHeader:    "application/json",
			mockReturnOrder: model.Order{},
			mockReturnError: storage.ErrUrlNotFound,
			wantStatus:      http.StatusNotFound,
			wantBody:        `"error":"not found"`,
		},
		{
			name:            "Internal Error",
			id:              "error123",
			acceptHeader:    "application/json",
			mockReturnOrder: model.Order{},
			mockReturnError: errors.New("some db error"),
			wantStatus:      http.StatusOK,
			wantBody:        `"error":"internal error"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockGetter := mocks.NewORDERGetter(t)
			mockGetter.On("GetOrderById", tc.id).Return(tc.mockReturnOrder, tc.mockReturnError)

			logger := slog.New(slog.NewTextHandler(io.Discard, nil))
			handler := order.GetOrder(logger, mockGetter)

			req := httptest.NewRequest("GET", "/order/"+tc.id, nil)
			if tc.acceptHeader != "" {
				req.Header.Set("Accept", tc.acceptHeader)
			}

			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("id", tc.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			w := httptest.NewRecorder()
			handler(w, req)

			assert.Equal(t, tc.wantStatus, w.Code)
			assert.Contains(t, w.Body.String(), tc.wantBody)
			mockGetter.AssertExpectations(t)
		})
	}
}

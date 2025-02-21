package order_test

import (
	"l0wb/internal/http-server/handlers/order"
	"l0wb/internal/http-server/handlers/order/mocks"
	"l0wb/internal/lib/logger/handlers/slogdiscard"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestOrder(t *testing.T) {
	cases := []struct {
		name      string
		id        string
		respError string
		mockError error
	}{
		{
			name: "happy test",
			id:   "b563feb7b2b84b6test",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			idOrderMock := mocks.NewORDERGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				idOrderMock.On("GetOrderById", tc.id).Return(tc.id, tc.mockError).Once()
			}

			r := chi.NewRouter()
			r.Get("/{order}", order.GetOrder(slogdiscard.NewDiscardLogger(), idOrderMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			// require.NoError(t, err)
			// докончить ................
		})
	}
}

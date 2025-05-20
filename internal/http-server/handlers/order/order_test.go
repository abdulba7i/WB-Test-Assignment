package order

// import (
// 	"l0/internal/http-server/handlers/order/mocks"
// 	"l0/internal/lib/logger/handlers/slogdiscard"
// 	"l0/internal/model"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/go-chi/chi"
// 	"github.com/go-chi/render"
// 	"github.com/stretchr/testify/assert"
// )

// func TestGetOrder(t *testing.T) {
// 	cases := []struct {
// 		name      string
// 		id        string
// 		mockOrder model.Order
// 		respError string
// 		mockError error
// 	}{
// 		{
// 			name: "Success",
// 			id:   "b563feb7b2b84b6test",
// 			mockOrder: model.Order{
// 				OrderUID: "b563feb7b2b84b6test",
// 			},
// 		},
// 	}

// 	for _, tc := range cases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			mockGetter := mocks.NewORDERGetter(t)

// 			expectedOrder := model.Order{
// 				OrderUID: tc.id,
// 			}

// 			if tc.respError == "" || tc.mockError != nil {
// 				mockGetter.On("GetOrderById", tc.id).
// 					Return(tc.mockOrder, tc.mockError).Once()

// 				r := chi.NewRouter()
// 				r.Get("/order/{id}", GetOrder(slogdiscard.NewDiscardLogger(), mockGetter))

// 				req, err := http.NewRequest("GET", "/order/b563feb7b2b84b6test", nil)
// 				assert.NoError(t, err)

// 				rr := httptest.NewRecorder()
// 				r.ServeHTTP(rr, req)

// 				assert.Equal(t, http.StatusOK, rr.Code)

// 				var response Response
// 				err = render.DecodeJSON(rr.Body, &response)
// 				assert.NoError(t, err)

// 				assert.Equal(t, expectedOrder, response.Order)
// 			}

// 		})
// 	}
// }

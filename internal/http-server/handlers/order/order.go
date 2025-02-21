package order

import (
	"errors"
	"log/slog"

	resp "l0wb/internal/lib/api/response"
	"l0wb/internal/storage"
	"l0wb/internal/storage/postgres"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Request struct {
	ID string `json:"id" validate:"required,id"`
}

type Response struct {
	resp.Response
	Order postgres.Order
}

//go:generate go run github.com/vektra/mockery/v2@v2.50.1 --name=OrderGetter
type ORDERGetter interface {
	GetOrderById(id string) (postgres.Order, error)
}

func GetOrder(logger *slog.Logger, OrderGetter ORDERGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var order postgres.Order
		order, err := OrderGetter.GetOrderById(id)

		if errors.Is(err, storage.ErrUrlNotFound) {
			render.JSON(w, r, "not found")
			return
		}

		if err != nil {
			render.JSON(w, r, resp.Error("intertanl error"))
			return
		}

		render.JSON(w, r, Response{Response: *resp.OK(), Order: order})
	}
}

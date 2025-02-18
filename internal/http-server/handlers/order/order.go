package order

import (
	"errors"

	// resp "l0/internal/lib/api/response"
	resp "l0wb/internal/lib/api/response"
	"l0wb/internal/storage"
	"l0wb/internal/storage/postgres"
	"net/http"

	"github.com/go-chi/render"
)

type Request struct {
	ID string `json:"id" validate:"required,id"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.50.1 --name=URLGetter
type ORDERGetter interface {
	GetOrderById(id string) (postgres.Order, error)
}

func GetOrder(id string, OrderGetter ORDERGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Get.GetOrder"

		// log := log.With(
		// 	slog.String("op", op),
		// 	slog.String("request_id", middleware.GetReqID(r.Context())),
		// )

		var order postgres.Order
		order, err := OrderGetter.GetOrderById(id)

		if errors.Is(err, storage.ErrUrlNotFound) {
			// log.Info("url not found", "alias", id)

			render.JSON(w, r, "not found")

			return
		}

		if err != nil {
			// log.Error("failed to get url", *sl.Err(err))

			render.JSON(w, r, resp.Error("intertanl error"))

			return
		}

		http.Redirect(w, r, order.Delivery.Address, http.StatusSeeOther)
	}
}

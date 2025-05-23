package order

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	resp "l0/internal/lib/api/response"
	"l0/internal/lib/storage"
	"l0/internal/model"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

type Request struct {
	ID string `json:"id" validate:"required,id"`
}

type Response struct {
	resp.Response
	Order model.Order
}

type HTMLResponse struct {
	Error string
	Order *model.Order
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.3 --name=ORDERGetter
type ORDERGetter interface {
	GetOrderById(id string) (model.Order, error)
}

func GetOrder(logger *slog.Logger, orderGetter ORDERGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			id = r.URL.Query().Get("id")
		}
		if id == "" {
			renderTemplate(w, HTMLResponse{Error: "ID не указан"})
			return
		}

		order, err := orderGetter.GetOrderById(id)

		accept := r.Header.Get("Accept")
		isJSON := strings.Contains(accept, "application/json")

		if err != nil {
			if errors.Is(err, storage.ErrUrlNotFound) {
				if isJSON {
					w.WriteHeader(http.StatusNotFound)
					render.JSON(w, r, resp.Error("not found"))
				} else {
					renderTemplate(w, HTMLResponse{Error: "Заказ не найден"})
				}
				return
			}

			logger.Error("ошибка получения заказа", slog.String("id", id), slog.String("err", err.Error()))

			if isJSON {
				render.JSON(w, r, resp.Error("internal error"))
			} else {
				renderTemplate(w, HTMLResponse{Error: "Заказ не найден"})
			}
			return
		}

		if isJSON {
			render.JSON(w, r, Response{Response: *resp.OK(), Order: order})
		} else {
			renderTemplate(w, HTMLResponse{Order: &order})
		}
	}
}

func renderTemplate(w http.ResponseWriter, data HTMLResponse) {
	tmpl, err := template.ParseFiles(filepath.Join("internal", "http-server", "templates", "order.html"))
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if data.Error != "" {
		w.WriteHeader(http.StatusNotFound)
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
	}
}

package order

import (
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"

	// "l0/garbage/storage/postgres"
	resp "l0/internal/lib/api/response"
	"l0/internal/lib/storage"
	"l0/internal/model"

	// "l0/internal/storage"
	// "l0/internal/storage/postgres"

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

func GetOrder(logger *slog.Logger, OrderGetter ORDERGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			id = r.URL.Query().Get("id")
		}

		// Если ID не указан, показываем пустую форму
		if id == "" {
			renderTemplate(w, HTMLResponse{})
			return
		}

		order, err := OrderGetter.GetOrderById(id)

		// Определяем формат ответа на основе заголовка Accept
		acceptHeader := r.Header.Get("Accept")
		if acceptHeader == "application/json" {
			if errors.Is(err, storage.ErrUrlNotFound) {
				render.JSON(w, r, resp.Error("not found"))
				return
			}

			if err != nil {
				render.JSON(w, r, resp.Error("internal error"))
				return
			}

			render.JSON(w, r, Response{Response: *resp.OK(), Order: order})
			return
		}

		// HTML ответ
		if err != nil {
			if errors.Is(err, storage.ErrUrlNotFound) {
				renderTemplate(w, HTMLResponse{Error: "Заказ не найден"})
				return
			}
			renderTemplate(w, HTMLResponse{Error: "Внутренняя ошибка сервера"})
			return
		}

		renderTemplate(w, HTMLResponse{Order: &order})
	}
}

func renderTemplate(w http.ResponseWriter, data HTMLResponse) {
	tmpl, err := template.ParseFiles(filepath.Join("internal", "http-server", "templates", "order.html"))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

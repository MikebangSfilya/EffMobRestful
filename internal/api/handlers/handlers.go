// internal/http/handlers/handlers.go
package handlers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	datatransfer "subscription/internal/api/dto"
	"subscription/internal/model"
	"subscription/internal/service"

	"github.com/gorilla/mux"
)

type HTTPRepository interface {
	HandleSubscribe(w http.ResponseWriter, r *http.Request)
	HandleGetInfoSubscribe(w http.ResponseWriter, r *http.Request)
	HandleGetAllInfoSubscribe(w http.ResponseWriter, r *http.Request)
	HandleDeleteSubscribe(w http.ResponseWriter, r *http.Request)
	HandleUpdateSubscribe(w http.ResponseWriter, r *http.Request)
	HandleSumInfo(w http.ResponseWriter, r *http.Request)
}

type HTTPHandlers struct {
	subscriptionStore service.ServiceRepository
}

func NewHTTPHandlers(subscriptionStore service.ServiceRepository) *HTTPHandlers {
	return &HTTPHandlers{
		subscriptionStore: subscriptionStore,
	}
}

// HandleSubscribe godoc
// @Summary      Create subscription
// @Description  Создать новую подписку
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      model.Subscription  true  "Subscription"
// @Success      201  {object}  model.Subscription
// @Failure      400  {object}  datatransfer.ErrorResponse
// @Failure      500  {object}  datatransfer.ErrorResponse
// @Router       /subscriptions [post]
func (h *HTTPHandlers) HandleSubscribe(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var DTOSubs datatransfer.DTOSubs
	if err := readJSON(r, &DTOSubs); err != nil {
		log.Printf("subscription bad request error: %v", err)
		datatransfer.WriteError(w, "invalid json body", http.StatusBadRequest)
		return
	}

	sub, err := h.subscriptionStore.Create(ctx, DTOSubs)
	if err != nil {
		log.Printf("subscription bad request error: %v", err)
		datatransfer.WriteError(w, "invalid json body", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := writeJSON(w, sub); err != nil {
		return
	}

	log.Printf("subscription add successfully")
}

// HandleGetInfoSubscribe godoc
// @Summary      Get subscription
// @Description  Получить информацию о подписке по её user_id
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  model.Subscription
// @Failure      404  {object}  datatransfer.ErrorResponse
// @Failure      500  {object}  datatransfer.ErrorResponse
// @Router       /subscriptions/{id} [get]
func (h *HTTPHandlers) HandleGetInfoSubscribe(w http.ResponseWriter, r *http.Request) {
	idSub := mux.Vars(r)["id"]
	ctx := r.Context()

	subs, err := h.subscriptionStore.GetInfo(ctx, idSub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("subscription not found: id=%s", idSub)
			datatransfer.WriteError(w, "subscription not found", http.StatusNotFound)
			return
		}
		log.Printf("db error while fetching subscription: %v", err)
		datatransfer.WriteError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := writeJSON(w, subs); err != nil {
		return
	}
	log.Printf("subscription retrieved successfully: id=%s", idSub)
}

// HandleGetAllSubscriptions godoc
// @Summary      Get all subscriptions
// @Description  Получить список всех подписок
// @Tags         subscriptions
// @Produce      json
// @Success      200  {array}   model.Subscription
// @Failure      500  {object}  datatransfer.ErrorResponse
// @Router       /subscriptions [get]
func (h *HTTPHandlers) HandleGetAllInfoSubscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	subs, err := h.subscriptionStore.GetAll(ctx)
	if err != nil {
		log.Printf("failed to get subs info: %v", err)
		datatransfer.WriteError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := writeJSON(w, subs); err != nil {
		return
	}
	log.Printf("subscription all info get successfully")

}

// HandleDeleteSubscription godoc
// @Summary      Delete subscription
// @Description  Удалить подписку по её user_id
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      204  "No Content"
// @Failure      404  {object}  datatransfer.ErrorResponse
// @Failure      500  {object}  datatransfer.ErrorResponse
// @Router       /subscriptions/{id} [delete]
func (h *HTTPHandlers) HandleDeleteSubscribe(w http.ResponseWriter, r *http.Request) {
	idSub := mux.Vars(r)["id"]
	ctx := r.Context()
	if err := h.subscriptionStore.Delete(ctx, idSub); err != nil {
		log.Printf("internal server error: %v", err)
		datatransfer.WriteError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// HandleUpdateSubscription godoc
// @Summary      Update subscription
// @Description  Обновить информацию о подписке
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id            path      string              true  "User ID"
// @Param        subscription  body      model.Subscription  true  "Updated Subscription"
// @Success      200  {object}  model.Subscription
// @Failure      400  {object}  datatransfer.ErrorResponse
// @Failure      404  {object}  datatransfer.ErrorResponse
// @Failure      500  {object}  datatransfer.ErrorResponse
// @Router       /subscriptions/{id} [put]
func (h *HTTPHandlers) HandleUpdateSubscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := mux.Vars(r)["id"]

	var dto datatransfer.DTOSubs
	if err := readJSON(r, &dto); err != nil {
		log.Printf("subscription bad request error: %v", err)
		datatransfer.WriteError(w, "invalid json body", http.StatusBadRequest)
		return
	}

	updatedSub, err := h.subscriptionStore.Update(ctx, userId, dto)
	if err != nil {
		log.Printf("failed to update subscription: %v", err)
		datatransfer.WriteError(w, "failed to update subscription", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	if err := writeJSON(w, updatedSub); err != nil {
		return
	}

	log.Printf("subscription update succefully")

}

// HandleSumInfo godoc
// @Summary      Calculate total subscription cost
// @Description  Подсчёт суммарной стоимости всех подписок за период с фильтрацией
// @Tags         subscriptions
// @Produce      json
// @Param        id            query     string  false  "User ID"
// @Param        service_name  query     string  false  "Service Name"
// @Param        from          query     string  true   "Start period (MM-YYYY)"
// @Param        to            query     string  true   "End period (MM-YYYY)"
// @Success      200  {object}  datatransfer.SumResponse
// @Failure      400  {object}  datatransfer.ErrorResponse
// @Failure      500  {object}  datatransfer.ErrorResponse
// @Router       /subscriptions/sum [get]
func (h *HTTPHandlers) HandleSumInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := r.URL.Query().Get("id")
	serviceName := r.URL.Query().Get("service_name")
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	// Проверка обязательных параметров
	if fromStr == "" || toStr == "" {
		datatransfer.WriteError(w, "missing 'from' or 'to' query parameter", http.StatusBadRequest)
		return
	}

	// Парсим даты в CustomDate
	var fromDate, toDate model.CustomDate
	if err := fromDate.UnmarshalJSON([]byte(`"` + fromStr + `"`)); err != nil {
		datatransfer.WriteError(w, "invalid 'from' date format", http.StatusBadRequest)
		return
	}
	if err := toDate.UnmarshalJSON([]byte(`"` + toStr + `"`)); err != nil {
		datatransfer.WriteError(w, "invalid 'to' date format", http.StatusBadRequest)
		return
	}

	// Получаем сумму
	sum, err := h.subscriptionStore.Sum(ctx, userID, serviceName, fromDate.Time, toDate.Time)
	if err != nil {
		log.Printf("failed to calculate sum: %v", err)
		datatransfer.WriteError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := datatransfer.SumResponse{TotalPrice: sum}

	if err := writeJSON(w, resp); err != nil {
		return
	}

	log.Printf("subscription sum calculated successfully: user_id=%s service_name=%s from=%s to=%s sum=%d",
		userID, serviceName, fromStr, toStr, sum)
}

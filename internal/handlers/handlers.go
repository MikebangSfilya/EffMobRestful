// handlers.go содержит наши http-"ручки"
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	datatransfer "subscription/internal/dto"
	"subscription/internal/model"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	subscriptionStore *model.SubscriptionStore
}

func NewHTTPHandlers(subscriptionStore *model.SubscriptionStore) *HTTPHandlers {
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
	if err := json.NewDecoder(r.Body).Decode(&DTOSubs); err != nil {
		log.Printf("subscription bad request error: %v", err)
		datatransfer.WriteError(w, "bad request", http.StatusBadRequest)
		return
	}

	subsNew := model.NewSubscription(DTOSubs.ServiceName, DTOSubs.Price)

	h.subscriptionStore.AddSub(ctx, subsNew)

	if err := json.NewEncoder(w).Encode(subsNew); err != nil {
		log.Printf("failed to encode subscription: %v", err)
		datatransfer.WriteError(w, "internal server error", http.StatusInternalServerError)
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
	subs, err := h.subscriptionStore.GetSubInfo(ctx, idSub)
	if err != nil {
		log.Printf("subscription not found for id: %s, error: %v", idSub, err)
		datatransfer.WriteError(w, "subscription not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(subs); err != nil {
		log.Printf("failed to encode subscription: %v", err)
		datatransfer.WriteError(w, "internal server error", http.StatusInternalServerError)
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
	subs, err := h.subscriptionStore.GetSubAllInfo(ctx)
	if err != nil {
		log.Printf("failed to get subs info: %v", err)
		datatransfer.WriteError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		log.Printf("failed to encode subscription: %v", err)
		datatransfer.WriteError(w, "internal server error", http.StatusInternalServerError)
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
	if err := h.subscriptionStore.DeleteInfo(ctx, idSub); err != nil {
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
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		log.Printf("subscription bad request error: %v", err)
		datatransfer.WriteError(w, "bad request", http.StatusBadRequest)
		return
	}

	oldSub, _ := h.subscriptionStore.GetSubInfo(ctx, userId)
	if (oldSub == model.Subscription{}) {
		log.Printf("subscription not found")
		datatransfer.WriteError(w, "not found", http.StatusNotFound)
		return
	}

	updatedSub := model.Subscription{
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserId:      oldSub.UserId,
		StartDate:   oldSub.StartDate,
	}

	if err := h.subscriptionStore.UpdateSub(ctx, userId, updatedSub); err != nil {
		log.Printf("failed to update subscription:")
		datatransfer.WriteError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(updatedSub); err != nil {
		log.Printf("failed to encode subscription: %v", err)
		datatransfer.WriteError(w, "internal server error", http.StatusInternalServerError)
		return
	}
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
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	sum, err := h.subscriptionStore.SumSubscriptions(ctx, userID, serviceName, from, to)
	if err != nil {
		log.Printf("failed to calculate sum: %v", err)
		datatransfer.WriteError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := map[string]int{"total_price": sum}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode sum response: %v", err)
		datatransfer.WriteError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("subscription sum calculated successfully: user_id=%s service_name=%s from=%s to=%s sum=%d",
		userID, serviceName, from, to, sum)
}

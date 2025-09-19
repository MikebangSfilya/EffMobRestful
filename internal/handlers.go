// package internal
// handlers.go содержит наши http-"ручки"
package internal

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	subscriptionStore *SubscriptionStore
}

func NewHTTPHandlers(subscriptionStore *SubscriptionStore) *HTTPHandlers {
	return &HTTPHandlers{
		subscriptionStore: subscriptionStore,
	}
}

// HandleSubscribe godoc
// @Summary      Create subscription
// @Description  Создать новую подписку (service_name, price, user_id, start_date)
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body   DTOSubs  true  "Subscription info"
// @Success      201  {object} Subscription
// @Failure      400  {object} ErrorResponse
// @Failure      409  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /subscriptions [post]
func (h *HTTPHandlers) HandleSubscribe(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	var DTOSubs DTOSubs
	if err := json.NewDecoder(r.Body).Decode(&DTOSubs); err != nil {
		log.Printf("subscription bad request error: %v", err)
		writeError(w, "bad request", http.StatusBadRequest)
		return
	}

	subsNew := NewSubscription(DTOSubs.ServiceName, DTOSubs.Price)

	h.subscriptionStore.AddSub(ctx, subsNew)

	if err := json.NewEncoder(w).Encode(subsNew); err != nil {
		log.Printf("failed to encode subscription: %v", err)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("subscription add successfully")
}

// HandleGetInfoSubscribe godoc
// @Summary      Get subscription by ID
// @Description  Получить информацию о подписке по её id
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      string  true  "Subscription ID"
// @Success      200  {object} Subscription
// @Failure      400  {object} ErrorResponse
// @Failure      404  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /subscriptions/{id} [get]
func (h *HTTPHandlers) HandleGetInfoSubscribe(w http.ResponseWriter, r *http.Request) {
	idSub := mux.Vars(r)["id"]
	ctx := r.Context()
	subs, err := h.subscriptionStore.GetSubInfo(ctx, idSub)
	if err != nil {
		log.Printf("subscription not found for id: %s, error: %v", idSub, err)
		writeError(w, "subscription not found", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(subs); err != nil {
		log.Printf("failed to encode subscription: %v", err)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("subscription retrieved successfully: id=%s", idSub)
}

// HandleGetAllInfoSubscribe godoc
// @Summary      Get all subscriptions
// @Description  Получить информацию обо всех подписках
// @Tags         subscriptions
// @Produce      json
// @Success      200  {array}  Subscription
// @Failure      500  {object} ErrorResponse
// @Router       /subscriptions [get]
func (h *HTTPHandlers) HandleGetAllInfoSubscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	subs, err := h.subscriptionStore.GetSubAllInfo(ctx)
	if err != nil {
		log.Printf("failed to get subs info: %v", err)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(subs); err != nil {
		log.Printf("failed to encode subscription: %v", err)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("subscription all info get successfully")

}

// HandleDeleteSubscribe godoc
// @Summary      Delete subscription
// @Description  Отменить подписку по её id
// @Tags         subscriptions
// @Param        id   path      string  true  "Subscription ID"
// @Success      204
// @Failure      400  {object} ErrorResponse
// @Failure      404  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /subscriptions/{id} [delete]
func (h *HTTPHandlers) HandleDeleteSubscribe(w http.ResponseWriter, r *http.Request) {
	idSub := mux.Vars(r)["id"]
	ctx := r.Context()
	if err := h.subscriptionStore.DeleteInfo(ctx, idSub); err != nil {
		log.Printf("internal server error: %v", err)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// HandleUpdateSubscribe godoc
// @Summary      Update subscription
// @Description  Обновить подписку по её id
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id            path    string  true  "Subscription ID"
// @Param        subscription  body    DTOSubs true  "Updated subscription info"
// @Success      200  {object} Subscription
// @Failure      400  {object} ErrorResponse
// @Failure      404  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /subscriptions/{id} [put]
func (h *HTTPHandlers) HandleUpdateSubscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId := mux.Vars(r)["id"]
	var dto DTOSubs
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		log.Printf("subscription bad request error: %v", err)
		writeError(w, "bad request", http.StatusBadRequest)
		return
	}

	oldSub, _ := h.subscriptionStore.GetSubInfo(ctx, userId)
	if (oldSub == Subscription{}) {
		log.Printf("subscription not found")
		writeError(w, "not found", http.StatusNotFound)
		return
	}

	updatedSub := Subscription{
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserId:      oldSub.UserId,
		StartDate:   oldSub.StartDate,
	}

	if err := h.subscriptionStore.UpdateSub(ctx, userId, updatedSub); err != nil {
		log.Printf("failed to update subscription:")
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(updatedSub); err != nil {
		log.Printf("failed to encode subscription: %v", err)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// HandleSumInfo godoc
// @Summary      Get total subscription cost
// @Description  Получить суммарную стоимость подписок с фильтрацией по id пользователя, названию подписки и периоду
// @Tags         subscriptions
// @Produce      json
// @Param        id            query   string  false  "User ID"
// @Param        service_name  query   string  false  "Service Name"
// @Param        from          query   string  false  "Start month-year, format 01-2006"
// @Param        to            query   string  false  "End month-year, format 01-2006"
// @Success      200  {object} map[string]int
// @Failure      400  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
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
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	resp := map[string]int{"total_price": sum}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode sum response: %v", err)
		writeError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("subscription sum calculated successfully: user_id=%s service_name=%s from=%s to=%s sum=%d",
		userID, serviceName, from, to, sum)
}

// package internal
// handlers.go содержит наши http-"ручки"
package internal

import (
	"encoding/json"
	"fmt"
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
// @Description  Создать новую подписку (service_name,   price, user_id, start_date).
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body   Subscription  true  "Subscription info"
// @Success      201  {object} Subscription
// @Failure      400  {object} ErrorResponse
// @Failure      409  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /subscriptions [post]
func (h *HTTPHandlers) HandleSubscribe(w http.ResponseWriter, r *http.Request) {
	var DTOSubs DTOSubs
	if err := json.NewDecoder(r.Body).Decode(&DTOSubs); err != nil {
		fmt.Println("errors is", err)
	}

	subsNew := NewSubscription(DTOSubs.ServiceName, DTOSubs.Price)

	h.subscriptionStore.AddSub(subsNew)

	b, err := json.MarshalIndent(subsNew, "", "	   ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http resonce", err)
		return
	}
}

// HandleGetInfoSubscribe godoc
// @Summary      Get subscription
// @Description  Получить информацию о подписке по её id.
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

	subs := h.subscriptionStore.GetSubInfo(idSub)

	b, err := json.MarshalIndent(subs, "", "	   ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http resonce", err)
		return
	}
}

// HandleGetInfoSubscribe godoc
// @Summary      Get subscription
// @Description  Получить информацию о всех подписках
// @Tags         subscriptions
// @Produce      json
// @Param        -
// @Success      200  {object} Subscription
// @Failure      400  {object} ErrorResponse
// @Failure      404  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /subscriptions [get]
func (h *HTTPHandlers) HandleGetAllInfoSubscribe(w http.ResponseWriter, r *http.Request) {
	subs := h.subscriptionStore.GetSubAllInfo()
	b, err := json.MarshalIndent(subs, "", "    ")
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http resonce", err)
		return
	}

}

// HandleDeleteSubscribe godoc
// @Summary      Delete subscription
// @Description  Отменить подписку по её id.
// @Tags         subscriptions
// @Param        id   path      string  true  "Subscription ID"
// @Success      204
// @Failure      400  {object} ErrorResponse
// @Failure      409  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /subscriptions/{id} [delete]
func (h *HTTPHandlers) HandleDeleteSubscribe(w http.ResponseWriter, r *http.Request) {
	idSub := mux.Vars(r)["id"]
	h.subscriptionStore.DeleteInfo(idSub)

	w.WriteHeader(http.StatusNoContent)

}

func (h *HTTPHandlers) HandleUpdateSubscribe(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["id"]

	var dto DTOSubs
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		fmt.Println("failed to decode request body:", err)
		return
	}

	oldSub := h.subscriptionStore.GetSubInfo(userId)
	if (oldSub == Subscription{}) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	updatedSub := Subscription{
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		UserId:      oldSub.UserId,
		StartDate:   oldSub.StartDate,
	}

	if err := h.subscriptionStore.UpdateSub(userId, updatedSub); err != nil {
		fmt.Println("failed to update subscription:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	b, err := json.MarshalIndent(updatedSub, "", "	   ")
	if err != nil {
		fmt.Println("failed to marshal response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

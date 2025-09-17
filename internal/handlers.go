// package internal
// handlers.go содержит наши http-"ручки"
package internal

import "net/http"

type HTTPHandlers struct {
	subscriptionStore *SubscriptionStore
}

func NewHTTPHandlers(subscriptionStore SubscriptionStore) *HTTPHandlers {
	return &HTTPHandlers{}
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is working fine"))
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

}

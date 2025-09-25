package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	datatransfer "subscription/internal/api/dto"
	"subscription/internal/api/handlers"
	"subscription/internal/model"
	"testing"
	"time"
)

type fakeService struct{}

func (f *fakeService) Create(ctx context.Context, dto datatransfer.DTOSubs) (model.Subscription, error) {
	return model.Subscription{}, nil
}
func (f *fakeService) GetInfo(ctx context.Context, idSub string) (model.Subscription, error) {
	return model.Subscription{}, nil
}
func (f *fakeService) GetAll(ctx context.Context) ([]model.Subscription, error) {
	return []model.Subscription{}, nil
}
func (f *fakeService) Delete(ctx context.Context, idSub string) error {
	return nil
}
func (f *fakeService) Update(ctx context.Context, id string, dto datatransfer.DTOSubs) (model.Subscription, error) {
	return model.Subscription{}, nil
}
func (f *fakeService) Sum(ctx context.Context, userId, serviceName string, from, to time.Time) (int, error) {
	return 1, nil
}

func TestHandleSubscribe_Unit(t *testing.T) {

	h := handlers.NewHTTPHandlers(&fakeService{})
	dto := datatransfer.DTOSubs{
		UserId:      "user1",
		ServiceName: "Netflix",
		Price:       10,
		StartDate:   "26.10.2025",
	}
	body, _ := json.Marshal(dto)

	req := httptest.NewRequest(http.MethodPost, "/subscriptions", bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.HandleSubscribe(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

}

func TestHandleDelete_Unit(t *testing.T) {

	h := handlers.NewHTTPHandlers(&fakeService{})

	req := httptest.NewRequest(http.MethodDelete, "/subscriptions/1", nil)
	w := httptest.NewRecorder()

	h.HandleDeleteSubscribe(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
	}

}

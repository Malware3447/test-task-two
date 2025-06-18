package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"test-task-two/internal/api/crut"
	"test-task-two/internal/models/request"
	"test-task-two/internal/models/response"
	"testing"
)

const testUUID = "f47ac10b-58cc-4372-a567-0e02b2c3d479"

type MockPgService struct {
	addAmountFunc      func(ctx context.Context, req *request.Request) (*response.Response, error)
	withdrawAmountFunc func(ctx context.Context, req *request.Request) (*response.Response, error)
	getAmountFunc      func(ctx context.Context, uuid string) (*response.Response, error)
}

func (m *MockPgService) AddAmount(ctx context.Context, req *request.Request) (*response.Response, error) {
	return m.addAmountFunc(ctx, req)
}

func (m *MockPgService) WithdrawAmount(ctx context.Context, req *request.Request) (*response.Response, error) {
	return m.withdrawAmountFunc(ctx, req)
}

func (m *MockPgService) GetAmount(ctx context.Context, uuid string) (*response.Response, error) {
	return m.getAmountFunc(ctx, uuid)
}

func TestAddAndWithdrawAmount(t *testing.T) {
	mockService := &MockPgService{
		addAmountFunc: func(ctx context.Context, req *request.Request) (*response.Response, error) {
			return &response.Response{Uuid: req.Uuid, Amount: req.Amount, Error: "OK"}, nil
		},
		withdrawAmountFunc: func(ctx context.Context, req *request.Request) (*response.Response, error) {
			return &response.Response{Uuid: req.Uuid, Amount: req.Amount, Error: "OK"}, nil
		},
	}

	crutHandler := crut.NewCrut(&crut.Params{RepoPg: mockService})

	testCases := []struct {
		name           string
		requestBody    request.Request
		expectedStatus int
	}{
		{
			name: "Successful deposit",
			requestBody: request.Request{
				Uuid:      testUUID,
				Amount:    100,
				Operation: "DEPOSIT",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Successful withdraw",
			requestBody: request.Request{
				Uuid:      testUUID,
				Amount:    50,
				Operation: "WITHDRAW",
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest("POST", "/amount", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			crutHandler.AddAndWithdrawAmount(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			var resp response.Response
			err := json.NewDecoder(w.Body).Decode(&resp)
			if err != nil {
				t.Fatalf("Error decoding response: %v", err)
			}

			if resp.Uuid != tc.requestBody.Uuid {
				t.Errorf("Expected UUID %s, got %s", tc.requestBody.Uuid, resp.Uuid)
			}

			if resp.Amount != tc.requestBody.Amount {
				t.Errorf("Expected amount %d, got %d", tc.requestBody.Amount, resp.Amount)
			}

			if resp.Error != "OK" {
				t.Errorf("Expected error 'OK', got '%s'", resp.Error)
			}
		})
	}
}

func TestGetAmount(t *testing.T) {
	mockService := &MockPgService{
		getAmountFunc: func(ctx context.Context, uuid string) (*response.Response, error) {
			return &response.Response{Uuid: uuid, Amount: 10000, Error: "OK"}, nil
		},
	}

	crutHandler := crut.NewCrut(&crut.Params{RepoPg: mockService})

	req := httptest.NewRequest("GET", "/amount/"+testUUID, nil)
	w := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, &chi.Context{
		URLParams: chi.RouteParams{Keys: []string{"uuid"}, Values: []string{testUUID}},
	})
	req = req.WithContext(ctx)

	crutHandler.GetAmount(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var resp response.Response
	err := json.NewDecoder(w.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if resp.Uuid != testUUID {
		t.Errorf("Expected UUID '%s', got '%s'", testUUID, resp.Uuid)
	}

	if resp.Amount != 10000 {
		t.Errorf("Expected amount 10000, got %d", resp.Amount)
	}

	if resp.Error != "OK" {
		t.Errorf("Expected error 'OK', got '%s'", resp.Error)
	}
}

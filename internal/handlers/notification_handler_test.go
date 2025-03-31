package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"notification-service/internal/models"
	"notification-service/internal/services"
	"testing"
	"time"
)

func TestNotificationHandler(t *testing.T) {
	// Setup
	factory := services.NewNotificationServiceFactory()
	defaultService, _ := factory.GetService(models.ChannelSlack)
	scheduler := services.NewSchedulerService(defaultService)
	scheduler.Start()
	defer scheduler.Stop()

	handler := NewNotificationHandler(factory, scheduler)

	tests := []struct {
		name          string
		request       SendNotificationRequest
		method        string
		expectedCode  int
		expectedBody  APIResponse
		validateExtra func(*testing.T, APIResponse)
	}{
		{
			name: "Successful immediate slack notification",
			request: SendNotificationRequest{
				Title:      "Test Slack",
				Content:    "Test content",
				Channel:    models.ChannelSlack,
				Recipients: []string{"user1"},
			},
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expectedBody: APIResponse{
				Success: true,
				Message: "Notification sent successfully",
			},
			validateExtra: func(t *testing.T, resp APIResponse) {
				if resp.Data == nil {
					t.Error("Expected data in response")
				}
			},
		},
		{
			name: "Successful scheduled email notification",
			request: SendNotificationRequest{
				Title:       "Test Email",
				Content:     "Test content",
				Channel:     models.ChannelEmail,
				Recipients:  []string{"test@example.com"},
				ScheduledAt: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			},
			method:       http.MethodPost,
			expectedCode: http.StatusAccepted,
			expectedBody: APIResponse{
				Success: true,
				Message: "Notification scheduled successfully",
			},
		},
		{
			name: "Missing required fields",
			request: SendNotificationRequest{
				Channel:    models.ChannelSlack,
				Recipients: []string{"user1"},
			},
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: APIResponse{
				Success: false,
				Message: "Title and content are required",
			},
		},
		{
			name: "Empty recipients",
			request: SendNotificationRequest{
				Title:      "Test",
				Content:    "Content",
				Channel:    models.ChannelSlack,
				Recipients: []string{},
			},
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: APIResponse{
				Success: false,
				Message: "At least one recipient is required",
			},
		},
		{
			name: "Invalid channel",
			request: SendNotificationRequest{
				Title:      "Test",
				Content:    "Content",
				Channel:    "invalid",
				Recipients: []string{"user1"},
			},
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: APIResponse{
				Success: false,
				Message: "Invalid notification channel: unsupported notification channel: invalid",
			},
		},
		{
			name: "Invalid scheduled time format",
			request: SendNotificationRequest{
				Title:       "Test",
				Content:     "Content",
				Channel:     models.ChannelEmail,
				Recipients:  []string{"test@example.com"},
				ScheduledAt: "invalid-time",
			},
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: APIResponse{
				Success: false,
				Message: "Invalid scheduled_at time format. Use RFC3339 format (e.g., 2024-03-31T21:20:00Z)",
			},
		},
		{
			name: "Past scheduled time",
			request: SendNotificationRequest{
				Title:       "Test",
				Content:     "Content",
				Channel:     models.ChannelEmail,
				Recipients:  []string{"test@example.com"},
				ScheduledAt: time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
			},
			method:       http.MethodPost,
			expectedCode: http.StatusBadRequest,
			expectedBody: APIResponse{
				Success: false,
				Message: "Scheduled time must be in the future",
			},
		},
		{
			name:         "Invalid HTTP method",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: APIResponse{
				Success: false,
				Message: "Method not allowed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			var err error
			if tt.method == http.MethodPost {
				reqBody, err = json.Marshal(tt.request)
				if err != nil {
					t.Fatalf("Failed to marshal request: %v", err)
				}
			}

			req := httptest.NewRequest(tt.method, "/notifications", bytes.NewBuffer(reqBody))
			if tt.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/json")
			}
			rr := httptest.NewRecorder()

			handler.SendNotification(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			var response APIResponse
			if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
				t.Fatalf("Failed to decode response: %v", err)
			}

			if response.Success != tt.expectedBody.Success {
				t.Errorf("Expected success %v, got %v", tt.expectedBody.Success, response.Success)
			}

			if response.Message != tt.expectedBody.Message {
				t.Errorf("Expected message %q, got %q", tt.expectedBody.Message, response.Message)
			}

			if tt.validateExtra != nil {
				tt.validateExtra(t, response)
			}
		})
	}
}

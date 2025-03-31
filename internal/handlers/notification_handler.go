package handlers

import (
	"encoding/json"
	"net/http"
	"notification-service/internal/models"
	"notification-service/internal/services"
	"time"

	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationFactory *services.NotificationServiceFactory
	schedulerService    *services.SchedulerService
}

func NewNotificationHandler(factory *services.NotificationServiceFactory, scheduler *services.SchedulerService) *NotificationHandler {
	return &NotificationHandler{
		notificationFactory: factory,
		schedulerService:    scheduler,
	}
}

type SendNotificationRequest struct {
	Title       string                     `json:"title"`
	Content     string                     `json:"content"`
	Channel     models.NotificationChannel `json:"channel"`
	Recipients  []string                   `json:"recipients"`
	ScheduledAt string                     `json:"scheduled_at,omitempty"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func generateID() string {
	return uuid.New().String()
}

func (h *NotificationHandler) SendNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSONResponse(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONResponse(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	// Validate required fields
	if req.Title == "" || req.Content == "" {
		sendJSONResponse(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Title and content are required",
		})
		return
	}

	if len(req.Recipients) == 0 {
		sendJSONResponse(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "At least one recipient is required",
		})
		return
	}

	// Get the service for the requested channel
	service, err := h.notificationFactory.GetService(req.Channel)
	if err != nil {
		sendJSONResponse(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid notification channel: " + err.Error(),
		})
		return
	}

	// Parse scheduled time if provided
	var scheduledTime *time.Time
	if req.ScheduledAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.ScheduledAt)
		if err != nil {
			sendJSONResponse(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Invalid scheduled_at time format. Use RFC3339 format (e.g., 2024-03-31T21:20:00Z)",
			})
			return
		}
		if parsedTime.Before(time.Now()) {
			sendJSONResponse(w, http.StatusBadRequest, APIResponse{
				Success: false,
				Message: "Scheduled time must be in the future",
			})
			return
		}
		scheduledTime = &parsedTime
	}

	// Create notification
	notification := &models.Notification{
		ID:          generateID(),
		Title:       req.Title,
		Content:     req.Content,
		Channel:     req.Channel,
		Recipients:  req.Recipients,
		ScheduledAt: scheduledTime,
		CreatedAt:   time.Now(),
	}

	// Handle scheduled vs immediate notifications
	if scheduledTime != nil {
		if err := h.schedulerService.ScheduleNotification(notification); err != nil {
			sendJSONResponse(w, http.StatusInternalServerError, APIResponse{
				Success: false,
				Message: "Failed to schedule notification: " + err.Error(),
			})
			return
		}

		sendJSONResponse(w, http.StatusAccepted, APIResponse{
			Success: true,
			Message: "Notification scheduled successfully",
			Data:    notification,
		})
		return
	}

	// Send immediate notification
	if err := service.Send(notification); err != nil {
		sendJSONResponse(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to send notification: " + err.Error(),
		})
		return
	}

	sendJSONResponse(w, http.StatusOK, APIResponse{
		Success: true,
		Message: "Notification sent successfully",
		Data:    notification,
	})
}

func sendJSONResponse(w http.ResponseWriter, status int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

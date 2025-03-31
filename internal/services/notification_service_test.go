package services

import (
	"notification-service/internal/models"
	"testing"
	"time"
)

func TestSlackNotificationService(t *testing.T) {
	service := &SlackNotificationService{}
	notification := &models.Notification{
		ID:         "test-1",
		Title:      "Test Slack Notification",
		Content:    "This is a test notification",
		Channel:    models.ChannelSlack,
		Recipients: []string{"test-user"},
		CreatedAt:  time.Now(),
	}

	err := service.Send(notification)
	if err != nil {
		t.Errorf("Failed to send Slack notification: %v", err)
	}
}

func TestEmailNotificationService(t *testing.T) {
	service := &EmailNotificationService{}
	notification := &models.Notification{
		ID:         "test-2",
		Title:      "Test Email Notification",
		Content:    "This is a test email",
		Channel:    models.ChannelEmail,
		Recipients: []string{"test@example.com"},
		CreatedAt:  time.Now(),
	}

	err := service.Send(notification)
	if err != nil {
		t.Errorf("Failed to send Email notification: %v", err)
	}
}

func TestMessageNotificationService(t *testing.T) {
	service := &MessageNotificationService{}
	notification := &models.Notification{
		ID:         "test-3",
		Title:      "Test SMS Notification",
		Content:    "This is a test SMS",
		Channel:    models.ChannelMessage,
		Recipients: []string{"+1234567890"},
		CreatedAt:  time.Now(),
	}

	err := service.Send(notification)
	if err != nil {
		t.Errorf("Failed to send SMS notification: %v", err)
	}
}

func TestNotificationServiceFactory(t *testing.T) {
	factory := NewNotificationServiceFactory()

	// Test getting Slack service
	slackService, err := factory.GetService(models.ChannelSlack)
	if err != nil {
		t.Errorf("Failed to get Slack service: %v", err)
	}
	if slackService == nil {
		t.Error("Slack service is nil")
	}

	// Test getting Email service
	emailService, err := factory.GetService(models.ChannelEmail)
	if err != nil {
		t.Errorf("Failed to get Email service: %v", err)
	}
	if emailService == nil {
		t.Error("Email service is nil")
	}

	// Test getting SMS service
	smsService, err := factory.GetService(models.ChannelMessage)
	if err != nil {
		t.Errorf("Failed to get SMS service: %v", err)
	}
	if smsService == nil {
		t.Error("SMS service is nil")
	}

	// Test invalid channel
	invalidService, err := factory.GetService("invalid-channel")
	if err == nil {
		t.Error("Expected error for invalid channel, got nil")
	}
	if invalidService != nil {
		t.Error("Expected nil service for invalid channel")
	}
}

func TestSchedulerService(t *testing.T) {
	// Create a test notification service
	testService := &SlackNotificationService{}
	scheduler := NewSchedulerService(testService)

	// Test scheduling a notification
	scheduledTime := time.Now().Add(2 * time.Second)
	notification := &models.Notification{
		ID:          "test-4",
		Title:       "Test Scheduled Notification",
		Content:     "This is a scheduled test",
		Channel:     models.ChannelSlack,
		Recipients:  []string{"test-user"},
		ScheduledAt: &scheduledTime,
		CreatedAt:   time.Now(),
	}

	err := scheduler.ScheduleNotification(notification)
	if err != nil {
		t.Errorf("Failed to schedule notification: %v", err)
	}

	// Start the scheduler
	scheduler.Start()
	defer scheduler.Stop()

	// Wait for the notification to be sent
	time.Sleep(3 * time.Second)
}

func TestMultipleScheduledNotifications(t *testing.T) {
	testService := &SlackNotificationService{}
	scheduler := NewSchedulerService(testService)
	scheduler.Start()
	defer scheduler.Stop()

	// Schedule multiple notifications with different delays
	notifications := []*models.Notification{
		{
			ID:          "test-5",
			Title:       "First Scheduled Notification",
			Content:     "This is the first scheduled test",
			Channel:     models.ChannelSlack,
			Recipients:  []string{"user1"},
			ScheduledAt: &time.Time{},
			CreatedAt:   time.Now(),
		},
		{
			ID:          "test-6",
			Title:       "Second Scheduled Notification",
			Content:     "This is the second scheduled test",
			Channel:     models.ChannelSlack,
			Recipients:  []string{"user2"},
			ScheduledAt: &time.Time{},
			CreatedAt:   time.Now(),
		},
	}

	// Set different delays
	*notifications[0].ScheduledAt = time.Now().Add(2 * time.Second)
	*notifications[1].ScheduledAt = time.Now().Add(4 * time.Second)

	// Schedule all notifications
	for _, notification := range notifications {
		err := scheduler.ScheduleNotification(notification)
		if err != nil {
			t.Errorf("Failed to schedule notification %s: %v", notification.ID, err)
		}
	}

	// Wait for all notifications to be sent
	time.Sleep(5 * time.Second)
}

func TestInvalidScheduledTime(t *testing.T) {
	testService := &SlackNotificationService{}
	scheduler := NewSchedulerService(testService)

	// Test with past scheduled time
	pastTime := time.Now().Add(-1 * time.Hour)
	notification := &models.Notification{
		ID:          "test-7",
		Title:       "Invalid Scheduled Notification",
		Content:     "This should fail",
		Channel:     models.ChannelSlack,
		Recipients:  []string{"test-user"},
		ScheduledAt: &pastTime,
		CreatedAt:   time.Now(),
	}

	err := scheduler.ScheduleNotification(notification)
	if err == nil {
		t.Error("Expected error for past scheduled time, got nil")
	}
}

func TestNilScheduledTime(t *testing.T) {
	testService := &SlackNotificationService{}
	scheduler := NewSchedulerService(testService)

	notification := &models.Notification{
		ID:         "test-8",
		Title:      "Nil Scheduled Time Notification",
		Content:    "This should fail",
		Channel:    models.ChannelSlack,
		Recipients: []string{"test-user"},
		CreatedAt:  time.Now(),
	}

	err := scheduler.ScheduleNotification(notification)
	if err == nil {
		t.Error("Expected error for nil scheduled time, got nil")
	}
}

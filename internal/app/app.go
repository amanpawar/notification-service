package app

import (
	"context"
	"fmt"
	"net/http"
	"notification-service/internal/config"
	"notification-service/internal/handlers"
	"notification-service/internal/models"
	"notification-service/internal/services"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	config              *config.Config
	notificationFactory *services.NotificationServiceFactory
	schedulerService    *services.SchedulerService
	server              *http.Server
}

func NewApp(cfg *config.Config) *App {
	notificationFactory := services.NewNotificationServiceFactory()
	defaultService, _ := notificationFactory.GetService(models.ChannelSlack)
	schedulerService := services.NewSchedulerService(defaultService)

	return &App{
		config:              cfg,
		notificationFactory: notificationFactory,
		schedulerService:    schedulerService,
	}
}

func (a *App) Run() error {
	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the scheduler service
	a.schedulerService.Start()
	defer a.schedulerService.Stop()

	fmt.Println("\nNotification service is running with the following examples:")
	fmt.Println("1. Immediate Slack notification to 3 users")
	fmt.Println("2. Email notification scheduled for 5 seconds from now")
	fmt.Println("3. Two SMS notifications scheduled for 10 and 15 seconds from now")
	fmt.Println("\nPress Ctrl+C to exit.")
	fmt.Println("\nSending notifications...\n")

	// Small delay to ensure messages are displayed
	time.Sleep(1 * time.Second)

	// Example 1: Immediate Slack notification to multiple users
	slackNotification := &models.Notification{
		ID:         "1",
		Title:      "Team Meeting Reminder",
		Content:    "Don't forget about the team meeting at 2 PM today!",
		Channel:    models.ChannelSlack,
		Recipients: []string{"user1", "user2", "user3"},
		CreatedAt:  time.Now(),
	}

	slackService, err := a.notificationFactory.GetService(slackNotification.Channel)
	if err != nil {
		return fmt.Errorf("failed to get slack service: %v", err)
	}

	if err := slackService.Send(slackNotification); err != nil {
		return fmt.Errorf("failed to send slack notification: %v", err)
	}

	// Example 2: Scheduled Email notification
	scheduledTime := time.Now().Add(5 * time.Second)
	emailNotification := &models.Notification{
		ID:          "2",
		Title:       "Weekly Report Ready",
		Content:     "Your weekly performance report is now available.",
		Channel:     models.ChannelEmail,
		Recipients:  []string{"manager@company.com", "hr@company.com"},
		ScheduledAt: &scheduledTime,
		CreatedAt:   time.Now(),
	}

	emailService, err := a.notificationFactory.GetService(emailNotification.Channel)
	if err != nil {
		return fmt.Errorf("failed to get email service: %v", err)
	}

	emailScheduler := services.NewSchedulerService(emailService)
	emailScheduler.Start()
	defer emailScheduler.Stop()

	if err := emailScheduler.ScheduleNotification(emailNotification); err != nil {
		return fmt.Errorf("failed to schedule email notification: %v", err)
	}

	// Example 3: Multiple scheduled SMS notifications with different delays
	smsService, err := a.notificationFactory.GetService(models.ChannelMessage)
	if err != nil {
		return fmt.Errorf("failed to get SMS service: %v", err)
	}

	smsScheduler := services.NewSchedulerService(smsService)
	smsScheduler.Start()
	defer smsScheduler.Stop()

	// Schedule multiple SMS notifications with different delays
	smsNotifications := []*models.Notification{
		{
			ID:          "3",
			Title:       "Appointment Reminder",
			Content:     "Your doctor's appointment is in 1 hour.",
			Channel:     models.ChannelMessage,
			Recipients:  []string{"+1234567890"},
			ScheduledAt: &time.Time{},
			CreatedAt:   time.Now(),
		},
		{
			ID:          "4",
			Title:       "Delivery Update",
			Content:     "Your package will arrive in 30 minutes.",
			Channel:     models.ChannelMessage,
			Recipients:  []string{"+1987654321"},
			ScheduledAt: &time.Time{},
			CreatedAt:   time.Now(),
		},
	}

	// Set different delays for SMS notifications
	smsNotifications[0].ScheduledAt = &time.Time{}
	*smsNotifications[0].ScheduledAt = time.Now().Add(10 * time.Second)
	smsNotifications[1].ScheduledAt = &time.Time{}
	*smsNotifications[1].ScheduledAt = time.Now().Add(15 * time.Second)

	for _, notification := range smsNotifications {
		if err := smsScheduler.ScheduleNotification(notification); err != nil {
			return fmt.Errorf("failed to schedule SMS notification: %v", err)
		}
	}

	// Create notification handler
	notificationHandler := handlers.NewNotificationHandler(a.notificationFactory, a.schedulerService)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/notifications", notificationHandler.SendNotification)

	// Create server
	a.server = &http.Server{
		Addr:    a.config.ServerPort,
		Handler: mux,
	}

	// Start HTTP server in a goroutine
	go func() {
		fmt.Printf("HTTP server listening on %s\n", a.config.ServerPort)
		if err := a.server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\nShutting down notification service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %v", err)
	}

	return nil
}

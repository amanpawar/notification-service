package services

import (
	"fmt"
	"notification-service/internal/models"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

type SchedulerService struct {
	cron                *cron.Cron
	notificationService NotificationService
	jobs                map[string]cron.EntryID
	mu                  sync.RWMutex
}

func NewSchedulerService(notificationService NotificationService) *SchedulerService {
	return &SchedulerService{
		cron:                cron.New(cron.WithSeconds()),
		notificationService: notificationService,
		jobs:                make(map[string]cron.EntryID),
	}
}

func (s *SchedulerService) Start() {
	s.cron.Start()
}

func (s *SchedulerService) Stop() {
	s.cron.Stop()
}

func (s *SchedulerService) ScheduleNotification(notification *models.Notification) error {
	if notification.ScheduledAt == nil {
		return fmt.Errorf("scheduled time is required")
	}

	delay := notification.ScheduledAt.Sub(time.Now())
	if delay <= 0 {
		return fmt.Errorf("scheduled time must be in the future")
	}

	// Create a one-time job that will run at the scheduled time
	job := func() {
		if err := s.notificationService.Send(notification); err != nil {
			fmt.Printf("Error sending notification: %v\n", err)
		}
		// Remove the job after execution
		s.mu.Lock()
		if entryID, exists := s.jobs[notification.ID]; exists {
			s.cron.Remove(entryID)
			delete(s.jobs, notification.ID)
		}
		s.mu.Unlock()
	}

	// Schedule the job
	entryID, err := s.cron.AddFunc("@every 1s", func() {
		now := time.Now()
		if now.After(*notification.ScheduledAt) || now.Equal(*notification.ScheduledAt) {
			job()
		}
	})

	if err != nil {
		return fmt.Errorf("failed to schedule notification: %v", err)
	}

	// Store the job ID
	s.mu.Lock()
	s.jobs[notification.ID] = entryID
	s.mu.Unlock()

	fmt.Printf("Scheduled notification for %s\n", notification.ScheduledAt)
	return nil
}

type notificationJob struct {
	notification *models.Notification
	service      NotificationService
}

func (j *notificationJob) Run() {
	if err := j.service.Send(j.notification); err != nil {
		fmt.Printf("Error sending notification: %v\n", err)
	}
}

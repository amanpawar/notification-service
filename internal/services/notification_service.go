package services

import (
	"fmt"
	"notification-service/internal/models"
)

type NotificationService interface {
	Send(notification *models.Notification) error
}

type SlackNotificationService struct{}

func (s *SlackNotificationService) Send(notification *models.Notification) error {
	fmt.Printf("[SLACK] Sending notification to %v: %s - %s\n",
		notification.Recipients,
		notification.Title,
		notification.Content)
	return nil
}

type EmailNotificationService struct{}

func (e *EmailNotificationService) Send(notification *models.Notification) error {
	fmt.Printf("[EMAIL] Sending notification to %v: %s - %s\n",
		notification.Recipients,
		notification.Title,
		notification.Content)
	return nil
}

type MessageNotificationService struct{}

func (m *MessageNotificationService) Send(notification *models.Notification) error {
	fmt.Printf("[MESSAGE] Sending notification to %v: %s - %s\n",
		notification.Recipients,
		notification.Title,
		notification.Content)
	return nil
}

type NotificationServiceFactory struct {
	services map[models.NotificationChannel]NotificationService
}

func NewNotificationServiceFactory() *NotificationServiceFactory {
	return &NotificationServiceFactory{
		services: map[models.NotificationChannel]NotificationService{
			models.ChannelSlack:   &SlackNotificationService{},
			models.ChannelEmail:   &EmailNotificationService{},
			models.ChannelMessage: &MessageNotificationService{},
		},
	}
}

func (f *NotificationServiceFactory) GetService(channel models.NotificationChannel) (NotificationService, error) {
	service, exists := f.services[channel]
	if !exists {
		return nil, fmt.Errorf("unsupported notification channel: %s", channel)
	}
	return service, nil
}

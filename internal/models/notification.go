package models

import "time"

type NotificationChannel string

const (
	ChannelSlack   NotificationChannel = "slack"
	ChannelEmail   NotificationChannel = "email"
	ChannelMessage NotificationChannel = "message"
)

type Notification struct {
	ID          string
	Title       string
	Content     string
	Channel     NotificationChannel
	Recipients  []string
	ScheduledAt *time.Time
	CreatedAt   time.Time
	SentAt      *time.Time
}

type User struct {
	ID       string
	Name     string
	Email    string
	SlackID  string
	Phone    string
	Metadata map[string]string
}

# Notification Service

A flexible notification service written in Go that supports multiple notification channels (Slack, Email, SMS) with scheduling capabilities.

## Features

- Multiple notification channels (Slack, Email, SMS)
- Scheduled notifications with configurable delays
- Factory pattern for notification services
- Clean architecture with separation of concerns
- Comprehensive test coverage
- Graceful shutdown handling

## Project Structure

```.
├── internal/
│   ├── app/          # Application setup and initialization
│   ├── config/       # Configuration management
│   ├── models/       # Data models
│   └── services/     # Business logic and services
├── go.mod           # Go module file
├── main.go          # Entry point
└── README.md        # This file
```

## Getting Started

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

3. Run tests:
   ```bash
   # Run all tests
   go test ./internal/services -v

   # Run specific test
   go test ./internal/services -v -run TestSlackNotificationService
   ```

## Usage Examples

The service currently supports three notification channels:
- Slack
- Email
- SMS

### Example Output
```
Starting notification service examples...
[SLACK] Sending notification to [user1 user2 user3]: Team Meeting Reminder - Don't forget about the team meeting at 2 PM today!
Scheduled notification for 2025-03-31 15:30:00 +0000 UTC
[EMAIL] Sending notification to [manager@company.com hr@company.com]: Weekly Report Ready - Your weekly performance report is now available.
[MESSAGE] Sending notification to [+1234567890]: Appointment Reminder - Your doctor's appointment is in 1 hour.
```

## Test Coverage

The service includes comprehensive tests covering:

1. **Basic Notification Services**
   - Slack notification sending
   - Email notification sending
   - SMS notification sending

2. **Factory Pattern**
   - Service creation for all channels
   - Error handling for invalid channels

3. **Scheduler Service**
   - Single notification scheduling
   - Multiple notifications with different delays
   - Invalid scheduling time handling
   - Nil scheduling time handling

### Running Tests

All tests can be run with:
```bash
go test ./internal/services -v
```

Test output will show:
```
=== RUN   TestSlackNotificationService
[SLACK] Sending notification to [test-user]: Test Slack Notification - This is a test notification
--- PASS: TestSlackNotificationService (0.00s)
=== RUN   TestEmailNotificationService
[EMAIL] Sending notification to [test@example.com]: Test Email Notification - This is a test email
--- PASS: TestEmailNotificationService (0.00s)
```

## Architecture

The service follows clean architecture principles:

1. **Models Layer**
   - Defines core data structures
   - Contains notification and user models
   - Defines notification channels

2. **Services Layer**
   - Implements notification sending logic
   - Handles scheduling and timing
   - Uses factory pattern for service creation

3. **Application Layer**
   - Coordinates between services
   - Handles application lifecycle
   - Manages graceful shutdown

## Future Improvements

1. **Channel Integration**
   - Implement actual Slack API integration
   - Add SMTP support for email
   - Integrate with SMS providers

2. **Features**
   - Add notification templates
   - Implement retry mechanisms
   - Add notification status tracking
   - Support for notification groups

3. **Testing**
   - Add integration tests
   - Implement mock services
   - Add performance benchmarks

## API Documentation

The service now provides a RESTful API for sending notifications:

### Send Notification

**Endpoint**: `POST /notifications`

**Request Body**:
```json
{
    "title": "Notification Title",
    "content": "Notification content",
    "channel": "slack|email|message",
    "recipients": ["user1", "user@example.com", "+1234567890"],
    "scheduled_at": "2025-03-31T15:30:00Z"
}
```

**Success Response** (200 OK for immediate, 202 Accepted for scheduled):
```json
{
    "success": true,
    "message": "Notification sent successfully",
    "data": {
        "id": "...",
        "title": "Notification Title",
        "content": "Notification content",
        "channel": "slack",
        "recipients": ["user1"],
        "created_at": "2025-03-31T15:30:00Z"
    }
}
```

**Error Response** (400 Bad Request):
```json
{
    "success": false,
    "message": "Error message describing the issue"
}
```

### Example API Usage

1. **Send immediate Slack notification**:
```bash
curl -X POST http://localhost:8080/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Team Meeting",
    "content": "Team meeting at 2 PM",
    "channel": "slack",
    "recipients": ["user1", "user2"]
  }'
```

2. **Schedule email notification**:
```bash
curl -X POST http://localhost:8080/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Weekly Report",
    "content": "Your weekly report is ready",
    "channel": "email",
    "recipients": ["user@example.com"],
    "scheduled_at": "2025-03-31T15:30:00Z"
  }'
```

3. **Send SMS notification**:
```bash
curl -X POST http://localhost:8080/notifications \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Appointment Reminder",
    "content": "Your appointment is in 1 hour",
    "channel": "message",
    "recipients": ["+1234567890"]
  }'
```

### API Error Cases

The API handles various error cases with appropriate status codes:

- `400 Bad Request`: Invalid input (missing fields, invalid channel, etc.)
- `405 Method Not Allowed`: Wrong HTTP method
- `500 Internal Server Error`: Server-side issues

Common error cases:
1. Missing required fields (title, content)
2. Empty recipients list
3. Invalid notification channel
4. Invalid scheduled time format
5. Past scheduled time
6. Invalid HTTP method

### API Testing

The API endpoints can be tested using:

```bash
# Run API tests
go test ./internal/handlers -v

# Run specific API test
go test ./internal/handlers -v -run TestNotificationHandler
```
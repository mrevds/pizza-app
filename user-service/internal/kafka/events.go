package kafka

import "time"

type UserRegisteredEvent struct {
    UserID      string    `json:"user_id"`
    PhoneNumber string    `json:"phone_number"`
    FirstName   string    `json:"first_name"`
    Timestamp   time.Time `json:"timestamp"`
}

type UserLoggedInEvent struct {
    UserID    string    `json:"user_id"`
    Timestamp time.Time `json:"timestamp"`
}
package models

type UserEvent struct {
	UserID  string
	EventID string
}

func (e *UserEvent) TableName() string {
	return "user_event"
}

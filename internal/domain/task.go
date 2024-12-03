package domain

import "time"

type Task struct {
	Id          uint64
	UserId      uint64
	Title       string
	Description string
	Status      Status
	Date        time.Time
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type Status string

const (
	NewTaskStatus       Status = "NEW"
	DoneTaskStatus      Status = "DONE"
	ImportantTaskStatus Status = "IMPORTANT"
	ExpiredTaskStatus   Status = "EXPIRED"
)

package requests

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type TaskRequest struct {
	Title       string        `json:"title" validate:"required"`
	Description string        `json:"description"`
	Status      domain.Status `json:"status" validate:"oneof=NEW IMPORTANT"`
	Date        int64         `json:"date"`
}

func (r TaskRequest) ToDomainModel() (interface{}, error) {
	return domain.Task{
		Title:       r.Title,
		Description: r.Description,
		Status:      r.Status,
		Date:        time.Unix(r.Date, 0),
	}, nil
}

package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type TaskDto struct {
	Id          uint64        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      domain.Status `json:"status"`
	Date        time.Time     `json:"date"`
}

func (d TaskDto) DomainToDto(t domain.Task) TaskDto {
	return TaskDto{
		Id:          t.Id,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Date:        t.Date,
	}
}

type ShortTaskDto struct {
	Id     uint64        `json:"id"`
	Title  string        `json:"title"`
	Status domain.Status `json:"status"`
	Date   time.Time     `json:"date"`
}

type TasksDto struct {
	Tasks []ShortTaskDto `json:"tasks"`
}

func (d TasksDto) DomainToDto(ts []domain.Task) TasksDto {
	tasks := make([]ShortTaskDto, len(ts))
	for i, t := range ts {
		tasks[i] = ShortTaskDto{}.DomainToDto(t)
	}
	return TasksDto{
		Tasks: tasks,
	}
}

func (d ShortTaskDto) DomainToDto(t domain.Task) ShortTaskDto {
	return ShortTaskDto{
		Id:     t.Id,
		Title:  t.Title,
		Status: t.Status,
		Date:   t.Date,
	}
}

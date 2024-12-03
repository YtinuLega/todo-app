package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}

func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController -> Save -> request.Bind: %s", err)
			BadRequest(w, err)
			return
		}
		user := r.Context().Value(UserKey).(domain.User)
		task.UserId = user.Id

		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController -> Save -> c.taskService.Save: %s", err)
			InternalServerError(w, err)
			return
		}

		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(task)
		Created(w, taskDto)
	}
}

func (c TaskController) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskIdStr := chi.URLParam(r, "taskId")
		taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
		if err != nil {
			log.Printf("TaskController -> Update -> strconv.ParseUint: %s", err)
			BadRequest(w, errors.New("invalid task ID"))
			return
		}

		var updateData requests.TaskRequest
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			log.Printf("TaskController -> Update -> json.Decode: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)

		task := domain.Task{
			Id:          taskId,
			UserId:      user.Id,
			Title:       updateData.Title,
			Description: updateData.Description,
			Status:      updateData.Status,
			Date:        time.Unix(updateData.Date, 0), // Преобразование Unix timestamp
		}

		updatedTask, err := c.taskService.Update(task)
		if err != nil {
			log.Printf("TaskController -> Update -> c.taskService.Update: %s", err)
			InternalServerError(w, err)
			return
		}

		taskDto := resources.TaskDto{}.DomainToDto(updatedTask)
		Success(w, taskDto)
	}
}

func (c TaskController) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskIdStr := chi.URLParam(r, "taskId")
		taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
		if err != nil {
			log.Printf("TaskController -> Delete -> strconv.ParseUint: %s", err)
			BadRequest(w, errors.New("invalid task ID"))
			return
		}

		user := r.Context().Value(UserKey).(domain.User)

		err = c.taskService.Delete(taskId, user.Id)
		if err != nil {
			log.Printf("TaskController -> Delete -> c.taskService.Delete: %s", err)
			InternalServerError(w, err)
			return
		}

		noContent(w)
	}
}

func (c TaskController) Find() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(chi.URLParam(r, "taskId"), 10, 64)
		if err != nil {
			log.Printf("TaskController -> ParseUint -> strconv.ParseUint(): %s", err)
			BadRequest(w, err)
			return
		}
		task, err := c.taskService.Find(id)
		if err != nil {
			log.Printf("TaskController -> Find -> c.taskService.Find: %s", err)
			InternalServerError(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		if user.Id != task.UserId {
			err = errors.New("access denied")
			Forbidden(w, err)
			return
		}

		var taskDto resources.TaskDto
		taskDto = taskDto.DomainToDto(task)
		Success(w, taskDto)
	}
}

func (c TaskController) FindList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserKey).(domain.User)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Получаем параметр status из query
		statusParam := r.URL.Query().Get("status")
		var statusFilter *domain.Status
		if statusParam != "" {
			status := domain.Status(statusParam)
			if !isValidStatus(status) {
				http.Error(w, "Invalid status parameter", http.StatusBadRequest)
				return
			}
			statusFilter = &status
		}

		// Передаем статус в сервис
		tasks, err := c.taskService.FindList(user.Id, statusFilter)
		if err != nil {
			log.Printf("TaskController -> FindList -> c.taskService.FindList: %s", err)
			InternalServerError(w, err)
			return
		}
		var tasksDto resources.TasksDto
		tasksDto = tasksDto.DomainToDto(tasks)
		Success(w, tasksDto)
	}
}

// Функция для проверки допустимости статуса
func isValidStatus(status domain.Status) bool {
	switch status {
	case domain.NewTaskStatus, domain.DoneTaskStatus, domain.ImportantTaskStatus, domain.ExpiredTaskStatus:
		return true
	default:
		return false
	}
}

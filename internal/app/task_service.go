package app

import (
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type TaskService struct {
	taskRepo database.TaskRepository
}

func NewTaskService(tr database.TaskRepository) TaskService {
	return TaskService{
		taskRepo: tr,
	}
}

func (s TaskService) Save(t domain.Task) (domain.Task, error) {
	tsk, err := s.taskRepo.Save(t)
	if err != nil {
		log.Printf("TaskService -> Save -> s.taskRepo.Save: %s", err)
		return domain.Task{}, err
	}

	return tsk, nil
}

func (s TaskService) Update(t domain.Task) (domain.Task, error) {
	tsk, err := s.taskRepo.Update(t)
	if err != nil {
		log.Printf("TaskService -> Update -> s.taskRepo.Update: %s", err)
		return domain.Task{}, err
	}

	return tsk, nil
}

func (s TaskService) Find(id uint64) (domain.Task, error) {
	tsk, err := s.taskRepo.Find(id)
	if err != nil {
		log.Printf("TaskService -> Find -> s.taskRepo.Find: %s", err)
		return domain.Task{}, err
	}

	return tsk, nil
}

func (s TaskService) FindList(uId uint64, status *domain.Status) ([]domain.Task, error) {
	tasks, err := s.taskRepo.FindList(uId, status)
	if err != nil {
		log.Printf("TaskService -> FindList -> s.taskRepo.FindList: %s", err)
		return nil, err
	}

	return tasks, nil
}

func (s TaskService) Delete(id uint64, userId uint64) error {
	err := s.taskRepo.Delete(id, userId)
	if err != nil {
		log.Printf("TaskService -> Delete -> s.taskRepo.Delete: %s", err)
		return err
	}
	return nil
}

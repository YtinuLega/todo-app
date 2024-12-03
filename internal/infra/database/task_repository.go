package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const TaskTableName = "tasks"

type task struct {
	Id          uint64        `db:"id,omitempty"`
	UserId      uint64        `db:"user_id"`
	Title       string        `db:"title"`
	Description string        `db:"description"`
	Status      domain.Status `db:"status"`
	Date        time.Time     `db:"date"`
	CreatedDate time.Time     `db:"created_date,omitempty"`
	UpdatedDate time.Time     `db:"updated_date,omitempty"`
	DeletedDate *time.Time    `db:"deleted_date,omitempty"`
}

type TaskRepository struct {
	coll db.Collection
	sess db.Session
}

func NewTaskRepository(sess db.Session) TaskRepository {
	return TaskRepository{
		coll: sess.Collection(TaskTableName),
		sess: sess,
	}
}

func (r TaskRepository) Save(t domain.Task) (domain.Task, error) {
	tsk := r.mapDomainToModel(t)
	tsk.CreatedDate = time.Now()
	tsk.UpdatedDate = time.Now()
	err := r.coll.InsertReturning(&tsk)
	if err != nil {
		return domain.Task{}, err
	}
	t = r.mapModelToDomain(tsk)
	return t, nil
}

func (r TaskRepository) Find(id uint64) (domain.Task, error) {
	var tsk task
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&tsk)
	if err != nil {
		return domain.Task{}, err
	}
	t := r.mapModelToDomain(tsk)
	return t, nil
}

func (r TaskRepository) FindList(uId uint64, status *domain.Status) ([]domain.Task, error) {
	var tasks []task
	cond := db.Cond{"user_id": uId, "deleted_date": nil}
	if status != nil {
		cond["status"] = *status
	}
	err := r.coll.Find(cond).All(&tasks)
	if err != nil {
		return nil, err
	}
	ts := r.mapModelToDomainCollection(tasks)
	return ts, nil
}

func (r TaskRepository) Update(d domain.Task) (domain.Task, error) {
	var existing task
	err := r.coll.Find(db.Cond{"id": d.Id, "user_id": d.UserId, "deleted_date": nil}).One(&existing)
	if err != nil {
		return domain.Task{}, err
	}

	existing.Title = d.Title
	existing.Description = d.Description
	existing.Status = d.Status
	existing.Date = d.Date
	existing.UpdatedDate = time.Now()

	err = r.coll.Find(existing.Id).Update(existing)
	if err != nil {
		return domain.Task{}, err
	}

	return r.mapModelToDomain(existing), nil
}

func (r TaskRepository) Delete(id uint64, userId uint64) error {
	var existing task
	err := r.coll.Find(db.Cond{"id": id, "user_id": userId, "deleted_date": nil}).One(&existing)
	if err != nil {
		return err
	}

	currentTime := time.Now()
	existing.DeletedDate = &currentTime

	err = r.coll.Find(existing.Id).Update(existing)
	if err != nil {
		return err
	}

	return nil
}

func (r TaskRepository) mapDomainToModel(d domain.Task) task {
	return task{
		Id:          d.Id,
		UserId:      d.UserId,
		Title:       d.Title,
		Description: d.Description,
		Status:      d.Status,
		Date:        d.Date,
		CreatedDate: d.CreatedDate,
		UpdatedDate: d.UpdatedDate,
		DeletedDate: d.DeletedDate,
	}
}

func (r TaskRepository) mapModelToDomain(m task) domain.Task {
	return domain.Task{
		Id:          m.Id,
		UserId:      m.UserId,
		Title:       m.Title,
		Description: m.Description,
		Status:      m.Status,
		Date:        m.Date,
		CreatedDate: m.CreatedDate,
		UpdatedDate: m.UpdatedDate,
		DeletedDate: m.DeletedDate,
	}
}

func (r TaskRepository) mapModelToDomainCollection(ts []task) []domain.Task {
	var tasks []domain.Task
	for _, t := range ts {
		tasks = append(tasks, r.mapModelToDomain(t))
	}
	return tasks
}

package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	results := []entity.Task{}
  	rows, err := r.db.WithContext(ctx).Table("tasks").Where("user_id =?", id).Select("*").Rows()
	if err != nil {
		return nil, err
	}
  	defer rows.Close()
  	for rows.Next() { // Next akan menyiapkan hasil baris berikutnya untuk dibaca dengan metode Scan.
	r.db.WithContext(ctx).ScanRows(rows, &results)
 	}
	return results, nil // TODO: replace this
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	 result := r.db.WithContext(ctx).Create(&task).Error
	 if result != nil {
		return 0, result
	 }
	 return task.ID, nil
}
	

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	hasil := entity.Task{}
	err := r.db.WithContext(ctx).Table("tasks").Where("id = ?", id).Take(&hasil).Error
	if err != nil {
		return entity.Task{}, err
	}
	return hasil, nil // TODO: replace this
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	results := []entity.Task{}
  	rows, err := r.db.WithContext(ctx).Table("tasks").Where("category_id =?", catId).Select("*").Rows()
	if err != nil {
		return nil, err
	}
  	defer rows.Close()
  	for rows.Next() { // Next akan menyiapkan hasil baris berikutnya untuk dibaca dengan metode Scan.
	r.db.WithContext(ctx).ScanRows(rows, &results)
 	}
	return results, nil // TODO: replace this
}


func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	err := r.db.WithContext(ctx).Model(&entity.Task{}).Where("id=?", task.ID).Updates(&task).Error
	if err != nil {
		return err
	}
	return nil// TODO: replace this
}


func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	cv := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Task{}).Error
	if cv != nil {
		return cv
	}
	return nil // TODO: replace this
}



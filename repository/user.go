package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	hasil := entity.User{}
	xs := r.db.WithContext(ctx).Table("users").Where("id = ?", id).Find(&hasil).Error
	if xs != nil {
		return entity.User{}, xs
	}
	return hasil, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	hasil := entity.User{}
	err := r.db.WithContext(ctx).Table("users").Where("email = ?", email).Find(&hasil).Error
	if err != nil {
		return entity.User{}, err
	}
	return hasil, nil 
}


func (r *userRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {

	tr := r.db.WithContext(ctx).Create(&user).Error
	if tr != nil {
		return entity.User{}, tr
	}
	return user, nil 
}

func (r *userRepository) UpdateUser(ctx context.Context, user entity.User) (entity.User, error) {
	ub := entity.User{}
	tx := r.db.WithContext(ctx).Model(&ub).Where("id =?", user.ID).Updates(user).Error
	if tx != nil {
		return entity.User{}, tx
	}

	return user, nil 
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	po := entity.User{}
	cv := r.db.WithContext(ctx).Where("id = ?", id).Delete(&po).Error
	if cv != nil {
		return cv
	}
	return nil // TODO: replace this
}

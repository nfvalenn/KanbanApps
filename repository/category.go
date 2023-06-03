package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error)
	StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error)
	StoreManyCategory(ctx context.Context, categories []entity.Category) error
	GetCategoryByID(ctx context.Context, id int) (entity.Category, error)
	UpdateCategory(ctx context.Context, category *entity.Category) error
	DeleteCategory(ctx context.Context, id int) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetCategoriesByUserId(ctx context.Context, id int) ([]entity.Category, error) {
	results := []entity.Category{}
  	rows, err := r.db.WithContext(ctx).Table("categories").Where("user_id =?", id).Select("*").Rows()
	if err != nil {
		return nil, err
	}
  	defer rows.Close()
  	for rows.Next() { // Next akan menyiapkan hasil baris berikutnya untuk dibaca dengan metode Scan.
	r.db.WithContext(ctx).ScanRows(rows, &results)
 	}	
	return results, nil // TODO: replace this
}



func (r *categoryRepository) StoreCategory(ctx context.Context, category *entity.Category) (categoryId int, err error) {
	// := entity.User{}
	tr := r.db.WithContext(ctx).Create(&category).Error
	if tr != nil {
		return 0, tr
	}
	return category.ID, nil // TODO: replace this
}

func (r *categoryRepository) StoreManyCategory(ctx context.Context, categories []entity.Category) error {
	//ox := []entity.User{}
	tr := r.db.WithContext(ctx).Create(&categories).Error
	if tr != nil {
		return tr
	}
	return nil // TODO: replace this
}

func (r *categoryRepository) GetCategoryByID(ctx context.Context, id int) (entity.Category, error) {
	hasil := entity.Category{}
	err := r.db.WithContext(ctx).Table("categories").Where("id = ?", id).Take(&hasil).Error
	if err != nil {
		return entity.Category{}, err
	}
	return hasil, nil // TODO: replace this
}

func (r *categoryRepository) UpdateCategory(ctx context.Context, category *entity.Category) error {
	err := r.db.WithContext(ctx).Model(&entity.Category{}).Where("id =?", category.ID).Updates(category).Error
	if err != nil {
		return err
	}
	return nil// TODO: replace this
}


func (r *categoryRepository) DeleteCategory(ctx context.Context, id int) error {
	
	cv := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Category{}).Error
	if cv != nil {
		return cv
	}
	return nil // TODO: replace this
}


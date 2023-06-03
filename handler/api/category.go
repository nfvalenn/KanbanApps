package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"

	"log"
	"net/http"
	"strconv"
)

type CategoryAPI interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTasks(w http.ResponseWriter, r *http.Request)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (c *categoryAPI) GetCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(string) //mengambil id dari middleware auth
	if userId == "" {
		w.WriteHeader(400)
		err := entity.ErrorResponse{Error: "invalid user id"}
		json, _ := json.Marshal(err)
		w.Write(json)
		return
	}
	IntUserId, _ := strconv.Atoi(userId)
	category, err := c.categoryService.GetCategories(r.Context(), IntUserId)
	if err != nil {
		w.WriteHeader(500)
		errResp := entity.ErrorResponse{Error: "error internal server"}
		json, _ := json.Marshal(errResp)
		w.Write(json)
		return
	}
	w.WriteHeader(200)
	json, _ := json.Marshal(category)
	w.Write(json)
}

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category entity.CategoryRequest

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	if category.Type == "" {
		w.WriteHeader(400)
		errResp := entity.ErrorResponse{Error: "invalid category request"}
		json, _ := json.Marshal(errResp)
		w.Write(json)
		return
	}
	userId := r.Context().Value("id").(string)
	if userId == "" {
		errResp := entity.ErrorResponse{Error: "invalid user id"}
		json, _ := json.Marshal(errResp)
		w.Write(json)
		return
	}
	UserIdInt,_ := strconv.Atoi(userId)
	Kategori := entity.Category{
		Type: category.Type,
		UserID: UserIdInt,
	}
	Cat, err := c.categoryService.StoreCategory(r.Context(), &Kategori)
	if err != nil {
		w.WriteHeader(500)
		errResp := entity.ErrorResponse{Error: "error internal server"}
		json, _ := json.Marshal(errResp)
		w.Write(json)
		return
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id" : UserIdInt,
		"category_id": Cat.ID,
		"message": "success create new category",
	})
}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(string)
	categoryID := r.URL.Query().Get("category_id")
	catIDInt, _ := strconv.Atoi(categoryID)
	del := c.categoryService.DeleteCategory(r.Context(), catIDInt)
	if del != nil {
		w.WriteHeader(500)
		err := entity.ErrorResponse{Error: "error internal server"}
		json, _ := json.Marshal(err)
		w.Write(json)
		return
	}
	IntUserId, _ := strconv.Atoi(userId)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id" : IntUserId,
		"category_id": catIDInt,
		"message": "success delete category",
	})
}

func (c *categoryAPI) GetCategoryWithTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category task", err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategoriesWithTasks(r.Context(), int(idLogin))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)

}

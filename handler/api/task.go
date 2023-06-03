package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(string)
	if userId == "" {
		w.WriteHeader(400)
		errResp := entity.ErrorResponse{Error: "invalid user id"}
		json, _ := json.Marshal(errResp)
		w.Write(json)
		return
	}
	UserIdInt, _ := strconv.Atoi(userId)
	taskId := r.URL.Query().Get("task_id")
	if taskId == "" {
		task, err := t.taskService.GetTasks(r.Context(), UserIdInt)
		if err != nil {
			w.WriteHeader(500)
			errResp := entity.ErrorResponse{Error: "error internal server"}
			json, _ := json.Marshal(errResp)
			w.Write(json)
			return
		}
		w.WriteHeader(200)
		json, _ := json.Marshal(task)
		w.Write(json)
	} else {
		taskIdInt, _ := strconv.Atoi(taskId)
		task, err := t.taskService.GetTaskByID(r.Context(), taskIdInt)
		if err != nil {
			w.WriteHeader(500)
			errResp := entity.ErrorResponse{Error: "error internal server"}
			json, _ := json.Marshal(errResp)
			w.Write(json)
			return
		}
		w.WriteHeader(200)
		json, _ := json.Marshal(task)
		w.Write(json)
	
	}
}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}
	if task.Title == "" || task.Description == "" || task.CategoryID == 0 {
		w.WriteHeader(400)
		err := entity.ErrorResponse{Error: "invalid task request"}
		json, _ := json.Marshal(err)
		w.Write(json)
		return
	}
	userId :=r.Context().Value("id").(string)
	if userId == "" {
		w.WriteHeader(400)
		err := entity.ErrorResponse{Error: "invalid user id"}
		json, _ := json.Marshal(err)
		w.Write(json)
		return
	}
	UserIdInt, _ := strconv.Atoi(userId)
	taks := entity.Task{
		ID: task.ID,
		Title: task.Title,
		Description: task.Description,
		CategoryID: task.CategoryID,
		UserID: UserIdInt,
	}
	taskE, err := t.taskService.StoreTask(r.Context(), &taks)
		if err != nil {
			w.WriteHeader(500)
			errResp := entity.ErrorResponse{Error: "error internal server"}
			json, _ := json.Marshal(errResp)
			w.Write(json)
			return
		}	
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"user_id": UserIdInt,
			"task_id": taskE.ID,
			"message": "success create new task",
	})
}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id").(string)
	taskID := r.URL.Query().Get("task_id")

	taskIdInt, _ := strconv.Atoi(taskID)
	task := t.taskService.DeleteTask(r.Context(), taskIdInt)
	if task != nil {
		w.WriteHeader(500)
		errResp := entity.ErrorResponse{Error: "error internal server"}
		json, _ := json.Marshal(errResp)
		w.Write(json)
		return
	}
	userIdInt, _ := strconv.Atoi(userId)
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id" : userIdInt,
		"task_id": taskIdInt,
		"message": "success delete task",
	})


}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id").(string)
	if userId == "" {
		w.WriteHeader(400)
		errResp := entity.ErrorResponse{Error: "invalid user id"}
		json, _ := json.Marshal(errResp)
		w.Write(json)
		return
	}
	UserIdInt, _ := strconv.Atoi(userId)
	ktask := entity.Task{
		ID: task.ID,
		Title: task.Title,
		Description: task.Description,
		CategoryID: task.CategoryID,
		UserID: UserIdInt,
	}
	taskk, err := t.taskService.UpdateTask(r.Context(), &ktask)
	if err != nil {
		w.WriteHeader(500)
		errResp := entity.ErrorResponse{Error: "error internal server"}
		json, _:= json.Marshal(errResp)
		w.Write(json)
		return
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id" : ktask.UserID,
		"task_id": taskk.ID,
		"message": "success update task",
	})
	return
}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}

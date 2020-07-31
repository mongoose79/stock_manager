package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"process_order_service"
	"utils"
)

// Request example: [POST] http://localhost:8080/api/v1/order	Body: ["pasta", "bread"]
func OrderHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Order request was received")
	processTask(w, r, true)
}

// Request example: [POST] http://localhost:8080/api/v1/supply	Body: ["bread"]
func SupplyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Supply request was received")
	processTask(w, r, false)
}

// Request example: [PUT] http://localhost:8080/api/v1/tasks/2/complete
func TaskCompleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Task complete request was received")
	vars := mux.Vars(r)
	taskId := vars["id"]
	if taskId == "" {
		errMsg := "Task ID is invalid"
		log.Println(errMsg)
		utils.WriteJSON(errMsg, w, http.StatusBadRequest)
		return
	}
	pos := process_order_service.NewProcessOrderService()
	err := pos.CompleteTask(taskId)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to complete the task %s", taskId)
		log.Println(errMsg)
		utils.WriteJSON(errMsg, w, http.StatusInternalServerError)
		return
	}
	log.Println("Task complete request was completed successfully")
}

// Request example: [GET] http://localhost:8080/api/v1/next-tasks
func NextTasksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get next tasks request was received")
	pos := process_order_service.NewProcessOrderService()
	tasks := pos.GetNextTasks()
	utils.WriteJSON(tasks, w, http.StatusOK)
	log.Println("Get next tasks request was completed successfully")
}

func processTask(w http.ResponseWriter, r *http.Request, isOrder bool) {
	var items []string
	bodyBytes, err := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(bodyBytes, &items)
	if err != nil {
		errMsg := "Failed to unmarshal request's body"
		log.Println(errMsg)
		utils.WriteJSON(errMsg, w, http.StatusBadRequest)
		return
	}
	pos := process_order_service.NewProcessOrderService()
	go pos.CreateNewTask(items, isOrder)
}

package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"models"
	"net/http"
	"process_order_service"
	"sort"
	"strings"
	"testing"
	"utils"
)

func TestOrderHandler(t *testing.T) {
	url := fmt.Sprintf("/order")
	request, err := http.NewRequest("POST", url, strings.NewReader(`["pasta", "bread"]`))
	assert.NoError(t, err)

	response := utils.InvokeRequest(request, OrderHandler, "/order")
	assert.Equal(t, http.StatusOK, response.Code)

	pos := process_order_service.NewProcessOrderService()
	assert.Equal(t, 2, len(pos.TasksMap))
	assert.Equal(t, "0", pos.TasksMap["0"].Id)
	assert.Equal(t, "pick_from_stock", pos.TasksMap["0"].Action)
	assert.Equal(t, "pasta", pos.TasksMap["0"].Product)
	assert.Equal(t, [2]int{9, 2}, pos.TasksMap["0"].Location)
	assert.Equal(t, "1", pos.TasksMap["1"].Id)
	assert.Equal(t, "pick_from_stock", pos.TasksMap["1"].Action)
	assert.Equal(t, "bread", pos.TasksMap["1"].Product)
	assert.Equal(t, [2]int{1, 5}, pos.TasksMap["1"].Location)
}

func TestSupplyHandler(t *testing.T) {
	url := fmt.Sprintf("/supply")
	request, err := http.NewRequest("POST", url, strings.NewReader(`["bread"]`))
	assert.NoError(t, err)

	response := utils.InvokeRequest(request, SupplyHandler, "/supply")
	assert.Equal(t, http.StatusOK, response.Code)

	pos := process_order_service.NewProcessOrderService()
	assert.Equal(t, "0", pos.TasksMap["0"].Id)
	assert.Equal(t, "put_to_stock", pos.TasksMap["0"].Action)
	assert.Equal(t, "bread", pos.TasksMap["0"].Product)
	assert.Equal(t, [2]int{1, 5}, pos.TasksMap["0"].Location)
}

func TestNextTasksHandler(t *testing.T) {
	url := fmt.Sprintf("/order")
	request, err := http.NewRequest("POST", url, strings.NewReader(`["pasta", "bread"]`))
	assert.NoError(t, err)

	response := utils.InvokeRequest(request, OrderHandler, "/order")
	assert.Equal(t, http.StatusOK, response.Code)

	url = fmt.Sprintf("/supply")
	request, err = http.NewRequest("POST", url, strings.NewReader(`["bread"]`))
	assert.NoError(t, err)

	response = utils.InvokeRequest(request, SupplyHandler, "/supply")
	assert.Equal(t, http.StatusOK, response.Code)

	url = fmt.Sprintf("/next-tasks")
	request, err = http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	response = utils.InvokeRequest(request, NextTasksHandler, "/next-tasks")
	assert.Equal(t, http.StatusOK, response.Code)

	var tasks []models.Task
	err = json.Unmarshal(response.Body.Bytes(), &tasks)
	assert.NoError(t, err)

	assert.Equal(t, 3, len(tasks))

	assert.Equal(t, "0", tasks[0].Id)
	assert.Equal(t, "pick_from_stock", tasks[0].Action)
	assert.Equal(t, "pasta", tasks[0].Product)
	assert.Equal(t, [2]int{9, 2}, tasks[0].Location)

	assert.Equal(t, "1", tasks[1].Id)
	assert.Equal(t, "pick_from_stock", tasks[1].Action)
	assert.Equal(t, "bread", tasks[1].Product)
	assert.Equal(t, [2]int{1, 5}, tasks[1].Location)

	assert.Equal(t, "2", tasks[2].Id)
	assert.Equal(t, "put_to_stock", tasks[2].Action)
	assert.Equal(t, "bread", tasks[2].Product)
	assert.Equal(t, [2]int{1, 5}, tasks[2].Location)
}

func TestTaskCompleteHandler(t *testing.T) {
	url := fmt.Sprintf("/order")
	request, err := http.NewRequest("POST", url, strings.NewReader(`["pasta", "bread"]`))
	assert.NoError(t, err)

	response := utils.InvokeRequest(request, OrderHandler, "/order")
	assert.Equal(t, http.StatusOK, response.Code)

	url = fmt.Sprintf("/supply")
	request, err = http.NewRequest("POST", url, strings.NewReader(`["bread"]`))
	assert.NoError(t, err)

	response = utils.InvokeRequest(request, SupplyHandler, "/supply")
	assert.Equal(t, http.StatusOK, response.Code)

	url = fmt.Sprintf("/tasks/0/complete")
	request, err = http.NewRequest("PUT", url, nil)
	assert.NoError(t, err)

	response = utils.InvokeRequest(request, TaskCompleteHandler, "/tasks/{id}/complete")
	assert.Equal(t, http.StatusOK, response.Code)

	url = fmt.Sprintf("/tasks/1/complete")
	request, err = http.NewRequest("PUT", url, nil)
	assert.NoError(t, err)

	response = utils.InvokeRequest(request, TaskCompleteHandler, "/tasks/{id}/complete")
	assert.Equal(t, http.StatusOK, response.Code)

	url = fmt.Sprintf("/tasks/2/complete")
	request, err = http.NewRequest("PUT", url, nil)
	assert.NoError(t, err)

	response = utils.InvokeRequest(request, TaskCompleteHandler, "/tasks/{id}/complete")
	assert.Equal(t, http.StatusOK, response.Code)

	url = fmt.Sprintf("/stock")
	request, err = http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	response = utils.InvokeRequest(request, StockHandler, "/stock")
	assert.Equal(t, http.StatusOK, response.Code)

	var products []models.Product
	err = json.Unmarshal(response.Body.Bytes(), &products)
	assert.NoError(t, err)

	sort.Sort(models.ByProductName(products))

	assert.Equal(t, 5, len(products))
	assert.Equal(t, "bread", products[0].Name)
	assert.Equal(t, 10, products[0].Amount)
	assert.Equal(t, "milk", products[1].Name)
	assert.Equal(t, 10, products[1].Amount)
	assert.Equal(t, "pasta", products[2].Name)
	assert.Equal(t, 9, products[2].Amount)
	assert.Equal(t, "salt", products[3].Name)
	assert.Equal(t, 10, products[3].Amount)
	assert.Equal(t, "soap", products[4].Name)
	assert.Equal(t, 10, products[4].Amount)
}

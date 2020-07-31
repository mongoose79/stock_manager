package process_order_service

import (
	"fmt"
	"github.com/palantir/stacktrace"
	"log"
	"models"
	"sort"
	"strconv"
	"sync"
)

var lock = sync.RWMutex{}

type ProductDetails struct {
	ShelfLocation [2]int
	Amount        int
}

type ProcessOrderService struct {
	Stock          map[string]ProductDetails
	TasksMap       map[string]models.Task
	CurrentOrderId int
}

var processOrderServiceInstance *ProcessOrderService
var processOrderServiceOnce sync.Once

func NewProcessOrderService() *ProcessOrderService {
	processOrderServiceOnce.Do(func() {
		processOrderServiceInstance = &ProcessOrderService{
			Stock:    make(map[string]ProductDetails),
			TasksMap: make(map[string]models.Task),
		}
		processOrderServiceInstance.initStock()
		processOrderServiceInstance.CurrentOrderId = -1
	})
	return processOrderServiceInstance
}

func (p *ProcessOrderService) CreateNewTask(items []string, isOrder bool) {
	if items == nil {
		log.Println("Failed to create the task. Items list is empty.")
		return
	}
	for _, item := range items {
		p.CurrentOrderId++
		id := strconv.Itoa(p.CurrentOrderId)
		task := models.Task{Id: id, Action: p.getAction(isOrder), Product: item,
			Location: p.getShelfLocation(p.Stock[item].ShelfLocation)}
		p.updateTask(task)
	}
}

func (p *ProcessOrderService) CompleteTask(taskId string) error {
	task, err := p.getTaskByTaskId(taskId)
	if err != nil {
		return stacktrace.Propagate(err, "Failed to complete the task %s", taskId)
	}
	lock.RLock()
	defer lock.RUnlock()
	if p.Stock[task.Product].Amount == 0 && task.Action == "pick_from_stock" {
		errMsg := fmt.Sprintf("Cannot pick from the stock %s because it is absent", task.Product)
		log.Println(errMsg)
		return stacktrace.NewError(errMsg)
	}
	product := p.Stock[task.Product]
	if task.Action == "pick_from_stock" {
		product.Amount--
	} else {
		product.Amount++
	}
	p.Stock[task.Product] = product
	delete(p.TasksMap, task.Id)
	return nil
}

func (p *ProcessOrderService) GetNextTasks() []models.Task {
	lock.RLock()
	defer lock.RUnlock()
	var taskIdArr []string
	for taskId, _ := range p.TasksMap {
		taskIdArr = append(taskIdArr, taskId)
	}
	sort.Strings(taskIdArr)
	var tasks []models.Task
	for _, taskId := range taskIdArr {
		tasks = append(tasks, p.TasksMap[taskId])
	}
	return tasks
}

func (p *ProcessOrderService) GetStockData() []models.Product {
	lock.RLock()
	defer lock.RUnlock()
	var products []models.Product
	for productName, productDetails := range p.Stock {
		product := models.Product{Name: productName, Amount: productDetails.Amount}
		products = append(products, product)
	}
	return products
}

func (p *ProcessOrderService) initStock() {
	p.Stock["milk"] = ProductDetails{ShelfLocation: [2]int{5, 6}, Amount: 10}
	p.Stock["bread"] = ProductDetails{ShelfLocation: [2]int{1, 5}, Amount: 10}
	p.Stock["salt"] = ProductDetails{ShelfLocation: [2]int{6, 5}, Amount: 10}
	p.Stock["soap"] = ProductDetails{ShelfLocation: [2]int{8, 9}, Amount: 10}
	p.Stock["pasta"] = ProductDetails{ShelfLocation: [2]int{9, 2}, Amount: 10}
}

func (p *ProcessOrderService) getShelfLocation(loc [2]int) [2]int {
	lock.RLock()
	defer lock.RUnlock()
	return [2]int{loc[0], loc[1]}
}

func (p *ProcessOrderService) getTaskByTaskId(taskId string) (models.Task, error) {
	lock.RLock()
	defer lock.RUnlock()
	var task models.Task
	var isExist bool
	if task, isExist = p.TasksMap[taskId]; !isExist {
		errMsg := fmt.Sprintf("Task %s does not exist", taskId)
		return task, stacktrace.NewError(errMsg)
	}
	return task, nil
}

func (p *ProcessOrderService) updateTask(task models.Task) {
	lock.RLock()
	defer lock.RUnlock()
	p.TasksMap[task.Id] = task
}

func (p *ProcessOrderService) getAction(isOrder bool) string {
	var action string
	if isOrder {
		action = "pick_from_stock"
	} else {
		action = "put_to_stock"
	}
	return action
}

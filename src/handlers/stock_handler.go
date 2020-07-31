package handlers

import (
	"log"
	"net/http"
	"process_order_service"
	"utils"
)

// Request example: [GET] http://localhost:8080/api/v1/stock
func StockHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Request for the current stock state was received")
	pos := process_order_service.NewProcessOrderService()
	products := pos.GetStockData()
	utils.WriteJSON(products, w, http.StatusOK)
	log.Println("Current state of stock request was completed successfully")
}

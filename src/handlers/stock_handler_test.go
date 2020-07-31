package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"models"
	"net/http"
	"sort"
	"testing"
	"utils"
)

func TestStockHandler(t *testing.T) {
	url := fmt.Sprintf("/stock")
	request, err := http.NewRequest("GET", url, nil)
	assert.NoError(t, err)

	response := utils.InvokeRequest(request, StockHandler, "/stock")
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
	assert.Equal(t, 10, products[2].Amount)
	assert.Equal(t, "salt", products[3].Name)
	assert.Equal(t, 10, products[3].Amount)
	assert.Equal(t, "soap", products[4].Name)
	assert.Equal(t, 10, products[4].Amount)
}

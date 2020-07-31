package models

type Task struct {
	Id       string `json:"id"`
	Action   string `json:"action"`
	Product  string `json:"product"`
	Location [2]int `json:"location"`
}

type Product struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

type ByProductName []Product

func (a ByProductName) Len() int           { return len(a) }
func (a ByProductName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByProductName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

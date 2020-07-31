package utils

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
)

func WriteJSON(model interface{}, w http.ResponseWriter, header int) {
	JSON, err := json.Marshal(model)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(header)
	w.Write(JSON)
}

//  for tests
func InvokeRequest(request *http.Request, handle func(w http.ResponseWriter, r *http.Request), path string) *httptest.ResponseRecorder {
	f := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handle(w, r)
	})
	m := mux.NewRouter()
	m.HandleFunc(path, f).Methods(request.Method)

	response := httptest.NewRecorder()
	m.ServeHTTP(response, request)
	return response
}

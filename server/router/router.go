package router

import (
	"github.com/dtsmith94/shared-expenses-tracker/server/middleware"
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {

	router := mux.NewRouter()

	router.HandleFunc("/api/expense", middleware.GetAllExpenses).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/expense", middleware.CreateExpense).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/expense/{id}", middleware.EditExpense).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/expense/{id}", middleware.DeleteExpense).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/expense", middleware.DeleteAllExpenses).Methods("DELETE", "OPTIONS")
	return router
}

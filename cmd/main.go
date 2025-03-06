package main

import (
	"log"
	"net/http"
	orch "github.com/schmalz302/Distributed_Calculator/internal/orchestrator"
)

// Запускаем сервер
func main() {
	queue := orch.NewExpressionQueue()

	http.HandleFunc("/api/v1/calculate", queue.CRUD_AddExpression)
	http.HandleFunc("/api/v1/expressions", queue.CRUD_GetExpressions)
	http.HandleFunc("/api/v1/expressions/{id}", queue.CRUD_GetExpression_id)
	http.HandleFunc("/internal/task", queue.CRUD_ProcessTask)

	log.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
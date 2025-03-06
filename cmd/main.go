package main

import (
	"fmt"
	"log"
	"net/http"
	orch "github.com/schmalz302/Distributed_Calculator/internal/orchestrator"
)

// Запускаем сервер
func main() {
	queue := orch.NewExpressionQueue()

	http.HandleFunc("/api/v1/calculate", queue.CRUD_AddExpression)
	http.HandleFunc("/api/v1/expressions", queue.CRUD_GetExpressions)
	http.HandleFunc("/api/v1/expressions/:id", queue.CRUD_GetExpression_id)
	http.HandleFunc("/internal/task", queue.CRUD_ProcessTask)

	fmt.Println("Orchestrator running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
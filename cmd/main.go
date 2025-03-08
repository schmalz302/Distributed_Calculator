package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	orch "github.com/schmalz302/Distributed_Calculator/internal/orchestrator"
	agent "github.com/schmalz302/Distributed_Calculator/internal/agent"
)

// Запускаем сервер
func main() {
	// запускаем сервак
	// создаем очередь выражений и роутер
	queue := orch.NewExpressionQueue()
	router := mux.NewRouter()
	// привызяваем эндпоинты к роутеру
	router.HandleFunc("/api/v1/calculate", queue.CRUD_AddExpression)
	router.HandleFunc("/api/v1/expressions", queue.CRUD_GetExpressions)
	router.HandleFunc("/api/v1/expressions/{id}", queue.CRUD_GetExpression_id)
	router.HandleFunc("/internal/task", queue.CRUD_ProcessTask)

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Println("Server is running on http://localhost:8080")
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Запускаем агента в отдельной горутине
	go agent.Start()

	// Блокируем main, чтобы он не завершился
	select {}
}
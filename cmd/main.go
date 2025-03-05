package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Expression struct{}

type ExpressionQueue struct {
	mu          sync.Mutex
	expressions map[string]*Expression
	results     map[string]float64
}

// Создаем новый TaskQueue
func NewExpressionQueue() *ExpressionQueue {
	return &ExpressionQueue{
		expressions: make(map[string]*Expression),
		results:     make(map[string]float64),
	}
}

// Запускаем сервер
func main() {
	queue := NewExpressionQueue()

	http.HandleFunc("/api/v1/calculate", queue.AddExpression)
	http.HandleFunc("/api/v1/expressions", queue.GetExpressions)
	http.HandleFunc("/api/v1/expressions/:id", queue.GetExpression_id)
	http.HandleFunc("/internal/task", queue.ProcessTask)

	fmt.Println("Orchestrator running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (q *ExpressionQueue) AddExpression(w http.ResponseWriter, r *http.Request) {}

func (q *ExpressionQueue) GetExpressions(w http.ResponseWriter, r *http.Request) {}

func (q *ExpressionQueue) GetExpression_id(w http.ResponseWriter, r *http.Request) {
	// r.PathValue("id")
	// на всякий, но вроде работате только с 1.22
}

func (q *ExpressionQueue) ProcessTask(w http.ResponseWriter, r *http.Request) {}

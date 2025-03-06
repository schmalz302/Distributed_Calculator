package orchestrator

import (
	"encoding/json"
	"net/http"
)

// Создание нового TaskQueue
func NewExpressionQueue() *ExpressionQueue {
	return &ExpressionQueue{
		expressions: make(map[string]*Expression),
		pool_task:   make(map[string]*Task),
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type Request_Expressions struct {
	Expressions []Expression `json:"expressions"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func (q *ExpressionQueue) CRUD_AddExpression(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	q.AddExpression(req.Expression)
}

func (q *ExpressionQueue) CRUD_GetExpressions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

func (q *ExpressionQueue) CRUD_GetExpression_id(w http.ResponseWriter, r *http.Request) {
	// r.PathValue("id")
	// на всякий, но вроде работате только с 1.22
}

func (q *ExpressionQueue) CRUD_ProcessTask(w http.ResponseWriter, r *http.Request) {}

package orchestrator

import (
	"encoding/json"

	"net/http"

	"github.com/gorilla/mux"
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

type Request_Expression struct {
	Expression string `json:"expression"`
}

type Response_Expression struct {
	Id string `json:"id"`
}

type Response_Expression_id struct {
	Expression Expression `json:"expression"`
}

type Response_Expression_List struct {
	Expressions []Expression `json:"expressions"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ProcessTaskRequest struct {
	Id             string
	Result         float64
	Operation_time int
}

func writeErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func (q *ExpressionQueue) CRUD_AddExpression(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request_Expression
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id := q.AddExpression(req.Expression)

	response := Response_Expression{Id: id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (q *ExpressionQueue) CRUD_GetExpressions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	expression_list := q.GetAllExpressions()

	respose_list := []Expression{}

	for _, exp := range expression_list {
		respose_list = append(respose_list, *exp)
	}

	response := Response_Expression_List{Expressions: respose_list}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (q *ExpressionQueue) CRUD_GetExpression_id(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	expression, _ := q.GetExpressionid(id)

	response := Response_Expression_id{Expression: *expression}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (q *ExpressionQueue) CRUD_ProcessTask(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		response := q.GetTask()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	case http.MethodPost:
		var req ProcessTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		q.SubmitResult(req.Id, req.Result)

	default:
		writeErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

package orchestrator

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Создание новой очереди выражений
func NewExpressionQueue() *ExpressionQueue {
	return &ExpressionQueue{
		expressions: make(map[string]*Expression),
		pool_task:   make(map[string]*Task),
	}
}

// шаблоны для запросов и ответов
type Request struct {
	Expression string `json:"expression"`
}

type Response struct {
	Status string `json:"status"`
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

// запись ошибки в ответ
func writeErrorResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	statusCode := 200
	switch message {
	case "Internal server error":
		statusCode = http.StatusInternalServerError
	case "Method not allowed":
		statusCode = http.StatusMethodNotAllowed
	case "Invalid data":
		statusCode = http.StatusUnprocessableEntity

	case "Not found":
		statusCode = http.StatusNotFound
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func writeResponse(w http.ResponseWriter, response_obj any, status_code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status_code)
	json.NewEncoder(w).Encode(response_obj)
}

// CRUD операции 
// эндпоинт для создания выражения
func (q *ExpressionQueue) CRUD_AddExpression(w http.ResponseWriter, r *http.Request) {
	// перехватываем любую необработанную панику и отпраляем код 500
	defer func() {
		if r := recover(); r != nil {
			writeErrorResponse(w, "Internal server error")
		}
	}()
	// проверка на метод
	if r.Method != http.MethodPost {
		writeErrorResponse(w, "Method not allowed")
		return
	}

	// проверка на десериализацию данных
	var req Request_Expression
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorResponse(w, "Invalid data")
		return
	}

	// создаем объект выражения
	id, err := q.AddExpression(req.Expression)
	if err != nil {
		writeErrorResponse(w, err.Error())
		return
	}

	// записываем ответ
	writeResponse(w, Response_Expression{Id: id}, 201)
}

// эндпоинт для получения всех выражений
func (q *ExpressionQueue) CRUD_GetExpressions(w http.ResponseWriter, r *http.Request) {
	// перехватываем любую необработанную панику и отпраляем код 500
	defer func() {
		if r := recover(); r != nil {
			writeErrorResponse(w, "Internal server error")
			return
		}
	}()
	// проверка на метод
	if r.Method != http.MethodGet {
		writeErrorResponse(w, "Method not allowed")
		return
	}

	// получаем список выражений 
	respose_list := q.GetAllExpressions()

	// записываем ответ
	writeResponse(w, Response_Expression_List{Expressions: respose_list}, 200)
}

// эндпоинт для получения конкретного выражения по id 
func (q *ExpressionQueue) CRUD_GetExpression_id(w http.ResponseWriter, r *http.Request) {
	// перехватываем любую необработанную панику и отпраляем код 500
	defer func() {
		if r := recover(); r != nil {
			writeErrorResponse(w, "Internal server error")
			return
		}
	}()
	// проверка на метод
	if r.Method != http.MethodGet {
		writeErrorResponse(w, "Method not allowed")
		return
	}
	// создаем словарь с параметрами пути
	vars := mux.Vars(r)
	// вытягиваем id
	id, err := vars["id"]
	if !err {
		panic("")
	}
	// вытягиваем инфу о выражении через id 
	expression, err2 := q.GetExpressionid(id)
	if err2 != nil {
		writeErrorResponse(w, "Not found")
		return
	}
	// записываем ответ
	writeResponse(w, Response_Expression_id{Expression: *expression}, 200)
}

// эндпоинт обработки задач (отдача задачи и получение ответов)
func (q *ExpressionQueue) CRUD_ProcessTask(w http.ResponseWriter, r *http.Request) {
	// перехватываем любую необработанную панику и отпраляем код 500
	defer func() {
		if r := recover(); r != nil {
			writeErrorResponse(w, "Internal server error")
			return
		}
	}()
	switch r.Method {
	case http.MethodGet:
		// выдаем задачу
		response := q.GetTask()
		// если задачи нет, выдаем 404
		if response == nil {
			writeErrorResponse(w, "Not Found")
			return
		}
		// если задача есть, записываем ответ
		writeResponse(w, *response, 200)
	case http.MethodPost:
		// создаем объект результата 
		var req ProcessTaskRequest
		// проверка на десериализацию данных
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorResponse(w, "Invalid data")
			return
		}
		// проверяем на наличие данного id 
		err := q.SubmitResult(req.Id, req.Result)
		if err != nil {
			writeErrorResponse(w, err.Error())
			return
		}
		// запись ответа: подтверждение записи результата
		writeResponse(w, Response{Status: "OK"}, 200)

	default:
		// любой другой метод помимо get и post 
		writeErrorResponse(w, "Method not allowed")
		return
	}
}

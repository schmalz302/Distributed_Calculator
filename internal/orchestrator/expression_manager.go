package orchestrator

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/google/uuid"
)

// структура задачи
type Task struct {
	ID            string   `json:"id"`
	Op            string   `json:"operation"`
	Arg1          string   `json:"arg1"`
	Arg2          string   `json:"arg2"`
	Status        int      `json:"-"` // 1 - принята, 2 - готова к обработке, 3 - обработка, 4 - завершена 
	Result        string `json:"-"`
	Expression_id string   `json:"-"`
}

// очередь выражений
type ExpressionQueue struct {
	mu          sync.Mutex
	expressions map[string]*Expression
	pool_task   map[string]*Task
}

// выражение
type Expression struct {
	ID          string `json:"id"`
	Status      string `json:"status"` // "pending", "in_progress", "done"
	Result      string `json:"result"`
	count_tasks int    `json:"-"`
}

// Добавляем задачи в очередь
func (q *ExpressionQueue) AddExpression(expression string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	// формируем дерево ast
	node := ParseExpression(expression)

	// распределяем его по задачам
	tasks := []Task{}
	SplitTasks(node, &tasks)

	// создаем объект выражения
	exp_obj := Expression{ID: uuid.New().String(), count_tasks: len(tasks), Status: "pending"}

	// закидываем задачи в пул задач
	for _, task := range tasks {
		q.pool_task[exp_obj.ID] = &task
	}
}

func (q *ExpressionQueue) GetExpressionid(id string) (*Expression, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	v, _ := q.expressions[id]
	return v, nil
}

func (q *ExpressionQueue) GetAllExpressions() []*Expression {
	expressions := []*Expression{}
	for _, expr := range q.expressions {
		expressions = append(expressions, expr)
	}
	return expressions
}

func (q *ExpressionQueue) Update_task(task *Task) {
	q.mu.Lock()
	defer q.mu.Unlock()
	var check int
	if !isNumber(task.Arg1) {
		arg1 := q.pool_task[task.Arg1].Result
		if arg1 != "" {
			task.Arg1 = task.Arg1
			check += 1
		}
	}
	if !isNumber(task.Arg2) {
		arg2 := q.pool_task[task.Arg2].Result
		if arg2 != "" {
			task.Arg2 = task.Arg2
			check += 1
		}
	}
	if check == 2 {
		task.Status = 2
		q.expressions[task.Expression_id].Status = "in_progress"
	}
}

// Отдаем агенту первую доступную задачу
func (q *ExpressionQueue) GetTask() *Task {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, task := range q.pool_task {
		q.Update_task(task)
		if task.Status == 2 {
			task.Status = 3
			return task
		}
	}
	return nil
}

// Получаем результат от агента
func (q *ExpressionQueue) SubmitResult(id string, result float64) {
	q.mu.Lock()
	defer q.mu.Unlock()

	task, exists := q.pool_task[id]

	if !exists {
		task.Result = fmt.Sprintf("%v", result)
		q.expressions[task.Expression_id].count_tasks -= 1
		task.Status = 4
	}
	if q.expressions[task.Expression_id].count_tasks == 0 {
		q.expressions[task.Expression_id].Result = fmt.Sprintf("%v", result)
		q.expressions[task.Expression_id].Status = "done"
	}
}

// Проверяем, является ли строка числом
func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// в этом файле будет происходить вся логика чтения, записи и обновления
// выражений и задач


package orchestrator

import (
	"fmt"
	"strconv"
	"sync"
	"github.com/google/uuid"
)

// структура задачи
type Task struct {
	ID             string `json:"id"`
	Op             string `json:"operation"`
	Arg1           string `json:"arg1"`
	Arg2           string `json:"arg2"`
	Operation_time int    `json:"operation_time"`
	Status         int    `json:"Status"` // 1 - принята, 2 - готова к обработке, 3 - обработка, 4 - завершена
	Result         string `json:"-"`
	Expression_id  string `json:"-"`
}

// очередь выражений (на самом деле очереди не будет, есть просто пул задач,
// которые можно выполнять параллельно)
type ExpressionQueue struct {
	mu          sync.Mutex
	expressions map[string]*Expression
	pool_task   map[string]*Task
}

// стрктура выражения
type Expression struct {
	ID          string `json:"id"`
	Status      string `json:"status"` // "pending", "in_progress", "done"
	Result      string `json:"result"`
	count_tasks int    `json:"-"`
}

// Добавляем задачи в очередь
func (q *ExpressionQueue) AddExpression(expression string) string {
	q.mu.Lock()
	defer q.mu.Unlock()
	// формируем дерево ast
	node := ParseExpression(expression)

	// распределяем его по задачам
	tasks := []Task{}
	SplitTasks(node, &tasks)

	// создаем объект выражения
	id_exp := uuid.New().String()
	exp_obj := Expression{ID: id_exp, count_tasks: len(tasks), Status: "pending"}

	q.expressions[id_exp] = &exp_obj

	// закидываем задачи в пул задач
	for _, task := range tasks {
		t := task 
		t.Expression_id = id_exp
		q.pool_task[task.ID] = &t
	}
	return id_exp
}

// получаем задачу по id
func (q *ExpressionQueue) GetExpressionid(id string) (*Expression, error) {
	v, _ := q.expressions[id]
	return v, nil
}

// получаем все задачи
func (q *ExpressionQueue) GetAllExpressions() []*Expression {
	expressions := []*Expression{}
	for _, expr := range q.expressions {
		expressions = append(expressions, expr)
	}
	return expressions
}


// обновляем задачу, если ее аргументы вычислены
func (q *ExpressionQueue) Update_task(task *Task) {
	var check int
	if !isNumber(task.Arg1) {
		t1, r := q.pool_task[task.Arg1]
		if r {	
			arg1 := t1.Result
			if arg1 != "" {
				task.Arg1 = arg1
				check += 1
			}
		} 
	} else {
		check += 1
	}
	if !isNumber(task.Arg2) {
		t2, r := q.pool_task[task.Arg2]
		if r {	
			arg2 := t2.Result
			if arg2 != "" {
				task.Arg2 = arg2
				check += 1
			}
		} 
	} else {
		check += 1
	}
	if check == 2 {
		task.Status = 2
		q.expressions[task.Expression_id].Status = "in_progress"
	}
}

// Отдаем агенту первую доступную задачу
func (q *ExpressionQueue) GetTask() *Task {
	for _, task := range q.pool_task {
		if task.Status == 1 {
			q.Update_task(task)
		}
		if task.Status == 2 {
			task.Status = 3
			return task
		}
	}
	return nil
}

// Получаем результат от агента
func (q *ExpressionQueue) SubmitResult(id string, result float64) {
	task, err := q.pool_task[id]
	if err {

		task.Result = fmt.Sprintf("%v", result)
		if task.Status == 3 {
			q.expressions[task.Expression_id].count_tasks -= 1
			task.Status = 4
		}
		if q.expressions[task.Expression_id].count_tasks == 0 {
			q.expressions[task.Expression_id].Result = fmt.Sprintf("%v", result)
			q.expressions[task.Expression_id].Status = "done"
			// удаляем все остальные таски
			exp_id := task.Expression_id
			for _, task := range q.pool_task {
				if task.Expression_id == exp_id {
					delete(q.pool_task, task.ID)
				}
		}

	}
	}
}

// Проверяем, является ли строка числом
func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

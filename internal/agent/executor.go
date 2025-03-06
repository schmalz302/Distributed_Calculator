package agent

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"sync"
// 	"time"

// 	orch "github.com/schmalz302/Distributed_Calculator/internal/orchestrator"
// )

// func getComputingPower() int {
// 	// получаем количество горутин
// 	computingPower := os.Getenv("COMPUTING_POWER")
// 	if computingPower == "" {
// 		return 4 // По умолчанию 4 горутины
// 	}
// 	num, err := strconv.Atoi(computingPower)
// 	if err != nil {
// 		return 4
// 	}
// 	return num
// }

// func getTask() (*orch.Task, error) {
// 	// отправляем запрос в оркестратор
// 	resp, err := http.Get("http://localhost:8080/internal/task")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode == http.StatusNotFound {
// 		return nil, nil // Нет задач
// 	}

// 	// создаем задачу
// 	var task orch.Task
// 	json.NewDecoder(resp.Body).Decode(&task)
// 	return &task, nil
// }

// func executeTask(task *orch.Task) orch.TaskResult {
// 	var result float64

// 	arg1 := to_float_64(task.Arg1)
// 	arg2 := to_float_64(task.Arg2)

// 	switch task.Op {
// 	case "+":
// 		result = arg1 + arg2
// 	case "-":
// 		result = arg1 - arg2
// 	case "/":
// 		result = arg1 / arg1
// 	case "*":
// 		result = arg1 * arg
// 	}
// 	return orch.TaskResult{ID: task.ID, Result: result}
// }

// func sendResult(result orch.TaskResult) {
// 	data, _ := json.Marshal(result)
// 	http.Post("http://localhost:8080/internal/task/", "application/json", bytes.NewBuffer(data))
// }

// func to_float_64(number string) float64 {
// 	n, _ := strconv.ParseFloat(number, 64);
// 	return n 
// }


// // Воркеры (горутины) для вычислений
// func worker() {
// 	for {
// 		// получаем задачу
// 		task, err := getTask()
// 		if err != nil {
// 			// ошибка получения задачи
// 			time.Sleep(2 * time.Second)
// 			continue
// 		}
// 		// если задач нет, ждем
// 		if task == nil {
// 			time.Sleep(2 * time.Second) 
// 			continue
// 		}
// 		// Выполняем задачу
// 		result := executeTask(task)

// 		// Отправляем результат обратно
// 		sendResult(result)
// 	}
// }

// func main() {
// 	// Количество допустимых горутин
// 	numWorkers := getComputingPower()

// 	var wg sync.WaitGroup

// 	// Запускаем горутины (пул воркеров)
// 	for i := 0; i < numWorkers; i++ {
// 		wg.Add(1)
// 		go func () {
// 			defer wg.Done()
// 			worker()
// 		}()
// 	}

// 	wg.Wait() // Ждем завершения всех горутин (но на практике агент будет работать бесконечно)
// }
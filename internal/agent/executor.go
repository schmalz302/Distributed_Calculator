package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"

	orch "github.com/schmalz302/Distributed_Calculator/internal/orchestrator"
)

func getComputingPower() int {
	// Загружаем переменные из .env
	err := godotenv.Load("../.env")
	if err != nil {
		return 4 // Значение по умолчанию
	}

	// Читаем значение переменной
	computingPower := os.Getenv("COMPUTING_POWER")

	if computingPower == "" {
		return 4 // Значение по умолчанию
	}

	// переводим в число
	num, err := strconv.Atoi(computingPower)
	if err != nil {
		return 4 // Значение по умолчанию
	}
	// возвращаем значение
	return num
}


func get_time_of_operation(oper_time_const string) int {
	// Загружаем переменные из .env
	err := godotenv.Load("../.env")
	if err != nil {
		return 0 // Значение по умолчанию
	}

	// Читаем значение переменной
	computingPower := os.Getenv(oper_time_const)

	if computingPower == "" {
		return 0 // Значение по умолчанию
	}

	// переводим в число
	num, err := strconv.Atoi(computingPower)
	if err != nil {
		return 0 // Значение по умолчанию
	}
	// возвращаем значение
	return num
}

func getTask() (*orch.Task, error) {
	// отправляем запрос в оркестратор
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil || resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Нет задач
	}

	// создаем задачу
	var task orch.Task
	json.NewDecoder(resp.Body).Decode(&task)
	if task.Status == 0 {
		// log.Println("Ошибка:", errors.New("Пустая задача"))
		return nil, errors.New("Пустая задача")
	}
	return &task, nil
}

func executeTask(task *orch.Task) orch.ProcessTaskRequest {
	var result float64	

	arg1 := to_float_64(task.Arg1)
	arg2 := to_float_64(task.Arg2)

	switch task.Op {
	case "+":
		result = arg1 + arg2
		// добавим искусственную нагрзку
		time_op := get_time_of_operation("TIME_ADDITION_MS")
		time.Sleep(time.Millisecond * time.Duration(time_op))
	case "-":
		result = arg1 - arg2
		time_op := get_time_of_operation("TIME_SUBTRACTION_MS")
		time.Sleep(time.Millisecond * time.Duration(time_op))
	case "/":
		result = arg1 / arg2
		time_op := get_time_of_operation("TIME_DIVISIONS_MS")
		time.Sleep(time.Millisecond * time.Duration(time_op))
	case "*":
		result = arg1 * arg2
		time_op := get_time_of_operation("TIME_MULTIPLICATIONS_MS")
		time.Sleep(time.Millisecond * time.Duration(time_op))
	}
	return orch.ProcessTaskRequest{Id: task.ID, Result: result}
}

func sendResult(result orch.ProcessTaskRequest) {
	data, err := json.Marshal(result)
	if err != nil {
		log.Println("Ошибка при сериализации JSON:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("Ошибка при отправке запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Ошибка! Код: %d, Ответ: %s", resp.StatusCode, string(body))
	}
}

func to_float_64(number string) float64 {
	n, _ := strconv.ParseFloat(number, 64)
	return n
}

// Воркеры (горутины) для вычислений
func worker() {
	for {
		// получаем задачу
		task, err := getTask()
		if err != nil {
			// ошибка получения задачи
			continue
		}
		// если задач нет, ждем
		if task == nil {
			continue
		}
		// Выполняем задачу
		result := executeTask(task)
		sendResult(result)
	}
}

func Start() {
	// Количество допустимых горутин
	numWorkers := getComputingPower()

	var wg sync.WaitGroup

	// Запускаем горутины (пул воркеров)
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker()
		}()
	}

	wg.Wait() // Ждем завершения всех горутин (но на практике агент будет работать бесконечно)
}

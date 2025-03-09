# Сервис для вычисления арифметических выражений
## Описание

Distributed Calculator – это распределенный калькулятор, который позволяет выполнять математические операции параллельно с помощью горутин. При запуске сервиса оркестратор начинает принимать запросы, агент в свою очередь запускает 
какое-то количество горутин, которые постоянно просят задачи для выполнения.

## Структура проекта
```
.
│   .gitignore
│   go.mod
│   go.sum
│   README.md
│
├───.vscode
│       launch.json
│
├───cmd
│       main.go
│
└───internal
    ├───agent
    │       executor.go
    │
    └───orchestrator
            ast.go 
            expression_manager.go
            service.go
            task_splitter.go
```

## Как запустить

1. Убедитесь, что у вас установлен Go версии 1.21.4.
2. Клонируйте этот репозиторий и перейдите в него:
   ```cmd
   git clone https://github.com/schmalz302/Distributed_Calculator
   cd Calc
   ```

3. Запустите сервис с помощью команды:
   ```cmd
   go run cmd/main.go
   ```

4. Сервис будет доступен по адресу: [http://localhost:8080](http://localhost:8080)

## Использование
Рекомендуется использовать `curl`, Postman или аналогичный инструмент для проверки работы сервиса. Проверьте все сценарии: корректные выражения, некорректные данные и симуляцию внутренних ошибок. Советую использовать Postman.

## Сценарии использования /api/v1/calculate

| **Request Method** | **Endpoint** | **Request Body**                                           | **Response Body**                                    | **HTTP Status Code** |
|--------------------|--------------|------------------------------------------------------------|------------------------------------------------------|----------------------|
| POST               | `/api/v1/calculate`  | `{ "expression": "2 + 2" }`                               | `{ "id": "какой-то id"}`                             | 201 OK               |
| POST               | `/api/v1/calculate`  | `{ "expression": "2 / 0" }`                               | `{"error": "Internal server error"}`                 | 500 Internal Server Error |
| POST               | `/api/v1/calculate`  | `{ "expression": "invalid expression" }`                  | `{ "error": "Invalid expression" }`                  | 422 Unprocessable Entity |
| GET                | `/api/v1/calculate`  | N/A                                                       | `{ "error": "Method not allowed" }`                  | 405 Method Not Allowed |


## Сценарии использования /api/v1/expressions

| **Request Method** | **Endpoint** | **Request Body**                                           | **Response Body**                                    | **HTTP Status Code** |
|--------------------|--------------|------------------------------------------------------------|------------------------------------------------------|----------------------|
| GET               | `/api/v1/expressions`  | N/A                                | `{"expressions": [{"id": <идентификатор выражения>, "status": <статус вычисления выражения> "result": <результат выражения>},{"id": <идентификатор выражения>, "status": <статус вычисления выражения> "result": <результат выражения>}]}`            | 200 OK                    |
| GET               | `/api/v1/expressions`  | N/A                              | `{"error": "Internal server error"}`               | 500 Internal Server Error | Entity |
| POST              | `/api/v1/expressions`  | N/A                              | `{ "error": "Method not allowed" }`                | 405 Method Not Allowed    |

## Сценарии использования /api/v1/expressions/:id

| **Request Method** | **Endpoint** | **Request Body**                                           | **Response Body**                                    | **HTTP Status Code** |
|--------------------|--------------|------------------------------------------------------------|------------------------------------------------------|----------------------|
| GET               | `/api/v1/expressions/:id`  | N/A                                | `{"expression": { "id": <идентификатор выражения>, "status": <статус вычисления выражения>, "result": <результат выражения>}}`                                    | 200 OK               |
| GET               | `/api/v1/expressions/:id`  | N/A                                | `{"error": "Internal server error"}`         | 500 Internal Server Error        |
| GET               | `/api/v1/expressions/:id`  | N/A                                | `{ "error": "Not Found" }`                   | 404 Not found                    |
| POST              | `/api/v1/expressions/:id`  | N/A                                | `{ "error": "Method not allowed" }`          | 405 Method Not Allowed           |

## Сценарии использования /internal/task

| **Request Method** | **Endpoint** | **Request Body**                                           | **Response Body**                                    | **HTTP Status Code** |
|--------------------|--------------|------------------------------------------------------------|------------------------------------------------------|----------------------|
| GET               | `/internal/task`  | N/A                                | `{"task":{"id": <идентификатор задачи>, "arg1": <имя первого аргумента>, "arg2": <имя второго аргумента>, "operation": <операция>,"operation_time": <время выполнения операции>}}`                                    | 200 OK               |
| GET               | `/internal/task`  | N/A                               | `{"error": "Internal server error"}`                 | 500 Internal Server Error |
| GET               | `/internal/task`  | N/A                          | `{ "error": "Not Found" }`                  | 404 Not found |
| POST               | `/internal/task`  | `{"id": 1, "result": 2.5}`   | `{"expressions":"OK"}`                  | 200 OK |
| POST               | `/internal/task`  | `{"id": 1, "result": 2.5}`   | `{ "error": "Not Found" }`                  | 404 Not found |
| POST               | `/internal/task`  | `{"id": 1, "result": 2.5}`   | `{ "error": "Invalid expression" }`                  | 422 Unprocessable Entity |
| POST               | `/internal/task`  | `{"id": 1, "result": 2.5}`   |  `{"error": "Internal server error"}`                  | 500 Internal Server Error |

## Коды ответов
- 200: Успешное вычисление
- 201: Выражение принято для вычисления
- 422: Ошибка в выражении (неверный формат) либо неправильно составленные данные запроса
- 404: Выражение/задача не найдена
- 405: Неверный метод запроса
- 500: Внутренняя ошибка сервера
## Тестирование

### Уточнения 

- тг ```@bll_nev_egor```
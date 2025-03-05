package main

import (
	"fmt"
	"strings"
)

// узел AST
type Node struct {
	Op    string
	Left  *Node
	Right *Node
}

// структура задачи
type Task struct {
	ID        string  `json:"id"`
	Op        string  `json:"operation"`
	Arg1      string  `json:"arg1"`
	Arg2      string  `json:"arg2"`
	Done      bool    `json:"-"`
	Result    *float64 `json:"-"`
}

// глобальный счетчик задач
var taskCounter = 1

// функция разбиения AST на задачи
func SplitTasks(node *Node, taskList *[]Task) string {
	if node == nil {
		return ""
	}

	// если это число — просто возвращаем его
	if node.Left == nil && node.Right == nil {
		return node.Op
	}

	// рекурсивно обрабатываем поддеревья
	leftVar := SplitTasks(node.Left, taskList)
	rightVar := SplitTasks(node.Right, taskList)

	// создаем уникальный идентификатор задачи cчетчиком
	taskID := fmt.Sprintf("id%d", taskCounter)
	taskCounter++

	// по сути это id данного узла, который будет являться аргументов узла выше
	taskID_child := fmt.Sprintf("id%d", taskCounter-1)

	// создаем задачу
	task := Task{
		ID:        taskID,
		Op:        node.Op,
		Arg1:      leftVar,
		Arg2:      rightVar,
	}

	// добавляем задачу в список
	*taskList = append(*taskList, task)

	return taskID_child
}

// функция парсинга выражения в AST
func ParseExpression(expression string) *Node {
	tokens := Tokenize(expression)
	pos := 0
	return ParseExpr(&tokens, &pos)
}

// функция токенизации (разбиение на числа, операторы, скобки)
func Tokenize(expression string) []string {
	expression = strings.ReplaceAll(expression, " ", "")
	tokens := []string{}
	num := ""
	// если цифра — добавляем к числу, если оператор — добавляем к токенам
	for _, char := range expression {
		switch {
		case char >= '0' && char <= '9':
			num += string(char)
		case char == '+' || char == '-' || char == '*' || char == '/' || char == '(' || char == ')':
			if num != "" {
				tokens = append(tokens, num)
				num = ""
			}
			tokens = append(tokens, string(char))
		}
	}
	if num != "" {
		tokens = append(tokens, num)
	}
	return tokens
}

// в дальнейших трех функция происходит тотальный вынос мозга
// если вы поняли как это работает, то вы молодец
// каждая функция представляет собой приоритет операций
// в конечном итоге получаем один узел, который связан указателями с другими узлами
// и по этим связям можно создавать задачи

// парсим выражение с учетом приоритетов и скобок
func ParseExpr(tokens *[]string, pos *int) *Node {
	node := ParseTerm(tokens, pos)

	for *pos < len(*tokens) {
		op := (*tokens)[*pos]
		if op != "+" && op != "-" {
			break
		}
		*pos++
		node = &Node{Op: op, Left: node, Right: ParseTerm(tokens, pos)}
	}
	return node
}

// учитываем умножение и деление
func ParseTerm(tokens *[]string, pos *int) *Node {
	node := ParseFactor(tokens, pos)

	for *pos < len(*tokens) {
		op := (*tokens)[*pos]
		if op != "*" && op != "/" {
			break
		}
		*pos++
		node = &Node{Op: op, Left: node, Right: ParseFactor(tokens, pos)}
	}
	return node
}

// числа и выражения в скобках
func ParseFactor(tokens *[]string, pos *int) *Node {
	token := (*tokens)[*pos]
	*pos++

	if token == "(" {
		node := ParseExpr(tokens, pos)
		*pos++ // пропускаем ')'
		return node
	}

	return &Node{Op: token}
}
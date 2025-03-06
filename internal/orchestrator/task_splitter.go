package orchestrator

import (

	"github.com/google/uuid"
)

// Функция разбиения AST на задачи
func SplitTasks(node *Node, taskList *[]Task) string {
	if node == nil {
		return ""
	}

	// Если это число — просто возвращаем его
	if node.Left == nil && node.Right == nil {
		return node.Op
	}

	// Рекурсивно обрабатываем поддеревья
	leftVar := SplitTasks(node.Left, taskList)
	rightVar := SplitTasks(node.Right, taskList)

	// Генерируем уникальный ID для текущей задачи
	taskID := uuid.New().String()

	// Создаем задачу
	task := Task{
		ID:   taskID,
		Op:   node.Op,
		Arg1: leftVar,
		Arg2: rightVar,
	}

	// Добавляем задачу в список
	*taskList = append(*taskList, task)

	// ID текущей задачи будет аргументом для родительского узла
	return taskID
}



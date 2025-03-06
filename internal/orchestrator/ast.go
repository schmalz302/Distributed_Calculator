package orchestrator

import (
	"strings"
)

// узел AST
type Node struct {
	Op    string
	Left  *Node
	Right *Node
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